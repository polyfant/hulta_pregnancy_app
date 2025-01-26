package models

import (
	"time"
)

// WeatherData represents environmental conditions that may affect horse health
type WeatherData struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        string    `json:"user_id"`
	LocationID    uint      `json:"location_id"`
	Timestamp     time.Time `json:"timestamp"`
	Temperature   float64   `json:"temperature"` // in Celsius
	Humidity      float64   `json:"humidity"`    // percentage
	WindSpeed     float64   `json:"wind_speed"`  // in m/s
	Precipitation float64   `json:"precipitation"` // in mm
	Pressure      float64   `json:"pressure"`    // in hPa
	RainAmount    float64   `json:"rain_amount"` // in mm
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// WeatherImpact represents the calculated impact of weather on horse health
type WeatherImpact struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	HorseID         uint      `json:"horse_id"`
	WeatherDataID   uint      `json:"weather_data_id"`
	WeatherData     WeatherData `json:"weather_data" gorm:"foreignKey:WeatherDataID"`
	StressLevel     float64   `json:"stress_level"`      // 0-100 scale
	ComfortIndex    float64   `json:"comfort_index"`     // 0-100 scale
	ExerciseSafe    bool      `json:"exercise_safe"`     // whether it's safe to exercise
	Recommendations []string  `json:"recommendations"`   // JSON array of recommendations
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Location represents a geographical location for weather tracking
type Location struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"` // in meters
}

// TemperatureThresholds defines safe temperature ranges for horses
type TemperatureThresholds struct {
	MinSafe     float64 // minimum safe temperature
	MaxSafe     float64 // maximum safe temperature
	MinCritical float64 // minimum critical temperature
	MaxCritical float64 // maximum critical temperature
}

// GetDefaultThresholds returns default temperature thresholds for horses
func GetDefaultThresholds() TemperatureThresholds {
	return TemperatureThresholds{
		MinSafe:     -5.0,  // Celsius
		MaxSafe:     25.0,  // Celsius
		MinCritical: -15.0, // Celsius
		MaxCritical: 30.0,  // Celsius
	}
}

// CalculateStressIndex calculates environmental stress on the horse
func (w *WeatherData) CalculateStressIndex() float64 {
	// Temperature-Humidity Index (THI) calculation
	// THI = (1.8 × T + 32) - (0.55 - 0.0055 × RH) × (1.8 × T - 26)
	// where T = temperature in Celsius, RH = relative humidity
	t := w.Temperature
	rh := w.Humidity
	
	tF := 1.8*t + 32
	thi := tF - (0.55-0.0055*rh)*(1.8*t-26)
	
	// Normalize to 0-100 scale
	normalizedTHI := (thi - 50) * 2
	if normalizedTHI < 0 {
		normalizedTHI = 0
	}
	if normalizedTHI > 100 {
		normalizedTHI = 100
	}
	
	return normalizedTHI
}
