package pregnancy

import (
	"context"
	"fmt"
	"time"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

// Service handles pregnancy-related business logic
type Service struct {
	pregnancyRepo repository.PregnancyRepository
}

// NewService creates a new pregnancy service
func NewService(pregnancyRepo repository.PregnancyRepository) *Service {
	return &Service{
		pregnancyRepo: pregnancyRepo,
	}
}

// DefaultGestationDays is the average number of days in a horse's pregnancy
const DefaultGestationDays = 340

// CalculateDueDate calculates the expected due date based on conception date and gestation days
func CalculateDueDate(conceptionDate time.Time, gestationDays int) time.Time {
	const averageHorseGestationDays = 340
	if gestationDays <= 0 {
		gestationDays = averageHorseGestationDays
	}
	return conceptionDate.AddDate(0, 0, gestationDays)
}

// CalculateGestationProgress calculates the pregnancy progress and days remaining
func CalculateGestationProgress(conceptionDate time.Time, gestationDays int) (float64, int) {
	const averageHorseGestationDays = 340
	if gestationDays <= 0 {
		gestationDays = averageHorseGestationDays
	}

	daysPregnant := time.Since(conceptionDate).Hours() / 24
	progress := (daysPregnant / float64(gestationDays)) * 100
	daysRemaining := int(float64(gestationDays) - daysPregnant)

	// Ensure progress doesn't exceed 100%
	if progress > 100 {
		progress = 100
		daysRemaining = 0
	}

	return progress, daysRemaining
}

// CalculatePregnancyProgress calculates the progress of a horse's pregnancy
func (s *Service) CalculatePregnancyProgress(pregnancy *models.Pregnancy) (float64, int, string) {
	const totalPregnancyDays = 340 // Average horse pregnancy duration

	daysPregnant := time.Since(pregnancy.StartDate).Hours() / 24
	progress := (daysPregnant / totalPregnancyDays) * 100
	daysRemaining := int(totalPregnancyDays - daysPregnant)

	var stage string
	switch {
	case daysPregnant < 120:
		stage = string(models.EarlyGestation)
	case daysPregnant < 240:
		stage = string(models.MidGestation)
	default:
		stage = string(models.LateGestation)
	}

	return progress, daysRemaining, stage
}

// GetPregnancyGuidelines returns guidelines for different pregnancy stages
func (s *Service) GetPregnancyGuidelines() []models.PregnancyGuideline {
	return []models.PregnancyGuideline{
		{
			Stage:       models.EarlyGestation,
			Category:    "Nutrition",
			Description: "Maintain regular diet, slight increase in nutrients",
		},
		{
			Stage:       models.MidGestation,
			Category:    "Nutrition",
			Description: "Increase nutrient intake, monitor weight",
		},
		{
			Stage:       models.LateGestation,
			Category:    "Nutrition",
			Description: "High-quality diet, prepare for foaling",
		},
	}
}

// GetPregnancyGuidelinesByStage returns guidelines for a specific pregnancy stage
func (s *Service) GetPregnancyGuidelinesByStage(stage models.PregnancyStage) []models.PregnancyGuideline {
	guidelines := s.GetPregnancyGuidelines()
	var stageGuidelines []models.PregnancyGuideline
	
	for _, guideline := range guidelines {
		if guideline.Stage == stage {
			stageGuidelines = append(stageGuidelines, guideline)
		}
	}
	
	return stageGuidelines
}

// GetPregnancyStage retrieves the current pregnancy stage for a specific horse
func (s *Service) GetPregnancyStage(ctx context.Context, horseID int64) (models.PregnancyStage, error) {
	// Fetch the current pregnancy for the horse
	pregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		return "", err
	}

	// If no active pregnancy, return an appropriate response
	if pregnancy == nil {
		return "", fmt.Errorf("no active pregnancy found for horse")
	}

	// Calculate the stage based on days pregnant
	_, _, stage := s.CalculatePregnancyProgress(pregnancy)
	return models.PregnancyStage(stage), nil
}
