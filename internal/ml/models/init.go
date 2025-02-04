package ml

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
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
		model, err := loadFromPath(path)
		if err != nil {
			return fmt.Errorf("failed to load model %s: %w", path, err)
		}

		modelMu.Lock()
		modelCache[model.Version] = model
		modelMu.Unlock()
	}

	return nil
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