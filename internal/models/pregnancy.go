package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Pregnancy struct {
	ID         uint       `gorm:"primaryKey"`
	HorseID    uint       `gorm:"not null"`
	UserID     string     `gorm:"index;not null"`
	StartDate  time.Time
	EndDate    *time.Time
	Status     string     `gorm:"size:50"`
	Events     []PregnancyEvent
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type PregnancyEvent struct {
	ID           uint      `gorm:"primaryKey"`
	PregnancyID  uint      `gorm:"not null"`
	Type         string    `gorm:"size:50"`
	Description  string
	Date         time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PreFoalingSign struct {
	ID          uint      `gorm:"primaryKey"`
	HorseID     uint      `gorm:"not null"`
	Description string
	Date        time.Time
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PreFoalingChecklistItem struct {
	ID          uint      `gorm:"primaryKey"`
	HorseID     uint      `gorm:"not null"`
	Description string
	IsCompleted bool      `gorm:"default:false"`
	DueDate     time.Time
	Priority    string    `gorm:"size:50"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Constants for pregnancy status
const (
	PregnancyStatusActive    = "ACTIVE"
	PregnancyStatusCompleted = "COMPLETED"
	PregnancyStatusLost     = "LOST"
)

// Constants for priority levels
const (
	PriorityHigh   = "HIGH"
	PriorityMedium = "MEDIUM"
	PriorityLow    = "LOW"
)

// Add these types and constants
type PregnancyStage string

const (
	EarlyGestation PregnancyStage = "EARLY_GESTATION"
	MidGestation   PregnancyStage = "MID_GESTATION"
	LateGestation  PregnancyStage = "LATE_GESTATION"
	PreFoaling    PregnancyStage = "PRE_FOALING"
	Foaling       PregnancyStage = "FOALING"
)

// Add these constants
var DefaultPreFoalingChecklist = []PreFoalingChecklistItem{
	{
		Description: "Prepare foaling kit",
		Priority:    PriorityHigh,
	},
	{
		Description: "Clean and prepare foaling stall",
		Priority:    PriorityHigh,
	},
	{
		Description: "Check emergency vet contacts",
		Priority:    PriorityHigh,
	},
	{
		Description: "Monitor mare's temperature",
		Priority:    PriorityMedium,
	},
	{
		Description: "Check udder development",
		Priority:    PriorityMedium,
	},
}

// Add these constants
const (
	EventFoaling = "FOALING"
	EventVetCheck = "VET_CHECK"
	EventUltrasound = "ULTRASOUND"
	EventVaccination = "VACCINATION"
	EventDeworming = "DEWORMING"
	EventComplication = "COMPLICATION"
)

// Add these methods to the Pregnancy struct
func (p *Pregnancy) IsActive() bool {
	return p.Status == PregnancyStatusActive
}

func (p *Pregnancy) DaysPregnant() int {
	if p.StartDate.IsZero() {
		return 0
	}
	return int(time.Since(p.StartDate).Hours() / 24)
}

func (p *Pregnancy) ExpectedDueDate() time.Time {
	return p.StartDate.AddDate(0, 0, 340) // Average gestation period
}

// Add these methods to PregnancyEvent

func (e *PregnancyEvent) Validate() error {
	if e.PregnancyID == 0 {
		return fmt.Errorf("pregnancy ID is required")
	}
	if e.Description == "" {
		return fmt.Errorf("description is required")
	}
	if e.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	return nil
}

func (e *PregnancyEvent) BeforeCreate(tx *gorm.DB) error {
	if e.Date.IsZero() {
		e.Date = time.Now()
	}
	return nil
} 