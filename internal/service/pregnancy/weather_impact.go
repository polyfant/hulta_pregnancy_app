package pregnancy

import (
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// WeatherImpact represents the impact of weather on pregnancy
type WeatherImpact struct {
	Temperature     float64 `json:"temperature"`      // in Celsius
	Humidity       float64 `json:"humidity"`         // percentage
	HeatIndex      float64 `json:"heat_index"`       // calculated heat index
	StressLevel    string  `json:"stress_level"`     // LOW, MODERATE, HIGH
	Recommendations []string `json:"recommendations"` // list of recommendations based on conditions
}

// WeatherCalculator handles weather-related calculations and recommendations
type WeatherCalculator struct{}

// NewWeatherCalculator creates a new weather impact calculator
func NewWeatherCalculator() *WeatherCalculator {
	return &WeatherCalculator{}
}

// CalculateWeatherImpact determines the impact of weather conditions on pregnancy
func (w *WeatherCalculator) CalculateWeatherImpact(temp, humidity float64, stage models.PregnancyStage) *WeatherImpact {
	impact := &WeatherImpact{
		Temperature: temp,
		Humidity:    humidity,
		HeatIndex:   w.calculateHeatIndex(temp, humidity),
	}

	impact.StressLevel = w.determineStressLevel(impact.HeatIndex, stage)
	impact.Recommendations = w.getRecommendations(impact.StressLevel, stage)

	return impact
}

// calculateHeatIndex calculates the heat index using a simplified formula
func (w *WeatherCalculator) calculateHeatIndex(temp, humidity float64) float64 {
	// For temperatures below 20°C, heat index equals temperature
	if temp < 20 {
		return temp
	}

	// Convert humidity to decimal
	rh := humidity / 100.0

	// Base heat index starts with temperature
	hi := temp

	// Add humidity effect that increases exponentially with temperature
	if temp >= 25 {
		// Temperature effect increases more rapidly at higher temps
		tempEffect := 0.7 * (temp - 25)
		if temp >= 30 {
			tempEffect *= 2.3
		}
		
		// Combine with humidity (stronger effect at higher humidity)
		humidityEffect := tempEffect * rh * (0.8 + 0.2*rh)
		hi += humidityEffect
	}

	return hi
}

// determineStressLevel determines stress level based on temperature and stage
func (w *WeatherCalculator) determineStressLevel(heatIndex float64, stage models.PregnancyStage) string {
	// Late stage pregnancies are more sensitive to heat
	isLateStage := stage == models.PregnancyStageLate || stage == models.PregnancyStageOverdue

	if isLateStage {
		if heatIndex >= 35 {
			return "HIGH"
		}
		if heatIndex >= 30 {
			return "MODERATE"
		}
	} else {
		if heatIndex >= 40 {
			return "HIGH"
		}
		if heatIndex >= 28 {
			return "MODERATE"
		}
	}
	return "LOW"
}

// getRecommendations provides stage-specific recommendations based on stress level
func (w *WeatherCalculator) getRecommendations(stressLevel string, stage models.PregnancyStage) []string {
	baseRecs := []string{
		"Ensure constant access to fresh, clean water",
		"Provide adequate shade in paddocks",
	}

	switch stressLevel {
	case "HIGH":
		return append(baseRecs,
			"Consider moving horse to climate-controlled stable",
			"Monitor water intake closely",
			"Avoid any exercise during peak heat hours",
			"Consider cold hosing if temperature exceeds 30°C",
			"Schedule vet check if heat stress signs appear")
	case "MODERATE":
		return append(baseRecs,
			"Limit outdoor exposure during peak heat hours",
			"Increase monitoring frequency",
			"Ensure good ventilation in stables",
			"Consider adding electrolytes to water")
	default:
		return append(baseRecs,
			"Maintain regular monitoring schedule",
			"Ensure normal exercise routine if comfortable")
	}
}

// IsWeatherAlert determines if current conditions warrant an alert
func (w *WeatherCalculator) IsWeatherAlert(impact *WeatherImpact, stage models.PregnancyStage) bool {
	isLateStage := stage == models.PregnancyStageLate || stage == models.PregnancyStageOverdue
	return impact.StressLevel == "HIGH" || (impact.StressLevel == "MODERATE" && isLateStage)
}
