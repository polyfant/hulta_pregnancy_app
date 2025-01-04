package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	gin.SetMode(cfg.Server.Mode)

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Path, cfg.Logger.Level); err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	// Initialize database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, 
		cfg.Database.User, cfg.Database.Password, 
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		logger.Error(err, "Failed to connect to database")
		os.Exit(1)
	}

	// Optional: Create database users
	if cfg.Features.StrictMode {
		_, err = database.CreateDatabaseUsers(db.DB())
		if err != nil {
			logger.Error(err, "Failed to create database users")
			os.Exit(1)
		}
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Horse{},
		&models.HealthRecord{},
		&models.Pregnancy{},
		&models.PregnancyEvent{},
		&models.PreFoalingSign{},
		&models.PreFoalingChecklistItem{},
		&models.Expense{},
		&models.RecurringExpense{},
		&models.BreedingRecord{},
	); err != nil {
		logger.Error(err, "Failed to auto-migrate database")
		os.Exit(1)
	}

	// Optional: Initialize memory cache if enabled
	var memoryCache *cache.MemoryCache
	if cfg.Features.EnableCaching {
		memoryCache = cache.NewMemoryCache()
		logger.Info("Memory cache initialized")
	} else {
		memoryCache = nil
	}

	// Setup Gin router with security middleware
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ConnectionIsolationMiddleware())

	// Determine frontend build path
	frontendBuildPath := filepath.Join(".", "frontend-react", "build")
	router.Use(middleware.StaticFileMiddleware(frontendBuildPath))

	// Initialize API handlers
	handler := api.NewHandler(db, memoryCache)

	// Setup routes
	api.SetupRouter(router, handler)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info(fmt.Sprintf("Starting server on %s", serverAddr))
	
	if err := router.Run(serverAddr); err != nil {
		logger.Error(err, "Server failed to start")
		os.Exit(1)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
