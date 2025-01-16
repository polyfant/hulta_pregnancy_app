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
		{"Early stage - just started", 1, models.PregnancyStageEarly},
		{"Early stage - middle", 50, models.PregnancyStageEarly},
		{"Early stage - end", 98, models.PregnancyStageEarly},
		{"Mid stage - start", 99, models.PregnancyStageMid},
		{"Mid stage - middle", 150, models.PregnancyStageMid},
		{"Mid stage - end", 196, models.PregnancyStageMid},
		{"Late stage - start", 197, models.PregnancyStageLate},
		{"Late stage - middle", 300, models.PregnancyStageLate},
		{"Late stage - end", 339, models.PregnancyStageLate},
		{"Overdue - start", 340, models.PregnancyStageOverdue},
		{"Overdue - middle", 355, models.PregnancyStageOverdue},
		{"Overdue - end", 369, models.PregnancyStageOverdue},
		{"High risk - start", 370, models.PregnancyStageHighRisk},
		{"High risk - middle", 380, models.PregnancyStageHighRisk},
		{"High risk - extreme", 400, models.PregnancyStageHighRisk},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stage := calc.GetPregnancyStage(tt.daysSinceConception)
			assert.Equal(t, tt.expectedStage, stage, "For %d days since conception", tt.daysSinceConception)
		})
	}
}

func TestGetStageInfo(t *testing.T) {
	calc := NewCalculator()
	now := time.Now()
	
	tests := []struct {
		name           string
		conceptionDate time.Time
		expect         struct {
			stage          models.PregnancyStage
			daysSoFar      int
			weeksSoFar     int
			daysRemaining  int
			weeksRemaining int
			progress       float64
			daysOverdue    int
			isOverdue      bool
		}
	}{
		{
			name: "Early pregnancy",
			conceptionDate: now.AddDate(0, 0, -50),
			expect: struct {
				stage          models.PregnancyStage
				daysSoFar      int
				weeksSoFar     int
				daysRemaining  int
				weeksRemaining int
				progress       float64
				daysOverdue    int
				isOverdue      bool
			}{
				stage:          models.PregnancyStageEarly,
				daysSoFar:      50,
				weeksSoFar:     7,
				daysRemaining:  290,
				weeksRemaining: 41,
				progress:       14.7,
				daysOverdue:    0,
				isOverdue:      false,
			},
		},
		{
			name: "Mid pregnancy",
			conceptionDate: now.AddDate(0, 0, -150),
			expect: struct {
				stage          models.PregnancyStage
				daysSoFar      int
				weeksSoFar     int
				daysRemaining  int
				weeksRemaining int
				progress       float64
				daysOverdue    int
				isOverdue      bool
			}{
				stage:          models.PregnancyStageMid,
				daysSoFar:      150,
				weeksSoFar:     21,
				daysRemaining:  190,
				weeksRemaining: 27,
				progress:       44.1,
				daysOverdue:    0,
				isOverdue:      false,
			},
		},
		{
			name: "Late pregnancy",
			conceptionDate: now.AddDate(0, 0, -300),
			expect: struct {
				stage          models.PregnancyStage
				daysSoFar      int
				weeksSoFar     int
				daysRemaining  int
				weeksRemaining int
				progress       float64
				daysOverdue    int
				isOverdue      bool
			}{
				stage:          models.PregnancyStageLate,
				daysSoFar:      300,
				weeksSoFar:     42,
				daysRemaining:  40,
				weeksRemaining: 5,
				progress:       88.2,
				daysOverdue:    0,
				isOverdue:      false,
			},
		},
		{
			name: "Overdue pregnancy",
			conceptionDate: now.AddDate(0, 0, -350),
			expect: struct {
				stage          models.PregnancyStage
				daysSoFar      int
				weeksSoFar     int
				daysRemaining  int
				weeksRemaining int
				progress       float64
				daysOverdue    int
				isOverdue      bool
			}{
				stage:          models.PregnancyStageOverdue,
				daysSoFar:      350,
				weeksSoFar:     50,
				daysRemaining:  0,
				weeksRemaining: 0,
				progress:       100.0,
				daysOverdue:    10,
				isOverdue:      true,
			},
		},
		{
			name: "High risk pregnancy",
			conceptionDate: now.AddDate(0, 0, -380),
			expect: struct {
				stage          models.PregnancyStage
				daysSoFar      int
				weeksSoFar     int
				daysRemaining  int
				weeksRemaining int
				progress       float64
				daysOverdue    int
				isOverdue      bool
			}{
				stage:          models.PregnancyStageHighRisk,
				daysSoFar:      380,
				weeksSoFar:     54,
				daysRemaining:  0,
				weeksRemaining: 0,
				progress:       100.0,
				daysOverdue:    40,
				isOverdue:      true,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := calc.GetStageInfo(tt.conceptionDate)
			
			assert.Equal(t, tt.expect.stage, info.Stage, "Stage mismatch")
			assert.InDelta(t, tt.expect.daysSoFar, info.DaysSoFar, 1, "Days so far mismatch")
			assert.InDelta(t, tt.expect.weeksSoFar, info.WeeksSoFar, 1, "Weeks so far mismatch")
			assert.InDelta(t, tt.expect.daysRemaining, info.DaysRemaining, 1, "Days remaining mismatch")
			assert.InDelta(t, tt.expect.weeksRemaining, info.WeeksRemaining, 1, "Weeks remaining mismatch")
			assert.InDelta(t, tt.expect.progress, info.Progress, 0.1, "Progress mismatch")
			assert.InDelta(t, tt.expect.daysOverdue, info.DaysOverdue, 1, "Days overdue mismatch")
			assert.Equal(t, tt.expect.isOverdue, info.IsOverdue, "Overdue status mismatch")
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