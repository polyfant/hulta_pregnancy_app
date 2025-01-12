package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
)

// PregnancyService handles pregnancy-related business logic
type PregnancyService struct {
	horseRepo     repository.HorseRepository
	pregnancyRepo repository.PregnancyRepository
	calculator    *pregnancy.Calculator
}

// NewPregnancyService creates a new pregnancy service instance
func NewPregnancyService(horseRepo repository.HorseRepository, pregnancyRepo repository.PregnancyRepository) PregnancyService {
	return PregnancyService{
		horseRepo:     horseRepo,
		pregnancyRepo: pregnancyRepo,
		calculator:    pregnancy.NewCalculator(),
	}
}

// GetPregnancy retrieves pregnancy information for a horse
func (s PregnancyService) GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	return s.pregnancyRepo.GetByHorseID(ctx, horseID)
}

// GetPregnancies retrieves all pregnancies for a user
func (s PregnancyService) GetPregnancies(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	pregnancies, err := s.pregnancyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancies: %w", err)
	}
	return pregnancies, nil
}

// StartTracking begins tracking a new pregnancy
func (s PregnancyService) StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error {
	pregnancy := &models.Pregnancy{
		HorseID:        horseID,
		StartDate:      start.ConceptionDate,
		Status:         models.PregnancyStatusActive,
		ConceptionDate: &start.ConceptionDate,
	}

	if err := s.pregnancyRepo.Create(ctx, pregnancy); err != nil {
		return fmt.Errorf("failed to create pregnancy: %w", err)
	}

	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	horse.IsPregnant = true
	if err := s.horseRepo.Update(ctx, horse); err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}

	return nil
}

// GetStatus retrieves the current pregnancy status
func (s PregnancyService) GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error) {
	pregnancy, err := s.pregnancyRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy: %w", err)
	}

	stage, err := s.GetPregnancyStage(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy stage: %w", err)
	}

	return &models.PregnancyStatus{
		Stage:        stage,
		DaysPregnant: int(time.Since(*pregnancy.ConceptionDate).Hours() / 24),
	}, nil
}

// GetPregnancyEvents retrieves all events for a pregnancy
func (s PregnancyService) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	return s.pregnancyRepo.GetEvents(ctx, horseID)
}

// GetGuidelines retrieves guidelines for a specific pregnancy stage
func (s PregnancyService) GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]models.Guideline, error) {
	guidelines := map[models.PregnancyStage][]models.Guideline{
		models.PregnancyStageEarly: {
			{
				Stage:       models.PregnancyStageEarly,
				Description: "Early pregnancy care",
				Tips: []string{
					"Schedule initial vet check",
					"Maintain regular exercise routine",
					"Monitor mare's appetite",
				},
			},
		},
		models.PregnancyStageMid: {
			{
				Stage:       models.PregnancyStageMid,
				Description: "Mid-pregnancy care",
				Tips: []string{
					"Increase feed gradually",
					"Continue moderate exercise",
					"Schedule vaccination updates",
				},
			},
		},
		models.PregnancyStageLate: {
			{
				Stage:       models.PregnancyStageLate,
				Description: "Late pregnancy care",
				Tips: []string{
					"Prepare foaling area",
					"Monitor for signs of labor",
					"Have vet contact ready",
				},
			},
		},
	}

	if g, ok := guidelines[stage]; ok {
		return g, nil
	}
	return nil, fmt.Errorf("invalid pregnancy stage: %s", stage)
}

func (s PregnancyService) GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	return s.pregnancyRepo.GetActive(ctx, userID)
}

func (s PregnancyService) UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	if err := s.pregnancyRepo.Update(ctx, pregnancy); err != nil {
		return fmt.Errorf("failed to update pregnancy: %w", err)
	}
	return nil
}

func (s PregnancyService) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	return s.pregnancyRepo.GetPreFoaling(ctx, horseID)
}

func (s PregnancyService) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return s.pregnancyRepo.AddPreFoaling(ctx, sign)
}

func (s PregnancyService) GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error) {
	pregnancy, err := s.pregnancyRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		return "", fmt.Errorf("failed to get pregnancy: %w", err)
	}

	if pregnancy.ConceptionDate == nil {
		return "", fmt.Errorf("pregnancy has no conception date")
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)

	switch {
	case daysPregnant < 120:
		return models.PregnancyStageEarly, nil
	case daysPregnant < 240:
		return models.PregnancyStageMid, nil
	case daysPregnant < 340:
		return models.PregnancyStageLate, nil
	default:
		return models.PregnancyStageOverdue, nil
	}
}

func (s PregnancyService) EndPregnancy(ctx context.Context, horseID uint, status string, date time.Time) error {
	pregnancy, err := s.pregnancyRepo.GetByHorseID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get pregnancy: %w", err)
	}

	pregnancy.Status = status
	pregnancy.EndDate = &date

	if err := s.pregnancyRepo.Update(ctx, pregnancy); err != nil {
		return fmt.Errorf("failed to update pregnancy: %w", err)
	}

	return nil
}

func (s PregnancyService) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	if err := s.pregnancyRepo.AddPregnancyEvent(ctx, event); err != nil {
		return fmt.Errorf("failed to add pregnancy event: %w", err)
	}
	return nil
}

func (s PregnancyService) GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	events, err := s.pregnancyRepo.GetEvents(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancy events: %w", err)
	}
	return events, nil
}

func (s PregnancyService) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	items, err := s.pregnancyRepo.GetPreFoalingChecklist(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pre-foaling checklist: %w", err)
	}
	return items, nil
}

func (s PregnancyService) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	if err := s.pregnancyRepo.AddPreFoalingChecklistItem(ctx, item); err != nil {
		return fmt.Errorf("failed to add pre-foaling checklist item: %w", err)
	}
	return nil
}
