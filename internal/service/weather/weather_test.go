package weather

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveWeatherData(ctx context.Context, data *WeatherData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockRepository) GetLatestWeatherData(ctx context.Context) (*WeatherData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*WeatherData), args.Error(1)
}

// MockHTTPClient implements a mock HTTP client for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetPregnancyWeatherAdvice(t *testing.T) {
	tests := []struct {
		name           string
		weatherData    *WeatherData
		stage          string
		expectedRisk   string
		minRecsCount   int
		expectError    bool
	}{
		{
			name: "Early stage, normal conditions",
			weatherData: &WeatherData{
				Temperature: 22.0,
				WindSpeed:   5.0,
				Description: "Clear",
				Conditions:  []string{"clear"},
			},
			stage:        "EARLY_GESTATION",
			expectedRisk: "LOW",
			minRecsCount: 2,
			expectError:  false,
		},
		{
			name: "Hot weather late stage",
			weatherData: &WeatherData{
				Temperature: 32.0,
				WindSpeed:   5.0,
				Description: "Hot and sunny",
				Conditions:  []string{"clear"},
			},
			stage:        "LATE_GESTATION",
			expectedRisk: "HIGH",
			minRecsCount: 7,
			expectError:  false,
		},
		{
			name: "Cold weather with wind",
			weatherData: &WeatherData{
				Temperature: 2.0,
				WindSpeed:   18.0,
				Description: "Cold and windy",
				Conditions:  []string{"clear"},
			},
			stage:        "MID_GESTATION",
			expectedRisk: "MODERATE", // Due to high wind speed
			minRecsCount: 5,
			expectError:  false,
		},
		{
			name: "Moderate temperature early stage",
			weatherData: &WeatherData{
				Temperature: 28.0,
				WindSpeed:   8.0,
				Description: "Pleasant",
				Conditions:  []string{"clear"},
			},
			stage:        "EARLY_GESTATION",
			expectedRisk: "MODERATE", // Due to temperature ≥28°C
			minRecsCount: 4,
			expectError:  false,
		},
		{
			name: "Severe weather conditions",
			weatherData: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   25.0,
				Description: "Stormy",
				Conditions:  []string{"thunderstorm"},
			},
			stage:        "LATE_GESTATION",
			expectedRisk: "HIGH",
			minRecsCount: 6,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetLatestWeatherData", mock.Anything).Return(tt.weatherData, nil)

			service := &Service{
				apiKey:     "test-key",
				httpClient: http.DefaultClient,
				repo:       mockRepo,
			}

			advice, err := service.GetPregnancyWeatherAdvice(context.Background(), tt.stage)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedRisk, advice.RiskLevel)
			assert.GreaterOrEqual(t, len(advice.Recommendations), tt.minRecsCount)
		})
	}
}

func TestGetCurrentWeather(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  string
		mockStatus    int
		mockError     error
		expectError   bool
		validateData  func(*testing.T, *WeatherData)
	}{
		{
			name: "Successful API response",
			mockResponse: `{
				"main": {"temp": 25.0},
				"wind": {"speed": 10.0},
				"weather": [{"description": "Partly cloudy", "main": "Clouds"}]
			}`,
			mockStatus: http.StatusOK,
			mockError: nil,
			expectError: false,
			validateData: func(t *testing.T, data *WeatherData) {
				assert.Equal(t, 25.0, data.Temperature)
				assert.Equal(t, 10.0, data.WindSpeed)
				assert.Equal(t, "Partly cloudy", data.Description)
				assert.Contains(t, data.Conditions, "clouds")
			},
		},
		{
			name: "API error",
			mockResponse: "",
			mockStatus: http.StatusForbidden,
			mockError: errors.New("API error"),
			expectError: true,
			validateData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockHTTP := new(MockHTTPClient)

			// Setup mock response
			if !tt.expectError {
				// Create mock response
				response := &http.Response{
					StatusCode: tt.mockStatus,
					Body: io.NopCloser(strings.NewReader(tt.mockResponse)),
				}
				mockHTTP.On("Do", mock.Anything).Return(response, tt.mockError)
				
				// Expect SaveWeatherData to be called with any context and WeatherData
				mockRepo.On("SaveWeatherData", mock.Anything, mock.MatchedBy(func(data *WeatherData) bool {
					return true // Accept any WeatherData for now
				})).Return(nil)
			} else {
				mockHTTP.On("Do", mock.Anything).Return(nil, tt.mockError)
			}

			service := &Service{
				apiKey:     "test-key",
				httpClient: mockHTTP,
				repo:       mockRepo,
			}

			data, err := service.GetCurrentWeather(context.Background(), 0, 0)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tt.validateData != nil {
				tt.validateData(t, data)
			}

			mockRepo.AssertExpectations(t)
			mockHTTP.AssertExpectations(t)
		})
	}
}

func TestHasSevereConditions(t *testing.T) {
	tests := []struct {
		name        string
		weatherData *WeatherData
		expectSevere bool
	}{
		{
			name: "Normal conditions",
			weatherData: &WeatherData{
				Temperature: 25.0,
				WindSpeed:   10.0,
				Description: "Clear",
				Conditions:  []string{"clear"},
			},
			expectSevere: false,
		},
		{
			name: "Severe wind",
			weatherData: &WeatherData{
				Temperature: 25.0,
				WindSpeed:   25.0, // Severe wind > 20 m/s
				Description: "Windy",
				Conditions:  []string{"wind"},
			},
			expectSevere: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			severe := tt.weatherData.HasSevereConditions()
			assert.Equal(t, tt.expectSevere, severe)
		})
	}
}

func TestCalculateRiskLevel(t *testing.T) {
	tests := []struct {
		name        string
		weatherData *WeatherData
		stage       string
		expectRisk  string
	}{
		{
			name: "Low risk conditions",
			weatherData: &WeatherData{
				Temperature: 22.0,
				WindSpeed:   5.0,
				Description: "Clear",
				Conditions:  []string{"clear"},
			},
			stage:      "EARLY_GESTATION",
			expectRisk: "LOW",
		},
		{
			name: "High risk conditions",
			weatherData: &WeatherData{
				Temperature: 35.0, // Very hot
				WindSpeed:   22.0, // Strong wind
				Description: "Hot and windy",
				Conditions:  []string{"hot", "wind"},
			},
			stage:      "LATE_GESTATION",
			expectRisk: "HIGH",
		},
		{
			name: "High risk - extreme heat late stage",
			weatherData: &WeatherData{
				Temperature: 35.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:      "LATE_GESTATION",
			expectRisk: "HIGH",
		},
		{
			name: "High risk - severe wind",
			weatherData: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   25.0,
				Conditions:  []string{"clear"},
			},
			stage:      "EARLY_GESTATION",
			expectRisk: "HIGH",
		},
		{
			name: "Moderate risk - high wind",
			weatherData: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   18.0,
				Conditions:  []string{"clear"},
			},
			stage:      "MID_GESTATION",
			expectRisk: "MODERATE",
		},
		{
			name: "Moderate risk - cold weather",
			weatherData: &WeatherData{
				Temperature: 2.0,
				WindSpeed:   10.0,
				Conditions:  []string{"clear"},
			},
			stage:      "EARLY_GESTATION",
			expectRisk: "MODERATE",
		},
		{
			name: "Moderate risk - warm weather",
			weatherData: &WeatherData{
				Temperature: 28.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:      "EARLY_GESTATION",
			expectRisk: "MODERATE",
		},
		{
			name: "Low risk - pleasant conditions",
			weatherData: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:      "MID_GESTATION",
			expectRisk: "LOW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				apiKey:     "test-key",
				httpClient: http.DefaultClient,
				repo:       new(MockRepository),
			}

			risk := service.calculateRiskLevel(tt.weatherData, tt.stage)
			assert.Equal(t, tt.expectRisk, risk)
		})
	}
}

func TestGetRecommendations(t *testing.T) {
	tests := []struct {
		name        string
		weatherData *WeatherData
		stage       string
		minRecs     int
		checkRecs   []string
	}{
		{
			name: "Basic recommendations",
			weatherData: &WeatherData{
				Temperature: 22.0,
				WindSpeed:   5.0,
				Description: "Clear",
				Conditions:  []string{"clear"},
			},
			stage:     "EARLY_GESTATION",
			minRecs:   2,
			checkRecs: []string{"water", "shelter"},
		},
		{
			name: "Late stage specific recommendations",
			weatherData: &WeatherData{
				Temperature: 32.0,
				WindSpeed:   15.0,
				Description: "Hot and windy",
				Conditions:  []string{"hot", "wind"},
			},
			stage:     "LATE_GESTATION",
			minRecs:   3,
			checkRecs: []string{"monitor", "shade", "hydration"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				apiKey:     "test-key",
				httpClient: http.DefaultClient,
				repo:       new(MockRepository),
			}

			recs := service.getRecommendations(tt.weatherData, tt.stage)
			assert.GreaterOrEqual(t, len(recs), tt.minRecs, "Should have minimum number of recommendations")
			
			for _, expected := range tt.checkRecs {
				found := false
				for _, rec := range recs {
					if strings.Contains(strings.ToLower(rec), expected) {
						found = true
						break
					}
				}
				assert.True(t, found, "Should contain recommendation with '%s'", expected)
			}
		})
	}
}
