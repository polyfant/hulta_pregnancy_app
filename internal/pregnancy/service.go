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
	repo repository.PregnancyRepository
}

// NewService creates a new pregnancy service
func NewService(repo repository.PregnancyRepository) *Service {
	return &Service{repo: repo}
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
		stage = string(models.PregnancyStageEarly)
	case daysPregnant < 240:
		stage = string(models.PregnancyStageMid)
	default:
		stage = string(models.PregnancyStageLate)
	}

	return progress, daysRemaining, stage
}

// GetPregnancyGuidelines returns guidelines for different pregnancy stages
func (s *Service) GetPregnancyGuidelines() []models.PregnancyGuideline {
	return []models.PregnancyGuideline{
		{
			Stage:       models.PregnancyStageEarly,
			Description: "Maintain regular diet, slight increase in nutrients",
		},
		{
			Stage:       models.PregnancyStageMid,
			Description: "Increase nutrient intake, monitor weight",
		},
		{
			Stage:       models.PregnancyStageLate,
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
func (s *Service) GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error) {
	pregnancy, err := s.repo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		return models.PregnancyStage("UNKNOWN"), fmt.Errorf("failed to get pregnancy: %w", err)
	}

	if pregnancy == nil {
		return models.PregnancyStage("UNKNOWN"), nil
	}

	return s.calculateStage(pregnancy), nil
}

func (s *Service) calculateStage(pregnancy *models.Pregnancy) models.PregnancyStage {
	if pregnancy.ConceptionDate == nil {
		return models.PregnancyStageEarly
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)
	progressPercentage := float64(daysPregnant) / float64(models.DefaultGestationDays)
	
	switch {
	case progressPercentage <= 0.33:
		return models.PregnancyStageEarly
	case progressPercentage <= 0.66:
		return models.PregnancyStageMid
	default:
		return models.PregnancyStageLate
	}
}

// GetPregnancy retrieves the pregnancy for a specific horse
func (s *Service) GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	pregnancy, err := s.repo.GetByHorseID(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy: %w", err)
	}
	return pregnancy, nil
}

func (s *Service) GetPregnancyStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error) {
	pregnancy, err := s.repo.GetByHorseID(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy: %w", err)
	}

	if pregnancy.ConceptionDate == nil {
		return nil, fmt.Errorf("pregnancy has no conception date")
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)
	dueDate := pregnancy.ConceptionDate.Add(340 * 24 * time.Hour)

	status := &models.PregnancyStatus{
		DaysPregnant: daysPregnant,
		DueDate:      dueDate,
	}

	return status, nil
}
