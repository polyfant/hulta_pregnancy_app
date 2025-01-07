package api

import (
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

// HandlerConfig defines the configuration for creating a new handler
type HandlerConfig struct {
	Database         *database.PostgresDB
	UserService      *service.UserService
	HorseService     *service.HorseService
	PregnancyService *service.PregnancyService
	HealthService    *service.HealthService
	Cache            *cache.MemoryCache
	HorseRepo        repository.HorseRepository
	BreedingRepo     repository.BreedingRepository
} 