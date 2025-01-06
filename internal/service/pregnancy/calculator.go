package pregnancy

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type Calculator struct {
    defaultGestationDays int
}

func NewCalculator() *Calculator {
    return &Calculator{
        defaultGestationDays: models.DefaultGestationDays,
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