package ml

import (
    "fmt"
    "math"
    "time"
)

// Feature preprocessing constants
const (
    maxAge         = 10950  // 30 years in days
    maxWeight      = 1000   // kg
    maxHeight      = 2.0    // meters
    maxTemperature = 45     // celsius
    maxExercise    = 8      // hours
)

func (s *ModelService) initializeModel(modelType string) *Model {
    // Create network architecture based on model type
    var network *gobrain.FeedForward
    switch modelType {
    case "GROWTH":
        network = &gobrain.FeedForward{
            InputNodes:    5,  // age, weight, height, temperature, exercise
            HiddenNodes:   8,
            OutputNodes:   2,  // weight_gain, height_gain
            LearningRate: 0.25,
            Momentum:     0.1,
        }
    case "HEALTH":
        network = &gobrain.FeedForward{
            InputNodes:    7,  // includes vital signs
            HiddenNodes:   10,
            OutputNodes:   3,  // health score, risk level, recommendations
            LearningRate: 0.2,
            Momentum:     0.1,
        }
    case "BREEDING":
        network = &gobrain.FeedForward{
            InputNodes:    10, // includes genetic markers
            HiddenNodes:   12,
            OutputNodes:   4,  // success probability, compatibility scores
            LearningRate: 0.15,
            Momentum:     0.1,
        }
    default:
        return nil
    }

    network.Init()

    return &Model{
        Network:     network,
        Version:     fmt.Sprintf("v%d", time.Now().Unix()),
        LastTrained: time.Now(),
        Accuracy:    0.0,
        Samples:     0,
    }
}

func (s *ModelService) preprocessFeatures(features map[string]float64) []float64 {
    // Convert map to normalized array in consistent order
    normalized := make([]float64, len(features))
    i := 0

    // Normalize each feature to [0,1] range
    if age, ok := features["age"]; ok {
        normalized[i] = math.Min(age/maxAge, 1.0)
        i++
    }
    if weight, ok := features["weight"]; ok {
        normalized[i] = math.Min(weight/maxWeight, 1.0)
        i++
    }
    if height, ok := features["height"]; ok {
        normalized[i] = math.Min(height/maxHeight, 1.0)
        i++
    }
    if temp, ok := features["temperature"]; ok {
        normalized[i] = math.Min(temp/maxTemperature, 1.0)
        i++
    }
    if exercise, ok := features["exercise"]; ok {
        normalized[i] = math.Min(exercise/maxExercise, 1.0)
        i++
    }

    return normalized
}

func (s *ModelService) postprocessOutput(output []float64) map[string]float64 {
    // Convert normalized output back to real values
    result := make(map[string]float64)
    
    // Denormalize based on model type
    if len(output) >= 2 {
        result["weight_gain"] = output[0] * 2.0  // max 2 kg/day
        result["height_gain"] = output[1] * 0.01 // max 1 cm/day
    }
    
    return result
}

func (s *ModelService) triggerRetraining(modelType string, samples []Sample) error {
    // Prepare training data
    trainingData := make([]TrainingData, len(samples))
    for i, sample := range samples {
        trainingData[i] = TrainingData{
            Features: sample.Features,
            Labels:   sample.Labels,
        }
    }

    // Train model
    if err := s.Train(modelType, trainingData); err != nil {
        return fmt.Errorf("retraining failed: %w", err)
    }

    return nil
} 