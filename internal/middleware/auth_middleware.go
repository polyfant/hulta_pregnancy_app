package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"	
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type key int

const (
	UserIDKey   key = iota
	UserRolesKey
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Roles []string `json:"permissions,omitempty"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// Auth0Config contains configuration for Auth0
type Auth0Config struct {
	Domain   string
	Audience string
	Issuer   string
}

// AuthMiddleware creates a gin middleware for Auth0 authentication
func AuthMiddleware(config Auth0Config) gin.HandlerFunc {
	// Initialize JWKS provider
	issuerURL := fmt.Sprintf("https://%s/", config.Domain)
	log.Printf("Auth0 Config - Domain: %s, Audience: %s, IssuerURL: %s\n", 
		config.Domain, config.Audience, issuerURL)

	parsedURL, err := url.Parse(issuerURL)
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(parsedURL, 5*time.Minute)

	// Set up the validator
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL,
		[]string{config.Audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
	)

	if err != nil {
		log.Fatalf("Failed to set up the validator: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(handleError),
	)

	return func(c *gin.Context) {
		log.Printf("Processing request to: %s", c.Request.URL.Path)
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
			return
		}
		log.Printf("Authorization header found: %s", authHeader[:15]+"...")

		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			
			log.Printf("Token validated successfully. Subject: %s", claims.RegisteredClaims.Subject)
			
			// Set user info in context
			if claims != nil && claims.RegisteredClaims.Subject != "" {
				log.Printf("Setting user_id in context: %s", claims.RegisteredClaims.Subject)
				c.Set("user_id", claims.RegisteredClaims.Subject)
			} else {
				log.Printf("No subject claim found in token")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No subject claim in token"})
				return
			}
			if claims.CustomClaims != nil {
				customClaims, ok := claims.CustomClaims.(*CustomClaims)
				if ok {
					log.Printf("Setting user_roles in context: %+v", customClaims.Roles)
					c.Set(string(UserRolesKey), customClaims.Roles)
				}
			}
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)

		if encounteredError {
			log.Println("Error encountered during JWT validation")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles middleware checks if the user has any of the required roles
// TODO: Implement proper role checking. Currently allows all authenticated users.
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, just check if user is authenticated
		_, exists := c.Get(string(UserIDKey))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		// TODO: Implement proper role checking here
		// For now, allow all authenticated users
		log.Printf("Role check bypassed for development. Required roles were: %v", roles)
		c.Next()
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Auth error: %v\n", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"error": "Invalid or missing token: %v"}`, err)))
}