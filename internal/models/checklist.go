package models

import (
	"time"
)

// Season represents different seasons for checklist customization
type Season string

const (
	SeasonSpring Season = "SPRING"
	SeasonSummer Season = "SUMMER"
	SeasonFall   Season = "FALL"
	SeasonWinter Season = "WINTER"
	SeasonAll    Season = "ALL"
)

// PreFoalingChecklistItem represents a single item in the pre-foaling checklist
type PreFoalingChecklistItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	HorseID     uint      `json:"horse_id"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	DueDate     time.Time `json:"due_date"`
	Notes       string    `json:"notes"`
	Completed   bool      `json:"completed"`
	Season      Season    `json:"season"`
	IsRequired  bool      `json:"is_required"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SeasonalChecklistTemplate defines season-specific checklist items
type SeasonalChecklistTemplate struct {
	Season      Season                    `json:"season"`
	Items       []PreFoalingChecklistItem `json:"items"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

// ChecklistCategory represents different categories of checklist items
type ChecklistCategory string

const (
	CategoryVeterinary  ChecklistCategory = "VETERINARY"
	CategoryNutrition   ChecklistCategory = "NUTRITION"
	CategoryEnvironment ChecklistCategory = "ENVIRONMENT"
	CategorySupplies    ChecklistCategory = "SUPPLIES"
	CategoryMonitoring  ChecklistCategory = "MONITORING"
)

// ChecklistProgress tracks completion status of a checklist
type ChecklistProgress struct {
	TotalItems     int     `json:"total_items"`
	CompletedItems int     `json:"completed_items"`
	Progress       float64 `json:"progress"` // Percentage complete
	RequiredItems  int     `json:"required_items"`
	CompletedRequired int  `json:"completed_required"`
}
