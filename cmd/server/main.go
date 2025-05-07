package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
	// Set Gin to release mode in production
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying *sql.DB for migrations
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("error getting underlying *sql.DB: %w", err)
	}
	defer sqlDB.Close()

	// Run database migrations
	if err := database.RunMigrations(sqlDB); err != nil {
		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	// Initialize repositories
	horseRepo := repository.NewHorseRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)
	pregnancyRepo := repository.NewPregnancyRepository(db.DB)
	healthRepo := repository.NewHealthRepository(db.DB)
	breedingRepo := repository.NewBreedingRepository(db.DB)

	// Initialize cache
	cache := cache.NewMemoryCache()

	// Initialize services
	userService := service.NewUserService(userRepo)
	horseService := service.NewHorseService(horseRepo)
	pregnancyService := service.NewPregnancyService(horseRepo, pregnancyRepo)
	healthService := service.NewHealthService(healthRepo)
	breedingService := service.NewBreedingService(breedingRepo)

	// Initialize growth service
	growthRepo := repository.NewGrowthRepository(db.DB)
	growthService := service.NewGrowthService(growthRepo, horseRepo)

	// Initialize expense service
	expenseRepo := repository.NewExpenseRepository(db.DB)
	recurringExpenseRepo := repository.NewRecurringExpenseRepository(db.DB)
	expenseService := service.NewExpenseService(expenseRepo, recurringExpenseRepo, horseRepo)

	// Initialize handler
	handler := api.NewHandler(api.HandlerConfig{
		Database:         db.DB,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		GrowthService:    growthService,
		ExpenseService:   expenseService,
		Cache:           cache,
		HorseRepo:       horseRepo,
		BreedingRepo:    breedingRepo,
		GrowthRepo:      growthRepo,
		Auth0:           cfg.Auth0,
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	return handler.Start(port)
}
