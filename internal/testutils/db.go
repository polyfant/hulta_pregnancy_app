package testutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	// Load test environment variables
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Logf("Warning: .env.test file not found, using default test values")
	}

	// Get database connection parameters from environment variables
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	port := getEnvOrDefault("TEST_DB_PORT", "5432")
	user := getEnvOrDefault("TEST_DB_USER", "postgres")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("TEST_DB_NAME", "horse_tracking_test")

	// Create database connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up and migrate tables
	tables := []interface{}{
		&models.PrivacyPreferences{},
		&models.WeatherData{},
		&models.HealthRecord{},
		&models.PregnancyEvent{},
	}

	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
		if err := db.AutoMigrate(table); err != nil {
			t.Fatalf("Failed to migrate table: %v", err)
		}
	}

	return db
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
