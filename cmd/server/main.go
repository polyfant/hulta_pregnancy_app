package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"path/filepath"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Path, cfg.Logger.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Construct database DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, 
		cfg.Database.Port, 
		cfg.Database.User, 
		cfg.Database.Password, 
		cfg.Database.DBName,
	)

	// Initialize database connection
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", 
			"error", err, 
			"dsn", dsn)
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

	// Initialize cache
	memoryCache := cache.NewMemoryCache()

	// Configure database backup if enabled
	if cfg.Backup.Enabled {
		backupUtil := database.NewDatabaseBackup(dsn, cfg.Backup.Directory)
		
		// Schedule backups
		backupUtil.ScheduleBackups(cfg.Backup.Interval)
		
		// Manage backup retention
		go func() {
			for {
				if err := backupUtil.ManageBackupRetention(cfg.Backup.MaxBackups); err != nil {
					logger.Error(err, "Backup retention management failed")
				}
				time.Sleep(24 * time.Hour)
			}
		}()

		logger.Info("Database backup service initialized", 
			"backup_dir", cfg.Backup.Directory, 
			"interval", cfg.Backup.Interval)
	}

	// Initialize API handlers
	handlers := api.NewHandler(db, memoryCache)

	// Setup Gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.ConnectionIsolationMiddleware())

	// Serve static files
	frontendBuildPath := filepath.Join(".", "frontend-react", "build")
	router.Use(middleware.StaticFileMiddleware(frontendBuildPath))

	// Setup routes
	api.SetupRouter(router, handlers)

	// Start server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	go func() {
		if err := router.Run(serverAddr); err != nil {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	logger.Info("Server started", 
		"address", serverAddr, 
		"environment", os.Getenv("APP_ENV"))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	// Perform any cleanup here
}
