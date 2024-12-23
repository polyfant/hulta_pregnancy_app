package health

import (
	"database/sql"
	"testing"
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockDB) AddHealthRecord(record *models.HealthRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

// Implement other required methods from DataStore interface
func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockDB) GetAllHorses() ([]models.Horse, error) {
	args := m.Called()
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockDB) GetHorse(id int64) (models.Horse, error) {
	args := m.Called(id)
	return args.Get(0).(models.Horse), args.Error(1)
}

func (m *MockDB) AddHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDB) UpdateHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDB) DeleteHorse(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) GetBreedingCosts(horseID int64) ([]models.BreedingCost, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.BreedingCost), args.Error(1)
}

func (m *MockDB) AddBreedingCost(cost *models.BreedingCost) error {
	args := m.Called(cost)
	return args.Error(0)
}

func (m *MockDB) DeleteBreedingCost(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockDB) GetLastSyncTime(userID int64) (time.Time, error) {
	args := m.Called(userID)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) GetPendingSyncCount(userID int64) (int, error) {
	args := m.Called(userID)
	return args.Get(0).(int), args.Error(1)
}

func TestGetVaccinationStatus(t *testing.T) {
	mockDB := new(MockDB)
	service := NewHealthService(mockDB)

	horse := models.Horse{ID: 1}
	
	// Test case 1: No vaccinations
	mockDB.On("GetHealthRecords", int64(1)).Return([]models.HealthRecord{}, nil)
	
	status := service.GetVaccinationStatus(horse)
	assert.False(t, status.IsUpToDate)
	assert.True(t, status.LastDate.IsZero())
	
	// Test case 2: Recent vaccination
	now := time.Now()
	mockDB.On("GetHealthRecords", int64(1)).Return([]models.HealthRecord{
		{
			Type: "Vaccination",
			Date: now.AddDate(0, -5, 0), // 5 months ago
		},
	}, nil)
	
	status = service.GetVaccinationStatus(horse)
	assert.True(t, status.IsUpToDate)
	assert.Equal(t, now.AddDate(0, -5, 0).Unix(), status.LastDate.Unix())
}

func TestGetHealthSummary(t *testing.T) {
	mockDB := new(MockDB)
	service := NewHealthService(mockDB)

	horse := models.Horse{ID: 1}
	now := time.Now()
	
	mockDB.On("GetHealthRecords", int64(1)).Return([]models.HealthRecord{
		{
			Type: "Checkup",
			Date: now.AddDate(0, -1, 0),
		},
		{
			Type: "Vaccination",
			Date: now.AddDate(0, -5, 0),
		},
	}, nil)
	
	summary := service.GetHealthSummary(horse)
	assert.Equal(t, 2, summary.TotalRecords)
	assert.Equal(t, now.AddDate(0, -1, 0).Unix(), summary.LastCheckup.Unix())
	assert.True(t, summary.VaccinationStatus.IsUpToDate)
}
