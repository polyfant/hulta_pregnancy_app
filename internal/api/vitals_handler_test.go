package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/vitals"
	"github.com/polyfant/hulta_pregnancy_app/internal/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockVitalsRepository struct {
	mock.Mock
}

func (m *MockVitalsRepository) SaveVitalSigns(ctx context.Context, vitals *models.VitalSigns) error {
	args := m.Called(ctx, vitals)
	return args.Error(0)
}

func (m *MockVitalsRepository) GetVitalSigns(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSigns, error) {
	args := m.Called(ctx, horseID, from, to)
	return args.Get(0).([]*models.VitalSigns), args.Error(1)
}

func (m *MockVitalsRepository) GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error) {
	args := m.Called(ctx, horseID)
	return args.Get(0).(*models.VitalSigns), args.Error(1)
}

func (m *MockVitalsRepository) SaveAlert(ctx context.Context, alert *models.VitalSignsAlert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

func (m *MockVitalsRepository) GetAlert(ctx context.Context, alertID uint) (*models.VitalSignsAlert, error) {
	args := m.Called(ctx, alertID)
	return args.Get(0).(*models.VitalSignsAlert), args.Error(1)
}

func (m *MockVitalsRepository) GetAlerts(ctx context.Context, horseID uint, includeAcknowledged bool) ([]*models.VitalSignsAlert, error) {
	args := m.Called(ctx, horseID, includeAcknowledged)
	return args.Get(0).([]*models.VitalSignsAlert), args.Error(1)
}

func (m *MockVitalsRepository) AcknowledgeAlert(ctx context.Context, alertID uint) error {
	args := m.Called(ctx, alertID)
	return args.Error(0)
}

func (m *MockVitalsRepository) SaveTrend(ctx context.Context, trend *models.VitalSignsTrend) error {
	args := m.Called(ctx, trend)
	return args.Error(0)
}

func (m *MockVitalsRepository) GetTrends(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSignsTrend, error) {
	args := m.Called(ctx, horseID, from, to)
	return args.Get(0).([]*models.VitalSignsTrend), args.Error(1)
}

func TestVitalsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(MockVitalsRepository)
	hub := websocket.NewHub()
	service := vitals.NewService(mockRepo, hub)
	handler := NewVitalsHandler(service)

	t.Run("HandleRecordVitalSigns", func(t *testing.T) {
		testCases := []struct {
			name           string
			vitals         *models.VitalSigns
			setupMock      func()
			expectedStatus int
		}{
			{
				name: "Valid vital signs",
				vitals: &models.VitalSigns{
					HorseID:     1,
					Temperature: 37.5,
					HeartRate:   32,
					Respiration: 10,
					RecordedAt:  time.Now(),
				},
				setupMock: func() {
					mockRepo.On("SaveVitalSigns", mock.Anything, mock.AnythingOfType("*models.VitalSigns")).
						Return(nil).Once()
				},
				expectedStatus: http.StatusOK,
			},
			{
				name: "Invalid temperature",
				vitals: &models.VitalSigns{
					HorseID:     1,
					Temperature: 45.0, // Invalid temperature
					HeartRate:   32,
					Respiration: 10,
					RecordedAt:  time.Now(),
				},
				setupMock:      func() {},
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.setupMock()

				// Create request
				body, err := json.Marshal(tc.vitals)
				require.NoError(t, err)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = httptest.NewRequest(http.MethodPost, "/api/vitals", bytes.NewBuffer(body))

				// Call handler
				handler.HandleRecordVitalSigns(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("HandleGetVitalSigns", func(t *testing.T) {
		// Setup test data
		horseID := uint(1)
		now := time.Now()
		from := now.Add(-24 * time.Hour)
		to := now

		vitals := []*models.VitalSigns{
			{
				HorseID:     horseID,
				Temperature: 37.5,
				HeartRate:   32,
				Respiration: 10,
				RecordedAt:  now.Add(-2 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 37.8,
				HeartRate:   34,
				Respiration: 12,
				RecordedAt:  now,
			},
		}

		mockRepo.On("GetVitalSigns", mock.Anything, horseID, from, to).
			Return(vitals, nil).Once()

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodGet, "/api/vitals/horse/1", nil)
		q := c.Request.URL.Query()
		q.Add("from", from.Format(time.RFC3339))
		q.Add("to", to.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		// Call handler
		handler.HandleGetVitalSigns(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.VitalSigns
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)

		mockRepo.AssertExpectations(t)
	})

	t.Run("HandleGetLatestVitalSigns", func(t *testing.T) {
		horseID := uint(1)
		vitals := &models.VitalSigns{
			HorseID:     horseID,
			Temperature: 37.8,
			HeartRate:   34,
			Respiration: 12,
			RecordedAt:  time.Now(),
		}

		mockRepo.On("GetLatestVitalSigns", mock.Anything, horseID).
			Return(vitals, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodGet, "/api/vitals/horse/1/latest", nil)

		handler.HandleGetLatestVitalSigns(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.VitalSigns
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, vitals.Temperature, response.Temperature)

		mockRepo.AssertExpectations(t)
	})

	t.Run("HandleGetAlerts", func(t *testing.T) {
		horseID := uint(1)
		alerts := []*models.VitalSignsAlert{
			{
				HorseID:   horseID,
				Type:      "temperature_high",
				Message:   "Temperature above normal range",
				CreatedAt: time.Now(),
				Severity:  "warning",
				Parameter: "temperature",
				Value:     39.5,
			},
		}

		mockRepo.On("GetAlerts", mock.Anything, horseID, false).
			Return(alerts, nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodGet, "/api/vitals/horse/1/alerts", nil)

		handler.HandleGetAlerts(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.VitalSignsAlert
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 1)

		mockRepo.AssertExpectations(t)
	})

	t.Run("HandleAcknowledgeAlert", func(t *testing.T) {
		alertID := uint(1)
		mockRepo.On("AcknowledgeAlert", mock.Anything, alertID).
			Return(nil).Once()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodPut, "/api/vitals/alerts/1/acknowledge", nil)

		handler.HandleAcknowledgeAlert(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockRepo.AssertExpectations(t)
	})

	t.Run("HandleGetTrends", func(t *testing.T) {
		horseID := uint(1)
		now := time.Now()
		from := now.Add(-24 * time.Hour)
		to := now

		// Mock vital signs data
		vitals := []*models.VitalSigns{
			{
				HorseID:     horseID,
				Temperature: 37.5,
				HeartRate:   32,
				Respiration: 10,
				RecordedAt:  now.Add(-2 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 37.8,
				HeartRate:   34,
				Respiration: 12,
				RecordedAt:  now,
			},
		}

		// Mock expected trends
		expectedTrends := []*models.VitalSignsTrend{
			{
				HorseID:    horseID,
				MetricType: "heart_rate",
				Direction:  "increasing",
				Magnitude:  1.0,
				StartTime:  from,
				EndTime:    to,
			},
			{
				HorseID:    horseID,
				MetricType: "temperature",
				Direction:  "increasing",
				Magnitude:  0.15,
				StartTime:  from,
				EndTime:    to,
			},
			{
				HorseID:    horseID,
				MetricType: "respiratory_rate",
				Direction:  "increasing",
				Magnitude:  1.0,
				StartTime:  from,
				EndTime:    to,
			},
		}

		// Setup mocks
		mockRepo.On("GetVitalSigns", mock.Anything, horseID, from, to).
			Return(vitals, nil).Once()

		mockRepo.On("GetTrends", mock.Anything, horseID, from, to).
			Return(expectedTrends, nil).Once()

		// Create request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{
			{Key: "id", Value: "1"},
		}
		c.Request = httptest.NewRequest(http.MethodGet, "/api/vitals/horse/1/trends", nil)
		q := c.Request.URL.Query()
		q.Add("from", from.Format(time.RFC3339))
		q.Add("to", to.Format(time.RFC3339))
		c.Request.URL.RawQuery = q.Encode()

		// Call handler
		handler.HandleGetTrends(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.VitalSignsTrend
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 3)

		// Verify trends
		for i, trend := range response {
			assert.Equal(t, expectedTrends[i].MetricType, trend.MetricType)
			assert.Equal(t, expectedTrends[i].Direction, trend.Direction)
			assert.Equal(t, expectedTrends[i].HorseID, trend.HorseID)
			assert.Equal(t, expectedTrends[i].Magnitude, trend.Magnitude)
		}

		mockRepo.AssertExpectations(t)
	})
}
