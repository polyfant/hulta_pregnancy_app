package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

// Implement all required methods from the DataStore interface
func (m *MockDB) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockDB) AddHealthRecord(record *models.HealthRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

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

func TestSyncData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)

	t.Run("successful sync", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Mock user authentication
		c.Set("userID", int64(1))

		// Test data
		syncData := models.SyncData{
			UserID: 1,
			Horses: []models.Horse{
				{
					ID:   1,
					Name: "Test Horse",
				},
			},
			Health: []models.HealthRecord{
				{
					ID:      1,
					HorseID: 1,
					Date:    time.Now(),
					Type:    "Checkup",
				},
			},
		}

		body, _ := json.Marshal(syncData)
		c.Request, _ = http.NewRequest("POST", "/sync", bytes.NewBuffer(body))

		mockDB.On("AddHorse", &syncData.Horses[0]).Return(nil)
		mockDB.On("AddHealthRecord", &syncData.Health[0]).Return(nil)

		handler.SyncData(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
	})

	t.Run("unauthorized sync", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		syncData := models.SyncData{
			UserID: 2, // Different user ID
			Horses: []models.Horse{
				{
					ID:   1,
					Name: "Test Horse",
				},
			},
		}

		body, _ := json.Marshal(syncData)
		c.Request, _ = http.NewRequest("POST", "/sync", bytes.NewBuffer(body))
		c.Set("userID", int64(1)) // Set different user ID in context

		handler.SyncData(c)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGetSyncStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)

	t.Run("successful status check", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", int64(1))

		lastSync := time.Now().Add(-24 * time.Hour)
		mockDB.On("GetLastSyncTime", int64(1)).Return(lastSync, nil)
		mockDB.On("GetPendingSyncCount", int64(1)).Return(5, nil)

		handler.GetSyncStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		
		var response SyncStatus
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, lastSync.Unix(), response.LastSync.Unix())
		assert.Equal(t, 5, response.PendingSync)
		assert.False(t, response.IsUpToDate)
	})

	t.Run("unauthorized status check", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Don't set userID in context

		handler.GetSyncStatus(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
