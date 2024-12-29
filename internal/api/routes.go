package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/service/pregnancy"
	"github.com/polyfant/horse_tracking/internal/models"
)

// SetupRouter sets up the routing for our API
func SetupRouter(h *Handler, store models.DataStore) *gin.Engine {
	router := gin.Default()

	// Add CORS middleware with more permissive settings for development
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow all origins in development
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	// Create pregnancy handler
	pregnancyHandler := NewPregnancyHandler(pregnancy.NewService(store), store)

	// API routes
	api := router.Group("/api")
	{
		// Dashboard route
		api.GET("/dashboard", h.GetDashboardStats)

		// Horse routes
		api.GET("/horses", h.ListHorses)
		api.POST("/horses", h.AddHorse)
		api.GET("/horses/:id", h.GetHorse)
		api.DELETE("/horses/:id", h.DeleteHorse)
		api.GET("/horses/:id/offspring", h.GetHorseOffspring)
		api.GET("/horses/:id/family", h.GetHorseFamilyTree)

		// Health routes
		api.GET("/horses/:id/health", h.GetHealthAssessment)
		api.POST("/horses/:id/health", h.AddHealthRecord)

		// Breeding cost routes
		api.GET("/horses/:id/breeding-costs", h.GetBreedingCosts)
		api.POST("/horses/:id/breeding-costs", h.AddBreedingCost)

		// Pregnancy routes
		api.GET("/horses/:id/pregnancy", pregnancyHandler.GetPregnancyStatus)
		api.POST("/horses/:id/pregnancy/start", pregnancyHandler.StartPregnancyTracking)
		api.POST("/horses/:id/pregnancy/end", pregnancyHandler.EndPregnancyTracking)
		api.GET("/horses/:id/pregnancy/events", pregnancyHandler.GetPregnancyEvents)
		api.POST("/horses/:id/pregnancy/events", pregnancyHandler.AddPregnancyEvent)
		api.GET("/horses/:id/pregnancy/guidelines", pregnancyHandler.GetPregnancyGuidelines)
		api.GET("/horses/:id/pregnancy/foaling-signs", pregnancyHandler.CheckPreFoalingSigns)
		api.POST("/horses/:id/pregnancy/foaling-signs", pregnancyHandler.RecordPreFoalingSign)
		api.GET("/horses/:id/pregnancy/foaling-checklist", pregnancyHandler.GetFoalingChecklist)
		api.GET("/horses/:id/pregnancy/post-foaling-checklist", pregnancyHandler.GetPostFoalingChecklist)
	}

	return router
}
