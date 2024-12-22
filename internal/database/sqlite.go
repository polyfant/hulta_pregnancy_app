package database

import (
	"database/sql"
	"fmt"
	"time"

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
		SELECT id, name, breed, date_of_birth, weight, conception_date, mother_id, father_id
		FROM horses WHERE id = ?`, id).Scan(
		&horse.ID, &horse.Name, &horse.Breed,
		&birthDateStr, &horse.Weight, &conceptionDateStr,
		&horse.MotherID, &horse.FatherID,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting horse: %v", err)
	}

	// Parse dates
	if birthDate, err := parseNullTime(birthDateStr); err != nil {
		return nil, err
	} else if birthDate != nil {
		horse.DateOfBirth = *birthDate
	}

	if conceptionDate, err := parseNullTime(conceptionDateStr); err != nil {
		return nil, err
	} else {
		horse.ConceptionDate = conceptionDate
	}

	return &horse, nil
}

func (s *SQLiteStore) GetHorseByName(name string) (*models.Horse, error) {
	var horse models.Horse
	var birthDateStr, conceptionDateStr sql.NullString

	err := s.db.QueryRow(`
		SELECT id, name, breed, date_of_birth, weight, conception_date, mother_id, father_id
		FROM horses WHERE name = ?`, name).Scan(
		&horse.ID, &horse.Name, &horse.Breed,
		&birthDateStr, &horse.Weight, &conceptionDateStr,
		&horse.MotherID, &horse.FatherID,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting horse by name: %v", err)
	}

	// Parse dates
	if birthDate, err := parseNullTime(birthDateStr); err != nil {
		return nil, err
	} else if birthDate != nil {
		horse.DateOfBirth = *birthDate
	}

	if conceptionDate, err := parseNullTime(conceptionDateStr); err != nil {
		return nil, err
	} else {
		horse.ConceptionDate = conceptionDate
	}

	return &horse, nil
}

func (s *SQLiteStore) GetAllHorses() ([]models.Horse, error) {
	rows, err := s.db.Query(`
		SELECT id, name, breed, date_of_birth, weight, conception_date, mother_id, father_id 
		FROM horses ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("error getting all horses: %v", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		var horse models.Horse
		var birthDateStr, conceptionDateStr sql.NullString

		err := rows.Scan(
			&horse.ID, &horse.Name, &horse.Breed,
			&birthDateStr, &horse.Weight, &conceptionDateStr,
			&horse.MotherID, &horse.FatherID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning horse row: %v", err)
		}

		// Parse dates
		if birthDate, err := parseNullTime(birthDateStr); err != nil {
			return nil, err
		} else if birthDate != nil {
			horse.DateOfBirth = *birthDate
		}

		if conceptionDate, err := parseNullTime(conceptionDateStr); err != nil {
			return nil, err
		} else {
			horse.ConceptionDate = conceptionDate
		}

		horses = append(horses, horse)
	}

	return horses, rows.Err()
}

func (s *SQLiteStore) GetUserHorses(userID int64) ([]models.Horse, error) {
	// TODO: Implement when user system is added
	return nil, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) AddHorse(horse *models.Horse) error {
	result, err := s.db.Exec(`
		INSERT INTO horses (name, breed, date_of_birth, weight, conception_date, mother_id, father_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		horse.Name, horse.Breed, horse.DateOfBirth.Format("2006-01-02"), horse.Weight,
		formatNullTime(horse.ConceptionDate), horse.MotherID, horse.FatherID)
	if err != nil {
		return fmt.Errorf("error adding horse: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	horse.ID = id
	return nil
}

func (s *SQLiteStore) UpdateHorse(horse *models.Horse) error {
	_, err := s.db.Exec(`
		UPDATE horses 
		SET name = ?, breed = ?, date_of_birth = ?, weight = ?, conception_date = ?, mother_id = ?, father_id = ?
		WHERE id = ?`,
		horse.Name, horse.Breed, horse.DateOfBirth.Format("2006-01-02"), horse.Weight,
		formatNullTime(horse.ConceptionDate), horse.MotherID, horse.FatherID, horse.ID)
	if err != nil {
		return fmt.Errorf("error updating horse: %v", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteHorse(id int64) error {
	_, err := s.db.Exec("DELETE FROM horses WHERE id = ?", id)
	return err
}

// Health record operations
func (s *SQLiteStore) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	rows, err := s.db.Query(`
		SELECT id, horse_id, date, type, notes
		FROM health_records 
		WHERE horse_id = ?
		ORDER BY date DESC`, horseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.HealthRecord
	for rows.Next() {
		var record models.HealthRecord
		var dateStr string
		err := rows.Scan(&record.ID, &record.HorseID, &dateStr, &record.Type, &record.Notes)
		if err != nil {
			return nil, err
		}
		record.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (s *SQLiteStore) GetUserHealthRecords(userID int64) ([]models.HealthRecord, error) {
	// TODO: Implement when user system is added
	return nil, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) AddHealthRecord(record *models.HealthRecord) error {
	result, err := s.db.Exec(`
		INSERT INTO health_records (horse_id, date, type, notes)
		VALUES (?, ?, ?, ?)`,
		record.HorseID, record.Date.Format("2006-01-02"), record.Type, record.Notes)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	record.ID = id
	return nil
}

func (s *SQLiteStore) UpdateHealthRecord(record *models.HealthRecord) error {
	_, err := s.db.Exec(`
		UPDATE health_records 
		SET horse_id = ?, date = ?, type = ?, notes = ?
		WHERE id = ?`,
		record.HorseID, record.Date.Format("2006-01-02"), record.Type, record.Notes, record.ID)
	return err
}

func (s *SQLiteStore) DeleteHealthRecord(id int64) error {
	_, err := s.db.Exec("DELETE FROM health_records WHERE id = ?", id)
	return err
}

// Pregnancy event operations
func (s *SQLiteStore) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	rows, err := s.db.Query(`
		SELECT id, horse_id, date, description
		FROM pregnancy_events 
		WHERE horse_id = ?
		ORDER BY date DESC`, horseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.PregnancyEvent
	for rows.Next() {
		var event models.PregnancyEvent
		var dateStr string
		err := rows.Scan(&event.ID, &event.HorseID, &dateStr, &event.Description)
		if err != nil {
			return nil, err
		}
		event.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

func (s *SQLiteStore) GetUserPregnancyEvents(userID int64) ([]models.PregnancyEvent, error) {
	// TODO: Implement when user system is added
	return nil, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) AddPregnancyEvent(event *models.PregnancyEvent) error {
	result, err := s.db.Exec(`
		INSERT INTO pregnancy_events (horse_id, date, description)
		VALUES (?, ?, ?)`,
		event.HorseID, event.Date.Format("2006-01-02"), event.Description)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = id
	return nil
}

func (s *SQLiteStore) UpdatePregnancyEvent(event *models.PregnancyEvent) error {
	_, err := s.db.Exec(`
		UPDATE pregnancy_events 
		SET horse_id = ?, date = ?, description = ?
		WHERE id = ?`,
		event.HorseID, event.Date.Format("2006-01-02"), event.Description, event.ID)
	return err
}

func (s *SQLiteStore) DeletePregnancyEvent(id int64) error {
	_, err := s.db.Exec("DELETE FROM pregnancy_events WHERE id = ?", id)
	return err
}

// Breeding cost operations
func (s *SQLiteStore) GetBreedingCosts(horseID int64) ([]models.BreedingCost, error) {
	rows, err := s.db.Query(`
		SELECT id, horse_id, description, amount, date
		FROM breeding_costs 
		WHERE horse_id = ?
		ORDER BY date DESC`, horseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var costs []models.BreedingCost
	for rows.Next() {
		var cost models.BreedingCost
		var dateStr string
		err := rows.Scan(&cost.ID, &cost.HorseID, &cost.Description, &cost.Amount, &dateStr)
		if err != nil {
			return nil, err
		}
		cost.Date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
		costs = append(costs, cost)
	}
	return costs, rows.Err()
}

func (s *SQLiteStore) AddBreedingCost(cost *models.BreedingCost) error {
	result, err := s.db.Exec(`
		INSERT INTO breeding_costs (horse_id, description, amount, date)
		VALUES (?, ?, ?, ?)`,
		cost.HorseID, cost.Description, cost.Amount, cost.Date.Format("2006-01-02"))
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	cost.ID = id
	return nil
}

func (s *SQLiteStore) UpdateBreedingCost(cost *models.BreedingCost) error {
	_, err := s.db.Exec(`
		UPDATE breeding_costs 
		SET horse_id = ?, description = ?, amount = ?, date = ?
		WHERE id = ?`,
		cost.HorseID, cost.Description, cost.Amount, cost.Date.Format("2006-01-02"), cost.ID)
	return err
}

func (s *SQLiteStore) DeleteBreedingCost(id int64) error {
	_, err := s.db.Exec("DELETE FROM breeding_costs WHERE id = ?", id)
	return err
}

// User operations
func (s *SQLiteStore) UpdateUserLastSync(userID int64, t time.Time) error {
	// TODO: Implement when user system is added
	return fmt.Errorf("not implemented")
}

func (s *SQLiteStore) GetUserLastSync(userID int64) (time.Time, error) {
	// TODO: Implement when user system is added
	return time.Time{}, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) Begin() (models.Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

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

func formatNullTime(t *time.Time) sql.NullString {
	if t == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{
		String: t.Format("2006-01-02"),
		Valid:  true,
	}
}
