package pregnancy

import (
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
)

const (
	// Average pregnancy duration for horses
	AveragePregnancyDays = 340

	// Stage durations in days
	EarlyGestationStart  = 0
	EarlyGestationEnd    = 113
	MidGestationStart    = 114
	MidGestationEnd      = 225
	LateGestationStart   = 226
	LateGestationEnd     = 310
	FinalGestationStart  = 311
	FinalGestationEnd    = 340

	// Important monitoring periods
	TemperatureMonitoringStart = 310  // Start daily temperature monitoring
	IntensiveMonitoringStart   = 320  // Start checking every 6 hours
	CriticalMonitoringStart    = 330  // Start checking every 2-3 hours
)

// PregnancyCalculator provides detailed pregnancy calculations
type PregnancyCalculator struct {
	ConceptionDate time.Time
	currentTime    time.Time
}

// NewPregnancyCalculator creates a new calculator instance
func NewPregnancyCalculator(conceptionDate time.Time) *PregnancyCalculator {
	return &PregnancyCalculator{
		ConceptionDate: conceptionDate,
		currentTime:    time.Now(),
	}
}

// GetDueDate calculates the expected due date
func (pc *PregnancyCalculator) GetDueDate() time.Time {
	return pc.ConceptionDate.AddDate(0, 0, AveragePregnancyDays)
}

// GetCurrentDay returns the current day of pregnancy (0-based)
func (pc *PregnancyCalculator) GetCurrentDay() int {
	duration := pc.currentTime.Sub(pc.ConceptionDate)
	days := int(duration.Hours() / 24)
	return days
}

// GetStage determines the current pregnancy stage
func (pc *PregnancyCalculator) GetStage() models.PregnancyStage {
	currentDay := pc.GetCurrentDay()

	switch {
	case currentDay <= EarlyGestationEnd:
		return models.EarlyGestation
	case currentDay <= MidGestationEnd:
		return models.MidGestation
	case currentDay <= LateGestationEnd:
		return models.LateGestation
	default:
		return models.FinalGestation
	}
}

// GetMonitoringSchedule returns the current monitoring requirements
func (pc *PregnancyCalculator) GetMonitoringSchedule() MonitoringSchedule {
	currentDay := pc.GetCurrentDay()

	switch {
	case currentDay >= CriticalMonitoringStart:
		return MonitoringSchedule{
			CheckFrequency:    2,  // hours
			TemperatureCheck:  true,
			BehaviorCheck:     true,
			UdderCheck:        true,
			VulvaCheck:        true,
			Priority:          "Critical",
		}
	case currentDay >= IntensiveMonitoringStart:
		return MonitoringSchedule{
			CheckFrequency:    6,  // hours
			TemperatureCheck:  true,
			BehaviorCheck:     true,
			UdderCheck:        true,
			VulvaCheck:        true,
			Priority:          "High",
		}
	case currentDay >= TemperatureMonitoringStart:
		return MonitoringSchedule{
			CheckFrequency:    24, // hours
			TemperatureCheck:  true,
			BehaviorCheck:     true,
			UdderCheck:        true,
			VulvaCheck:        false,
			Priority:          "Medium",
		}
	default:
		return MonitoringSchedule{
			CheckFrequency:    24, // hours
			TemperatureCheck:  false,
			BehaviorCheck:     true,
			UdderCheck:        false,
			VulvaCheck:        false,
			Priority:          "Normal",
		}
	}
}

// MonitoringSchedule defines what needs to be monitored and how frequently
type MonitoringSchedule struct {
	CheckFrequency    int    // How often to check in hours
	TemperatureCheck  bool   // Whether to check temperature
	BehaviorCheck     bool   // Whether to monitor behavior
	UdderCheck        bool   // Whether to check udder development
	VulvaCheck        bool   // Whether to check vulva
	Priority          string // Monitoring priority level
}

// GetRemainingDays calculates days remaining until due date
func (pc *PregnancyCalculator) GetRemainingDays() int {
	dueDate := pc.GetDueDate()
	duration := dueDate.Sub(pc.currentTime)
	days := int(duration.Hours() / 24)
	return days
}

// IsInCriticalPeriod determines if the mare is in the critical pre-foaling period
func (pc *PregnancyCalculator) IsInCriticalPeriod() bool {
	return pc.GetCurrentDay() >= CriticalMonitoringStart
}

// GetProgressPercentage calculates the percentage of pregnancy completed
func (pc *PregnancyCalculator) GetProgressPercentage() float64 {
	currentDay := pc.GetCurrentDay()
	if currentDay < 0 {
		return 0.0
	}
	percentage := float64(currentDay) / float64(AveragePregnancyDays) * 100
	// Round to 2 decimal places
	return float64(int(percentage*100)) / 100
}

// SetCurrentTime sets the current time for testing purposes
func (pc *PregnancyCalculator) SetCurrentTime(t time.Time) {
	pc.currentTime = t
}
