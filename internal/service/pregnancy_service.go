package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

// PregnancyService handles business logic for horse pregnancies
type PregnancyService struct {
	horseRepo      repository.HorseRepository
	pregnancyRepo  repository.PregnancyRepository
}

// NewPregnancyService creates a new pregnancy service instance
func NewPregnancyService(horseRepo repository.HorseRepository, pregnancyRepo repository.PregnancyRepository) *PregnancyService {
	return &PregnancyService{
		horseRepo:      horseRepo,
		pregnancyRepo:  pregnancyRepo,
	}
}

// Essential methods for pregnancy tracking
func (s *PregnancyService) GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	return s.pregnancyRepo.GetByHorseID(ctx, horseID)
}

func (s *PregnancyService) StartPregnancy(ctx context.Context, horseID uint, userID string, conceptionDate time.Time) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	if horse.UserID != userID {
		return fmt.Errorf("unauthorized: horse does not belong to user")
	}

	horse.IsPregnant = true
	horse.ConceptionDate = &conceptionDate

	if err := s.horseRepo.Update(ctx, horse); err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}

	pregnancy := &models.Pregnancy{
		HorseID:        horseID,
		ConceptionDate: &conceptionDate,
		Status:         string(models.PregnancyStatusActive),
	}

	return s.pregnancyRepo.Create(ctx, pregnancy)
}

func (s *PregnancyService) EndPregnancy(ctx context.Context, horseID uint, status string) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	horse.IsPregnant = false
	horse.ConceptionDate = nil

	if err := s.horseRepo.Update(ctx, horse); err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}

	pregnancy, err := s.pregnancyRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get pregnancy: %w", err)
	}

	pregnancy.Status = status
	now := time.Now()
	pregnancy.EndDate = &now

	return s.pregnancyRepo.Update(ctx, pregnancy)
}

func (s *PregnancyService) GetPregnancyStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error) {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return nil, err
	}

	if !horse.IsPregnant || horse.ConceptionDate == nil {
		return &models.PregnancyStatus{
			IsPregnant: false,
		}, nil
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	progress := float64(daysPregnant) / float64(models.DefaultGestationDays) * 100
	
	return &models.PregnancyStatus{
		IsPregnant:          true,
		ConceptionDate:      *horse.ConceptionDate,
		DaysPregnant:        daysPregnant,
		PregnancyPercentage: progress,
		Stage:               s.calculateStage(&models.Pregnancy{ConceptionDate: horse.ConceptionDate}),
	}, nil
}

// Helper method
func (s *PregnancyService) calculateStage(pregnancy *models.Pregnancy) models.PregnancyStage {
	if pregnancy == nil || pregnancy.ConceptionDate == nil {
		return models.PregnancyStageUnknown
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)
	progressPercentage := float64(daysPregnant) / float64(models.DefaultGestationDays)
	
	switch {
	case progressPercentage <= 0.33:
		return models.PregnancyStageEarlyGestation
	case progressPercentage <= 0.66:
		return models.PregnancyStageMidGestation
	default:
		return models.PregnancyStageLateGestation
	}
}

// Additional features - to be implemented later
func (s *PregnancyService) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	// TODO: Implement pregnancy event tracking
	return []models.PregnancyEvent{}, nil
}

func (s *PregnancyService) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	// TODO: Implement event adding
	return nil
}

func (s *PregnancyService) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	// TODO: Implement pre-foaling checklist
	return []models.PreFoalingChecklistItem{}, nil
}

func (s *PregnancyService) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	// TODO: Implement pre-foaling signs tracking
	return []models.PreFoalingSign{}, nil
}

func (s *PregnancyService) GetPregnancyGuidelinesByStage(stage models.PregnancyStage) ([]models.PregnancyGuideline, error) {
	// Basic guidelines implementation - can be enhanced later
	guidelines := map[models.PregnancyStage][]models.PregnancyGuideline{
		models.PregnancyStageEarlyGestation: {
			{
				Stage:       models.PregnancyStageEarlyGestation,
				Category:    "Nutrition",
				Description: "Ensure balanced diet with proper nutrients",
			},
		},
		models.PregnancyStageMidGestation: {
			{
				Stage:       models.PregnancyStageMidGestation,
				Category:    "Health Monitoring",
				Description: "Regular veterinary check-ups",
			},
		},
		models.PregnancyStageLateGestation: {
			{
				Stage:       models.PregnancyStageLateGestation,
				Category:    "Preparation",
				Description: "Prepare foaling area and emergency kit",
			},
		},
	}
	
	if _, exists := guidelines[stage]; !exists {
		return nil, fmt.Errorf("invalid pregnancy stage: %v", stage)
	}
	
	return guidelines[stage], nil
}

// Helper types and methods that might be useful later
func (s *PregnancyService) UpdateHorsePregnancyStatus(ctx context.Context, horseID uint, userID string, isPregnant bool) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	if horse.UserID != userID {
		return fmt.Errorf("unauthorized: horse does not belong to user")
	}

	horse.IsPregnant = isPregnant
	if !isPregnant {
		horse.ConceptionDate = nil
	}

	return s.horseRepo.Update(ctx, horse)
}

/* 
// These methods can be uncommented and implemented when needed
func (s *PregnancyService) AddPreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	return nil
}

func (s *PregnancyService) GetPreFoalingSign(sign *models.PreFoalingSign) error {
	return nil
}
*/

// Add missing method for getting pregnancy stage
func (s *PregnancyService) GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error) {
	pregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		return "UNKNOWN", fmt.Errorf("failed to get pregnancy: %w", err)
	}

	if pregnancy == nil {
		return "UNKNOWN", nil
	}

	return s.calculateStage(pregnancy), nil
}
