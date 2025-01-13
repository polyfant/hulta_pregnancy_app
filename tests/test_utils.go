package tests

import (
	"os"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
)

func SetupTestDB(t *testing.T) *database.PostgresDB {
	// Set test environment variables
	os.Setenv("DB_HOST", "172.31.112.63")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "jonas")
	os.Setenv("DB_PASSWORD", "infernal")
	os.Setenv("DB_NAME", "jonas")
	os.Setenv("APP_ENV", "test")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	return db
}

func TeardownTestEnvironment(db *database.PostgresDB) {
	if err := db.Close(); err != nil {
		// Log or handle the error as needed
	}
}
