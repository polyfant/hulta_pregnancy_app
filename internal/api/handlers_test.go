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

// MockDataStore implements models.DataStore for testing
type MockDataStore struct {
	mock.Mock
}

func (m *MockDataStore) GetHorse(id int64) (*models.Horse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Horse), args.Error(1)
}

func (m *MockDataStore) GetHorseByName(name string) (*models.Horse, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Horse), args.Error(1)
}

func (m *MockDataStore) GetAllHorses() ([]models.Horse, error) {
	args := m.Called()
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockDataStore) AddHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDataStore) UpdateHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDataStore) DeleteHorse(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDataStore) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func (m *MockDataStore) AddHealthRecord(record *models.HealthRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockDataStore) GetBreedingCosts(horseID int64) ([]models.BreedingCost, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.BreedingCost), args.Error(1)
}

func (m *MockDataStore) AddBreedingCost(cost *models.BreedingCost) error {
	args := m.Called(cost)
	return args.Error(0)
}

func (m *MockDataStore) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.PregnancyEvent), args.Error(1)
}

func (m *MockDataStore) AddPregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockDataStore) {
	gin.SetMode(gin.TestMode)
	mockDB := new(MockDataStore)
	handler := NewHandler(mockDB)
	router := SetupRouter(handler)
	return router, mockDB
}

func TestGetHorses(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful get all horses", func(t *testing.T) {
		horses := []models.Horse{
			{ID: 1, Name: "Thunder", Breed: "Arabian"},
			{ID: 2, Name: "Storm", Breed: "Friesian"},
		}
		mockDB.On("GetAllHorses").Return(horses, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []models.Horse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, "Thunder", response[0].Name)
	})
}

func TestGetHorse(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful get horse", func(t *testing.T) {
		horse := &models.Horse{
			ID:    1,
			Name:  "Thunder",
			Breed: "Arabian",
		}
		mockDB.On("GetHorse", int64(1)).Return(horse, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Horse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Thunder", response.Name)
	})

	t.Run("horse not found", func(t *testing.T) {
		mockDB.On("GetHorse", int64(999)).Return(nil, models.ErrNotFound).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCreateHorse(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful create horse", func(t *testing.T) {
		horse := models.Horse{
			Name:        "Thunder",
			Breed:       "Arabian",
			DateOfBirth: time.Now(),
		}
		mockDB.On("AddHorse", mock.AnythingOfType("*models.Horse")).Return(nil).Once()

		horseJSON, _ := json.Marshal(horse)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/horses", bytes.NewBuffer(horseJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.Horse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, horse.Name, response.Name)
	})

	t.Run("invalid horse data", func(t *testing.T) {
		horse := models.Horse{
			// Missing required fields
			Breed: "Arabian",
		}
		horseJSON, _ := json.Marshal(horse)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/horses", bytes.NewBuffer(horseJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetHealthAssessment(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful get health assessment", func(t *testing.T) {
		horse := &models.Horse{
			ID:          1,
			Name:        "Thunder",
			DateOfBirth: time.Now().AddDate(-5, 0, 0),
		}
		mockDB.On("GetHorse", int64(1)).Return(horse, nil).Once()
		mockDB.On("GetHealthRecords", int64(1)).Return([]models.HealthRecord{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses/1/health", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetPregnancyGuidelines(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful get pregnancy guidelines", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -3, 0)
		horse := &models.Horse{
			ID:             1,
			Name:           "Thunder",
			ConceptionDate: &conceptionDate,
		}
		mockDB.On("GetHorse", int64(1)).Return(horse, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses/1/pregnancy-guidelines", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("horse not pregnant", func(t *testing.T) {
		horse := &models.Horse{
			ID:   1,
			Name: "Thunder",
		}
		mockDB.On("GetHorse", int64(1)).Return(horse, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/horses/1/pregnancy-guidelines", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddHealthRecord(t *testing.T) {
	router, mockDB := setupTestRouter()

	t.Run("successful add health record", func(t *testing.T) {
		record := models.HealthRecord{
			HorseID: 1,
			Type:    "Checkup",
			Date:    time.Now(),
			Notes:   "Regular checkup",
		}
		mockDB.On("AddHealthRecord", mock.AnythingOfType("*models.HealthRecord")).Return(nil).Once()

		recordJSON, _ := json.Marshal(record)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/horses/1/health-records", bytes.NewBuffer(recordJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
