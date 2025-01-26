package models

import "fmt"

// FeedRequirements represents the daily feed requirements for a horse
type FeedRequirements struct {
    Hay      float64 `json:"hay"`      // Daily hay in kg
    Grain    float64 `json:"grain"`    // Daily grain in kg
    Minerals float64 `json:"minerals"` // Daily minerals in kg
    Water    float64 `json:"water"`    // Daily water in liters
}

// Validate checks if the feed requirements are within reasonable bounds
func (f *FeedRequirements) Validate() error {
    if f.Hay < 0 || f.Hay > 20 {
        return fmt.Errorf("hay amount must be between 0 and 20 kg")
    }
    if f.Grain < 0 || f.Grain > 5 {
        return fmt.Errorf("grain amount must be between 0 and 5 kg")
    }
    if f.Minerals < 0 || f.Minerals > 0.5 {
        return fmt.Errorf("minerals amount must be between 0 and 0.5 kg")
    }
    if f.Water < 0 || f.Water > 100 {
        return fmt.Errorf("water amount must be between 0 and 100 liters")
    }
    return nil
} 