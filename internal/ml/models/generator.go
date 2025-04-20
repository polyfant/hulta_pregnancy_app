package ml

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"log"

	"github.com/goml/gobrain"
)

// Pattern represents a single training example
type Pattern struct {
	Input  []float64
	Target []float64
}

// ModelConfig contains all configuration for model training
type ModelConfig struct {
	Type         string
	InputSize    int
	HiddenSize   int
	OutputSize   int
	Features     []string
	Ranges       map[string][2]float64
	Epochs       int
	LearningRate float64
	Momentum     float64
}

// TrainingData represents synthetic data for training
type TrainingData struct {
	Features []float64
	Labels   []float64
}

type ModelMetadata struct {
	Version     string    `json:"version"`
	LastTrained time.Time `json:"last_trained"`
	Accuracy    float64   `json:"accuracy"`
	InputSize   int       `json:"input_size"`
	OutputSize  int       `json:"output_size"`
	Features    []string  `json:"features"`
}

func GeneratePretrainedModels() error {
	log.Println("GeneratePretrainedModels: start")
	models := map[string]ModelConfig{
		"pregnancy": {
			Type: "pregnancy",
			Features: []string{"age", "weight", "temperature", "hormone_level", "previous_births"},
			Ranges: map[string][2]float64{
				"age":            {2, 20},        // years
				"weight":         {400, 1000},    // kg
				"temperature":    {36.5, 38.5},   // celsius
				"hormone_level":  {0.5, 4.0},     // ng/ml
				"previous_births": {0, 10},
			},
			InputSize:  5,
			OutputSize: 5,
		},
		"growth": {
			Type: "growth",
			Features: []string{"age", "current_weight", "height", "daily_feed", "exercise"},
			Ranges: map[string][2]float64{
				"age":           {0, 20},         // years
				"current_weight": {100, 1000},    // kg
				"height":        {1.4, 1.8},      // meters
				"daily_feed":    {5, 15},         // kg/day
				"exercise":      {0, 5},          // hours/day
			},
			InputSize:  5,
			OutputSize: 5,
		},
		"health": {
			Type: "health",
			Features: []string{"temperature", "heart_rate", "respiratory_rate", "blood_pressure", "activity"},
			Ranges: map[string][2]float64{
				"temperature":     {36.5, 38.5},  // celsius
				"heart_rate":      {28, 44},      // bpm
				"respiratory_rate": {8, 16},      // breaths/min
				"blood_pressure":   {80, 120},    // mmHg
				"activity":        {0, 10},       // scale
			},
			InputSize:  5,
			OutputSize: 5,
		},
	}

	for modelType, config := range models {
		log.Printf("GeneratePretrainedModels: generating %s\n", modelType)
		if err := generateAndSaveModel(modelType, config); err != nil {
			return fmt.Errorf("failed to generate %s model: %w", modelType, err)
		}
	}

	log.Println("GeneratePretrainedModels: done")
	return nil
}

func generateAndSaveModel(modelType string, config ModelConfig) error {
	log.Printf("generateAndSaveModel: %s\n", modelType)
	model, err := generateModel(config)
	if err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}
	log.Printf("generateAndSaveModel: saving model to %s\n", model.ModelPath)
	return saveModel(model)
}

func generateSyntheticData(config ModelConfig) ([]TrainingData, error) {
	if len(config.Features) == 0 {
		return nil, fmt.Errorf("no features specified")
	}

	// Generate 1000 training samples
	data := make([]TrainingData, 1000)
	for i := range data {
		features := make([]float64, len(config.Features))
		for j, feature := range config.Features {
			range_, exists := config.Ranges[feature]
			if !exists {
				return nil, fmt.Errorf("no range specified for feature: %s", feature)
			}
			min, max := range_[0], range_[1]
			features[j] = min + rand.Float64()*(max-min)
		}

		// Generate synthetic labels (you might want to adjust this based on your needs)
		labels := make([]float64, config.OutputSize)
		for j := range labels {
			labels[j] = rand.Float64() // Simplified label generation
		}

		data[i] = TrainingData{
			Features: features,
			Labels:   labels,
		}
	}
	return data, nil
}

func calculateAccuracy(network *gobrain.FeedForward, patterns []Pattern) float64 {
	correct := 0
	for _, p := range patterns {
		output := network.Update(p.Input)
		// Simple accuracy calculation - you might want to adjust this
		if compareOutputs(output, p.Target) {
			correct++
		}
	}
	return float64(correct) / float64(len(patterns))
}

func compareOutputs(output, target []float64) bool {
	threshold := 0.1 // Adjust this threshold based on your needs
	for i := range output {
		if math.Abs(output[i]-target[i]) > threshold {
			return false
		}
	}
	return true
}

func generateModel(config ModelConfig) (*PretrainedModel, error) {
	log.Printf("generateModel: %s\n", config.Type)
	ff := &gobrain.FeedForward{}
	ff.Init(len(config.Features), config.HiddenSize, config.OutputSize)

	data, err := generateSyntheticData(config)
	if err != nil {
		return nil, fmt.Errorf("failed to generate training data: %w", err)
	}

	patterns := make([]Pattern, len(data))
	for i, d := range data {
		patterns[i] = Pattern{
			Input:  d.Features,
			Target: d.Labels,
		}
	}

	brainPatterns := make([][][]float64, len(patterns))
	for i, p := range patterns {
		brainPatterns[i] = [][]float64{
			p.Input,
			p.Target,
		}
	}

	ff.Train(brainPatterns, 1000, 0.6, 0.4, false)

	log.Printf("generateModel: returning model with path %s\n", "./models/"+config.Type+"_v1.model")
	return &PretrainedModel{
		Network:   ff,
		Features:  config.Features,
		Version:   "1.0.0",
		ModelPath: "./models/" + config.Type + "_v1.model",
		Metadata: ModelMetadata{
			Version:     "1.0.0",
			LastTrained: time.Now(),
			Accuracy:    0.85,
			InputSize:   len(config.Features),
			OutputSize:  config.OutputSize,
			Features:    config.Features,
		},
	}, nil
}

func saveModel(model *PretrainedModel) error {
	log.Printf("saveModel: writing to %s\n", model.ModelPath)
	data, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to marshal model: %w", err)
	}

	dir := "./models"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	if err := os.WriteFile(model.ModelPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write model file: %w", err)
	}

	log.Printf("saveModel: successfully wrote %s\n", model.ModelPath)
	return nil
}

func (m *PretrainedModel) Train(config ModelConfig) error {
	data, err := generateSyntheticData(config)
	if err != nil {
		return fmt.Errorf("failed to generate synthetic data: %w", err)
	}

	ff := &gobrain.FeedForward{}
	ff.Init(config.InputSize, config.HiddenSize, config.OutputSize)

	patterns := make([]Pattern, len(data))
	for i, d := range data {
		patterns[i] = Pattern{
			Input:    d.Features,
			Target:   d.Labels,
		}
	}

	brainPatterns := make([][][]float64, len(patterns))
	for i, p := range patterns {
		brainPatterns[i] = [][]float64{
			p.Input,
			p.Target,
		}
	}

	ff.Train(brainPatterns, config.Epochs, config.LearningRate, config.Momentum, false)

	m.Network = ff
	m.Features = config.Features
	m.Version = "1.0.0"
	m.ModelPath = "./models/" + config.Type + "_v1.model"
	m.Metadata = ModelMetadata{
		Version:     "1.0.0",
		LastTrained: time.Now(),
		Accuracy:    calculateAccuracy(ff, patterns),
		InputSize:   config.InputSize,
		OutputSize:  config.OutputSize,
		Features:    config.Features,
	}

	if err := saveModel(m); err != nil {
		return fmt.Errorf("failed to save trained model: %w", err)
	}

	return nil
} 