package health

import (
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

var timeNow = time.Now // Allow overriding in tests

// ActivityLevel represents the horse's current activity level
type ActivityLevel int

const (
	Maintenance ActivityLevel = iota
	LightWork
	ModerateWork
	HeavyWork
)

// Season represents the current season
type Season int

const (
	Spring Season = iota
	Summer
	Autumn
	Winter
)

// NutritionService handles the calculation and management of horse nutrition requirements.
// It takes into account factors such as:
// - Horse weight and size
// - Activity level
// - Pregnancy status
// - Environmental conditions
type NutritionService struct {
	healthRepo repository.HealthRepository
	horseRepo  repository.HorseRepository
}

func NewNutritionService(healthRepo repository.HealthRepository, horseRepo repository.HorseRepository) *NutritionService {
	return &NutritionService{
		healthRepo: healthRepo,
		horseRepo:  horseRepo,
	}
}

// CalculateDailyFeedRequirements calculates feed requirements based on horse's condition
func (s *NutritionService) CalculateDailyFeedRequirements(horse models.Horse, activity ActivityLevel) (models.FeedRequirements, error) {
	if err := activity.Validate(); err != nil {
		return models.FeedRequirements{}, err
	}

	baseRequirement := s.calculateBaseFeedRequirement(horse)
	
	// Adjust for activity level
	baseRequirement = s.adjustForActivity(baseRequirement, activity)
	
	// Adjust for pregnancy if applicable
	if horse.IsPregnant {
		stage := s.getPregnancyStage(horse)
		baseRequirement = s.adjustForPregnancy(baseRequirement, stage)
	}

	// Adjust for seasonal changes
	baseRequirement = s.adjustForSeason(baseRequirement, getCurrentSeason())

	// Validate the final requirements
	if err := baseRequirement.Validate(); err != nil {
		return models.FeedRequirements{}, fmt.Errorf("invalid feed requirements: %w", err)
	}

	return baseRequirement, nil
}

func (s *NutritionService) calculateBaseFeedRequirement(horse models.Horse) models.FeedRequirements {
	// Base calculations using weight and activity level
	weight := horse.Weight
	if weight == 0 {
		weight = 500 // Default weight if not specified
		logger.Warn("Using default weight of 500kg for horse", map[string]interface{}{
			"horseID": horse.ID,
			"horseName": horse.Name,
		})
	}

	return models.FeedRequirements{
		Hay:      weight * 0.02, // 2% of body weight
		Grain:    weight * 0.005, // 0.5% of body weight
		Minerals: 0.1,            // 100g minerals
		Water:    weight * 0.05,  // 5% of body weight
	}
}

func (s *NutritionService) adjustForPregnancy(base models.FeedRequirements, stage models.PregnancyStage) models.FeedRequirements {
	switch stage {
	case models.PregnancyStageEarlyGestation:
		// Minimal increase in early pregnancy
		base.Hay *= 1.1
		base.Minerals *= 1.2
	case models.PregnancyStageMidGestation:
		// Moderate increase
		base.Hay *= 1.2
		base.Grain *= 1.1
		base.Minerals *= 1.3
	case models.PregnancyStageLateGestation:
		// Significant increase
		base.Hay *= 1.3
		base.Grain *= 1.2
		base.Minerals *= 1.5
		base.Water *= 1.2
	}
	return base
}

func (s *NutritionService) getPregnancyStage(horse models.Horse) models.PregnancyStage {
	if !horse.IsPregnant || horse.ConceptionDate == nil {
		return ""
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	switch {
	case daysPregnant <= 113:
		return models.PregnancyStageEarlyGestation
	case daysPregnant <= 226:
		return models.PregnancyStageMidGestation
	default:
		return models.PregnancyStageLateGestation
	}
}

func (s *NutritionService) adjustForActivity(base models.FeedRequirements, activity ActivityLevel) models.FeedRequirements {
	switch activity {
	case LightWork:
		base.Hay *= 1.1
		base.Grain *= 1.25
		base.Water *= 1.2
	case ModerateWork:
		base.Hay *= 1.2
		base.Grain *= 1.5
		base.Water *= 1.4
	case HeavyWork:
		base.Hay *= 1.3
		base.Grain *= 2.0
		base.Water *= 1.6
	}
	return base
}

func getCurrentSeason() Season {
	month := timeNow().Month()
	switch {
	case month >= 3 && month <= 5:
		return Spring
	case month >= 6 && month <= 8:
		return Summer
	case month >= 9 && month <= 11:
		return Autumn
	default:
		return Winter
	}
}

func (s *NutritionService) adjustForSeason(base models.FeedRequirements, season Season) models.FeedRequirements {
	switch season {
	case Winter:
		// Increase hay for warmth and energy
		base.Hay *= 1.15
		base.Grain *= 1.1
	case Summer:
		// Increase water for hydration
		base.Water *= 1.3
		// Slightly decrease hay due to available pasture
		base.Hay *= 0.9
	}
	return base
}

func (a ActivityLevel) Validate() error {
	if a < Maintenance || a > HeavyWork {
		return fmt.Errorf("invalid activity level: %d", a)
	}
	return nil
}
