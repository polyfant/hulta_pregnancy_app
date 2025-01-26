package wire

import (
	"github.com/google/wire"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"gorm.io/gorm"
)

// ProvideHandlerConfig sets up the configuration for the API handler
func ProvideHandlerConfig(
	db *gorm.DB,
	userService service.UserService,
	horseService service.HorseService,
	pregnancyService service.PregnancyService,
	healthService service.HealthService,
	breedingService service.BreedingService,
	growthService service.GrowthService,
	cacheService cache.Cache,
	horseRepo repository.HorseRepository,
	breedingRepo repository.BreedingRepository,
	growthRepo repository.GrowthRepository,
	auth0Config config.Auth0Config,
) api.HandlerConfig {
	return api.HandlerConfig{
		Database:         db,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		GrowthService:    growthService,
		Cache:            cacheService,
		HorseRepo:        horseRepo,
		BreedingRepo:     breedingRepo,
		GrowthRepo:       growthRepo,
		Auth0:            auth0Config,
	}
}

// ProvideGrowthService sets up the growth service
func ProvideGrowthService(
	growthRepo repository.GrowthRepository,
	horseRepo repository.HorseRepository,
) service.GrowthService {
	return service.NewGrowthService(growthRepo, horseRepo)
}

// WireSet for API dependencies
var WireSet = wire.NewSet(
	ProvideHandlerConfig,
	ProvideGrowthService,
	api.NewHandler,
	api.NewGrowthHandler,
)
