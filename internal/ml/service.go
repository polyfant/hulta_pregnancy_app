package ml

import (
	"fmt"
	"time"

	"github.com/goml/gobrain"
)

type ModelService struct {
    models       map[string]*Model
    trainingData map[string][]TrainingData
    config       *Config
}

type ModelMetadata struct {
    Type      string
    Version   string
    Accuracy  float64
    UpdatedAt string
}

type Model struct {
    Network     *gobrain.FeedForward
    Version     string
    LastTrained time.Time
    Accuracy    float64
    Samples     int
    Metadata    ModelMetadata
}

type NeuralNetwork struct {
    Layers []int
    Weights [][]float64
}

type Config struct {
    ModelPaths map[string]string
    BatchSize  int
    Threshold  float64
}

type TrainingData struct {
    Features map[string]float64
    Labels   map[string]float64
}

func NewModelService(config *Config) *ModelService {
    return &ModelService{
        models:       make(map[string]*Model),
        trainingData: make(map[string][]TrainingData),
        config:      config,
    }
}

func (s *ModelService) Predict(modelType string, features map[string]float64) (map[string]float64, error) {
    model, exists := s.models[modelType]
    if !exists {
        return nil, fmt.Errorf("model %s not found", modelType)
    }

    // Convert features to network input
    input := s.preprocessFeatures(features)
    
    // Make prediction
    output := model.Network.Update(input)
    
    return s.postprocessOutput(output), nil
}

func (s *ModelService) Train(modelType string, data []TrainingData) error {
    model, exists := s.models[modelType]
    if !exists {
        // Handle error from initializeModel
        ff, err := s.initializeModel(10, 5, 2) // Adjust these numbers based on your needs
        if err != nil {
            return fmt.Errorf("failed to initialize model: %w", err)
        }
        model = &Model{Network: ff}
        s.models[modelType] = model
    }

    // Convert to [][][]float64 format that gobrain expects
    patterns := make([][][]float64, len(data))
    for i, d := range data {
        patterns[i] = [][]float64{
            s.preprocessFeatures(d.Features),
            s.preprocessLabels(d.Labels),
        }
    }

    // Train the model
    err := s.trainModel(model, patterns)
    if err != nil {
        return fmt.Errorf("training failed: %w", err)
    }

    // Update model metadata
    model.Metadata.Version = fmt.Sprintf("v%d", time.Now().Unix())
    model.Samples += len(data)
    
    return nil
}

func (s *ModelService) trainModel(model *Model, patterns [][][]float64) error {
    if model.Network == nil {
        return fmt.Errorf("network not initialized")
    }
    
    model.Network.Train(patterns, 1000, 0.6, 0.4, false)
    model.LastTrained = time.Now()
    return nil
}

func (s *ModelService) initializeModel(inputSize, hiddenSize, outputSize int) (*gobrain.FeedForward, error) {
    ff := &gobrain.FeedForward{}
    ff.Init(inputSize, hiddenSize, outputSize)
    return ff, nil
}

func (s *ModelService) normalizeFeatures(modelType string, features map[string]float64) ([]float64, error) {
    config, exists := DefaultModels[modelType]
    if !exists {
        return nil, fmt.Errorf("unknown model type: %s", modelType)
    }

    normalized := make([]float64, len(config.Features))
    for i, feature := range config.Features {
        value, exists := features[feature]
        if !exists {
            return nil, fmt.Errorf("missing required feature: %s", feature)
        }
        
        minVal, maxVal := config.Ranges[feature][0], config.Ranges[feature][1]
        normalized[i] = (value - minVal) / (maxVal - minVal)
    }

    return normalized, nil
}

func (s *ModelService) getModelMetadata(modelType string) (*ModelMetadata, error) {
    model, exists := s.models[modelType]
    if !exists {
        return nil, fmt.Errorf("model type %s not found", modelType)
    }
    return &model.Metadata, nil
}

func (s *ModelService) preprocessLabels(labels map[string]float64) []float64 {
    // Implementation
    normalized := make([]float64, len(labels))
    // Add normalization logic
    return normalized
} 