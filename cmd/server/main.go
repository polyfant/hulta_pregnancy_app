package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

func main() {
	logger.Init()

	// Get database configuration from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPass := getEnvOrDefault("DB_PASSWORD", "your_password_here")
	dbName := getEnvOrDefault("DB_NAME", "HE_horse_db")

	// Construct database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	// Initialize database
	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		logger.Error(err, "Failed to initialize database")
		os.Exit(1)
	}

	// Initialize API handlers
	handler := api.NewHandler(db)

	// Set up router
	router := gin.Default()
	api.SetupRouter(handler, db)

	// Start server
	port := getEnvOrDefault("PORT", "8080")
	if err := router.Run(":" + port); err != nil {
		logger.Error(err, "Failed to start server")
		os.Exit(1)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
