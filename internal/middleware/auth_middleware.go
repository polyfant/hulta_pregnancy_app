package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type key int

const (
	UserIDKey key = iota
	UserRolesKey
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Scope string `json:"scope"`
	Roles []string `json:"https://hulta-pregnancy.com/roles"`
}

// Validate implements the validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	if len(c.Roles) == 0 {
		return errors.New("missing required roles claim")
	}
	return nil
}

// Auth0Config holds the configuration for Auth0
type Auth0Config struct {
	Domain    string
	Audience  string
	Issuer    string
}

// AuthMiddleware creates the Auth0 middleware
func AuthMiddleware(config Auth0Config) gin.HandlerFunc {
	// Initialize JWKS provider
	issuerURL, err := url.Parse(fmt.Sprintf("https://%s/", config.Domain))
	if err != nil {
		panic(err)
	}

	provider := jwks.NewCachingProvider(issuerURL, time.Duration(5)*time.Minute)

	// Create the validator
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		config.Issuer,
		[]string{config.Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
	)
	if err != nil {
		panic(err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(handleError),
	)

	return func(c *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			// Extract claims
			claims := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			customClaims := claims.CustomClaims.(*CustomClaims)
			
			// Set user info in context
			ctx := context.WithValue(r.Context(), UserRolesKey, customClaims.Roles)
			// Extract sub claim for user ID
			if claims.RegisteredClaims.Subject != "" {
				ctx = context.WithValue(ctx, UserIDKey, claims.RegisteredClaims.Subject)
			}
			
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)

		if encounteredError {
			c.Abort()
		}
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error": "invalid token"}`))
}

// GetUserID extracts the user ID from the context
func GetUserID(ctx context.Context) (string, error) {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID, nil
	}
	return "", errors.New("no user ID found in context")
}

// GetUserRoles extracts the user roles from the context
func GetUserRoles(ctx context.Context) ([]string, error) {
	if roles, ok := ctx.Value(UserRolesKey).([]string); ok {
		return roles, nil
	}
	return nil, errors.New("no roles found in context")
}

// RequireRoles middleware checks if the user has the required roles
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, err := GetUserRoles(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: no roles found"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if strings.EqualFold(requiredRole, userRole) {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}