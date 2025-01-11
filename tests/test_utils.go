package tests

import (
	"fmt"
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

	cfg := config.LoadConfig()
	
	// Construct test DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	db, err := database.NewPostgresDB(dsn)
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
