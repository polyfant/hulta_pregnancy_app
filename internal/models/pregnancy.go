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

// RiskLevel represents the risk assessment for a pregnancy
type RiskLevel string

const (
	LowRisk       RiskLevel = "LOW"
	MediumRisk    RiskLevel = "MEDIUM"
	HighRisk      RiskLevel = "HIGH"
	CriticalRisk  RiskLevel = "CRITICAL"
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

// Additional stage constants for internal use
const (
	EarlyStage   PregnancyStage = "EARLY_GESTATION"
	MidStage     PregnancyStage = "MID_GESTATION"
	LateStage    PregnancyStage = "LATE_GESTATION"
	OverdueStage PregnancyStage = "OVERDUE"
)

// PregnancyStageInfo provides comprehensive pregnancy stage details
type PregnancyStageInfo struct {
	Stage             PregnancyStage `json:"stage"`
	Progress          float64        `json:"progress"`
	Description       string         `json:"description"`
	RiskLevel         RiskLevel      `json:"risk_level"`
	DaysSoFar         int            `json:"days_so_far"`
	WeeksSoFar        int            `json:"weeks_so_far"`
	DaysRemaining     int            `json:"days_remaining"`
	WeeksRemaining    int            `json:"weeks_remaining"`
	DaysOverdue       int            `json:"days_overdue"`
	IsOverdue         bool           `json:"is_overdue"`
	NutritionAdvice   string         `json:"nutrition_advice"`
	MonitoringAdvice  string         `json:"monitoring_advice"`
}

// PregnancyStatus represents the status of a pregnancy with enhanced tracking
type PregnancyStatus struct {
	DaysPregnant     int            `json:"daysPregnant"`
	DueDate          time.Time      `json:"dueDate"`
	NextCheckDate    time.Time      `json:"nextCheckDate"`
	Stage            PregnancyStage `json:"stage"`
	RiskLevel        RiskLevel      `json:"riskLevel"`
	Progress         float64        `json:"progress"`
	StageDescription string         `json:"stageDescription"`
}

// Pregnancy model is updated to include more comprehensive tracking
type Pregnancy struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	HorseID            uint           `json:"horseID"`
	StartDate          time.Time      `json:"startDate"`
	EndDate            *time.Time     `json:"endDate,omitempty"`
	Status             string         `json:"status"`
	ConceptionDate     *time.Time     `json:"conceptionDate,omitempty"`
	CurrentStage       PregnancyStage `json:"currentStage"`
	RiskLevel          RiskLevel      `json:"riskLevel"`
	ProgressPercentage float64        `json:"progressPercentage"`
	Notes              string         `json:"notes,omitempty"`
	CreatedAt          time.Time      `json:"createdAt"`
	UpdatedAt          time.Time      `json:"updatedAt"`
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

// PregnancyEvent represents a pregnancy event
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

// PreFoalingSign represents a pre-foaling sign
type PreFoalingSign struct {
	ID          uint      `gorm:"primaryKey"`
	HorseID     uint      `gorm:"not null"`
	Description string
	Date        time.Time
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PreFoalingChecklistItem represents a pre-foaling checklist item
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

// Guideline represents a guideline
type Guideline struct {
	Stage       PregnancyStage `json:"stage"`
	Description string         `json:"description"`
	Tips        []string       `json:"tips"`
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

// Additional methods for Pregnancy model
func (p *Pregnancy) GetStageInfo() *PregnancyStageInfo {
	if p.ConceptionDate == nil {
		return nil
	}

	daysSinceConception := int(time.Since(*p.ConceptionDate).Hours() / 24)
	
	var stage PregnancyStage
	var riskLevel RiskLevel
	var description string

	switch {
	case daysSinceConception <= 120:
		stage = EarlyStage
		riskLevel = LowRisk
		description = "Early pregnancy: Embryo development and initial health monitoring"
	case daysSinceConception <= 270:
		stage = MidStage
		riskLevel = MediumRisk
		description = "Mid pregnancy: Fetal growth and critical nutritional period"
	case daysSinceConception <= 340:
		stage = LateStage
		riskLevel = HighRisk
		description = "Late pregnancy: Preparing for foaling, increased monitoring required"
	default:
		stage = OverdueStage
		riskLevel = CriticalRisk
		description = "Overdue: Immediate veterinary consultation recommended"
	}

	return &PregnancyStageInfo{
		Stage:             stage,
		Progress:          p.ProgressPercentage,
		Description:       description,
		RiskLevel:         riskLevel,
		DaysSoFar:         daysSinceConception,
		WeeksSoFar:        daysSinceConception / 7,
		DaysRemaining:     340 - daysSinceConception,
		WeeksRemaining:    (340 - daysSinceConception) / 7,
		DaysOverdue:       max(0, daysSinceConception - 340),
		IsOverdue:         daysSinceConception > 340,
		NutritionAdvice:   getNutritionAdvice(stage),
		MonitoringAdvice:  getMonitoringAdvice(stage),
	}
}

// Helper functions for nutrition and monitoring advice
func getNutritionAdvice(stage PregnancyStage) string {
	switch stage {
	case EarlyStage:
		return "Focus on balanced diet, maintain body condition. Supplement with folic acid and minerals."
	case MidStage:
		return "Increase protein and energy intake. Monitor weight gain carefully."
	case LateStage:
		return "High-quality protein, increased calories. Prepare for increased nutritional demands."
	case OverdueStage:
		return "Consult veterinarian. Specialized nutrition may be required."
	default:
		return "Maintain standard pregnancy nutrition protocol."
	}
}

func getMonitoringAdvice(stage PregnancyStage) string {
	switch stage {
	case EarlyStage:
		return "Weekly health checks. Monitor for early signs of complications."
	case MidStage:
		return "Bi-weekly detailed health assessments. Ultrasound recommended."
	case LateStage:
		return "Weekly veterinary check-ups. Monitor for pre-foaling signs."
	case OverdueStage:
		return "Immediate and frequent veterinary monitoring. Prepare for potential intervention."
	default:
		return "Standard pregnancy monitoring protocol."
	}
}

// Utility function to get max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}