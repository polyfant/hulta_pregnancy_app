package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// VitalSigns represents vital signs measurements for a horse
type VitalSigns struct {
	ID          uint      `json:"id"`
	HorseID     uint      `json:"horse_id"`
	Temperature float64   `json:"temperature"`
	HeartRate   uint      `json:"heart_rate"`
	Respiration uint      `json:"respiration"`
	Hydration   float64   `json:"hydration"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// VitalSignsAlert represents an alert generated from vital signs monitoring
type VitalSignsAlert struct {
	ID          uint      `json:"id"`
	HorseID     uint      `json:"horse_id"`
	AlertType   string    `json:"alert_type"`
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Severity    string    `json:"severity"`
	Parameter   string    `json:"parameter"`
	Value       float64   `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	Acknowledged bool     `json:"acknowledged"`
}

// VitalSignsTrend represents a trend analysis of vital signs
type VitalSignsTrend struct {
	ID         uint      `json:"id"`
	HorseID    uint      `json:"horse_id"`
	MetricType string    `json:"metric_type"`
	Direction  string    `json:"direction"`
	Magnitude  float64   `json:"magnitude"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	CreatedAt  time.Time `json:"created_at"`
	DataPoints int       `json:"data_points"`
	Average    float64   `json:"average"`
}

// VitalSignsPrediction represents ML predictions for vital signs
type VitalSignsPrediction struct {
	gorm.Model
	HorseID            uint           `json:"horse_id" gorm:"not null"`
	PredictedFoaling   time.Time      `json:"predicted_foaling"`
	FoalingProbability float64        `json:"foaling_probability"`
	RiskLevel          string         `json:"risk_level"`
	Alerts             pq.StringArray `json:"alerts" gorm:"type:text[]"`
	CreatedAt          time.Time      `json:"created_at"`
}

// VitalSignsThresholds defines normal ranges for vital signs
type VitalSignsThresholds struct {
	TemperatureMin float64
	TemperatureMax float64
	HeartRateMin   int
	HeartRateMax   int
	RespirationMin int
	RespirationMax int
}

// DefaultThresholds returns the default thresholds for vital signs
func DefaultThresholds() VitalSignsThresholds {
	return VitalSignsThresholds{
		TemperatureMin: 37.2, // Celsius
		TemperatureMax: 38.3,
		HeartRateMin:   28, // Beats per minute
		HeartRateMax:   44,
		RespirationMin: 8,  // Breaths per minute
		RespirationMax: 16,
	}
}
