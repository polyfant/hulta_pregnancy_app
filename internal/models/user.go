package models

import "time"

type User struct {
	ID             string         `json:"id" db:"id"`
    Email          string         `json:"email" db:"email"`
    HashedPassword string         `json:"-" db:"hashed_password"` // Already properly hidden
    LastLogin      *time.Time     `json:"last_login" db:"last_login"`
	LastSync       time.Time      `json:"last_sync" db:"last_sync"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	IsActive       bool       `json:"is_active" db:"is_active"`
	WeatherSettings WeatherSettings `json:"weather_settings" gorm:"embedded"`
}

// WeatherSettings stores user preferences for weather tracking
type WeatherSettings struct {
	NotificationsEnabled bool    `json:"notifications_enabled" gorm:"default:false"`
	DefaultLatitude     float64 `json:"default_latitude"`
	DefaultLongitude    float64 `json:"default_longitude"`
	UpdateFrequency     string  `json:"update_frequency" gorm:"default:'hourly'"` // hourly, daily, realtime
	ForecastAlerts      bool    `json:"forecast_alerts" gorm:"default:false"`
}

type UserSettings struct {
	WeatherNotificationsEnabled bool `json:"weatherNotificationsEnabled"`
	WeatherPreferences         WeatherPreferences `json:"weatherPreferences"`
}

type WeatherPreferences struct {
	Enabled        bool    `json:"enabled"`
	ForecastAlerts bool    `json:"forecastAlerts"`
	DefaultLatitude  float64 `json:"defaultLatitude"`
	DefaultLongitude float64 `json:"defaultLongitude"`
	UpdateFrequency string  `json:"updateFrequency"`
}

type SyncData struct {
	UserID    string          `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	Horses    []Horse         `json:"horses"`
	Health    []HealthRecord  `json:"health"`
	Events    []PregnancyEvent `json:"events"`
}

type UserDashboard struct {
	TotalHorses          int    `json:"total_horses"`
	PregnantHorses       int    `json:"pregnant_horses"`
	TotalExpenses        string `json:"total_expenses"`
	UpcomingEvents       int    `json:"upcoming_events"`
	ActivePregnancies    int    `json:"active_pregnancies"`
}
