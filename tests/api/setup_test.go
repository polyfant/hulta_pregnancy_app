package api_test

import (
	"context"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/mocks"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/breeding"
)

// setupTestHandler creates a new API handler with mock repositories for testing.
// It returns the handler and all mock repositories to allow test cases to set expectations.
func setupTestHandler() (*api.Handler, *mocks.MockHorseRepository, *mocks.MockUserRepository, *mocks.PregnancyRepository, *mocks.MockHealthRepository, *mocks.MockBreedingRepository) {
	mockHorseRepo := new(mocks.MockHorseRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockPregnancyRepo := new(mocks.PregnancyRepository)
	mockHealthRepo := new(mocks.MockHealthRepository)
	mockBreedingRepo := new(mocks.MockBreedingRepository)

	// Initialize services with mock repositories
	userService := service.NewUserService(mockUserRepo)
	pregnancyService := service.NewPregnancyService(mockHorseRepo, mockPregnancyRepo)
	healthService := service.NewHealthService(mockHealthRepo)
	breedingService := breeding.NewBreedingService(mockBreedingRepo)
	horseService := service.NewHorseService(mockHorseRepo)
	// Initialize cache
	cache := cache.NewMemoryCache()

	// Create and return the handler
	handler := api.NewHandler(api.HandlerConfig{
		UserService:      userService,
		HorseService:    horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		Cache:            cache,
		HorseRepo:        mockHorseRepo,
		BreedingRepo:     mockBreedingRepo,
		Auth0:            config.Auth0Config{Domain: "test"},
	})
	return handler, mockHorseRepo, mockUserRepo, mockPregnancyRepo, mockHealthRepo, mockBreedingRepo
}

// setupTestContext creates a new context for testing.
// This can be extended in the future to add test-specific context values if needed.
func setupTestContext(_ *testing.T) context.Context {
    return context.Background()
}