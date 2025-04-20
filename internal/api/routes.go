package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
)

func setupCORS() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowCredentials = true
	return config
}

func testBypassAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", "test_user")
		c.Next()
	}
}

func SetupRouter(router *gin.Engine, h *Handler) *gin.Engine {
	// Set trusted proxies
	router.SetTrustedProxies([]string{"127.0.0.1"})

	// CORS middleware
	router.Use(cors.New(setupCORS()))

	// Rate limiting middleware
	rateLimitConfig := middleware.RateLimitConfig{
		RequestsPerSecond: 10, // 10 requests per second
		BurstLimit:        30, // Allow bursts up to 30 requests
	}
	router.Use(middleware.RateLimitMiddleware(rateLimitConfig))

	var auth gin.HandlerFunc
	if h.config.Auth0.Domain == "test" {
		auth = testBypassAuthMiddleware()
	} else {
		auth = middleware.AuthMiddleware(middleware.Auth0Config{
			Domain:   h.config.Auth0.Domain,
			Audience: h.config.Auth0.Audience,
			Issuer:   h.config.Auth0.Issuer,
		})
	}

	// Public routes
	public := router.Group("/api/v1")
	{
		public.GET("/health", HealthCheck)
		public.GET("/version", Version)
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(auth)
	{
		// User routes
		protected.GET("/user/profile", h.GetUserProfile)
		protected.PUT("/user/profile", h.UpdateUserProfile)

		// Horse routes
		protected.GET("/horses", h.ListHorses)
		protected.POST("/horses", h.AddHorse)
		protected.GET("/horses/:id", h.GetHorse)
		protected.PUT("/horses/:id", h.UpdateHorse)
		protected.DELETE("/horses/:id", h.DeleteHorse)

		// Health routes
		protected.GET("/horses/:id/health", h.GetHealthRecords)
		protected.POST("/horses/:id/health", h.AddHealthRecord)
		protected.PUT("/horses/:id/health/:recordId", h.UpdateHealthRecord)
		protected.DELETE("/horses/:id/health/:recordId", h.DeleteHealthRecord)

		// Pregnancy routes
		protected.GET("/horses/:id/pregnancy", h.GetPregnancy)
		protected.POST("/horses/:id/pregnancy/start", h.StartPregnancyTracking)
		protected.GET("/horses/:id/pregnancy/status", h.GetPregnancyStatus)
		protected.GET("/horses/:id/pregnancy/events", h.GetPregnancyEvents)
		protected.POST("/horses/:id/pregnancy/events", h.AddPregnancyEvent)
		protected.GET("/horses/:id/pregnancy/guidelines", h.GetPregnancyGuidelines)

		// Breeding routes
		protected.GET("/horses/:id/breeding", h.GetBreedingRecords)
		protected.POST("/horses/:id/breeding", h.AddBreedingRecord)
		protected.PUT("/horses/:id/breeding/:recordId", h.UpdateBreedingRecord)
		protected.DELETE("/horses/:id/breeding/:recordId", h.DeleteBreedingRecord)

		// Growth routes
		protected.POST("/horses/:id/growth", h.growthHandler.RecordGrowthMeasurement)
		protected.GET("/horses/:id/growth", h.growthHandler.GetFoalGrowthData)
		protected.GET("/horses/:id/growth/analysis", h.growthHandler.AnalyzeGrowthTrends)

		// Dashboard route
		protected.GET("/dashboard", h.GetDashboardStats)
	}

	return router
}

// RequireRoles middleware wrapper for gin
func RequireRoles(roles ...string) gin.HandlerFunc {
	return middleware.RoleMiddleware(roles...)
}

// HealthCheck handles GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// Version handles GET /version
func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": "1.0.0", // You can make this dynamic by using build flags or environment variables
	})
}
