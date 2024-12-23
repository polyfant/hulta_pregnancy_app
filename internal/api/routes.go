package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all the routes for the application
func SetupRouter(handler *Handler) *gin.Engine {
	router := gin.Default()

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Horse management endpoints
	horses := router.Group("/horses")
	{
		horses.GET("", handler.ListHorses)
		horses.POST("", handler.AddHorse)
		horses.GET("/:id", handler.GetHorse)

		// Health routes
		horses.GET("/:id/health", handler.GetHealthAssessment)
		horses.POST("/:id/health-records", handler.AddHealthRecord)

		// Pregnancy routes
		horses.GET("/:id/pregnancy-guidelines", handler.GetPregnancyGuidelines)

		// Breeding routes
		horses.GET("/:id/breeding-costs", handler.GetBreedingCosts)
		horses.POST("/:id/breeding-costs", handler.AddBreedingCost)
	}

	return router
}
