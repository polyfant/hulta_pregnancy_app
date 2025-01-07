package integration

import (
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

var testDSN = "host=172.31.112.63 port=5432 user=jonas dbname=jonas password=infernal sslmode=disable"

func TestDatabaseConnection(t *testing.T) {
	db, err := database.NewPostgresDB(testDSN)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Test query
	result := db.DB.Exec("SELECT 1")
	if result.Error != nil {
		t.Fatalf("Failed to execute test query: %v", result.Error)
	}

	// Test migration
	err = db.AutoMigrate(&models.Horse{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}
