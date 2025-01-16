package shop

import (
	"context"
	"time"
)

// ExampleRecommendationService shows how the shop integration would work
type ExampleRecommendationService struct{}

func (s *ExampleRecommendationService) GetWeatherBasedRecs(ctx context.Context, temp float64, conditions []string) ([]ProductRecommendation, error) {
	recs := []ProductRecommendation{}

	// Hot weather recommendations
	if temp >= 25 {
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "ELECTRO-500",
				Name:        "Premium Electrolyte Supplement",
				Description: "Keep your horse hydrated in hot weather with our premium electrolyte mix",
				Category:    "Supplements",
				Price:       29.99,
				ImageURL:    "https://yoursisters.shop/images/electrolytes.jpg",
				StoreURL:    "https://yoursisters.shop/product/premium-electrolytes",
			},
			ReasonForRec: "High temperature - essential for maintaining hydration",
			Type:         WeatherBased,
			Priority:     1,
			ValidUntil:   time.Now().Add(24 * time.Hour), // Valid for next 24 hours
		})

		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "COOL-BLANKET-1",
				Name:        "Cooling Horse Blanket",
				Description: "Lightweight, breathable blanket that helps regulate temperature",
				Category:    "Equipment",
				Price:       89.99,
				ImageURL:    "https://yoursisters.shop/images/cooling-blanket.jpg",
				StoreURL:    "https://yoursisters.shop/product/cooling-blanket",
			},
			ReasonForRec: "Protect from sun while staying cool",
			Type:         WeatherBased,
			Priority:     2,
		})
	}

	// Cold weather recommendations
	if temp <= 5 {
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "WINTER-FEED-1",
				Name:        "Winter Wellness Feed Mix",
				Description: "High-energy feed mix perfect for cold weather",
				Category:    "Feed",
				Price:       45.99,
				ImageURL:    "https://yoursisters.shop/images/winter-feed.jpg",
				StoreURL:    "https://yoursisters.shop/product/winter-feed",
			},
			ReasonForRec: "Extra energy needed in cold weather",
			Type:         WeatherBased,
			Priority:     1,
		})
	}

	return recs, nil
}

func (s *ExampleRecommendationService) GetPregnancyStageRecs(ctx context.Context, stage string, daysPregnant int) ([]ProductRecommendation, error) {
	recs := []ProductRecommendation{}

	switch stage {
	case "LATE_GESTATION":
		// Essential foaling supplies
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "FOALING-KIT-PRO",
				Name:        "Professional Foaling Kit",
				Description: "Complete kit with everything needed for foaling",
				Category:    "Foaling Supplies",
				Price:       199.99,
				ImageURL:    "https://yoursisters.shop/images/foaling-kit.jpg",
				StoreURL:    "https://yoursisters.shop/product/pro-foaling-kit",
			},
			ReasonForRec: "Essential supplies for upcoming foaling",
			Type:         PregnancyStage,
			Priority:     1,
		})

		// Special late-stage nutrition
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "PREG-NUTRITION-3",
				Name:        "Late-Stage Pregnancy Feed",
				Description: "Specially formulated for the final trimester",
				Category:    "Feed",
				Price:       65.99,
				ImageURL:    "https://yoursisters.shop/images/late-stage-feed.jpg",
				StoreURL:    "https://yoursisters.shop/product/late-stage-feed",
			},
			ReasonForRec: "Optimal nutrition for final pregnancy stage",
			Type:         PregnancyStage,
			Priority:     2,
		})

	case "MID_GESTATION":
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "PREG-VITAMIN-2",
				Name:        "Mid-Stage Pregnancy Vitamins",
				Description: "Essential vitamins for developing foal",
				Category:    "Supplements",
				Price:       42.99,
				ImageURL:    "https://yoursisters.shop/images/pregnancy-vitamins.jpg",
				StoreURL:    "https://yoursisters.shop/product/pregnancy-vitamins",
			},
			ReasonForRec: "Support healthy foal development",
			Type:         PregnancyStage,
			Priority:     2,
		})
	}

	return recs, nil
}

func (s *ExampleRecommendationService) GetEmergencyKitRecs(ctx context.Context) ([]ProductRecommendation, error) {
	return []ProductRecommendation{
		{
			Product: Product{
				ID:          "EMERG-KIT-1",
				Name:        "Emergency First Aid Kit",
				Description: "Comprehensive first aid kit for horses",
				Category:    "Emergency",
				Price:       149.99,
				ImageURL:    "https://yoursisters.shop/images/first-aid-kit.jpg",
				StoreURL:    "https://yoursisters.shop/product/emergency-kit",
			},
			ReasonForRec: "Essential emergency supplies",
			Type:         EmergencySupplies,
			Priority:     1,
		},
		{
			Product: Product{
				ID:          "EMERG-BLANKET-1",
				Name:        "Emergency Thermal Blanket",
				Description: "High-performance blanket for emergency situations",
				Category:    "Emergency",
				Price:       79.99,
				ImageURL:    "https://yoursisters.shop/images/thermal-blanket.jpg",
				StoreURL:    "https://yoursisters.shop/product/thermal-blanket",
			},
			ReasonForRec: "Emergency temperature regulation",
			Type:         EmergencySupplies,
			Priority:     2,
		},
	}, nil
}

func (s *ExampleRecommendationService) GetSeasonalRecs(ctx context.Context, season string) ([]ProductRecommendation, error) {
	recs := []ProductRecommendation{}

	switch season {
	case "SUMMER":
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "FLY-SPRAY-1",
				Name:        "Natural Fly Spray",
				Description: "Chemical-free fly protection",
				Category:    "Care",
				Price:       24.99,
				ImageURL:    "https://yoursisters.shop/images/fly-spray.jpg",
				StoreURL:    "https://yoursisters.shop/product/fly-spray",
			},
			ReasonForRec: "Essential summer fly protection",
			Type:         SeasonalCare,
			Priority:     1,
		})
	case "WINTER":
		recs = append(recs, ProductRecommendation{
			Product: Product{
				ID:          "HOOF-OIL-WINTER",
				Name:        "Winter Hoof Oil",
				Description: "Extra protection for winter conditions",
				Category:    "Care",
				Price:       19.99,
				ImageURL:    "https://yoursisters.shop/images/hoof-oil.jpg",
				StoreURL:    "https://yoursisters.shop/product/winter-hoof-oil",
			},
			ReasonForRec: "Protect hooves in wet winter conditions",
			Type:         SeasonalCare,
			Priority:     1,
		})
	}

	return recs, nil
}
