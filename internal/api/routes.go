package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the routing for our API
func SetupRouter(router *gin.Engine, h *Handler) {
	// Add CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// Create API group
	api := router.Group("/api")
	{
		// Horse routes
		api.GET("/horses", h.ListHorses)
		api.POST("/horses", h.AddHorse)
		api.GET("/horses/:id", h.GetHorse)
		api.PUT("/horses/:id", h.UpdateHorse)
		api.DELETE("/horses/:id", h.DeleteHorse)

		// Health routes
		api.GET("/horses/:id/health", h.GetHealthRecords)
		api.POST("/horses/:id/health", h.AddHealthRecord)
		api.PUT("/horses/:id/health/:recordId", h.UpdateHealthRecord)
		api.DELETE("/horses/:id/health/:recordId", h.DeleteHealthRecord)

		// Pregnancy routes
		api.GET("/horses/:id/pregnancy", h.GetPregnancy)
		api.POST("/horses/:id/pregnancy/start", h.StartPregnancyTracking)
		api.GET("/horses/:id/pregnancy/status", h.GetPregnancyStatus)
		api.GET("/horses/:id/pregnancy/events", h.GetPregnancyEvents)
		api.POST("/horses/:id/pregnancy/events", h.AddPregnancyEvent)
		api.GET("/horses/:id/pregnancy/guidelines", h.GetPregnancyGuidelines)

		// Breeding routes
		api.GET("/horses/:id/breeding", h.GetBreedingRecords)
		api.POST("/horses/:id/breeding", h.AddBreedingRecord)
		api.PUT("/horses/:id/breeding/:recordId", h.UpdateBreedingRecord)
		api.DELETE("/horses/:id/breeding/:recordId", h.DeleteBreedingRecord)

		// User routes
		api.GET("/user/profile", h.GetUserProfile)
		api.PUT("/user/profile", h.UpdateUserProfile)

		// Dashboard route
		api.GET("/dashboard", h.GetDashboardStats)
	}
}
