package api

import (
	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

// HandlerConfig defines the configuration for creating a new handler
type HandlerConfig struct {
	Database         *database.PostgresDB
	UserService      *service.UserService
	HorseService     *service.HorseService
	PregnancyService *service.PregnancyService
	HealthService    *service.HealthService
	Cache           *cache.MemoryCache
}

// Handler interface defines all the HTTP handlers
type Handler interface {
	// Horse methods
	ListHorses(c *gin.Context)
	GetHorse(c *gin.Context)
	AddHorse(c *gin.Context)
	DeleteHorse(c *gin.Context)
	
	// Pregnancy methods
	GetPregnancies(c *gin.Context)
	GetPregnancyStage(c *gin.Context)
	GetPregnancyStatus(c *gin.Context)
	UpdateHorsePregnancyStatus(c *gin.Context)
	GetPregnancyEvents(c *gin.Context)
	AddPregnancyEvent(c *gin.Context)
	
	// Health methods
	GetHealthAssessment(c *gin.Context)
	AddHealthRecord(c *gin.Context)
} 