package integration

import (
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/config"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestDatabaseConnection(t *testing.T) {
	cfg := config.DatabaseConfig{
		Host:     "172.31.112.63",
		Port:     5432,
		User:     "jonas",
		Password: "xxxxx",
		DBName:   "jonas",
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
