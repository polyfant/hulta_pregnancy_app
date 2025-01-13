package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
)

type PostgresDB struct {
	*gorm.DB
}

func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting underlying database: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresDB{DB: db}, nil
}

func (db *PostgresDB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return fmt.Errorf("error getting underlying database: %w", err)
	}
	return sqlDB.Close()
}