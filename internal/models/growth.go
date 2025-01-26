package models

import "time"

type GrowthData struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	FoalID          uint      `json:"foalId"`
	Age             int       `json:"age"`
	Weight          float64   `json:"weight"`
	Height          float64   `json:"height"`
	ExpectedWeight  float64   `json:"expectedWeight"`
	ExpectedHeight  float64   `json:"expectedHeight"`
	MeasurementDate time.Time `json:"measurementDate"`
}

type BodyCondition struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FoalID       uint      `json:"foalId"`
	Score        float64   `json:"score"`
	Neck         float64   `json:"neck"`
	Withers      float64   `json:"withers"`
	Loin         float64   `json:"loin"`
	Tailhead     float64   `json:"tailhead"`
	Ribs         float64   `json:"ribs"`
	Shoulder     float64   `json:"shoulder"`
	LastUpdated  time.Time `json:"lastUpdated"`
}
