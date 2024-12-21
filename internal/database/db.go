package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./horse_tracker.db")
	if err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS horses (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			breed TEXT NOT NULL,
			birth_date DATETIME,
			conception_date DATETIME,
			mother_id INTEGER,
			father_id INTEGER,
			FOREIGN KEY(mother_id) REFERENCES horses(id),
			FOREIGN KEY(father_id) REFERENCES horses(id)
		)`,
		`CREATE TABLE IF NOT EXISTS health_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			horse_id INTEGER,
			date DATETIME NOT NULL,
			type TEXT NOT NULL,
			notes TEXT,
			FOREIGN KEY(horse_id) REFERENCES horses(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pregnancy_events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			horse_id INTEGER,
			date DATETIME NOT NULL,
			description TEXT NOT NULL,
			FOREIGN KEY(horse_id) REFERENCES horses(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_horse_name ON horses(name)`,
		`CREATE INDEX IF NOT EXISTS idx_horse_breed ON horses(breed)`,
		`CREATE INDEX IF NOT EXISTS idx_horse_conception_date ON horses(conception_date)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}
