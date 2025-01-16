package pregnancy

import (
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateWeatherImpact(t *testing.T) {
	calc := NewWeatherCalculator()

	tests := []struct {
		name           string
		temp          float64
		humidity      float64
		stage         models.PregnancyStage
		wantStress    string
		wantAlertFlag bool
	}{
		{
			name:      "Normal conditions early stage",
			temp:      20.0,
			humidity:  50.0,
			stage:     models.PregnancyStageEarly,
			wantStress: "LOW",
			wantAlertFlag: false,
		},
		{
			name:      "High heat late stage",
			temp:      35.0,
			humidity:  70.0,
			stage:     models.PregnancyStageLate,
			wantStress: "HIGH",
			wantAlertFlag: true,
		},
		{
			name:      "Moderate conditions mid stage",
			temp:      28.0,
			humidity:  60.0,
			stage:     models.PregnancyStageMid,
			wantStress: "MODERATE",
			wantAlertFlag: false,
		},
		{
			name:      "Borderline conditions overdue",
			temp:      30.0,
			humidity:  65.0,
			stage:     models.PregnancyStageOverdue,
			wantStress: "MODERATE",
			wantAlertFlag: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impact := calc.CalculateWeatherImpact(tt.temp, tt.humidity, tt.stage)
			assert.Equal(t, tt.wantStress, impact.StressLevel)
			
			alert := calc.IsWeatherAlert(impact, tt.stage)
			assert.Equal(t, tt.wantAlertFlag, alert)

			// Verify recommendations are not empty
			assert.NotEmpty(t, impact.Recommendations)
		})
	}
}

func TestHeatIndexCalculation(t *testing.T) {
	calc := NewWeatherCalculator()

	tests := []struct {
		name     string
		temp     float64
		humidity float64
		wantHI   float64
		delta    float64 // Acceptable difference for floating point comparison
	}{
		{
			name:     "Normal conditions",
			temp:     25.0,
			humidity: 50.0,
			wantHI:   25.9,
			delta:    1.0,
		},
		{
			name:     "High heat conditions",
			temp:     35.0,
			humidity: 70.0,
			wantHI:   45.5,
			delta:    1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impact := calc.CalculateWeatherImpact(tt.temp, tt.humidity, models.PregnancyStageEarly)
			assert.InDelta(t, tt.wantHI, impact.HeatIndex, tt.delta)
		})
	}
}
