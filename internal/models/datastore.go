package models

import (
	"time"
)

type Horse struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index;not null"`
	Name      string    `gorm:"size:255;not null"`
	Breed     string    `gorm:"size:255"`
	Gender    string    `gorm:"size:50"`
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Pregnancy struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index;not null"`
	HorseID   uint      `gorm:"not null"`
	StartDate time.Time
	EndDate   *time.Time
	Status    string    `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Add other models as needed

// DataStore defines the interface for database operations
type DataStore interface {
	// Horse operations
	GetHorse(id int64) (Horse, error)
	GetAllHorses() ([]Horse, error)
	AddHorse(horse *Horse) error
	UpdateHorse(horse *Horse) error
	DeleteHorse(id int64) error

	// Health record operations
	GetHealthRecords(horseID int64) ([]HealthRecord, error)
	AddHealthRecord(record *HealthRecord) error

	// Pregnancy operations
	GetPregnancyEvents(horseID int64) ([]PregnancyEvent, error)
	AddPregnancyEvent(event *PregnancyEvent) error

	// Pre-foaling operations
	GetPreFoalingSigns(horseID int64) ([]PreFoalingSign, error)
	AddPreFoalingSign(sign *PreFoalingSign) error

	// Pre-foaling checklist operations
	GetPreFoalingChecklist(horseID int64) ([]PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	UpdatePreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	DeletePreFoalingChecklistItem(id int64) error
}
