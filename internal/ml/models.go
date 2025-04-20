package ml

type ModelConfig struct {
    Features []string
    Ranges   map[string][2]float64
}

var DefaultModels = map[string]ModelConfig{
    "pregnancy": {
        Features: []string{"age", "weight", "temperature", "hormone_level", "previous_births"},
        Ranges: map[string][2]float64{
            "age":            {2, 20},        // years
            "weight":         {400, 1000},    // kg
            "temperature":    {36.5, 38.5},   // celsius
            "hormone_level":  {0.5, 4.0},     // ng/ml
            "previous_births": {0, 10},
        },
    },
    "growth": {
        Features: []string{"age", "current_weight", "height", "daily_feed", "exercise"},
        Ranges: map[string][2]float64{
            "age":           {0, 20},         // years
            "current_weight": {100, 1000},    // kg
            "height":        {1.4, 1.8},      // meters
            "daily_feed":    {5, 15},         // kg/day
            "exercise":      {0, 5},          // hours/day
        },
    },
    "health": {
        Features: []string{"temperature", "heart_rate", "respiratory_rate", "blood_pressure", "activity"},
        Ranges: map[string][2]float64{
            "temperature":     {36.5, 38.5},  // celsius
            "heart_rate":      {28, 44},      // bpm
            "respiratory_rate": {8, 16},      // breaths/min
            "blood_pressure":   {80, 120},    // mmHg
            "activity":        {0, 10},       // scale
        },
    },
} 