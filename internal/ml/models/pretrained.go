package ml

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/goml/gobrain"
)

type PretrainedModel struct {
	Network     *gobrain.FeedForward
	Metadata    ModelMetadata
	Features    []string
	Version     string
	ModelPath   string
}

func LoadPretrainedModel(modelType string) (*PretrainedModel, error) {
	// Load smaller, pre-trained models instead of training from scratch
	switch modelType {
	case "PREGNANCY":
		return loadFromPath("models/pregnancy_v1.model")
	case "GROWTH":
		return loadFromPath("models/growth_v1.model")
	case "HEALTH":
		return loadFromPath("models/health_v1.model")
	}
	return nil, fmt.Errorf("unknown model type: %s", modelType)
}

func loadFromPath(path string) (*PretrainedModel, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open model file: %w", err)
	}
	defer file.Close()

	var model PretrainedModel
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&model); err != nil {
		return nil, fmt.Errorf("failed to decode model: %w", err)
	}

	return &model, nil
}

func (m *PretrainedModel) Predict(features map[string]float64) (map[string]float64, error) {
	// Convert map to slice in correct order
	inputData := make([]float64, len(m.Features))
	for i, feature := range m.Features {
		val, ok := features[feature]
		if !ok {
			return nil, fmt.Errorf("missing feature: %s", feature)
		}
		inputData[i] = val
	}

	output := m.Network.Update(inputData)
	
	result := make(map[string]float64)
	for i, val := range output {
		result[m.Features[i]] = val
	}
	return result, nil
} 