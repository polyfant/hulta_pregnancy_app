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

// GetPregnancyStage determines the stage based on days since conception with advanced logic
func (c *Calculator) GetPregnancyStage(daysSinceConception int) models.PregnancyStage {
	const (
		earlyStageEnd     = 120 // First 4 months
		midStageEnd       = 270 // 9 months
		lateStageEnd      = 340 // Standard gestation
		overdueThreshold  = 370 // Beyond standard gestation
	)

	switch {
	case daysSinceConception <= earlyStageEnd:
		return models.EarlyStage
	case daysSinceConception <= midStageEnd:
		return models.MidStage
	case daysSinceConception <= lateStageEnd:
		return models.LateStage
	case daysSinceConception > overdueThreshold:
		return models.OverdueStage
	default:
		return models.LateStage
	}
}

// CalculateProgress calculates the pregnancy progress as a percentage with advanced heuristics
func (c *Calculator) CalculateProgress(conceptionDate time.Time) float64 {
	daysSinceConception := time.Since(conceptionDate).Hours() / 24
	progress := (daysSinceConception / float64(defaultGestationDays)) * 100
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}
	return progress
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

	progress := c.CalculateProgress(conceptionDate)
	stage := c.GetPregnancyStage(daysSinceConception)

	var stageDescription string
	var riskLevel models.RiskLevel = models.LowRisk

	switch stage {
	case models.EarlyStage:
		stageDescription = "Early pregnancy: Embryo development and initial health monitoring"
		riskLevel = models.LowRisk
	case models.MidStage:
		stageDescription = "Mid pregnancy: Fetal growth and critical nutritional period"
		riskLevel = models.MediumRisk
	case models.LateStage:
		stageDescription = "Late pregnancy: Preparing for foaling, increased monitoring required"
		riskLevel = models.HighRisk
	case models.OverdueStage:
		stageDescription = "Overdue: Immediate veterinary consultation recommended"
		riskLevel = models.CriticalRisk
	}

	return &models.PregnancyStageInfo{
		Stage:             stage,
		Progress:          progress,
		Description:       stageDescription,
		RiskLevel:         riskLevel,
		DaysSoFar:         daysSinceConception,
		WeeksSoFar:        daysSinceConception / 7,
		DaysRemaining:     daysRemaining,
		WeeksRemaining:    daysRemaining / 7,
		DaysOverdue:       daysOverdue,
		IsOverdue:         isOverdue,
		NutritionAdvice:   c.getNutritionAdvice(stage),
		MonitoringAdvice:  c.getMonitoringAdvice(stage),
	}
}

// calculateDaysRemaining estimates remaining days in pregnancy
func (c *Calculator) calculateDaysRemaining(conceptionDate time.Time) int {
	const averageGestationDays = 340
	daysSinceConception := int(time.Since(conceptionDate).Hours() / 24)
	remaining := averageGestationDays - daysSinceConception

	// Ensure non-negative
	if remaining < 0 {
		return 0
	}
	return remaining
}

// getNutritionAdvice provides stage-specific nutrition recommendations
func (c *Calculator) getNutritionAdvice(stage models.PregnancyStage) string {
	switch stage {
	case models.EarlyStage:
		return "Focus on balanced diet, maintain body condition. Supplement with folic acid and minerals."
	case models.MidStage:
		return "Increase protein and energy intake. Monitor weight gain carefully."
	case models.LateStage:
		return "High-quality protein, increased calories. Prepare for increased nutritional demands."
	case models.OverdueStage:
		return "Consult veterinarian. Specialized nutrition may be required."
	default:
		return "Maintain standard pregnancy nutrition protocol."
	}
}

// getMonitoringAdvice provides stage-specific health monitoring recommendations
func (c *Calculator) getMonitoringAdvice(stage models.PregnancyStage) string {
	switch stage {
	case models.EarlyStage:
		return "Weekly health checks. Monitor for early signs of complications."
	case models.MidStage:
		return "Bi-weekly detailed health assessments. Ultrasound recommended."
	case models.LateStage:
		return "Weekly veterinary check-ups. Monitor for pre-foaling signs."
	case models.OverdueStage:
		return "Immediate and frequent veterinary monitoring. Prepare for potential intervention."
	default:
		return "Standard pregnancy monitoring protocol."
	}
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
		return models.EarlyStage
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