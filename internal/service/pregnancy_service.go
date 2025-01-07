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
	repo           repository.PregnancyRepository
}

// NewPregnancyService creates a new pregnancy service instance
func NewPregnancyService(horseRepo repository.HorseRepository, pregnancyRepo repository.PregnancyRepository) *PregnancyService {
	return &PregnancyService{
		horseRepo:      horseRepo,
		pregnancyRepo:  pregnancyRepo,
		repo:           pregnancyRepo,
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
	switch stage {
	case models.PregnancyStageEarlyGestation:
		return []models.PregnancyGuideline{
			{Stage: stage, Description: "Monitor mare's appetite and weight"},
			{Stage: stage, Description: "Continue regular exercise routine"},
			{Stage: stage, Description: "Schedule initial pregnancy check"},
		}, nil
	case models.PregnancyStageMidGestation:
		return []models.PregnancyGuideline{
			{Stage: stage, Description: "Adjust feed for increasing nutritional needs"},
			{Stage: stage, Description: "Monitor for any signs of discomfort"},
			{Stage: stage, Description: "Schedule mid-term pregnancy check"},
		}, nil
	case models.PregnancyStageLateGestation:
		return []models.PregnancyGuideline{
			{Stage: stage, Description: "Prepare foaling area"},
			{Stage: stage, Description: "Monitor for signs of approaching labor"},
			{Stage: stage, Description: "Have vet contact information ready"},
		}, nil
	default:
		return nil, fmt.Errorf("invalid pregnancy stage: %s", stage)
	}
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

func (s *PregnancyService) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return s.repo.AddPreFoalingChecklistItem(ctx, item)
}

func (s *PregnancyService) UpdatePreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return s.repo.UpdatePreFoalingChecklistItem(ctx, item)
}

func (s *PregnancyService) GetPreFoalingChecklistItem(ctx context.Context, itemID uint) (*models.PreFoalingChecklistItem, error) {
	return s.repo.GetPreFoalingChecklistItem(ctx, itemID)
}

func (s *PregnancyService) DeletePreFoalingChecklistItem(ctx context.Context, itemID uint) error {
	return s.repo.DeletePreFoalingChecklistItem(ctx, itemID)
}

func (s *PregnancyService) InitializePreFoalingChecklist(ctx context.Context, horseID uint) error {
	return s.repo.InitializePreFoalingChecklist(ctx, horseID)
}

func (s *PregnancyService) UpdatePregnancyStatus(ctx context.Context, horseID uint, isPregnant bool, conceptionDate *time.Time) error {
	return s.repo.UpdatePregnancyStatus(ctx, horseID, isPregnant, conceptionDate)
}

func (s *PregnancyService) GetPregnancies(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	return s.pregnancyRepo.GetByUserID(ctx, userID)
}

func (s *PregnancyService) GetHorse(ctx context.Context, horseID uint) (*models.Horse, error) {
	return s.horseRepo.GetByID(ctx, horseID)
}

func (s *PregnancyService) UpdateHorse(ctx context.Context, horse *models.Horse) error {
	return s.horseRepo.Update(ctx, horse)
}

func (s *PregnancyService) UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	return s.pregnancyRepo.Update(ctx, pregnancy)
}

func (s *PregnancyService) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return s.pregnancyRepo.AddPreFoalingSign(ctx, sign)
}
