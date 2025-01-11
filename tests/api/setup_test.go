package api_test

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository/postgres"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/breeding"
)

func setupTestData(db *database.PostgresDB) error {
	// Clean up existing test data
	if err := db.GetDB().Exec("DELETE FROM breeding_records WHERE horse_id = ?", 1).Error; err != nil {
		return fmt.Errorf("failed to clean up breeding records: %v", err)
	}
	if err := db.GetDB().Exec("DELETE FROM horses WHERE user_id = ?", "test_user_id").Error; err != nil {
		return fmt.Errorf("failed to clean up test data: %v", err)
	}

	// Create test horse
	horse := &models.Horse{
		ID:      1,
		UserID:  "test_user_id",
		Name:    "Test Horse",
		Breed:   "Test Breed",
		Gender:  models.GenderMare,
	}

	if err := db.GetDB().Create(horse).Error; err != nil {
		return fmt.Errorf("failed to create test horse: %v", err)
	}

	return nil
}

func setupTestRouter(db *database.PostgresDB) *gin.Engine {
	if err := setupTestData(db); err != nil {
		panic(fmt.Sprintf("Failed to set up test data: %v", err))
	}
	
	gin.SetMode(gin.TestMode)

	// Initialize repositories
	horseRepo := postgres.NewHorseRepository(db.GetDB())
	userRepo := postgres.NewUserRepository(db.GetDB())
	pregnancyRepo := postgres.NewPregnancyRepository(db.GetDB())
	healthRepo := postgres.NewHealthRepository(db.GetDB())
	breedingRepo := postgres.NewBreedingRepository(db.GetDB())

	// Initialize services
	userService := service.NewUserService(userRepo)
	horseService := service.NewHorseService(horseRepo)
	pregnancyService := service.NewPregnancyService(horseRepo, pregnancyRepo)
	healthService := service.NewHealthService(healthRepo)
	breedingService := breeding.NewBreedingService(breedingRepo, horseRepo)

	// Initialize handler config
	handlerConfig := api.HandlerConfig{
		Database:         db,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		Cache:           cache.NewMemoryCache(),
		HorseRepo:       horseRepo,
		BreedingRepo:    breedingRepo,
	}

	handler := api.NewHandler(handlerConfig)
	router := gin.New()
	
	// Add test middleware to bypass auth
	router.Use(func(c *gin.Context) {
		// Mock authenticated user
		c.Set("user_id", "test_user_id")
		c.Next()
	})

	api.SetupRouter(router, handler)

	return router
} 