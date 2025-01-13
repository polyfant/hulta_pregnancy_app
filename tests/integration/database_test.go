package integration

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestDatabaseConnection(t *testing.T) {
	// Skip if not in integration test environment
	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run.")
	}

	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		t.Skipf("Could not load .env file: %v", err)
	}

	cfg := config.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432, // Default postgres port
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Test query
	result := db.WithContext(nil).Exec("SELECT 1")
	if result.Error != nil {
		t.Fatalf("Failed to execute test query: %v", result.Error)
	}

	// Test migration
	if err := db.AutoMigrate(&models.Horse{}); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}
