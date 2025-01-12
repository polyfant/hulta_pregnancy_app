package api_test

import (
	"context"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/mocks"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/breeding"
)

// Mock repositories for testing
var (
	mockHorseRepo     *mocks.MockHorseRepository
	mockPregnancyRepo *mocks.MockPregnancyRepository
	mockUserRepo      *mocks.MockUserRepository
	mockHealthRepo    *mocks.MockHealthRepository
	mockBreedingRepo  *mocks.MockBreedingRepository
)

func setupTestHandler(t *testing.T) *api.Handler {
	mockDB := &database.PostgresDB{}
	mockCache := cache.NewMemoryCache()

	// Initialize mock repositories with testify/mock
	mockHorseRepo = new(mocks.MockHorseRepository)
	mockPregnancyRepo = new(mocks.MockPregnancyRepository)
	mockUserRepo = new(mocks.MockUserRepository)
	mockHealthRepo = new(mocks.MockHealthRepository)
	mockBreedingRepo = new(mocks.MockBreedingRepository)

	// Initialize services
	horseService := service.NewHorseService(mockHorseRepo)
	userService := service.NewUserService(mockUserRepo)
	pregnancyService := service.NewPregnancyService(mockPregnancyRepo)
	healthService := service.NewHealthService(mockHealthRepo)
	breedingService := breeding.NewBreedingService(mockBreedingRepo)

	return api.NewHandler(api.HandlerConfig{
		Database:         mockDB,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		Cache:            mockCache,
		HorseRepo:        mockHorseRepo,
		BreedingRepo:     mockBreedingRepo,
	})
}

func setupTestContext(t *testing.T) context.Context {
	return context.Background()
} 