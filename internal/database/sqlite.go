package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/polyfant/horse_tracking/internal/logger"
	"github.com/polyfant/horse_tracking/internal/models"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// Horse operations
func (s *SQLiteStore) GetHorse(id int64) (*models.Horse, error) {
	var horse models.Horse
	var birthDateStr, conceptionDateStr sql.NullString

	err := s.db.QueryRow(`
		SELECT id, name, breed, birth_date, conception_date, mother_id, father_id
		FROM horses WHERE id = ?`, id).Scan(
		&horse.ID, &horse.Name, &horse.Breed,
		&birthDateStr, &conceptionDateStr,
		&horse.MotherID, &horse.FatherID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get horse: %w", err)
	}

	if birthDateStr.Valid {
		date, err := time.Parse("2006-01-02", birthDateStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid birth date format: %w", err)
		}
		horse.DateOfBirth = date
	}

	if conceptionDateStr.Valid {
		date, err := time.Parse("2006-01-02", conceptionDateStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid conception date format: %w", err)
		}
		horse.ConceptionDate = &date
	}

	return &horse, nil
}

func (s *SQLiteStore) GetHorseByName(name string) (*models.Horse, error) {
	var horse models.Horse
	var birthDateStr, conceptionDateStr sql.NullString

	err := s.db.QueryRow(`
		SELECT id, name, breed, birth_date, conception_date, mother_id, father_id
		FROM horses WHERE name = ?`, name).Scan(
		&horse.ID, &horse.Name, &horse.Breed,
		&birthDateStr, &conceptionDateStr,
		&horse.MotherID, &horse.FatherID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get horse by name: %w", err)
	}

	if birthDateStr.Valid {
		date, err := time.Parse("2006-01-02", birthDateStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid birth date format: %w", err)
		}
		horse.DateOfBirth = date
	}

	if conceptionDateStr.Valid {
		date, err := time.Parse("2006-01-02", conceptionDateStr.String)
		if err != nil {
			return nil, fmt.Errorf("invalid conception date format: %w", err)
		}
		horse.ConceptionDate = &date
	}

	return &horse, nil
}

func (s *SQLiteStore) GetAllHorses() ([]models.Horse, error) {
	rows, err := s.db.Query(`
		SELECT id, name, breed, birth_date, conception_date, mother_id, father_id
		FROM horses`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all horses: %w", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		var horse models.Horse
		var birthDateStr, conceptionDateStr sql.NullString

		err := rows.Scan(
			&horse.ID, &horse.Name, &horse.Breed,
			&birthDateStr, &conceptionDateStr,
			&horse.MotherID, &horse.FatherID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan horse row: %w", err)
		}

		if birthDateStr.Valid {
			date, err := time.Parse("2006-01-02", birthDateStr.String)
			if err != nil {
				logger.Error(err, "Invalid birth date format", map[string]interface{}{
					"horseID": horse.ID,
					"date":    birthDateStr.String,
				})
				continue
			}
			horse.DateOfBirth = date
		}

		if conceptionDateStr.Valid {
			date, err := time.Parse("2006-01-02", conceptionDateStr.String)
			if err != nil {
				logger.Error(err, "Invalid conception date format", map[string]interface{}{
					"horseID": horse.ID,
					"date":    conceptionDateStr.String,
				})
				continue
			}
			horse.ConceptionDate = &date
		}

		horses = append(horses, horse)
	}

	return horses, nil
}

func (s *SQLiteStore) AddHorse(horse *models.Horse) error {
	result, err := s.db.Exec(`
		INSERT INTO horses (name, breed, birth_date, conception_date, mother_id, father_id)
		VALUES (?, ?, ?, ?, ?, ?)`,
		horse.Name, horse.Breed,
		horse.DateOfBirth.Format("2006-01-02"),
		horse.ConceptionDate.Format("2006-01-02"),
		horse.MotherID, horse.FatherID,
	)
	if err != nil {
		return fmt.Errorf("failed to add horse: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	horse.ID = id
	return nil
}

func (s *SQLiteStore) UpdateHorse(horse *models.Horse) error {
	_, err := s.db.Exec(`
		UPDATE horses
		SET name = ?, breed = ?, birth_date = ?, conception_date = ?, mother_id = ?, father_id = ?
		WHERE id = ?`,
		horse.Name, horse.Breed,
		horse.DateOfBirth.Format("2006-01-02"),
		horse.ConceptionDate.Format("2006-01-02"),
		horse.MotherID, horse.FatherID,
		horse.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteHorse(id int64) error {
	_, err := s.db.Exec("DELETE FROM horses WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete horse: %w", err)
	}
	return nil
}

// Health record operations
func (s *SQLiteStore) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	rows, err := s.db.Query(`
		SELECT id, horse_id, date, type, notes
		FROM health_records
		WHERE horse_id = ?`, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health records: %w", err)
	}
	defer rows.Close()

	var records []models.HealthRecord
	for rows.Next() {
		var record models.HealthRecord
		var dateStr string

		err := rows.Scan(
			&record.ID, &record.HorseID,
			&dateStr, &record.Type, &record.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan health record row: %w", err)
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			logger.Error(err, "Invalid health record date format", map[string]interface{}{
				"recordID": record.ID,
				"date":     dateStr,
			})
			continue
		}
		record.Date = date

		records = append(records, record)
	}

	return records, nil
}

func (s *SQLiteStore) AddHealthRecord(record *models.HealthRecord) error {
	result, err := s.db.Exec(`
		INSERT INTO health_records (horse_id, date, type, notes)
		VALUES (?, ?, ?, ?)`,
		record.HorseID, record.Date.Format("2006-01-02"),
		record.Type, record.Notes,
	)
	if err != nil {
		return fmt.Errorf("failed to add health record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	record.ID = id
	return nil
}

// Implement other methods similarly...

// Helper functions
func parseNullTime(s sql.NullString) (*time.Time, error) {
	if !s.Valid {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s.String)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
