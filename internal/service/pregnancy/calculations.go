package pregnancy

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

const (
	// Average pregnancy duration for horses
	AveragePregnancyDays = 340

	// DefaultGestationDays is the average gestation period for horses (340 days)
	DefaultGestationDays = 340

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
		return models.Foaling
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

// CalculateGuidelines returns pregnancy guidelines based on the current stage
func (pc *PregnancyCalculator) CalculateGuidelines() map[string]string {
	stage := pc.GetStage()
	guidelines := make(map[string]string)

	switch stage {
	case models.EarlyGestation:
		guidelines["Feeding"] = "Maintain normal feeding routine"
		guidelines["Exercise"] = "Continue regular exercise"
		guidelines["Monitoring"] = "Watch for signs of morning sickness"
	case models.MidGestation:
		guidelines["Feeding"] = "Increase feed quality, not quantity"
		guidelines["Exercise"] = "Light to moderate exercise"
		guidelines["Monitoring"] = "Regular vet checks recommended"
	case models.LateGestation:
		guidelines["Feeding"] = "Gradual increase in feed quality and quantity"
		guidelines["Exercise"] = "Light exercise only"
		guidelines["Monitoring"] = "Weekly vet checks, prepare for foaling"
	}

	return guidelines
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

// CalculateDueDate calculates the expected foaling date based on conception date
// gestationDays parameter allows for breed-specific or mare-specific adjustments
func CalculateDueDate(conceptionDate time.Time, gestationDays int) time.Time {
	if gestationDays <= 0 {
		gestationDays = DefaultGestationDays
	}
	return conceptionDate.AddDate(0, 0, gestationDays)
}

// CalculateGestationProgress returns the current progress as a percentage and days remaining
func CalculateGestationProgress(conceptionDate time.Time, gestationDays int) (float64, int) {
	if gestationDays <= 0 {
		gestationDays = DefaultGestationDays
	}
	
	daysPregnant := int(time.Since(conceptionDate).Hours() / 24)
	progress := float64(daysPregnant) / float64(gestationDays) * 100
	daysRemaining := gestationDays - daysPregnant

	return progress, daysRemaining
}
