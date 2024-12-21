package database

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/polyfant/horse_tracking/internal/models"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// Set up test database
	var err error
	testDB, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	// Create tables
	err = createTables(testDB)
	if err != nil {
		panic(err)
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestCreateHorse(t *testing.T) {
	// Create a test horse
	horse := &models.Horse{
		Name:        "TestHorse",
		Breed:       "TestBreed",
		DateOfBirth: time.Now().AddDate(-5, 0, 0), // 5 years old
	}

	// Insert horse
	result, err := testDB.Exec(
		"INSERT INTO horses (name, breed, birth_date) VALUES (?, ?, ?)",
		horse.Name, horse.Breed, horse.DateOfBirth,
	)
	assert.NoError(t, err)

	// Verify insertion
	id, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Retrieve horse
	var retrievedHorse models.Horse
	err = testDB.QueryRow(
		"SELECT id, name, breed, birth_date FROM horses WHERE id = ?",
		id,
	).Scan(&retrievedHorse.ID, &retrievedHorse.Name, &retrievedHorse.Breed, &retrievedHorse.DateOfBirth)
	
	assert.NoError(t, err)
	assert.Equal(t, horse.Name, retrievedHorse.Name)
	assert.Equal(t, horse.Breed, retrievedHorse.Breed)
}

func TestCreateHealthRecord(t *testing.T) {
	// First create a horse
	horse := &models.Horse{
		Name:        "TestHorse2",
		Breed:       "TestBreed2",
		DateOfBirth: time.Now().AddDate(-3, 0, 0),
	}

	result, err := testDB.Exec(
		"INSERT INTO horses (name, breed, birth_date) VALUES (?, ?, ?)",
		horse.Name, horse.Breed, horse.DateOfBirth,
	)
	assert.NoError(t, err)
	horseID, err := result.LastInsertId()
	assert.NoError(t, err)

	// Create health record
	record := &models.HealthRecord{
		HorseID: horseID,
		Date:    time.Now(),
		Type:    "Vaccination",
		Notes:   "Annual vaccination",
	}

	result, err = testDB.Exec(
		"INSERT INTO health_records (horse_id, date, type, notes) VALUES (?, ?, ?, ?)",
		record.HorseID, record.Date, record.Type, record.Notes,
	)
	assert.NoError(t, err)

	id, err := result.LastInsertId()
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Verify record
	var retrievedRecord models.HealthRecord
	err = testDB.QueryRow(
		"SELECT id, horse_id, date, type, notes FROM health_records WHERE id = ?",
		id,
	).Scan(&retrievedRecord.ID, &retrievedRecord.HorseID, &retrievedRecord.Date, 
		&retrievedRecord.Type, &retrievedRecord.Notes)
	
	assert.NoError(t, err)
	assert.Equal(t, record.HorseID, retrievedRecord.HorseID)
	assert.Equal(t, record.Type, retrievedRecord.Type)
	assert.Equal(t, record.Notes, retrievedRecord.Notes)
}
