package pregnancy

import (
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPregnancyStage(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name              string
		daysSinceConception int
		expectedStage     models.PregnancyStage
	}{
		{"Early stage", 50, models.PregnancyStageEarly},
		{"Mid stage", 150, models.PregnancyStageMid},
		{"Late stage", 300, models.PregnancyStageLate},
		{"Overdue", 345, models.PregnancyStageOverdue},
		{"High risk", 380, models.PregnancyStageHighRisk},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stage := calc.GetPregnancyStage(tt.daysSinceConception)
			assert.Equal(t, tt.expectedStage, stage)
		})
	}
}

func TestGetStageInfo(t *testing.T) {
	calc := NewCalculator()
	now := time.Now()
	
	tests := []struct {
		name           string
		conceptionDate time.Time
		expectOverdue  bool
		expectDaysOverdue int
	}{
		{
			name: "Normal pregnancy",
			conceptionDate: now.AddDate(0, 0, -300),
			expectOverdue: false,
			expectDaysOverdue: 0,
		},
		{
			name: "Overdue pregnancy",
			conceptionDate: now.AddDate(0, 0, -350),
			expectOverdue: true,
			expectDaysOverdue: 10,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := calc.GetStageInfo(tt.conceptionDate)
			assert.Equal(t, tt.expectOverdue, info.IsOverdue)
			assert.Equal(t, tt.expectDaysOverdue, info.DaysOverdue)
		})
	}
}

func TestCalculator(t *testing.T) {
	calc := NewCalculator()
	now := time.Now()

	// Test progress calculation
	t.Run("CalculateProgress", func(t *testing.T) {
		tests := []struct {
			name              string
			conceptionDate    time.Time
			expectProgress    float64
			expectDaysRemaining int
		}{
			{
				name: "Just started",
				conceptionDate: now,
				expectProgress: 0,
				expectDaysRemaining: defaultGestationDays,
			},
			{
				name: "Half way",
				conceptionDate: now.AddDate(0, 0, -170),
				expectProgress: 50,
				expectDaysRemaining: 170,
			},
			{
				name: "Completed",
				conceptionDate: now.AddDate(0, 0, -340),
				expectProgress: 100,
				expectDaysRemaining: 0,
			},
			{
				name: "Over completed",
				conceptionDate: now.AddDate(0, 0, -400),
				expectProgress: 100,
				expectDaysRemaining: 0,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				progress := calc.CalculateProgress(tt.conceptionDate)
				assert.InDelta(t, tt.expectProgress, progress, 0.1)
			})
		}
	})
}

func TestCalculateDueDateInfo(t *testing.T) {
	calc := NewCalculator()
	now := time.Now()

	tests := []struct {
		name           string
		conceptionDate time.Time
		expect        struct {
			daysUntilDue    int
			weeksUntilDue   int
			isInDueWindow   bool
			daysSpread      int  // Days between earliest and latest
		}
	}{
		{
			name: "Normal pregnancy - mid term",
			conceptionDate: now.AddDate(0, 0, -170), // Half way through
			expect: struct {
				daysUntilDue    int
				weeksUntilDue   int
				isInDueWindow   bool
				daysSpread      int
			}{
				daysUntilDue:  170,
				weeksUntilDue: 24,
				isInDueWindow: false,
				daysSpread:    28, // 14 days before and after
			},
		},
		{
			name: "Almost due",
			conceptionDate: now.AddDate(0, 0, -319), // One day before earliest
			expect: struct {
				daysUntilDue    int
				weeksUntilDue   int
				isInDueWindow   bool
				daysSpread      int
			}{
				daysUntilDue:  21,
				weeksUntilDue: 3,
				isInDueWindow: false,
				daysSpread:    28,
			},
		},
		{
			name: "In due window",
			conceptionDate: now.AddDate(0, 0, -330), // Between min and max
			expect: struct {
				daysUntilDue    int
				weeksUntilDue   int
				isInDueWindow   bool
				daysSpread      int
			}{
				daysUntilDue:  10,
				weeksUntilDue: 1,
				isInDueWindow: true,
				daysSpread:    28,
			},
		},
		{
			name: "Overdue",
			conceptionDate: now.AddDate(0, 0, -380), // Past latest
			expect: struct {
				daysUntilDue    int
				weeksUntilDue   int
				isInDueWindow   bool
				daysSpread      int
			}{
				daysUntilDue:  -40,
				weeksUntilDue: -5,
				isInDueWindow: false,
				daysSpread:    28,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := calc.CalculateDueDateInfo(tt.conceptionDate)

			// Test days and weeks until due
			assert.Equal(t, tt.expect.daysUntilDue, info.DaysUntilDue)
			assert.Equal(t, tt.expect.weeksUntilDue, info.WeeksUntilDue)
			
			// Test due window status
			assert.Equal(t, tt.expect.isInDueWindow, info.IsInDueWindow)
			
			// Test due date ranges
			daysSpread := int(info.LatestDueDate.Sub(info.EarliestDueDate).Hours() / 24)
			assert.Equal(t, tt.expect.daysSpread, daysSpread)
			
			// Test date ordering
			assert.True(t, info.EarliestDueDate.Before(info.ExpectedDueDate))
			assert.True(t, info.ExpectedDueDate.Before(info.LatestDueDate))
			
			// Validate due dates are calculated correctly
			expectedDueDate := tt.conceptionDate.AddDate(0, 0, defaultGestationDays)
			earliestDueDate := tt.conceptionDate.AddDate(0, 0, defaultGestationDays-14)
			latestDueDate := tt.conceptionDate.AddDate(0, 0, defaultGestationDays+14)
			
			assert.Equal(t, expectedDueDate, info.ExpectedDueDate)
			assert.Equal(t, earliestDueDate, info.EarliestDueDate)
			assert.Equal(t, latestDueDate, info.LatestDueDate)
		})
	}
}