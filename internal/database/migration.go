package database

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

type MigrationConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
}

func getMigrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "migrations")
}

func RunMigrations(db *sql.DB) error {
	migrationsPath := getMigrationsPath()
	log.Printf("Looking for migrations in: %s", migrationsPath)

	// Set goose binary path
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %v", err)
	}

	// Run migrations
	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Printf("All migrations completed successfully")
	return nil
}
