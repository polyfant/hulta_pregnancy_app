package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)


type Pregnancy struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	HorseID        uint       `json:"horse_id"`
	StartDate      time.Time  `json:"start_date"`
	ConceptionDate *time.Time `json:"conception_date"`
	EndDate        *time.Time `json:"end_date"`
	Status         string     `json:"status"`
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
	Priority    Priority  `gorm:"size:50"`
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PregnancyGuideline struct {
	Stage       PregnancyStage `json:"stage"`
	Description string         `json:"description"`
}

type PregnancyStatus struct {
	IsPregnant          bool              `json:"is_pregnant"`
	ConceptionDate      time.Time         `json:"conception_date,omitempty"`
	DaysPregnant        int               `json:"days_pregnant"`
	PregnancyPercentage float64           `json:"pregnancy_percentage"`
	Stage               PregnancyStage    `json:"stage"`
	LastEvent           *PregnancyEvent   `json:"last_event,omitempty"`
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
	return p.StartDate.AddDate(0, 0, 340)
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