package models

import (
	"database/sql"
	"time"
)

type Transaction interface {
	Commit() error
	Rollback() error
}

type DataStore interface {
	Begin() (*sql.Tx, error)

	// Horse operations
	GetHorse(id int64) (Horse, error)
	GetHorseByName(name string) (*Horse, error)
	GetAllHorses() ([]Horse, error)
	GetUserHorses(userID int64) ([]Horse, error)
	AddHorse(horse *Horse) error
	UpdateHorse(horse *Horse) error
	UpdateHorsePregnancyStatus(horseID int64, isPregnant bool, conceptionDate time.Time) error
	DeleteHorse(id int64) error

	// Health record operations
	GetHealthRecords(horseID int64) ([]HealthRecord, error)
	GetUserHealthRecords(userID int64) ([]HealthRecord, error)
	AddHealthRecord(record *HealthRecord) error
	UpdateHealthRecord(record *HealthRecord) error
	DeleteHealthRecord(id int64) error

	// Pregnancy operations
	AddPregnancyEvent(event *PregnancyEvent) error
	GetPregnancyEvents(horseID int64) ([]PregnancyEvent, error)
	GetUserPregnancyEvents(userID int64) ([]PregnancyEvent, error)
	UpdatePregnancyEvent(event *PregnancyEvent) error
	DeletePregnancyEvent(id int64) error
	GetPreFoalingSigns(horseID int64) ([]PreFoalingSign, error)
	UpdatePreFoalingSign(sign *PreFoalingSign) error
	AddPreFoalingSign(sign *PreFoalingSign) error

	// Breeding cost operations
	GetBreedingCosts(horseID int64) ([]BreedingCost, error)
	AddBreedingCost(cost *BreedingCost) error
	UpdateBreedingCost(cost *BreedingCost) error
	DeleteBreedingCost(id int64) error

	// User operations
	UpdateUserLastSync(userID int64, time time.Time) error
	GetUserLastSync(userID int64) (time.Time, error)
	GetLastSyncTime(userID int64) (time.Time, error)
	GetPendingSyncCount(userID int64) (int, error)

	// Horse pregnancy operations
	UpdateHorseConceptionDate(horseID int64, conceptionDate time.Time) error

	// Pre-foaling checklist operations
	GetPreFoalingChecklist(horseID int64) ([]PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	UpdatePreFoalingChecklistItem(item *PreFoalingChecklistItem) error
	DeletePreFoalingChecklistItem(id int64) error
}

type BreedingCost struct {
	ID          int64     `json:"id" db:"id"`
	HorseID     int64     `json:"horseId" db:"horse_id"`
	Description string    `json:"description" db:"description"`
	Amount      float64   `json:"amount" db:"amount"`
	Date        time.Time `json:"date" db:"date"`
}

var ErrNotFound = sql.ErrNoRows
