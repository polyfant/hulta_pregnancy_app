package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"

	"github.com/polyfant/hulta_pregnancy_app/internal/repository/postgres"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

func main() {
	// Load configuration
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.LoadConfig()

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Path, cfg.Logger.Level); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Log the DSN we're about to use (mask password)
	logDsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
	)
	logger.Info("Attempting database connection", "dsn", logDsn)

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
			"dsn", logDsn)
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

	// Initialize handlers
	handlerConfig := api.HandlerConfig{
		Database:         db,
		UserService:      userService,
		HorseService:     horseService,
		PregnancyService: pregnancyService,
		HealthService:    healthService,
		Cache:           memoryCache,
		HorseRepo:       horseRepo,
		BreedingRepo:    breedingRepo,
	}

	handler := api.NewHandler(handlerConfig)

	// Setup Gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.ConnectionIsolationMiddleware())

	// Serve static files
	frontendBuildPath := filepath.Join(".", "frontend-react", "build")
	router.Use(middleware.StaticFileMiddleware(frontendBuildPath))

	// Setup routes
	api.SetupRouter(router, handler)

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
