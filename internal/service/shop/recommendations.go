package shop

import (
	"context"
	"time"
)

// Product represents a shop product with its details
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	StoreURL    string    `json:"store_url"`
}

// RecommendationType indicates why a product is being recommended
type RecommendationType string

const (
	WeatherBased     RecommendationType = "WEATHER"
	PregnancyStage   RecommendationType = "PREGNANCY_STAGE"
	SeasonalCare     RecommendationType = "SEASONAL"
	EmergencySupplies RecommendationType = "EMERGENCY"
)

// ProductRecommendation represents a recommended product with context
type ProductRecommendation struct {
	Product       Product           `json:"product"`
	ReasonForRec  string           `json:"reason"`
	Type          RecommendationType `json:"type"`
	Priority      int              `json:"priority"` // 1-5, with 1 being highest
	ValidUntil    time.Time        `json:"valid_until,omitempty"`
}

// RecommendationService handles product recommendations
type RecommendationService interface {
	// GetWeatherBasedRecs returns products recommended for current weather
	GetWeatherBasedRecs(ctx context.Context, temp float64, conditions []string) ([]ProductRecommendation, error)

	// GetPregnancyStageRecs returns products recommended for pregnancy stage
	GetPregnancyStageRecs(ctx context.Context, stage string, daysPregnant int) ([]ProductRecommendation, error)

	// GetEmergencyKitRecs returns recommended items for emergency preparation
	GetEmergencyKitRecs(ctx context.Context) ([]ProductRecommendation, error)

	// GetSeasonalRecs returns season-appropriate recommendations
	GetSeasonalRecs(ctx context.Context, season string) ([]ProductRecommendation, error)
}

// Example usage for future implementation:
/*
recommendations := []ProductRecommendation{
	{
		Product: Product{
			ID:          "FOALING-KIT-001",
			Name:        "Premium Foaling Kit",
			Description: "Complete kit with all essential items for foaling",
			Category:    "Foaling Supplies",
			Price:       149.99,
			ImageURL:    "https://example.com/images/foaling-kit.jpg",
			StoreURL:    "https://yoursisters.shop/foaling-kit",
		},
		ReasonForRec: "Approaching due date - ensure you're prepared with essential foaling supplies",
		Type:         PregnancyStage,
		Priority:     1,
	},
	{
		Product: Product{
			ID:          "ELECTRO-PLUS-500",
			Name:        "Electrolyte Plus Supplement",
			Description: "Premium electrolyte supplement for hot weather",
			Category:    "Supplements",
			Price:       34.99,
			ImageURL:    "https://example.com/images/electrolytes.jpg",
			StoreURL:    "https://yoursisters.shop/electrolytes",
		},
		ReasonForRec: "High temperature forecast - maintain proper hydration",
		Type:         WeatherBased,
		Priority:     2,
	},
}
*/
