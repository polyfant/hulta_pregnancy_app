package database

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
}

func RunMigrations(cfg MigrationConfig) error {
	migrationURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", 
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Schema)
	
	m, err := migrate.New(
		"file://internal/database/migrations", 
		migrationURL,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
