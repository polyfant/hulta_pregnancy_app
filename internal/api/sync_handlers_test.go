package api

import (
	"bytes"
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

func (m *MockDB) Begin() (models.Transaction, error) {
	args := m.Called()
	return args.Get(0).(models.Transaction), args.Error(1)
}

func (m *MockDB) AddHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDB) GetUserHorses(userID int64) ([]models.Horse, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockDB) AddHealthRecord(record *models.HealthRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockDB) GetUserHealthRecords(userID int64) ([]models.HealthRecord, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockDB) GetUserPregnancyEvents(userID int64) ([]models.PregnancyEvent, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.PregnancyEvent), args.Error(1)
}

func (m *MockDB) UpdateUserLastSync(userID int64, time time.Time) error {
	args := m.Called(userID, time)
	return args.Error(0)
}

func (m *MockDB) GetUserLastSync(userID int64) (time.Time, error) {
	args := m.Called(userID)
	return args.Get(0).(time.Time), args.Error(1)
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func setupTestSyncRouter() (*gin.Engine, *MockDB) {
	gin.SetMode(gin.TestMode)
	mockDB := new(MockDB)
	handler := NewHandler(mockDB)
	router := gin.Default()

	router.POST("/sync", handler.SyncData)
	router.GET("/sync/status", handler.GetSyncStatus)
	router.GET("/sync/restore", handler.RestoreData)

	return router, mockDB
}

func TestSyncData(t *testing.T) {
	router, mockDB := setupTestSyncRouter()

	t.Run("successful sync", func(t *testing.T) {
		// Setup test data
		syncData := models.SyncData{
			UserID:    1,
			Timestamp: time.Now(),
			Horses: []models.Horse{
				{ID: 1, Name: "Thunder"},
			},
			Health: []models.HealthRecord{
				{ID: 1, HorseID: 1, Type: "Checkup"},
			},
			Events: []models.PregnancyEvent{
				{ID: 1, HorseID: 1, Type: "Conception"},
			},
		}

		// Setup mock transaction
		mockTx := new(MockTransaction)
		mockTx.On("Commit").Return(nil)
		mockTx.On("Rollback").Return(nil)
		mockDB.On("Begin").Return(mockTx, nil)

		// Setup expectations
		mockDB.On("AddHorse", mock.AnythingOfType("*models.Horse")).Return(nil)
		mockDB.On("AddHealthRecord", mock.AnythingOfType("*models.HealthRecord")).Return(nil)
		mockDB.On("AddPregnancyEvent", mock.AnythingOfType("*models.PregnancyEvent")).Return(nil)
		mockDB.On("UpdateUserLastSync", int64(1), mock.AnythingOfType("time.Time")).Return(nil)

		// Create request
		jsonData, _ := json.Marshal(syncData)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/sync", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Set user ID in context
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("userID", int64(1))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockDB.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("unauthorized sync", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/sync", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRestoreData(t *testing.T) {
	router, mockDB := setupTestSyncRouter()

	t.Run("successful restore", func(t *testing.T) {
		// Setup test data
		horses := []models.Horse{{ID: 1, Name: "Thunder"}}
		health := []models.HealthRecord{{ID: 1, HorseID: 1, Type: "Checkup"}}
		events := []models.PregnancyEvent{{ID: 1, HorseID: 1, Type: "Conception"}}

		// Setup expectations
		mockDB.On("GetUserHorses", int64(1)).Return(horses, nil)
		mockDB.On("GetUserHealthRecords", int64(1)).Return(health, nil)
		mockDB.On("GetUserPregnancyEvents", int64(1)).Return(events, nil)

		// Create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sync/restore", nil)

		// Set user ID in context
		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("userID", int64(1))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.SyncData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response.Horses))
		assert.Equal(t, 1, len(response.Health))
		assert.Equal(t, 1, len(response.Events))
		mockDB.AssertExpectations(t)
	})
}

func TestGetSyncStatus(t *testing.T) {
	router, mockDB := setupTestSyncRouter()

	t.Run("successful status check", func(t *testing.T) {
		lastSync := time.Now()
		mockDB.On("GetUserLastSync", int64(1)).Return(lastSync, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/sync/status", nil)

		ctx, _ := gin.CreateTestContext(w)
		ctx.Set("userID", int64(1))
		router.HandleContext(ctx)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
		mockDB.AssertExpectations(t)
	})
}
