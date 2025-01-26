package validation

import (
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HorseValidator provides validation for horse-related operations
type HorseValidator struct {
	minAge time.Duration
	maxAge time.Duration
}

// NewHorseValidator creates a new horse validator with default age limits
func NewHorseValidator() *HorseValidator {
	return &HorseValidator{
		minAge: 2 * 365 * 24 * time.Hour,  // 2 years
		maxAge: 25 * 365 * 24 * time.Hour, // 25 years
	}
}

// ValidateHorse performs all validations on a horse
func (v *HorseValidator) ValidateHorse(horse *models.Horse) error {
	if err := v.validateBasicInfo(horse); err != nil {
		return err
	}
	return v.validateBreedingEligibility(horse)
}

func (v *HorseValidator) validateBasicInfo(horse *models.Horse) error {
	if horse.Name == "" {
		return fmt.Errorf("horse name is required")
	}

	if horse.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	if err := v.validateAge(horse); err != nil {
		return err
	}

	return v.validatePhysicalAttributes(horse)
}

func (v *HorseValidator) validateAge(horse *models.Horse) error {
	if horse.BirthDate.IsZero() {
		return fmt.Errorf("birth date is required")
	}

	age := time.Since(horse.BirthDate)
	if age < v.minAge {
		return fmt.Errorf("horse is too young (minimum age: 2 years)")
	}
	if age > v.maxAge {
		return fmt.Errorf("horse age exceeds maximum (maximum age: 25 years)")
	}

	return nil
}

func (v *HorseValidator) validatePhysicalAttributes(horse *models.Horse) error {
	if horse.Weight < 0 || horse.Weight > 1200 {
		return fmt.Errorf("invalid weight: must be between 0 and 1200 kg")
	}

	if horse.Height < 0 || horse.Height > 200 {
		return fmt.Errorf("invalid height: must be between 0 and 200 cm")
	}

	return nil
}

func (v *HorseValidator) validateBreedingEligibility(horse *models.Horse) error {
	if !horse.IsPregnant {
		return nil
	}

	if horse.Gender != models.GenderMare {
		return fmt.Errorf("only mares can be pregnant")
	}

	if horse.ConceptionDate == nil {
		return fmt.Errorf("conception date is required for pregnant horse")
	}

	return nil
}
