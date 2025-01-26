package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeatherData represents weather conditions
type WeatherData struct {
	Temperature float64
	WindSpeed   float64
	Description string
	Conditions  []string
}

// HasSevereConditions checks if there are any severe weather conditions
func (w *WeatherData) HasSevereConditions() bool {
	if w.WindSpeed > 20 {
		return true
	}
	for _, condition := range w.Conditions {
		switch condition {
		case "thunderstorm", "tornado", "hurricane", "blizzard":
			return true
		}
	}
	return false
}

// Repository interface for weather data persistence
type Repository interface {
	SaveWeatherData(ctx context.Context, data *WeatherData) error
	GetLatestWeatherData(ctx context.Context) (*WeatherData, error)
}

// Service provides weather-related functionality
type Service struct {
	apiKey     string
	httpClient *http.Client
	repo       Repository
}

// ServiceConfig holds configuration for the weather service
type ServiceConfig struct {
	APIKey      string
	HTTPTimeout time.Duration
}

// NewService creates a new weather service
func NewService(cfg ServiceConfig, repo Repository) *Service {
	return &Service{
		apiKey: cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
		repo: repo,
	}
}

// GetCurrentWeather retrieves current weather data for a location
func (s *Service) GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*WeatherData, error) {
	// Build API URL with API key and coordinates
	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/current.json?key=%s&q=%.6f,%.6f",
		s.apiKey,
		latitude,
		longitude,
	)

	// Make HTTP request
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

	// Parse response
	var apiResp struct {
		Current struct {
			TempC     float64 `json:"temp_c"`
			WindKph   float64 `json:"wind_kph"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Convert wind speed from km/h to m/s
	windSpeed := apiResp.Current.WindKph * 0.277778

	weatherData := &WeatherData{
		Temperature: apiResp.Current.TempC,
		WindSpeed:   windSpeed,
		Description: apiResp.Current.Condition.Text,
		Conditions:  []string{apiResp.Current.Condition.Text},
	}

	// Save weather data if repository is available
	if s.repo != nil {
		if err := s.repo.SaveWeatherData(ctx, weatherData); err != nil {
			// Log error but don't fail the request
			fmt.Printf("failed to save weather data: %v\n", err)
		}
	}

	return weatherData, nil
}

// GetLatestWeatherData retrieves the most recent weather data from storage
func (s *Service) GetLatestWeatherData(ctx context.Context) (*WeatherData, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("no repository configured")
	}

	return s.repo.GetLatestWeatherData(ctx)
}
