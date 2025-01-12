package main

import (
	"fmt"
	"log"
	"os"

	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/breeding"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run() error {
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
	defer db.Close()

	// Initialize repositories
	horseRepo := repository.NewHorseRepository(db)
	userRepo := repository.NewUserRepository(db)
	pregnancyRepo := repository.NewPregnancyRepository(db)
	healthRepo := repository.NewHealthRepository(db)
	breedingRepo := repository.NewBreedingRepository(db)

	// Initialize cache
	cache := cache.NewMemoryCache()

	// Initialize services
	userService := service.NewUserService(userRepo)
	horseService := service.NewHorseService(horseRepo)
	pregnancyService := service.NewPregnancyService(horseRepo, pregnancyRepo)
	healthService := service.NewHealthService(healthRepo)
	breedingService := breeding.NewBreedingService(breedingRepo)

	// Initialize handler
	handler := api.NewHandler(api.HandlerConfig{
		Database:         db,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		BreedingService:  breedingService,
		Cache:           cache,
		HorseRepo:       horseRepo,
		BreedingRepo:    breedingRepo,
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	return handler.Start(port)
}
