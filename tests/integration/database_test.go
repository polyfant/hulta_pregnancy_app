package integration

import (
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/tests"
)

func TestDatabaseOperations(t *testing.T) {
	db := tests.SetupTestEnvironment(t)
	defer tests.TeardownTestEnvironment(db)

	// Test model creation
	horse := &models.Horse{
		Name: "Test Horse",
		Breed: "Thoroughbred",
	}

	result := db.Create(horse)
	if result.Error != nil {
		t.Fatalf("Failed to create horse: %v", result.Error)
	}

	// Verify creation
	var retrievedHorse models.Horse
	result = db.First(&retrievedHorse, horse.ID)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve horse: %v", result.Error)
	}

	if retrievedHorse.Name != horse.Name {
		t.Errorf("Expected name %s, got %s", horse.Name, retrievedHorse.Name)
	}
}
