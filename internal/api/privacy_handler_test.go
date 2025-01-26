package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

// MockPrivacyService mocks the privacy service for testing
type MockPrivacyService struct {
	mock.Mock
}

func (m *MockPrivacyService) GetUserPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PrivacyPreferences), args.Error(1)
}

func (m *MockPrivacyService) UpdatePreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error {
	args := m.Called(ctx, userID, prefs)
	return args.Error(0)
}

func (m *MockPrivacyService) CheckFeatureEnabled(ctx context.Context, userID string, feature string) (bool, error) {
	args := m.Called(ctx, userID, feature)
	return args.Bool(0), args.Error(1)
}

func (m *MockPrivacyService) DeleteUserData(ctx context.Context, userID string, dataType string) error {
	args := m.Called(ctx, userID, dataType)
	return args.Error(0)
}

// setupTestRouter creates a test router with authentication middleware
func setupTestRouter(userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	// Add authentication middleware that sets the user ID
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	
	return r
}

func TestPrivacyHandler_GetPrivacySettings(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupAuth      bool
		mockPrefs      *models.PrivacyPreferences
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful retrieval",
			userID:    "user123",
			setupAuth: true,
			mockPrefs: &models.PrivacyPreferences{
				UserID:                "user123",
				WeatherTrackingEnabled: true,
				LocationSharingEnabled: false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "no authentication",
			userID:         "",
			setupAuth:      false,
			mockPrefs:      nil,
			mockError:      nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody: types.ErrorResponse{
				Error: "User not authenticated",
			},
		},
		{
			name:           "service error",
			userID:         "user123",
			setupAuth:      true,
			mockPrefs:      nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: types.ErrorResponse{
				Error: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPrivacyService)
			if tt.setupAuth {
				mockService.On("GetUserPreferences", mock.Anything, tt.userID).
					Return(tt.mockPrefs, tt.mockError)
			}

			handler := NewPrivacyHandler(mockService)
			router := setupTestRouter(tt.userID)
			router.GET("/privacy/settings", handler.GetPrivacySettings)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/privacy/settings", nil)
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response types.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(types.ErrorResponse).Error, response.Error)
			} else if tt.mockPrefs != nil {
				var response models.PrivacyPreferences
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPrefs.UserID, response.UserID)
				assert.Equal(t, tt.mockPrefs.WeatherTrackingEnabled, response.WeatherTrackingEnabled)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPrivacyHandler_UpdatePrivacySettings(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupAuth      bool
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful update",
			userID:    "user123",
			setupAuth: true,
			requestBody: models.PrivacyPreferences{
				WeatherTrackingEnabled: true,
				LocationSharingEnabled: false,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid request body",
			userID:         "user123",
			setupAuth:      true,
			requestBody:    "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: types.ErrorResponse{
				Error: "Invalid request body",
			},
		},
		{
			name:      "service error",
			userID:    "user123",
			setupAuth: true,
			requestBody: models.PrivacyPreferences{
				WeatherTrackingEnabled: true,
			},
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: types.ErrorResponse{
				Error: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPrivacyService)
			if tt.setupAuth && tt.name != "invalid request body" {
				mockService.On("UpdatePreferences", mock.Anything, tt.userID, mock.AnythingOfType("*models.PrivacyPreferences")).
					Return(tt.mockError)
			}

			handler := NewPrivacyHandler(mockService)
			router := setupTestRouter(tt.userID)
			router.PUT("/privacy/settings", handler.UpdatePrivacySettings)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/privacy/settings", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response types.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(types.ErrorResponse).Error, response.Error)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPrivacyHandler_RequestDataDeletion(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		setupAuth      bool
		requestBody    interface{}
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:      "successful deletion",
			userID:    "user123",
			setupAuth: true,
			requestBody: struct {
				DataTypes []string `json:"data_types"`
			}{
				DataTypes: []string{"weather_data", "health_data"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:      "partial deletion failure",
			userID:    "user123",
			setupAuth: true,
			requestBody: struct {
				DataTypes []string `json:"data_types"`
			}{
				DataTypes: []string{"weather_data"},
			},
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: types.ErrorResponse{
				Error: assert.AnError.Error(),
			},
		},
		{
			name:           "invalid request body",
			userID:         "user123",
			setupAuth:      true,
			requestBody:    "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: types.ErrorResponse{
				Error: "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPrivacyService)
			if tt.setupAuth {
				if req, ok := tt.requestBody.(struct {
					DataTypes []string `json:"data_types"`
				}); ok {
					for _, dataType := range req.DataTypes {
						mockService.On("DeleteUserData", mock.Anything, tt.userID, dataType).
							Return(tt.mockError)
					}
				}
			}

			handler := NewPrivacyHandler(mockService)
			router := setupTestRouter(tt.userID)
			router.POST("/privacy/data-deletion", handler.RequestDataDeletion)

			// Create request
			body, _ := json.Marshal(tt.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/privacy/data-deletion", bytes.NewBuffer(body))
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var response types.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(types.ErrorResponse).Error, response.Error)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPrivacyHandler_GetDataSharingStatus(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		feature        string
		setupAuth      bool
		mockEnabled    bool
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "feature enabled",
			userID:         "user123",
			feature:        "weather_tracking",
			setupAuth:      true,
			mockEnabled:    true,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"feature": "weather_tracking",
				"enabled": true,
			},
		},
		{
			name:           "feature disabled",
			userID:         "user123",
			feature:        "location_sharing",
			setupAuth:      true,
			mockEnabled:    false,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"feature": "location_sharing",
				"enabled": false,
			},
		},
		{
			name:           "service error",
			userID:         "user123",
			feature:        "health_data",
			setupAuth:      true,
			mockEnabled:    false,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody: types.ErrorResponse{
				Error: assert.AnError.Error(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockPrivacyService)
			if tt.setupAuth {
				mockService.On("CheckFeatureEnabled", mock.Anything, tt.userID, tt.feature).
					Return(tt.mockEnabled, tt.mockError)
			}

			handler := NewPrivacyHandler(mockService)
			router := setupTestRouter(tt.userID)
			router.GET("/privacy/sharing-status/:feature", handler.GetDataSharingStatus)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/privacy/sharing-status/"+tt.feature, nil)
			router.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if errResp, ok := tt.expectedBody.(types.ErrorResponse); ok {
				var response types.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, errResp.Error, response.Error)
			} else if status, ok := tt.expectedBody.(gin.H); ok {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, status["feature"], response["feature"])
				assert.Equal(t, status["enabled"], response["enabled"])
			}

			mockService.AssertExpectations(t)
		})
	}
}
