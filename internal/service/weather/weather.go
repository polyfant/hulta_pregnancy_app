// Package weather provides weather-related functionality for horse pregnancy management
package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Service defines the interface for weather operations
type Service interface {
	GetCurrentWeather(ctx context.Context, location string) (*models.WeatherData, error)
	GetForecast(ctx context.Context, location string, days int) ([]*models.WeatherData, error)
	GetWeatherAlerts(ctx context.Context, location string) ([]*models.WeatherAlert, error)
	GetRecommendations(ctx context.Context, weatherData *models.WeatherData) ([]*models.WeatherRecommendation, error)
	GetLatestWeatherData(ctx context.Context) (*models.WeatherData, error)
	GetPregnancyWeatherAdvice(ctx context.Context, stage string) (*models.PregnancyWeatherAdvice, error)
}

// Repository interface for weather data persistence
type Repository interface {
	SaveWeatherData(ctx context.Context, data *models.WeatherData) error
	GetLatestWeatherData(ctx context.Context) (*models.WeatherData, error)
	SaveWeatherAlert(ctx context.Context, alert *models.WeatherAlert) error
	GetWeatherAlerts(ctx context.Context, location string) ([]*models.WeatherAlert, error)
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ServiceConfig holds configuration for the weather service
type ServiceConfig struct {
	APIKey      string
	HTTPTimeout time.Duration
}

// service implements the Service interface
type service struct {
	apiKey     string
	httpClient HTTPClient
	repo       Repository
}

// NewService creates a new weather service
func NewService(cfg ServiceConfig, repo Repository, httpClient HTTPClient) *service {
	return &service{
		apiKey:     cfg.APIKey,
		httpClient: httpClient,
		repo:       repo,
	}
}

// GetCurrentWeather retrieves current weather data for a location
func (s *service) GetCurrentWeather(ctx context.Context, lat, lon float64) (*models.WeatherData, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, s.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var apiResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	weatherData := &models.WeatherData{
		Temperature: apiResp.Main.Temp,
		WindSpeed:  apiResp.Wind.Speed,
		Timestamp:  time.Now(),
	}

	// Save the weather data
	if err := s.repo.SaveWeatherData(ctx, weatherData); err != nil {
		return nil, fmt.Errorf("saving weather data: %w", err)
	}

	return weatherData, nil
}

// GetLatestWeatherData retrieves the most recent weather data from storage
func (s *service) GetLatestWeatherData(ctx context.Context) (*models.WeatherData, error) {
	return s.repo.GetLatestWeatherData(ctx)
}

// GetPregnancyWeatherAdvice provides tailored recommendations based on:
// - Current weather conditions
// - Pregnancy stage (EARLY_GESTATION, MID_GESTATION, LATE_GESTATION, OVERDUE)
// - Risk factors (temperature, wind, severe conditions)
func (s *service) GetPregnancyWeatherAdvice(ctx context.Context, stage string) (*models.PregnancyWeatherAdvice, error) {
	weather, err := s.GetLatestWeatherData(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting latest weather: %w", err)
	}

	riskLevel := s.calculateRiskLevel(weather, stage)
	recommendations := s.getRecommendations(weather, stage)

	return &models.PregnancyWeatherAdvice{
		WeatherData:     weather,
		Recommendations: recommendations,
		RiskLevel:       riskLevel,
		StageSpecific:   true,
	}, nil
}

// WeatherAPIResponse represents the structure of the weather API response
type WeatherAPIResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
}

// calculateRiskLevel determines risk based on weather and pregnancy stage:
// HIGH risk:
// - Temperature ≥ 30°C in late stage
// - Temperature ≥ 35°C any stage
// - Wind speed > 20 m/s
// MODERATE risk:
// - Wind speed > 15 m/s
// - Temperature ≥ 28°C or ≤ 5°C
// - Temperature ≥ 25°C in late stage
func (s *service) calculateRiskLevel(weather *models.WeatherData, stage string) string {
	if weather.WindSpeed > 20 {
		return "HIGH"
	}

	if weather.Temperature >= 35 {
		return "HIGH"
	}

	if stage == "LATE_GESTATION" || stage == "OVERDUE" {
		if weather.Temperature >= 30 {
			return "HIGH"
		}
		if weather.Temperature >= 25 {
			return "MODERATE"
		}
	}

	if weather.WindSpeed > 15 || weather.Temperature >= 28 || weather.Temperature <= 5 {
		return "MODERATE"
	}

	return "LOW"
}

// getRecommendations generates advice based on conditions:
// - Base recommendations (water, food)
// - Temperature-specific (heat/cold management)
// - Wind-specific (shelter, protection)
// - Stage-specific (late pregnancy care)
// - Emergency recommendations for severe conditions
func (s *service) getRecommendations(data *models.WeatherData, stage string) []string {
	var recommendations []string

	// Base recommendations
	recommendations = append(recommendations, "Ensure fresh water is always available")
	recommendations = append(recommendations, "Monitor food and water intake")

	// Temperature recommendations
	if data.Temperature >= 28 {
		recommendations = append(recommendations,
			"Provide shade and ventilation",
			"Consider using fans or misting systems",
			"Monitor for signs of heat stress",
		)
	} else if data.Temperature <= 5 {
		recommendations = append(recommendations,
			"Provide adequate shelter from cold",
			"Consider using blankets if necessary",
			"Ensure adequate feed for warmth",
		)
	}

	// Wind recommendations
	if data.WindSpeed > 15 {
		recommendations = append(recommendations,
			"Ensure access to wind shelter",
			"Monitor for signs of stress",
		)
	}

	// Stage-specific recommendations
	if stage == "LATE_GESTATION" || stage == "OVERDUE" {
		recommendations = append(recommendations,
			"Monitor more frequently",
			"Ensure easy access to shelter",
			"Keep stress levels minimal",
		)
	}

	// Emergency conditions
	if data.WindSpeed > 20 || data.Temperature >= 35 {
		recommendations = append(recommendations,
			"URGENT: Move to protected shelter",
			"Contact veterinarian if signs of distress",
			"Monitor vital signs frequently",
		)
	}

	return recommendations
}

// PregnancyWeatherAdvice contains recommendations based on weather and pregnancy stage
type PregnancyWeatherAdvice struct {
	WeatherData     *models.WeatherData `json:"weather_data"`
	Recommendations []string             `json:"recommendations"`
	RiskLevel       string               `json:"risk_level"`    // HIGH, MODERATE, or LOW
	StageSpecific   bool                 `json:"stage_specific"`
}
