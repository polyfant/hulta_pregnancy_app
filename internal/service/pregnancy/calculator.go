package pregnancy

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Calculator handles pregnancy stage and due date calculations
type Calculator struct{}

// NewCalculator creates a new pregnancy calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// CalculateStage determines the current stage of pregnancy
func (c *Calculator) CalculateStage(pregnancy *models.Pregnancy) models.PregnancyStage {
	if pregnancy.ConceptionDate == nil {
		return models.PregnancyStageEarly
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)
	progressPercentage := float64(daysPregnant) / float64(models.DefaultGestationDays)

	switch {
	case progressPercentage <= 0.33:
		return models.PregnancyStageEarly
	case progressPercentage <= 0.66:
		return models.PregnancyStageMid
	default:
		return models.PregnancyStageLate
	}
}

// CalculateDueDate calculates the expected due date
func (c *Calculator) CalculateDueDate(conceptionDate time.Time) time.Time {
	return conceptionDate.AddDate(0, 0, models.DefaultGestationDays)
}

// CalculateDaysPregnant calculates the number of days pregnant
func (c *Calculator) CalculateDaysPregnant(conceptionDate time.Time) int {
	return int(time.Since(conceptionDate).Hours() / 24)
}

// CalculateWeeksPregnant calculates the number of weeks pregnant
func (c *Calculator) CalculateWeeksPregnant(conceptionDate time.Time) int {
	return c.CalculateDaysPregnant(conceptionDate) / 7
}

// CalculateIsInDueWindow checks if the pregnancy is in the due window
func (c *Calculator) CalculateIsInDueWindow(conceptionDate time.Time) bool {
	daysPregnant := c.CalculateDaysPregnant(conceptionDate)
	return daysPregnant >= models.DefaultGestationDays-14 && daysPregnant <= models.DefaultGestationDays+14
}

// CalculateDueWindow returns the earliest and latest due dates
func (c *Calculator) CalculateDueWindow(conceptionDate time.Time) (time.Time, time.Time) {
	earliest := conceptionDate.AddDate(0, 0, models.DefaultGestationDays-14)
	latest := conceptionDate.AddDate(0, 0, models.DefaultGestationDays+14)
	return earliest, latest
} 