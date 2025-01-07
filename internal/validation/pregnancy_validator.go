package validation

import (
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type PregnancyValidator struct{}

func NewPregnancyValidator() *PregnancyValidator {
    return &PregnancyValidator{}
}

func (v *PregnancyValidator) ValidatePregnancyStart(pregnancy *models.Pregnancy) error {
	if pregnancy.HorseID == 0 {
		return fmt.Errorf("horse ID is required")
	}

	if pregnancy.ConceptionDate == nil {
		return fmt.Errorf("conception date is required")
	}

	if pregnancy.ConceptionDate.After(time.Now()) {
		return fmt.Errorf("conception date cannot be in the future")
	}

	// Typical horse pregnancy is around 340 days
	if time.Since(*pregnancy.ConceptionDate) > 365*24*time.Hour {
		return fmt.Errorf("invalid conception date: pregnancy duration exceeds one year")
	}

	return nil
}

func (v *PregnancyValidator) ValidatePreFoalingSign(sign *models.PreFoalingSign) error {
	if sign.HorseID == 0 {
		return fmt.Errorf("horse ID is required")
	}

	if sign.Description == "" {
		return fmt.Errorf("description is required")
	}

	if sign.Date.After(time.Now()) {
		return fmt.Errorf("sign date cannot be in the future")
	}

	return nil
} 