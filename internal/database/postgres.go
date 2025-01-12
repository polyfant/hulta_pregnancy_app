package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(cfg config.DatabaseConfig) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

func (db *PostgresDB) Close() error {
	return db.db.Close()
}