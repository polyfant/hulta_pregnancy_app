package models

import (
	"fmt"
)

// FeedRequirements represents the daily feed requirements for a horse
type FeedRequirements struct {
	HayRequirementKg       float64 `json:"hay_requirement_kg" validate:"required,min=0,max=50"`
	ConcentrateRequirementKg float64 `json:"concentrate_requirement_kg" validate:"required,min=0,max=20"`
	HorseID                int     `json:"horse_id" validate:"required"`
}

// Validate checks if the feed requirements are within reasonable bounds
func (f *FeedRequirements) Validate() error {
	// Validate total daily intake
	totalIntake := f.HayRequirementKg + f.ConcentrateRequirementKg
	if totalIntake <= 0 {
		return fmt.Errorf("total feed intake must be positive")
	}

	// Optional: Add a reasonable maximum intake limit
	const maxDailyIntake = 50.0 // kg
	if totalIntake > maxDailyIntake {
		return fmt.Errorf("total daily feed intake exceeds recommended maximum of %.1f kg", maxDailyIntake)
	}

	return nil
}