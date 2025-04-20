package ml

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	modelCache = make(map[string]*PretrainedModel)
	modelMu    sync.RWMutex
)

func InitializeModels() error {
	models := []string{
		"models/pregnancy_v1.model",
		"models/growth_v1.model",
		"models/health_v1.model",
	}

	for _, path := range models {
		model, err := loadModel(path)
		if err != nil {
			return fmt.Errorf("failed to load model %s: %w", path, err)
		}

		modelMu.Lock()
		modelCache[model.Version] = model
		modelMu.Unlock()
	}

	return nil
}

func loadModel(path string) (*PretrainedModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read model file: %w", err)
	}

	var model PretrainedModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, fmt.Errorf("failed to unmarshal model: %w", err)
	}

	return &model, nil
} 