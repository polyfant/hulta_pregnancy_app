package models

import (
	"time"

	"gorm.io/gorm"
)

// Recommendation represents a personalized suggestion or product recommendation
type Recommendation struct {
	gorm.Model
	UserID        string    `json:"user_id"`
	HorseID       uint      `json:"horse_id"`
	Type          string    `json:"type"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Priority      string    `json:"priority"`
	Category      string    `json:"category"`
	URL           string    `json:"url,omitempty"`
	ImageURL      string    `json:"image_url,omitempty"`
	ExpiresAt     time.Time `json:"expires_at,omitempty"`
	Relevance     float64   `json:"relevance"`
	IsPersonal    bool      `json:"is_personal"`
	IsActionable  bool      `json:"is_actionable"`
}

// RecommendationType defines the types of recommendations
const (
	RecommendationTypeProduct    = "product"
	RecommendationTypeHealth     = "health"
	RecommendationTypeTraining   = "training"
	RecommendationTypeNutrition  = "nutrition"
	RecommendationTypePregnancy  = "pregnancy"
	RecommendationTypeGeneral    = "general"
)

// RecommendationPriority defines the priority levels
const (
	RecommendationPriorityLow      = "low"
	RecommendationPriorityMedium   = "medium"
	RecommendationPriorityHigh     = "high"
	RecommendationPriorityUrgent   = "urgent"
)
