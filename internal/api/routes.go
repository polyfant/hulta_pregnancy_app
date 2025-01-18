package api

import (
	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
)

// RegisterRoutes sets up all API routes
func RegisterRoutes(router *gin.Engine, handlers *Handlers) {
	// Apply global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// API routes group
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", handlers.HealthHandler.HandleHealthCheck)

		// Pregnancy routes
		pregnancy := api.Group("/pregnancy")
		{
			pregnancy.POST("", handlers.PregnancyHandler.HandleCreatePregnancy)
			pregnancy.GET("/:id", handlers.PregnancyHandler.HandleGetPregnancy)
			pregnancy.PUT("/:id", handlers.PregnancyHandler.HandleUpdatePregnancy)
			pregnancy.DELETE("/:id", handlers.PregnancyHandler.HandleDeletePregnancy)
			pregnancy.GET("", handlers.PregnancyHandler.HandleListPregnancies)
		}

		// Vital signs routes
		vitals := api.Group("/vitals")
		{
			vitals.POST("", handlers.VitalsHandler.HandleRecordVitalSigns)
			vitals.GET("/horse/:id", handlers.VitalsHandler.HandleGetVitalSigns)
			vitals.GET("/horse/:id/latest", handlers.VitalsHandler.HandleGetLatestVitalSigns)
			vitals.GET("/horse/:id/alerts", handlers.VitalsHandler.HandleGetAlerts)
			vitals.GET("/alerts/:id", handlers.VitalsHandler.HandleGetAlert)
			vitals.PUT("/alerts/:id/acknowledge", handlers.VitalsHandler.HandleAcknowledgeAlert)
			vitals.GET("/horse/:id/trends", handlers.VitalsHandler.HandleGetTrends)
		}

		// Weather routes
		weather := api.Group("/weather")
		{
			weather.GET("/current", handlers.WeatherHandler.HandleGetCurrentWeather)
			weather.GET("/recommendations", handlers.WeatherHandler.HandleGetRecommendations)
		}

		// Feedback routes
		feedback := api.Group("/feedback")
		{
			feedback.POST("/feature", handlers.FeedbackHandler.HandleFeatureRequest)
			feedback.GET("/features", handlers.FeedbackHandler.HandleListFeatures)
			feedback.POST("/features/:id/vote", handlers.FeedbackHandler.HandleFeatureVote)
			feedback.GET("/votes", handlers.FeedbackHandler.HandleGetUserVotes)
		}

		// Privacy routes
		privacy := api.Group("/privacy")
		{
			privacy.GET("/preferences", handlers.PrivacyHandler.HandleGetPreferences)
			privacy.PUT("/preferences", handlers.PrivacyHandler.HandleUpdatePreferences)
		}
	}

	// WebSocket endpoint
	router.GET("/ws", gin.WrapF(handlers.WebSocketHub.ServeHTTP))
}
