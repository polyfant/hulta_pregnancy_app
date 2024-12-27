package database

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/polyfant/horse_tracking/internal/models"
)

//go:embed schema.sql
var schemaFS embed.FS

type SQLiteStore struct {
	db *sql.DB
}

// UpdatePreFoalingSign implements models.DataStore.
func (s *SQLiteStore) UpdatePreFoalingSign(sign *models.PreFoalingSign) error {
	panic("unimplemented")
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Initialize schema
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, fmt.Errorf("error reading schema: %v", err)
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, fmt.Errorf("error initializing schema: %v", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func (s *SQLiteStore) parseDate(str string) (time.Time, error) {
	return time.Parse("2006-01-02", str)
}

func (s *SQLiteStore) GetHorse(id int64) (models.Horse, error) {
	query := `
		SELECT h.id, h.name, h.breed, h.gender, h.date_of_birth, h.weight, 
		       h.conception_date, h.mother_id, h.father_id, 
		       h.external_mother, h.external_father,
		       m.name as mother_name, f.name as father_name
		FROM horses h
		LEFT JOIN horses m ON h.mother_id = m.id
		LEFT JOIN horses f ON h.father_id = f.id
		WHERE h.id = ?`

	var horse models.Horse
	var dateOfBirthStr string
	var conceptionDateStr sql.NullString
	var motherID, fatherID sql.NullInt64
	var externalMother, externalFather sql.NullString
	var motherName, fatherName sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&horse.ID, &horse.Name, &horse.Breed, &horse.Gender,
		&dateOfBirthStr, &horse.Weight, &conceptionDateStr,
		&motherID, &fatherID, &externalMother, &externalFather,
		&motherName, &fatherName,
	)
	if err != nil {
		return horse, fmt.Errorf("error getting horse: %v", err)
	}

	// Parse dates
	dateOfBirth, err := s.parseDate(dateOfBirthStr)
	if err != nil {
		return horse, fmt.Errorf("error parsing birth date: %v", err)
	}
	horse.DateOfBirth = dateOfBirth

	if conceptionDateStr.Valid {
		conceptionDate, err := s.parseDate(conceptionDateStr.String)
		if err != nil {
			return horse, fmt.Errorf("error parsing conception date: %v", err)
		}
		horse.ConceptionDate = &conceptionDate
	}

	if motherID.Valid {
		horse.MotherID = &motherID.Int64
	}
	if fatherID.Valid {
		horse.FatherID = &fatherID.Int64
	}

	if externalMother.Valid {
		horse.ExternalMother = externalMother.String
	}
	if externalFather.Valid {
		horse.ExternalFather = externalFather.String
	}

	// Calculate age
	horse.Age = horse.CalculateAge(time.Now())

	return horse, nil
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
	if birthDate, err := s.parseDate(birthDateStr.String); err != nil {
		return nil, err
	} else if birthDateStr.Valid {
		horse.DateOfBirth = birthDate
	}

	if conceptionDate, err := s.parseDate(conceptionDateStr.String); err != nil {
		return nil, err
	} else if conceptionDateStr.Valid {
		horse.ConceptionDate = &conceptionDate
	}

	return &horse, nil
}

func (s *SQLiteStore) GetAllHorses() ([]models.Horse, error) {
	rows, err := s.db.Query(`
		SELECT id, name, breed, gender, date_of_birth, weight, conception_date, mother_id, father_id 
		FROM horses ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("error getting all horses: %v", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		var horse models.Horse
		var birthDateStr string
		var conceptionDateStr sql.NullString
		var motherID, fatherID sql.NullInt64

		err := rows.Scan(
			&horse.ID, &horse.Name, &horse.Breed, &horse.Gender,
			&birthDateStr, &horse.Weight, &conceptionDateStr,
			&motherID, &fatherID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning horse row: %v", err)
		}

		// Parse dates
		birthDate, err := s.parseDate(birthDateStr)
		if err != nil {
			return nil, err
		}
		horse.DateOfBirth = birthDate

		if conceptionDateStr.Valid {
			conceptionDate, err := s.parseDate(conceptionDateStr.String)
			if err != nil {
				return nil, err
			}
			horse.ConceptionDate = &conceptionDate
		}

		if motherID.Valid {
			horse.MotherID = &motherID.Int64
		}
		if fatherID.Valid {
			horse.FatherID = &fatherID.Int64
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
	query := `
		INSERT INTO horses (
			name, breed, gender, date_of_birth, weight, 
			conception_date, mother_id, father_id, 
			external_mother, external_father
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	dateOfBirth := s.formatDate(horse.DateOfBirth)
	var conceptionDate *string
	if horse.ConceptionDate != nil {
		formatted := s.formatDate(*horse.ConceptionDate)
		conceptionDate = &formatted
	}

	result, err := s.db.Exec(query,
		horse.Name, horse.Breed, horse.Gender, dateOfBirth,
		horse.Weight, conceptionDate, horse.MotherID, horse.FatherID,
		nullString(horse.ExternalMother), nullString(horse.ExternalFather),
	)
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
		SET name = ?, breed = ?, gender = ?, date_of_birth = ?, weight = ?, 
		conception_date = ?, mother_id = ?, father_id = ?, 
		external_mother = ?, external_father = ?
		WHERE id = ?`,
		horse.Name, horse.Breed, horse.Gender, s.formatDate(horse.DateOfBirth), horse.Weight,
		s.formatDate(*horse.ConceptionDate), horse.MotherID, horse.FatherID,
		nullString(horse.ExternalMother), nullString(horse.ExternalFather), horse.ID)
	if err != nil {
		return fmt.Errorf("error updating horse: %v", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteHorse(id int64) error {
	_, err := s.db.Exec("DELETE FROM horses WHERE id = ?", id)
	return err
}

func (s *SQLiteStore) GetOffspring(horseID int64) ([]models.Horse, error) {
	query := `
		SELECT id, name, breed, gender, date_of_birth, weight, conception_date, mother_id, father_id 
		FROM horses 
		WHERE mother_id = ? OR father_id = ?
		ORDER BY date_of_birth DESC
	`
	rows, err := s.db.Query(query, horseID, horseID)
	if err != nil {
		return nil, fmt.Errorf("error getting offspring: %v", err)
	}
	defer rows.Close()

	var offspring []models.Horse
	for rows.Next() {
		var h models.Horse
		var dateOfBirthStr string
		var conceptionDateStr sql.NullString
		var motherID, fatherID sql.NullInt64

		err := rows.Scan(&h.ID, &h.Name, &h.Breed, &h.Gender, &dateOfBirthStr, &h.Weight, &conceptionDateStr, &motherID, &fatherID)
		if err != nil {
			return nil, fmt.Errorf("error scanning offspring: %v", err)
		}

		dateOfBirth, err := s.parseDate(dateOfBirthStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing birth date: %v", err)
		}
		h.DateOfBirth = dateOfBirth

		if conceptionDateStr.Valid {
			conceptionDate, err := s.parseDate(conceptionDateStr.String)
			if err != nil {
				return nil, fmt.Errorf("error parsing conception date: %v", err)
			}
			h.ConceptionDate = &conceptionDate
		}

		if motherID.Valid {
			h.MotherID = &motherID.Int64
		}
		if fatherID.Valid {
			h.FatherID = &fatherID.Int64
		}

		offspring = append(offspring, h)
	}

	return offspring, nil
}

func (s *SQLiteStore) GetFamilyTree(horseID int64) (models.FamilyTree, error) {
	tree := models.FamilyTree{}

	// Get the main horse
	horse, err := s.GetHorse(horseID)
	if err != nil {
		return tree, fmt.Errorf("error getting horse: %v", err)
	}
	tree.Horse = horse

	// Get mother information
	if horse.MotherID != nil {
		mother, err := s.GetHorse(*horse.MotherID)
		if err == nil {
			tree.Mother = &models.FamilyMember{
				ID:          mother.ID,
				Name:        mother.Name,
				Breed:       mother.Breed,
				Gender:      mother.Gender,
				DateOfBirth: mother.DateOfBirth,
				Age:         mother.CalculateAge(time.Now()),
				IsExternal:  false,
			}
		}
	} else if horse.ExternalMother != "" {
		tree.Mother = &models.FamilyMember{
			Name:           horse.ExternalMother,
			IsExternal:     true,
			ExternalSource: "External Mare",
		}
	}

	// Get father information
	if horse.FatherID != nil {
		father, err := s.GetHorse(*horse.FatherID)
		if err == nil {
			tree.Father = &models.FamilyMember{
				ID:          father.ID,
				Name:        father.Name,
				Breed:       father.Breed,
				Gender:      father.Gender,
				DateOfBirth: father.DateOfBirth,
				Age:         father.CalculateAge(time.Now()),
				IsExternal:  false,
			}
		}
	} else if horse.ExternalFather != "" {
		tree.Father = &models.FamilyMember{
			Name:           horse.ExternalFather,
			IsExternal:     true,
			ExternalSource: "External Stallion",
		}
	}

	// Get offspring
	offspring, err := s.GetOffspring(horseID)
	if err == nil && len(offspring) > 0 {
		for _, child := range offspring {
			tree.Offspring = append(tree.Offspring, models.FamilyMember{
				ID:          child.ID,
				Name:        child.Name,
				Breed:       child.Breed,
				Gender:      child.Gender,
				DateOfBirth: child.DateOfBirth,
				Age:         child.CalculateAge(time.Now()),
				IsExternal:  false,
			})
		}
	}

	// Get siblings (same mother or father)
	if horse.MotherID != nil || horse.FatherID != nil {
		query := `
			SELECT DISTINCT h.id, h.name, h.breed, h.gender, h.date_of_birth
			FROM horses h
			WHERE h.id != ? AND (
				(h.mother_id = ? AND ? IS NOT NULL) OR
				(h.father_id = ? AND ? IS NOT NULL)
			)
			ORDER BY h.date_of_birth
		`
		rows, err := s.db.Query(query, horseID,
			horse.MotherID, horse.MotherID,
			horse.FatherID, horse.FatherID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var sibling models.FamilyMember
				var dateOfBirthStr string
				err := rows.Scan(&sibling.ID, &sibling.Name, &sibling.Breed, &sibling.Gender, &dateOfBirthStr)
				if err == nil {
					if dateOfBirth, err := s.parseDate(dateOfBirthStr); err == nil {
						sibling.DateOfBirth = dateOfBirth
						sibling.Age = (&models.Horse{DateOfBirth: dateOfBirth}).CalculateAge(time.Now())
					}
					sibling.IsExternal = false
					tree.Siblings = append(tree.Siblings, sibling)
				}
			}
		}
	}

	return tree, nil
}

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
		record.Date, err = s.parseDate(dateStr)
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
		record.HorseID, s.formatDate(record.Date), record.Type, record.Notes)
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
		record.HorseID, s.formatDate(record.Date), record.Type, record.Notes, record.ID)
	return err
}

func (s *SQLiteStore) DeleteHealthRecord(id int64) error {
	_, err := s.db.Exec("DELETE FROM health_records WHERE id = ?", id)
	return err
}

func (s *SQLiteStore) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	query := `
		SELECT id, horse_id, date, type, description, notes
		FROM pregnancy_events
		WHERE horse_id = ?
		ORDER BY date DESC
	`

	rows, err := s.db.Query(query, horseID)
	if err != nil {
		return nil, fmt.Errorf("error querying pregnancy events: %v", err)
	}
	defer rows.Close()

	var events []models.PregnancyEvent
	for rows.Next() {
		var event models.PregnancyEvent
		var dateStr string
		err := rows.Scan(
			&event.ID,
			&event.HorseID,
			&dateStr,
			&event.Type,
			&event.Description,
			&event.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning pregnancy event: %v", err)
		}

		date, err := s.parseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		event.Date = date

		events = append(events, event)
	}

	return events, nil
}

func (s *SQLiteStore) GetUserPregnancyEvents(userID int64) ([]models.PregnancyEvent, error) {
	// TODO: Implement when user system is added
	return nil, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) AddPregnancyEvent(event *models.PregnancyEvent) error {
	query := `
		INSERT INTO pregnancy_events (horse_id, date, type, description, notes)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := s.db.Exec(query,
		event.HorseID,
		event.Date,
		event.Type,
		event.Description,
		nullString(event.Notes),
	)
	if err != nil {
		return fmt.Errorf("error adding pregnancy event: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}

	event.ID = id
	return nil
}

func (s *SQLiteStore) UpdatePregnancyEvent(event *models.PregnancyEvent) error {
	_, err := s.db.Exec(`
		UPDATE pregnancy_events 
		SET horse_id = ?, date = ?, type = ?, description = ?, notes = ?
		WHERE id = ?`,
		event.HorseID, s.formatDate(event.Date), event.Type, event.Description, event.Notes, event.ID)
	return err
}

func (s *SQLiteStore) DeletePregnancyEvent(id int64) error {
	_, err := s.db.Exec("DELETE FROM pregnancy_events WHERE id = ?", id)
	return err
}

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
		cost.Date, err = s.parseDate(dateStr)
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
		cost.HorseID, cost.Description, cost.Amount, s.formatDate(cost.Date))
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
		cost.HorseID, cost.Description, cost.Amount, s.formatDate(cost.Date), cost.ID)
	return err
}

func (s *SQLiteStore) DeleteBreedingCost(id int64) error {
	_, err := s.db.Exec("DELETE FROM breeding_costs WHERE id = ?", id)
	return err
}

func (s *SQLiteStore) GetLastSyncTime(userID int64) (time.Time, error) {
	// For now, we'll just use GetUserLastSync since they serve the same purpose
	return s.GetUserLastSync(userID)
}

func (s *SQLiteStore) GetPendingSyncCount(userID int64) (int, error) {
	// This is a placeholder implementation
	return 0, nil
}

func (s *SQLiteStore) UpdateUserLastSync(userID int64, t time.Time) error {
	// TODO: Implement when user system is added
	return fmt.Errorf("not implemented")
}

func (s *SQLiteStore) GetUserLastSync(userID int64) (time.Time, error) {
	// TODO: Implement when user system is added
	return time.Time{}, fmt.Errorf("not implemented")
}

func (s *SQLiteStore) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *SQLiteStore) GetDashboardStats() (models.DashboardStats, error) {
	stats := models.DashboardStats{}
	var err error

	// Get total horses
	err = s.db.QueryRow("SELECT COUNT(*) FROM horses").Scan(&stats.TotalHorses)
	if err != nil {
		return stats, fmt.Errorf("error getting total horses: %v", err)
	}

	// Get pregnant mares (horses with conception_date set and no due date passed)
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM horses 
		WHERE conception_date IS NOT NULL 
		AND date(conception_date, '+340 days') >= date('now')`).Scan(&stats.PregnantMares)
	if err != nil {
		return stats, fmt.Errorf("error getting pregnant mares: %v", err)
	}

	// Get upcoming due dates
	rows, err := s.db.Query(`
		SELECT h.id, h.name, date(h.conception_date, '+340 days') as due_date,
		   julianday(date(h.conception_date, '+340 days')) - julianday('now') as days_remaining
		FROM horses h
		WHERE h.conception_date IS NOT NULL 
		AND date(h.conception_date, '+340 days') >= date('now')
		ORDER BY due_date ASC
		LIMIT 5`)
	if err != nil {
		return stats, fmt.Errorf("error getting upcoming due dates: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var summary models.PregnancySummary
		var dueDateStr string
		err := rows.Scan(&summary.HorseID, &summary.HorseName, &dueDateStr, &summary.DaysRemaining)
		if err != nil {
			return stats, fmt.Errorf("error scanning pregnancy summary: %v", err)
		}

		date, err := s.parseDate(dueDateStr)
		if err != nil {
			return stats, fmt.Errorf("error parsing due date: %v", err)
		}
		summary.DueDate = date

		stats.UpcomingDueDates = append(stats.UpcomingDueDates, summary)
	}

	// Get recent health events
	rows, err = s.db.Query(`
		SELECT h.id, h.name, hr.type, hr.date, hr.notes
		FROM health_records hr
		JOIN horses h ON h.id = hr.horse_id
		WHERE date(hr.date) >= date('now', '-30 days')
		ORDER BY hr.date DESC
		LIMIT 5`)
	if err != nil {
		return stats, fmt.Errorf("error getting recent health events: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var summary models.HealthSummary
		var dateStr string
		err := rows.Scan(&summary.HorseID, &summary.HorseName, &summary.Type, &dateStr, &summary.Notes)
		if err != nil {
			return stats, fmt.Errorf("error scanning health summary: %v", err)
		}

		date, err := s.parseDate(dateStr)
		if err != nil {
			return stats, fmt.Errorf("error parsing health date: %v", err)
		}
		summary.Date = date

		stats.RecentHealth = append(stats.RecentHealth, summary)
	}

	// Get cost summary
	err = s.db.QueryRow(`
		SELECT 
			COALESCE(SUM(amount), 0) as total_costs,
			COALESCE(AVG(amount), 0) as monthly_average,
			COALESCE((
				SELECT SUM(amount) 
				FROM breeding_costs 
				WHERE date >= date('now', '-30 days')
			), 0) as recent_costs
		FROM breeding_costs`).Scan(
		&stats.Costs.TotalCosts,
		&stats.Costs.MonthlyAverage,
		&stats.Costs.RecentCosts)
	if err != nil {
		return stats, fmt.Errorf("error getting cost summary: %v", err)
	}

	return stats, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func (s *SQLiteStore) UpdateHorseConceptionDate(horseID int64, conceptionDate time.Time) error {
	query := `UPDATE horses SET conception_date = ? WHERE id = ?`

	var dateStr *string
	if !conceptionDate.IsZero() {
		str := s.formatDate(conceptionDate)
		dateStr = &str
	}

	_, err := s.db.Exec(query, dateStr, horseID)
	if err != nil {
		return fmt.Errorf("error updating conception date: %v", err)
	}

	return nil
}

func (s *SQLiteStore) GetPreFoalingSigns(horseID int64) ([]models.PreFoalingSign, error) {
	query := `
		SELECT id, horse_id, name, observed, date_observed, notes
		FROM pre_foaling_signs
		WHERE horse_id = ?
		ORDER BY date_observed DESC
	`

	rows, err := s.db.Query(query, horseID)
	if err != nil {
		return nil, fmt.Errorf("error querying pre-foaling signs: %v", err)
	}
	defer rows.Close()

	var signs []models.PreFoalingSign
	for rows.Next() {
		var sign models.PreFoalingSign
		var dateStr string
		err := rows.Scan(
			&sign.ID,
			&sign.HorseID,
			&sign.Name,
			&sign.Observed,
			&dateStr,
			&sign.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning pre-foaling sign: %v", err)
		}

		date, err := s.parseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}
		sign.DateObserved = date

		signs = append(signs, sign)
	}

	return signs, nil
}

func (s *SQLiteStore) AddPreFoalingSign(sign *models.PreFoalingSign) error {
	query := `
		INSERT INTO pre_foaling_signs (horse_id, name, observed, date_observed, notes)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := s.db.Exec(query,
		sign.HorseID,
		sign.Name,
		sign.Observed,
		s.formatDate(sign.DateObserved),
		sign.Notes,
	)

	if err != nil {
		return fmt.Errorf("error inserting pre-foaling sign: %v", err)
	}

	return nil
}

func (s *SQLiteStore) UpdateHorsePregnancyStatus(horseID int64, isPregnant bool, conceptionDate time.Time) error {
	query := `
        UPDATE horses 
        SET is_pregnant = ?,
            conception_date = ?
        WHERE id = ?`

	formatted := s.formatDate(conceptionDate)
	_, err := s.db.Exec(query, isPregnant, formatted, horseID)
	if err != nil {
		return fmt.Errorf("error updating horse pregnancy status: %v", err)
	}

	return nil
}

func (s *SQLiteStore) GetPreFoalingChecklist(horseID int64) ([]models.PreFoalingChecklistItem, error) {
	query := `SELECT id, horse_id, description, is_completed, due_date, priority, notes 
			 FROM pre_foaling_checklist 
			 WHERE horse_id = ? 
			 ORDER BY due_date ASC, priority DESC`
	
	rows, err := s.db.Query(query, horseID)
	if err != nil {
		return nil, fmt.Errorf("error querying pre-foaling checklist: %w", err)
	}
	defer rows.Close()

	var items []models.PreFoalingChecklistItem
	for rows.Next() {
		var item models.PreFoalingChecklistItem
		err := rows.Scan(
			&item.ID,
			&item.HorseID,
			&item.Description,
			&item.IsCompleted,
			&item.DueDate,
			&item.Priority,
			&item.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning pre-foaling checklist row: %w", err)
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pre-foaling checklist rows: %w", err)
	}
	return items, nil
}

func (s *SQLiteStore) AddPreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	query := `INSERT INTO pre_foaling_checklist 
			 (horse_id, description, is_completed, due_date, priority, notes) 
			 VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := s.db.Exec(query, 
		item.HorseID,
		item.Description,
		item.IsCompleted,
		item.DueDate,
		item.Priority,
		item.Notes)
	
	if err != nil {
		return fmt.Errorf("error adding pre-foaling checklist item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}
	item.ID = id
	return nil
}

func (s *SQLiteStore) UpdatePreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	query := `UPDATE pre_foaling_checklist 
			 SET description = ?, is_completed = ?, due_date = ?, priority = ?, notes = ? 
			 WHERE id = ? AND horse_id = ?`
	
	result, err := s.db.Exec(query,
		item.Description,
		item.IsCompleted,
		item.DueDate,
		item.Priority,
		item.Notes,
		item.ID,
		item.HorseID)
	
	if err != nil {
		return fmt.Errorf("error updating pre-foaling checklist item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("checklist item not found")
	}
	return nil
}

func (s *SQLiteStore) DeletePreFoalingChecklistItem(id int64) error {
	query := `DELETE FROM pre_foaling_checklist WHERE id = ?`
	
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting pre-foaling checklist item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("checklist item not found")
	}
	return nil
}
