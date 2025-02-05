package ml

import (
	"testing"
	"time"
)

func TestModelPrediction(t *testing.T) {
	config := &Config{
		BatchSize:  32,
		Threshold:  0.8,
	}
	
	service := NewModelService(config)
	collector := NewDataCollector(service)

	// Test cases
	tests := []struct {
		name      string
		modelType string
		features  map[string]float64
		expected  map[string]float64
	}{
		{
			name:      "Growth Prediction",
			modelType: "GROWTH",
			features: map[string]float64{
				"age":         365,  // days
				"weight":      400,  // kg
				"height":      1.5,  // meters
				"temperature": 20,   // celsius
				"exercise":    2,    // hours/day
			},
			expected: map[string]float64{
				"weight_gain": 0.5,  // kg/day
				"height_gain": 0.001, // meters/day
			},
		},
		// Add more test cases
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make prediction
			pred, err := service.Predict(tt.modelType, tt.features)
			if err != nil {
				t.Fatalf("Prediction failed: %v", err)
			}

			// Collect sample
			collector.CollectSample(tt.modelType, Sample{
				HorseID:    "test_horse",
				Timestamp:  time.Now(),
				Features:   tt.features,
				Labels:     tt.expected,
				Prediction: pred,
				Error:     calculateError(pred, tt.expected),
			})

			// Verify predictions are within acceptable range
			for key, expected := range tt.expected {
				if pred[key] < expected*0.5 || pred[key] > expected*1.5 {
					t.Errorf("Prediction for %s outside acceptable range: got %v, want %v", 
						key, pred[key], expected)
				}
			}
		})
	}
}

func calculateError(pred, actual map[string]float64) float64 {
	var totalError float64
	for key, actualVal := range actual {
		if predVal, exists := pred[key]; exists {
			// Mean squared error
			diff := predVal - actualVal
			totalError += diff * diff
		}
	}
	return totalError / float64(len(actual))
} 