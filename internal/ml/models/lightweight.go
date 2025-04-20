package ml

// Instead of deep neural networks, use:
// - Random Forests
// - Gradient Boosting
// - Simple regression models
type LightweightModel struct {
    Algorithm string // "rf", "gbm", "regression"
    Features  []string
    Weights   map[string]float64
} 