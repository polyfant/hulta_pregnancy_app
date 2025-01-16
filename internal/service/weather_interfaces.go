package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// WeatherService defines the interface for weather-related operations
type WeatherService interface {
	GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*models.WeatherData, error)
	GetLatestWeatherData(ctx context.Context) (*models.WeatherData, error)
	GetPregnancyWeatherAdvice(ctx context.Context, stage string) (*models.PregnancyWeatherAdvice, error)
}
