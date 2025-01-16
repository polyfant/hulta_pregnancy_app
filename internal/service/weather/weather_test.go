package weather

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetPregnancyWeatherAdvice(t *testing.T) {
	tests := []struct {
		name           string
		weatherData    *WeatherData
		stage          string
		expectedRisk   string
		minRecsCount   int
		expectSevere   bool
		expectedAdvice []string
	}{
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
			expectedAdvice: []string{
				"Provide extra water sources",
				"Add electrolytes to water if needed",
				"Monitor more frequently",
			},
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
			expectedAdvice: []string{
				"Provide extra hay for warmth",
				"Ensure wind breaks are available",
			},
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
			expectedAdvice: []string{
				"Ensure constant access to fresh water",
				"Monitor food intake",
			},
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
			expectSevere: true,
			expectedAdvice: []string{
				"Consider indoor housing",
				"Contact vet if signs of distress",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("GetLatestWeatherData", mock.Anything).Return(tt.weatherData, nil)

			service := &Service{
				repo: mockRepo,
			}

			advice, err := service.GetPregnancyWeatherAdvice(context.Background(), tt.stage)
			assert.NoError(t, err)
			assert.NotNil(t, advice)

			// Check risk level
			assert.Equal(t, tt.expectedRisk, advice.RiskLevel, "Risk level should match")

			// Check minimum number of recommendations
			assert.GreaterOrEqual(t, len(advice.Recommendations), tt.minRecsCount, 
				"Should have minimum number of recommendations")

			// Check for expected advice items
			for _, expectedAdvice := range tt.expectedAdvice {
				found := false
				for _, actualAdvice := range advice.Recommendations {
					if strings.Contains(strings.ToLower(actualAdvice), 
						strings.ToLower(expectedAdvice)) {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected advice not found: %s", expectedAdvice)
			}

			// Verify weather data
			assert.Equal(t, tt.weatherData, advice.WeatherData)

			// Check if severe conditions are properly detected
			if tt.expectSevere {
				assert.True(t, tt.weatherData.HasSevereConditions())
			}
		})
	}
}

func TestCalculateRiskLevel(t *testing.T) {
	service := &Service{}
	
	tests := []struct {
		name         string
		weather      *WeatherData
		stage        string
		expectedRisk string
	}{
		{
			name: "High risk - extreme heat late stage",
			weather: &WeatherData{
				Temperature: 35.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:        "LATE_GESTATION",
			expectedRisk: "HIGH",
		},
		{
			name: "High risk - severe wind",
			weather: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   25.0,
				Conditions:  []string{"clear"},
			},
			stage:        "EARLY_GESTATION",
			expectedRisk: "HIGH",
		},
		{
			name: "Moderate risk - high wind",
			weather: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   18.0,
				Conditions:  []string{"clear"},
			},
			stage:        "MID_GESTATION",
			expectedRisk: "MODERATE",
		},
		{
			name: "Moderate risk - cold weather",
			weather: &WeatherData{
				Temperature: 2.0,
				WindSpeed:   10.0,
				Conditions:  []string{"clear"},
			},
			stage:        "EARLY_GESTATION",
			expectedRisk: "MODERATE",
		},
		{
			name: "Moderate risk - warm weather",
			weather: &WeatherData{
				Temperature: 28.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:        "EARLY_GESTATION",
			expectedRisk: "MODERATE",
		},
		{
			name: "Low risk - pleasant conditions",
			weather: &WeatherData{
				Temperature: 20.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage:        "MID_GESTATION",
			expectedRisk: "LOW",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			risk := service.calculateRiskLevel(tt.weather, tt.stage)
			assert.Equal(t, tt.expectedRisk, risk, 
				"Risk level should match for %s", tt.name)
		})
	}
}

func TestGetRecommendations(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name               string
		weather           *WeatherData
		stage             string
		expectedContains  []string
		notExpectedContains []string
	}{
		{
			name: "Hot weather recommendations",
			weather: &WeatherData{
				Temperature: 30.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage: "MID_GESTATION",
			expectedContains: []string{
				"water",
				"shade",
				"electrolytes",
			},
			notExpectedContains: []string{
				"frozen",
				"warm mash",
			},
		},
		{
			name: "Cold weather recommendations",
			weather: &WeatherData{
				Temperature: 2.0,
				WindSpeed:   8.0,
				Conditions:  []string{"clear"},
			},
			stage: "EARLY_GESTATION",
			expectedContains: []string{
				"hay",
				"warm",
			},
			notExpectedContains: []string{
				"shade",
				"cool water",
			},
		},
		{
			name: "Late stage specific recommendations",
			weather: &WeatherData{
				Temperature: 22.0,
				WindSpeed:   5.0,
				Conditions:  []string{"clear"},
			},
			stage: "LATE_GESTATION",
			expectedContains: []string{
				"Monitor",
				"vet",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recs := service.getRecommendations(tt.weather, tt.stage)

			// Check for expected content
			for _, expected := range tt.expectedContains {
				found := false
				for _, rec := range recs {
					if strings.Contains(strings.ToLower(rec), 
						strings.ToLower(expected)) {
						found = true
						break
					}
				}
				assert.True(t, found, 
					"Should contain recommendation with: %s", expected)
			}

			// Check that unwanted recommendations are not present
			for _, notExpected := range tt.notExpectedContains {
				for _, rec := range recs {
					assert.False(t, strings.Contains(strings.ToLower(rec), 
						strings.ToLower(notExpected)),
						"Should not contain recommendation with: %s", notExpected)
				}
			}
		})
	}
}
