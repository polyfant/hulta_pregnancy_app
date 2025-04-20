package ml

import (
	"os"
	"testing"
)

func TestPretrainedModels(t *testing.T) {
	if err := GeneratePretrainedModels(); err != nil {
		t.Fatalf("Failed to generate models: %v", err)
	}

	if _, err := os.Stat("models/pregnancy_v1.model"); os.IsNotExist(err) {
		t.Fatalf("Model file models/pregnancy_v1.model was not created")
	}
	if _, err := os.Stat("models/growth_v1.model"); os.IsNotExist(err) {
		t.Fatalf("Model file models/growth_v1.model was not created")
	}
	if _, err := os.Stat("models/health_v1.model"); os.IsNotExist(err) {
		t.Fatalf("Model file models/health_v1.model was not created")
	}

	if err := InitializeModels(); err != nil {
		t.Fatalf("Failed to initialize models: %v", err)
	}

	tests := []struct {
		name      string
		modelType string
		features  map[string]float64
		wantErr   bool
	}{
		{
			name:      "Test Pregnancy Prediction",
			modelType: "PREGNANCY",
			features: map[string]float64{
				"age":            8,
				"weight":         550,
				"temperature":    37.5,
				"hormone_level":  2.5,
				"previous_births": 2,
			},
			wantErr: false,
		},
		{
			name:      "Test Growth Prediction",
			modelType: "GROWTH",
			features: map[string]float64{
				"age":           3,
				"current_weight": 450,
				"height":        1.5,
				"daily_feed":    10,
				"exercise":      2,
			},
			wantErr: false,
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := LoadPretrainedModel(tt.modelType)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("LoadPretrainedModel() error = %v", err)
				}
				return
			}

			prediction, err := model.Predict(tt.features)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Predict() error = %v", err)
				}
				return
			}

			// Verify prediction values are within expected ranges
			for key, value := range prediction {
				if value < 0 || value > 1 {
					t.Errorf("Prediction %s = %v, want between 0 and 1", key, value)
				}
			}
		})
	}
} 