package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	PregnancyStatusActive   = "ACTIVE"
	PregnancyStatusComplete = "COMPLETE"
	PregnancyStatusLost     = "LOST"
	PregnancyStatusAborted  = "ABORTED"
)

// PregnancyStage represents different stages of pregnancy
type PregnancyStage string

const (
	PregnancyStageEarly   PregnancyStage = "EARLY_GESTATION"
	PregnancyStageMid     PregnancyStage = "MID_GESTATION"
	PregnancyStageLate    PregnancyStage = "LATE_GESTATION"
	PregnancyStageOverdue PregnancyStage = "OVERDUE"
	PregnancyStageHighRisk PregnancyStage = "HIGH_RISK"
)

// PregnancyStatus represents the status of a pregnancy
type PregnancyStatus struct {
	DaysPregnant  int       `json:"daysPregnant"`
	DueDate       time.Time `json:"dueDate"`
	NextCheckDate time.Time `json:"nextCheckDate"`
	Stage         PregnancyStage `json:"stage"`
}

// PregnancyStart represents the data needed to start pregnancy tracking
type PregnancyStart struct {
	ConceptionDate time.Time `json:"conceptionDate" binding:"required"`
}

// PregnancyGuideline represents guidelines for different pregnancy stages
type PregnancyGuideline struct {
	Stage       PregnancyStage `json:"stage"`
	Description string         `json:"description"`
	Tips        []string       `json:"tips,omitempty"`
}

// Pregnancy represents a pregnancy record
type Pregnancy struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	HorseID        uint           `json:"horseID"`
	StartDate      time.Time      `json:"startDate"`
	EndDate        *time.Time     `json:"endDate,omitempty"`
	Status         string         `json:"status"`
	ConceptionDate *time.Time     `json:"conceptionDate,omitempty"`
	Notes          string         `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type PregnancyEvent struct {
	ID           uint      `json:"id"`
	PregnancyID  uint      `json:"pregnancy_id"`
	UserID       string    `json:"user_id"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	Date         time.Time `json:"date"`
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

type PregnancyStageInfo struct {
	Stage           PregnancyStage `json:"stage"`
	DaysSoFar      int            `json:"days_so_far"`
	WeeksSoFar     int            `json:"weeks_so_far"`
	DaysRemaining  int            `json:"days_remaining"`
	WeeksRemaining int            `json:"weeks_remaining"`
	Progress       float64        `json:"progress"`
	DaysOverdue    int            `json:"days_overdue"`
	IsOverdue      bool           `json:"is_overdue"`
}

type Guideline struct {
	Stage       PregnancyStage `json:"stage"`
	Description string         `json:"description"`
	Tips        []string       `json:"tips"`
}

// Define the default checklist using the constants from constants.go
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

// Keep the methods
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
	return p.StartDate.AddDate(0, 0, DefaultGestationDays)
}

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

// DueDateInfo contains comprehensive due date information
type DueDateInfo struct {
	ExpectedDueDate time.Time `json:"expected_due_date"`
	EarliestDueDate time.Time `json:"earliest_due_date"`
	LatestDueDate   time.Time `json:"latest_due_date"`
	DaysUntilDue    int       `json:"days_until_due"`
	WeeksUntilDue   int       `json:"weeks_until_due"`
	IsInDueWindow   bool      `json:"is_in_due_window"`
} 