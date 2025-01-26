package pregnancy

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

const (
	defaultGestationDays = models.DefaultGestationDays
	overdueDays          = 340
	highRiskDays         = 370
)

// Calculator handles pregnancy stage and due date calculations
type Calculator struct{}

// NewCalculator creates a new pregnancy calculator
func NewCalculator() *Calculator {
	return &Calculator{}
}

// GetPregnancyStage determines the stage based on days since conception
func (c *Calculator) GetPregnancyStage(daysSinceConception int) models.PregnancyStage {
	switch {
	case daysSinceConception >= highRiskDays:
		return models.PregnancyStageHighRisk
	case daysSinceConception >= overdueDays:
		return models.PregnancyStageOverdue
	case daysSinceConception >= 197:
		return models.PregnancyStageLate
	case daysSinceConception >= 99:
		return models.PregnancyStageMid
	default:
		return models.PregnancyStageEarly
	}
}

// GetStageInfo returns detailed information about the current pregnancy stage
func (c *Calculator) GetStageInfo(conceptionDate time.Time) *models.PregnancyStageInfo {
	daysSinceConception := int(time.Since(conceptionDate).Hours() / 24)
	daysRemaining := defaultGestationDays - daysSinceConception
	daysOverdue := 0
	isOverdue := false

	if daysSinceConception > defaultGestationDays {
		daysOverdue = daysSinceConception - defaultGestationDays
		daysRemaining = 0
		isOverdue = true
	}

	return &models.PregnancyStageInfo{
		Stage:          c.GetPregnancyStage(daysSinceConception),
		DaysSoFar:      daysSinceConception,
		WeeksSoFar:     daysSinceConception / 7,
		DaysRemaining:  daysRemaining,
		WeeksRemaining: daysRemaining / 7,
		Progress:       float64(daysSinceConception) / float64(defaultGestationDays) * 100,
		DaysOverdue:    daysOverdue,
		IsOverdue:      isOverdue,
	}
}

// CalculateProgress calculates the pregnancy progress as a percentage
func (c *Calculator) CalculateProgress(conceptionDate time.Time) float64 {
	daysSinceConception := int(time.Since(conceptionDate).Hours() / 24)
	if daysSinceConception >= defaultGestationDays {
		return 100.0
	}
	return float64(daysSinceConception) / float64(defaultGestationDays) * 100
}

// CalculateDueDateInfo returns comprehensive information about the due date
func (c *Calculator) CalculateDueDateInfo(conceptionDate time.Time) *models.DueDateInfo {
	expectedDueDate := conceptionDate.AddDate(0, 0, defaultGestationDays)
	earliestDueDate := conceptionDate.AddDate(0, 0, defaultGestationDays-14)
	latestDueDate := conceptionDate.AddDate(0, 0, defaultGestationDays+14)
	
	now := time.Now()
	// Round to whole days by truncating to midnight
	expectedDueDateMidnight := expectedDueDate.Truncate(24 * time.Hour)
	nowMidnight := now.Truncate(24 * time.Hour)
	
	daysUntilDue := int(expectedDueDateMidnight.Sub(nowMidnight).Hours() / 24)
	weeksUntilDue := daysUntilDue / 7
	
	isInDueWindow := now.After(earliestDueDate) && now.Before(latestDueDate)
	
	return &models.DueDateInfo{
		ExpectedDueDate: expectedDueDate,
		EarliestDueDate: earliestDueDate,
		LatestDueDate:   latestDueDate,
		DaysUntilDue:    daysUntilDue,
		WeeksUntilDue:   weeksUntilDue,
		IsInDueWindow:   isInDueWindow,
	}
}

// CalculateStage determines the current stage of pregnancy
func (c *Calculator) CalculateStage(pregnancy *models.Pregnancy) models.PregnancyStage {
	if pregnancy.ConceptionDate == nil {
		return models.PregnancyStageEarly
	}

	daysSinceConception := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)
	return c.GetPregnancyStage(daysSinceConception)
}

// CalculateDueDate calculates the expected due date
func (c *Calculator) CalculateDueDate(conceptionDate time.Time) time.Time {
	return conceptionDate.AddDate(0, 0, defaultGestationDays)
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
	dueDate := c.CalculateDueDate(conceptionDate)
	earliestDue := dueDate.AddDate(0, 0, -14)
	latestDue := dueDate.AddDate(0, 0, 14)
	now := time.Now()
	return now.After(earliestDue) && now.Before(latestDue)
}

// CalculateDueWindow returns the earliest and latest due dates
func (c *Calculator) CalculateDueWindow(conceptionDate time.Time) (time.Time, time.Time) {
	dueDate := c.CalculateDueDate(conceptionDate)
	return dueDate.AddDate(0, 0, -14), dueDate.AddDate(0, 0, 14)
}