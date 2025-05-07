package models

import "time"

type GrowthData struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	FoalID          uint      `json:"foal_id" gorm:"not null" validate:"required"`
	Age             int       `json:"age" validate:"required,gte=0"`
	Weight          float64   `json:"weight" validate:"required,gt=0"`
	Height          float64   `json:"height" validate:"required,gt=0"`
	ExpectedWeight  float64   `json:"expected_weight,omitempty"`
	ExpectedHeight  float64   `json:"expected_height,omitempty"`
	MeasurementDate time.Time `json:"measurement_date" validate:"required"`
}

type BodyCondition struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FoalID       uint      `json:"foal_id" gorm:"not null" validate:"required"`
	Score        float64   `json:"score" validate:"required,gte=1,lte=9"`
	Neck         float64   `json:"neck,omitempty" validate:"omitempty,gte=0"`
	Withers      float64   `json:"withers,omitempty" validate:"omitempty,gte=0"`
	Loin         float64   `json:"loin,omitempty" validate:"omitempty,gte=0"`
	Tailhead     float64   `json:"tailhead,omitempty" validate:"omitempty,gte=0"`
	Ribs         float64   `json:"ribs,omitempty" validate:"omitempty,gte=0"`
	Shoulder     float64   `json:"shoulder,omitempty" validate:"omitempty,gte=0"`
	LastUpdated  time.Time `json:"last_updated" validate:"required"`
}