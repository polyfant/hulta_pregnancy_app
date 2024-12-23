package pregnancy

import (
	"math"
	"testing"
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/stretchr/testify/assert"
)

// floatEquals checks if two floats are equal within a small epsilon
func floatEquals(a, b float64) bool {
	epsilon := 0.0001
	return math.Abs(a-b) < epsilon
}

func TestPregnancyCalculator(t *testing.T) {
	// Setup a fixed current time for testing
	now := time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		conceptionDate time.Time
		currentTime    time.Time
		expectedTests  struct {
			currentDay    int
			stage        models.PregnancyStage
			daysRemaining int
			progress     float64
			monitoring   MonitoringSchedule
		}
	}{
		{
			name:           "Early Gestation",
			conceptionDate: now.AddDate(0, 0, -50),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    50,
				stage:        models.EarlyGestation,
				daysRemaining: 290,
				progress:     14.71,
				monitoring: MonitoringSchedule{
					CheckFrequency:    24,
					TemperatureCheck:  false,
					BehaviorCheck:     true,
					UdderCheck:        false,
					VulvaCheck:        false,
					Priority:          "Normal",
				},
			},
		},
		{
			name:           "Mid Gestation",
			conceptionDate: now.AddDate(0, 0, -150),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    150,
				stage:        models.MidGestation,
				daysRemaining: 190,
				progress:     44.12,
				monitoring: MonitoringSchedule{
					CheckFrequency:    24,
					TemperatureCheck:  false,
					BehaviorCheck:     true,
					UdderCheck:        false,
					VulvaCheck:        false,
					Priority:          "Normal",
				},
			},
		},
		{
			name:           "Late Gestation",
			conceptionDate: now.AddDate(0, 0, -250),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    250,
				stage:        models.LateGestation,
				daysRemaining: 90,
				progress:     73.53,
				monitoring: MonitoringSchedule{
					CheckFrequency:    24,
					TemperatureCheck:  false,
					BehaviorCheck:     true,
					UdderCheck:        false,
					VulvaCheck:        false,
					Priority:          "Normal",
				},
			},
		},
		{
			name:           "Final Gestation - Temperature Monitoring",
			conceptionDate: now.AddDate(0, 0, -315),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    315,
				stage:        models.FinalGestation,
				daysRemaining: 25,
				progress:     92.65,
				monitoring: MonitoringSchedule{
					CheckFrequency:    24,
					TemperatureCheck:  true,
					BehaviorCheck:     true,
					UdderCheck:        true,
					VulvaCheck:        false,
					Priority:          "Medium",
				},
			},
		},
		{
			name:           "Final Gestation - Intensive Monitoring",
			conceptionDate: now.AddDate(0, 0, -325),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    325,
				stage:        models.FinalGestation,
				daysRemaining: 15,
				progress:     95.59,
				monitoring: MonitoringSchedule{
					CheckFrequency:    6,
					TemperatureCheck:  true,
					BehaviorCheck:     true,
					UdderCheck:        true,
					VulvaCheck:        true,
					Priority:          "High",
				},
			},
		},
		{
			name:           "Final Gestation - Critical Monitoring",
			conceptionDate: now.AddDate(0, 0, -335),
			currentTime:    now,
			expectedTests: struct {
				currentDay    int
				stage        models.PregnancyStage
				daysRemaining int
				progress     float64
				monitoring   MonitoringSchedule
			}{
				currentDay:    335,
				stage:        models.FinalGestation,
				daysRemaining: 5,
				progress:     98.53,
				monitoring: MonitoringSchedule{
					CheckFrequency:    2,
					TemperatureCheck:  true,
					BehaviorCheck:     true,
					UdderCheck:        true,
					VulvaCheck:        true,
					Priority:          "Critical",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calculator := NewPregnancyCalculator(tt.conceptionDate)
			calculator.SetCurrentTime(tt.currentTime)

			currentDay := calculator.GetCurrentDay()
			assert.Equal(t, tt.expectedTests.currentDay, currentDay, "Current day calculation incorrect")

			stage := calculator.GetStage()
			assert.Equal(t, tt.expectedTests.stage, stage, "Stage calculation incorrect")

			daysRemaining := calculator.GetRemainingDays()
			assert.Equal(t, tt.expectedTests.daysRemaining, daysRemaining, "Days remaining calculation incorrect")

			progress := calculator.GetProgressPercentage()
			assert.True(t, floatEquals(tt.expectedTests.progress, progress), 
				"Progress percentage calculation incorrect: expected %v, got %v", 
				tt.expectedTests.progress, progress)

			monitoring := calculator.GetMonitoringSchedule()
			assert.Equal(t, tt.expectedTests.monitoring, monitoring, "Monitoring schedule incorrect")

			// Test IsInCriticalPeriod
			expectedCritical := tt.expectedTests.currentDay >= CriticalMonitoringStart
			assert.Equal(t, expectedCritical, calculator.IsInCriticalPeriod(), 
				"Critical period detection incorrect")

			// Test GetDueDate
			expectedDueDate := tt.conceptionDate.AddDate(0, 0, AveragePregnancyDays)
			assert.Equal(t, expectedDueDate, calculator.GetDueDate(), 
				"Due date calculation incorrect")
		})
	}
}

func TestEdgeCases(t *testing.T) {
	now := time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC)

	t.Run("Zero Days Pregnant", func(t *testing.T) {
		calculator := NewPregnancyCalculator(now)
		calculator.SetCurrentTime(now)
		assert.Equal(t, 0, calculator.GetCurrentDay())
		assert.Equal(t, models.EarlyGestation, calculator.GetStage())
		assert.Equal(t, AveragePregnancyDays, calculator.GetRemainingDays())
		assert.True(t, floatEquals(0.0, calculator.GetProgressPercentage()), 
			"Progress percentage calculation incorrect: expected 0.0, got %v", 
			calculator.GetProgressPercentage())
	})

	t.Run("Over Due Date", func(t *testing.T) {
		calculator := NewPregnancyCalculator(now.AddDate(0, 0, -(AveragePregnancyDays + 5)))
		calculator.SetCurrentTime(now)
		assert.Equal(t, AveragePregnancyDays+5, calculator.GetCurrentDay())
		assert.Equal(t, models.FinalGestation, calculator.GetStage())
		assert.Equal(t, -5, calculator.GetRemainingDays())
		assert.True(t, calculator.GetProgressPercentage() > 100.0)
	})

	t.Run("Far Future Conception", func(t *testing.T) {
		calculator := NewPregnancyCalculator(now.AddDate(0, 0, 30))
		calculator.SetCurrentTime(now)
		assert.Equal(t, -30, calculator.GetCurrentDay())
		assert.Equal(t, models.EarlyGestation, calculator.GetStage())
		assert.Equal(t, AveragePregnancyDays+30, calculator.GetRemainingDays())
		assert.True(t, floatEquals(0.0, calculator.GetProgressPercentage()), 
			"Progress percentage calculation incorrect: expected 0.0, got %v", 
			calculator.GetProgressPercentage())
	})
}
