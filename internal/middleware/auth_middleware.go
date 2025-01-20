package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"os"
)

type key int

const (
	UserIDKey   key = iota
	UserRolesKey
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	jwt.StandardClaims
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}

// Validate satisfies validator.CustomClaims interface
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// Auth0Config contains configuration for Auth0
type Auth0Config struct {
	Domain   string
	Audience string
	Issuer   string
}

// User represents a system user
type User struct {
	ID            string
	IsFarmManager bool
}

// Predefined roles matching frontend
const (
	RoleUser         = "user"
	RoleAdmin        = "admin"
	RoleOwner        = "owner"
	RoleFarmManager  = "farm_manager"
)

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

		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			log.Printf("Token validated successfully. Subject: %s", claims.RegisteredClaims.Subject)

			// Set user info in context
			if claims != nil && claims.RegisteredClaims.Subject != "" {
				log.Printf("Setting user_id in context: %s", claims.RegisteredClaims.Subject)
				c.Set(fmt.Sprintf("user_id_%d", UserIDKey), claims.RegisteredClaims.Subject)
			}

			// Extract custom claims
			if claims.CustomClaims != nil {
				customClaims, ok := claims.CustomClaims.(*CustomClaims)
				if ok {
					// Default to user role if not specified
					role := RoleUser
					if customClaims.Role == RoleAdmin {
						role = RoleAdmin
					}
					
					log.Printf("Setting user_role in context: %s", role)
					c.Set(fmt.Sprintf("user_roles_%d", UserRolesKey), role)
				}
			}
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)
		c.Next()
	}
}

// AdminMiddleware restricts access to admin-only routes
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleKey := fmt.Sprintf("user_roles_%d", UserRolesKey)
		userRole, exists := c.Get(userRoleKey)
		
		if !exists || userRole != RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// Example route setup
func setupRoutes(r *gin.Engine) {
	// Public routes
	r.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Public endpoint"})
	})

	// User routes (require authentication)
	r.GET("/horses", AuthMiddleware(Auth0Config{}))

	// Admin-only routes
	r.GET("/admin/dashboard", AuthMiddleware(Auth0Config{}), AdminMiddleware())
	r.POST("/admin/settings", AuthMiddleware(Auth0Config{}), AdminMiddleware())
}

// Simplified token generation
func generateJWTToken(user User) (string, error) {
	role := determineUserRole(user)

	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   user.ID,
		},
		Role:   role,
		UserID: user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Determine user role based on account type and permissions
func determineUserRole(user User) string {
	switch {
	case user.IsFarmManager:
		return RoleFarmManager
	default:
		return RoleOwner
	}
}

// RoleMiddleware provides flexible role-based access control
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Validate token and extract claims
		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// If no roles specified, just validate token
		if len(requiredRoles) == 0 {
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
			c.Next()
			return
		}

		// Check role against required roles
		for _, role := range requiredRoles {
			if claims.Role == role {
				c.Set("user_id", claims.UserID)
				c.Set("user_role", claims.Role)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":           "Insufficient permissions",
			"required_roles":  requiredRoles,
			"user_role":       claims.Role,
		})
		c.Abort()
	}
}

// Extract token from Authorization header
func extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}

	// Bearer token format
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Auth error: %v\n", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf(`{"error": "Invalid or missing token: %v"}`, err)))
}