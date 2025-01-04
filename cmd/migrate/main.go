package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
)

func main() {
	// Use new logger initialization
	if err := logger.InitLogger("./logs", "info"); err != nil {
		// Handle error
		panic(err)
	}

	// Source database (Java app)
	srcDB, err := sql.Open("sqlite3", filepath.Join("..", "..", "..", "horse_tracking_app", "horse_tracker.db"))
	if err != nil {
		logger.Error(err, "Failed to open source database")
		os.Exit(1)
	}
	defer srcDB.Close()

	// Destination database (Go app)
	destDB, err := sql.Open("sqlite3", "./horse_tracker.db")
	if err != nil {
		logger.Error(err, "Failed to open destination database")
		os.Exit(1)
	}
	defer destDB.Close()

	// Start migration
	if err := migrateData(srcDB, destDB); err != nil {
		logger.Error(err, "Migration failed")
		os.Exit(1)
	}

	logger.Info("Migration completed successfully")
}

func migrateData(srcDB, destDB *sql.DB) error {
	// Start a transaction
	tx, err := destDB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Migrate horses
	if err := migrateHorses(srcDB, tx); err != nil {
		return fmt.Errorf("failed to migrate horses: %w", err)
	}

	// Migrate health records
	if err := migrateHealthRecords(srcDB, tx); err != nil {
		return fmt.Errorf("failed to migrate health records: %w", err)
	}

	// Migrate pregnancy events
	if err := migratePregnancyEvents(srcDB, tx); err != nil {
		return fmt.Errorf("failed to migrate pregnancy events: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func migrateHorses(srcDB *sql.DB, tx *sql.Tx) error {
	rows, err := srcDB.Query(`SELECT id, name, breed, birth_date, conception_date, mother_id FROM horses`)
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := tx.Prepare(`
		INSERT INTO horses (id, name, breed, birth_date, conception_date, mother_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id, motherId                   sql.NullInt64
			name, breed                    string
			birthDate, conceptionDate      sql.NullString
			parsedBirth, parsedConception  time.Time
			birthErr, conceptionErr        error
		)

		if err := rows.Scan(&id, &name, &breed, &birthDate, &conceptionDate, &motherId); err != nil {
			return err
		}

		if birthDate.Valid {
			parsedBirth, birthErr = time.Parse("2006-01-02", birthDate.String)
			if birthErr != nil {
				logger.Error(birthErr, "Invalid birth date", map[string]interface{}{
					"horseID": id.Int64,
					"date":    birthDate.String,
				})
				continue
			}
		}

		if conceptionDate.Valid {
			parsedConception, conceptionErr = time.Parse("2006-01-02", conceptionDate.String)
			if conceptionErr != nil {
				logger.Error(conceptionErr, "Invalid conception date", map[string]interface{}{
					"horseID": id.Int64,
					"date":    conceptionDate.String,
				})
				continue
			}
		}

		_, err = stmt.Exec(
			id.Int64,
			name,
			breed,
			parsedBirth,
			parsedConception,
			motherId.Int64,
		)
		if err != nil {
			return err
		}
		count++
	}

	logger.Info("Migrated horses", map[string]interface{}{"count": count})
	return nil
}

func migrateHealthRecords(srcDB *sql.DB, tx *sql.Tx) error {
	rows, err := srcDB.Query(`SELECT id, horse_id, date, type, notes FROM health_records`)
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := tx.Prepare(`
		INSERT INTO health_records (id, horse_id, date, type, notes)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id, horseId          sql.NullInt64
			dateStr              sql.NullString
			recordType, notes    sql.NullString
			parsedDate           time.Time
			dateErr              error
		)

		if err := rows.Scan(&id, &horseId, &dateStr, &recordType, &notes); err != nil {
			return err
		}

		if dateStr.Valid {
			parsedDate, dateErr = time.Parse("2006-01-02", dateStr.String)
			if dateErr != nil {
				logger.Error(dateErr, "Invalid health record date", map[string]interface{}{
					"recordID": id.Int64,
					"date":     dateStr.String,
				})
				continue
			}
		}

		_, err = stmt.Exec(
			id.Int64,
			horseId.Int64,
			parsedDate,
			recordType.String,
			notes.String,
		)
		if err != nil {
			return err
		}
		count++
	}

	logger.Info("Migrated health records", map[string]interface{}{"count": count})
	return nil
}

func migratePregnancyEvents(srcDB *sql.DB, tx *sql.Tx) error {
	rows, err := srcDB.Query(`SELECT id, horse_id, date, description FROM pregnancy_events`)
	if err != nil {
		return err
	}
	defer rows.Close()

	stmt, err := tx.Prepare(`
		INSERT INTO pregnancy_events (id, horse_id, date, description)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id, horseId          sql.NullInt64
			dateStr              sql.NullString
			description          sql.NullString
			parsedDate           time.Time
			dateErr              error
		)

		if err := rows.Scan(&id, &horseId, &dateStr, &description); err != nil {
			return err
		}

		if dateStr.Valid {
			parsedDate, dateErr = time.Parse("2006-01-02", dateStr.String)
			if dateErr != nil {
				logger.Error(dateErr, "Invalid pregnancy event date", map[string]interface{}{
					"eventID": id.Int64,
					"date":    dateStr.String,
				})
				continue
			}
		}

		_, err = stmt.Exec(
			id.Int64,
			horseId.Int64,
			parsedDate,
			description.String,
		)
		if err != nil {
			return err
		}
		count++
	}

	logger.Info("Migrated pregnancy events", map[string]interface{}{"count": count})
	return nil
}
