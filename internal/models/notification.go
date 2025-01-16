package models

import "time"

// Notification represents a system notification
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"index"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PregnancyAlert represents a pregnancy-related alert
type PregnancyAlert struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	HorseID   uint      `json:"horse_id" gorm:"index"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Stage     string    `json:"stage"`
	CreatedAt time.Time `json:"created_at"`
}

// NotificationType represents different types of notifications
type NotificationType string

const (
	NotificationTypeWeather    NotificationType = "WEATHER"
	NotificationTypeVitalSigns NotificationType = "VITAL_SIGNS"
	NotificationTypePregnancy  NotificationType = "PREGNANCY"
	NotificationTypeSystem     NotificationType = "SYSTEM"
)

// NotificationSeverity represents the severity level of notifications
type NotificationSeverity string

const (
	NotificationSeverityInfo     NotificationSeverity = "INFO"
	NotificationSeverityWarning  NotificationSeverity = "WARNING"
	NotificationSeverityCritical NotificationSeverity = "CRITICAL"
)
