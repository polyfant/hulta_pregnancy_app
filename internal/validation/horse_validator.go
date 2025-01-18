package validation

import (
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HorseValidator provides validation for Horse entities
type HorseValidator struct{}

// NewHorseValidator creates a new instance of HorseValidator
func NewHorseValidator() *HorseValidator {
	return &HorseValidator{}
}

// ValidateHorse performs comprehensive validation on a Horse
func (v *HorseValidator) ValidateHorse(horse *models.Horse) error {
	// Use the generic struct validation first
	if err := ValidateStruct(horse); err != nil {
		return err
	}

	// Additional custom validations
	if err := v.validateBasicInfo(horse); err != nil {
		return err
	}

	if err := v.validateAge(horse); err != nil {
		return err
	}

	if err := v.validatePhysicalAttributes(horse); err != nil {
		return err
	}

	if err := v.validateBreedingEligibility(horse); err != nil {
		return err
	}

	return nil
}

func (v *HorseValidator) validateBasicInfo(horse *models.Horse) error {
	if horse.Name == "" {
		return fmt.Errorf("horse name cannot be empty")
	}
	return nil
}

func (v *HorseValidator) validateAge(horse *models.Horse) error {
	if horse.BirthDate == nil {
		return fmt.Errorf("birth date is required")
	}

	age := time.Since(*horse.BirthDate)
	if age.Hours() < 0 {
		return fmt.Errorf("birth date cannot be in the future")
	}

	// Assuming horses can live up to 40 years
	if age.Hours() > 40*365*24 {
		return fmt.Errorf("horse age seems unrealistic")
	}

	return nil
}

func (v *HorseValidator) validatePhysicalAttributes(horse *models.Horse) error {
	if horse.Weight < 0 {
		return fmt.Errorf("horse weight cannot be negative")
	}

	if horse.Height < 0 {
		return fmt.Errorf("horse height cannot be negative")
	}

	return nil
}

func (v *HorseValidator) validateBreedingEligibility(horse *models.Horse) error {
	// Typical breeding age for horses is between 2-20 years
	if horse.BirthDate != nil {
		age := time.Since(*horse.BirthDate).Hours() / (365.25 * 24)
		if age < 2 {
			return fmt.Errorf("horse is too young for breeding")
		}
		if age > 20 {
			return fmt.Errorf("horse is too old for breeding")
		}
	}

	return nil
}
