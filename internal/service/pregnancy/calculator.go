package pregnancy

import (
	"math"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type Calculator struct {
    defaultGestationDays int
    stageThresholds     map[models.PregnancyStage]int
    maxOverdueDays      int  // Maximum number of days we consider normal for overdue
    minGestationDays    int
    maxGestationDays    int
}

func NewCalculator() *Calculator {
    return &Calculator{
        defaultGestationDays: models.DefaultGestationDays,
        stageThresholds: map[models.PregnancyStage]int{
            models.PregnancyStageEarlyGestation: 98,  // 14 weeks
            models.PregnancyStageMidGestation:   196, // 28 weeks
            models.PregnancyStageLateGestation:  340, // 48 weeks
        },
        maxOverdueDays: 30,  // Up to 30 days overdue is considered normal
        minGestationDays: 320, // Minimum viable gestation
        maxGestationDays: 370, // Maximum expected gestation
    }
}

func (c *Calculator) CalculateDueDate(conceptionDate time.Time) time.Time {
    return conceptionDate.AddDate(0, 0, c.defaultGestationDays)
}

func (c *Calculator) CalculateProgress(conceptionDate time.Time) (float64, int) {
    now := time.Now()
    daysSinceConception := int(now.Sub(conceptionDate).Hours() / 24)
    
    progress := float64(daysSinceConception) / float64(c.defaultGestationDays) * 100
    if progress > 100 {
        progress = 100
    }
    
    daysRemaining := c.defaultGestationDays - daysSinceConception
    if daysRemaining < 0 {
        daysRemaining = 0
    }
    
    return progress, daysRemaining
}

func (c *Calculator) GetPregnancyStage(daysSinceConception int) models.PregnancyStage {
    switch {
    case daysSinceConception <= c.stageThresholds[models.PregnancyStageEarlyGestation]:
        return models.PregnancyStageEarlyGestation
    case daysSinceConception <= c.stageThresholds[models.PregnancyStageMidGestation]:
        return models.PregnancyStageMidGestation
    case daysSinceConception <= c.stageThresholds[models.PregnancyStageLateGestation]:
        return models.PregnancyStageLateGestation
    case daysSinceConception <= c.stageThresholds[models.PregnancyStageLateGestation] + c.maxOverdueDays:
        return models.PregnancyStageOverdue
    default:
        return models.PregnancyStageHighRisk
    }
}

// GetStageInfo returns detailed information about the current pregnancy stage
func (c *Calculator) GetStageInfo(conceptionDate time.Time) models.PregnancyStageInfo {
    daysSinceConception := int(time.Since(conceptionDate).Hours() / 24)
    currentStage := c.GetPregnancyStage(daysSinceConception)
    
    // Calculate days overdue if applicable
    daysOverdue := 0
    if daysSinceConception > c.defaultGestationDays {
        daysOverdue = daysSinceConception - c.defaultGestationDays
    }
    
    // Calculate progress, capping at 100%
    progress := float64(daysSinceConception) / float64(c.defaultGestationDays) * 100
    if progress > 100 {
        progress = 100
    }
    
    return models.PregnancyStageInfo{
        Stage:           currentStage,
        DaysSoFar:      daysSinceConception,
        WeeksSoFar:     daysSinceConception / 7,
        DaysRemaining:  c.defaultGestationDays - daysSinceConception,
        WeeksRemaining: (c.defaultGestationDays - daysSinceConception) / 7,
        Progress:       progress,
        DaysOverdue:    daysOverdue,
        IsOverdue:      daysOverdue > 0,
    }
}

// DueDateInfo provides comprehensive due date information
func (c *Calculator) CalculateDueDateInfo(conceptionDate time.Time) models.DueDateInfo {
    expectedDueDate := conceptionDate.AddDate(0, 0, c.defaultGestationDays)
    earliestDueDate := conceptionDate.AddDate(0, 0, c.minGestationDays)
    latestDueDate := conceptionDate.AddDate(0, 0, c.maxGestationDays)
    
    now := time.Now()
    // Round to nearest day to avoid time-of-day differences
    daysUntilDue := int(math.Round(expectedDueDate.Sub(now).Hours() / 24))
    
    return models.DueDateInfo{
        ExpectedDueDate: expectedDueDate,
        EarliestDueDate: earliestDueDate,
        LatestDueDate:   latestDueDate,
        DaysUntilDue:    daysUntilDue,
        // For weeks, we want to floor negative numbers and ceil positive numbers
        WeeksUntilDue:   daysUntilDue / 7,
        IsInDueWindow:   now.After(earliestDueDate) && now.Before(latestDueDate),
    }
} 