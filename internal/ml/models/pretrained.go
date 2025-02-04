package ml

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
)

type PretrainedModel struct {
	BaseModel *ensemble.RandomForest  // Lightweight model
	ModelPath string
	Version   string
	Features  []string
	Metadata  ModelMetadata
}

type ModelMetadata struct {
	Accuracy      float64
	LastUpdated   string
	SampleSize    int
	FeatureRanges map[string][2]float64  // min/max ranges for normalization
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
	// Normalize input features
	normalizedFeatures := make([]float64, len(m.Features))
	for i, feature := range m.Features {
		value := features[feature]
		minVal, maxVal := m.Metadata.FeatureRanges[feature][0], m.Metadata.FeatureRanges[feature][1]
		normalizedFeatures[i] = (value - minVal) / (maxVal - minVal)
	}

	// Make prediction using the random forest
	prediction, err := m.BaseModel.Predict(normalizedFeatures)
	if err != nil {
		return nil, fmt.Errorf("prediction failed: %w", err)
	}

	return prediction, nil
} 