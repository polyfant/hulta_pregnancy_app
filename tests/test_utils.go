package tests

import (
	"fmt"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func SetupTestEnvironment(t *testing.T) *database.PostgresDB {
	cfg := config.LoadTestConfig()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, 
		cfg.Database.User, cfg.Database.Password, 
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := database.NewPostgresDB(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.Horse{},
		&models.HealthRecord{},
		&models.Pregnancy{},
		&models.PregnancyEvent{},
	); err != nil {
		t.Fatalf("Failed to auto-migrate test database: %v", err)
	}

	return db
}

func TeardownTestEnvironment(db *database.PostgresDB) {
	if err := db.Close(); err != nil {
		// Log or handle the error as needed
	}
}
