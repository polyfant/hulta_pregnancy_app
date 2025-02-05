package ml

import (
    "sync"
    "time"
)

type DataCollector struct {
    mu           sync.RWMutex
    samples      map[string][]Sample
    modelService *ModelService
}

type Sample struct {
    HorseID    string
    Timestamp  time.Time
    Features   map[string]float64
    Labels     map[string]float64
    Prediction map[string]float64
    Error      float64
}

func NewDataCollector(ms *ModelService) *DataCollector {
    return &DataCollector{
        samples:      make(map[string][]Sample),
        modelService: ms,
    }
}

func (dc *DataCollector) CollectSample(modelType string, sample Sample) {
    dc.mu.Lock()
    defer dc.mu.Unlock()

    dc.samples[modelType] = append(dc.samples[modelType], sample)

    // If we have enough samples, trigger model evaluation
    if len(dc.samples[modelType]) >= 100 {
        go dc.evaluateModel(modelType)
    }
}

func (dc *DataCollector) evaluateModel(modelType string) {
    dc.mu.RLock()
    samples := dc.samples[modelType]
    dc.mu.RUnlock()

    // Calculate model accuracy
    var totalError float64
    for _, sample := range samples {
        totalError += sample.Error
    }
    accuracy := 1 - (totalError / float64(len(samples)))

    // If accuracy drops below threshold, trigger retraining
    if accuracy < 0.8 { // 80% accuracy threshold
        dc.triggerRetraining()
    }
}

// Add specialized collectors for each type
func (dc *DataCollector) CollectPregnancyData(horseId string) {
    // Collect vital signs, hormone levels, etc.
}

func (dc *DataCollector) CollectFoalData(horseId string) {
    // Collect growth rates, health markers, etc.
}

func (dc *DataCollector) triggerRetraining() error {
    // Implementation
    return nil
} 