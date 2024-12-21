package models

import "time"

type Horse struct {
	ID             int64      `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Breed          string     `json:"breed" db:"breed"`
	DateOfBirth    time.Time  `json:"date_of_birth" db:"date_of_birth"`
	Weight         float64    `json:"weight" db:"weight"`
	ConceptionDate *time.Time `json:"conception_date,omitempty" db:"conception_date"`
	MotherID       *int64     `json:"motherId,omitempty" db:"mother_id"`
	FatherID       *int64     `json:"fatherId,omitempty" db:"father_id"`
}

type PregnancyStage string

const (
	EarlyGestation  PregnancyStage = "EARLY_GESTATION"
	MidGestation    PregnancyStage = "MID_GESTATION"
	LateGestation   PregnancyStage = "LATE_GESTATION"
	FinalGestation  PregnancyStage = "FINAL_GESTATION"
)

type HealthRecord struct {
	ID      int64     `json:"id" db:"id"`
	HorseID int64     `json:"horseId" db:"horse_id"`
	Date    time.Time `json:"date" db:"date"`
	Type    string    `json:"type" db:"type"`
	Notes   string    `json:"notes" db:"notes"`
}

type PregnancyEvent struct {
	ID          int64     `json:"id" db:"id"`
	HorseID     int64     `json:"horseId" db:"horse_id"`
	Date        time.Time `json:"date" db:"date"`
	Description string    `json:"description" db:"description"`
}
