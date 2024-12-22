package api

import (
	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Horse routes
	horses := r.Group("/horses")
	{
		horses.GET("", handler.GetHorses)
		horses.POST("", handler.CreateHorse)
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

	return r
}

// SetupRoutes initializes the router with all routes and the database
func SetupRoutes(r *gin.Engine, db models.DataStore) {
	handler := NewHandler(db)
	
	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Horse routes
	horses := r.Group("/horses")
	{
		horses.GET("", handler.GetHorses)
		horses.POST("", handler.CreateHorse)
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
}
