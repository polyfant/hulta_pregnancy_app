package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
	"gorm.io/gorm"
)

// PregnancyServiceImpl handles pregnancy-related business logic
type PregnancyServiceImpl struct {
	horseRepo     repository.HorseRepository
	pregnancyRepo repository.PregnancyRepository
	calculator    *pregnancy.Calculator
}

// NewPregnancyService creates a new pregnancy service instance
func NewPregnancyService(horseRepo repository.HorseRepository, pregnancyRepo repository.PregnancyRepository) *PregnancyServiceImpl {
	return &PregnancyServiceImpl{
		horseRepo:     horseRepo,
		pregnancyRepo: pregnancyRepo,
		calculator:    pregnancy.NewCalculator(),
	}
}

// GetPregnancy retrieves pregnancy information for a horse
func (s *PregnancyServiceImpl) GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	return s.pregnancyRepo.GetByHorseID(ctx, horseID)
}

// GetPregnancies retrieves all pregnancies for a user
func (s *PregnancyServiceImpl) GetPregnancies(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	pregnancies, err := s.pregnancyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancies: %w", err)
	}
	return pregnancies, nil
}

// StartTracking begins tracking a new pregnancy
func (s *PregnancyServiceImpl) StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error {
	_, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err == nil {
		return fmt.Errorf("horse %d already has an active pregnancy", horseID)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing active pregnancy for horse %d: %w", horseID, err)
	}

	pregnancy := &models.Pregnancy{
		HorseID:        horseID,
		StartDate:      start.ConceptionDate,
		Status:         models.PregnancyStatusActive,
		ConceptionDate: &start.ConceptionDate,
	}

	if err := s.pregnancyRepo.Create(ctx, pregnancy); err != nil {
		return fmt.Errorf("failed to create pregnancy: %w", err)
	}

	if err := s.pregnancyRepo.UpdatePregnancyStatus(ctx, horseID, true, &start.ConceptionDate); err != nil {
		fmt.Printf("WARN: Failed to update horse %d pregnancy status after starting tracking: %v\n", horseID, err)
	}

	return nil
}

// GetStatus retrieves the current pregnancy status
func (s *PregnancyServiceImpl) GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error) {
	pregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no active pregnancy found for horse ID %d", horseID)
		}
		return nil, fmt.Errorf("failed to get active pregnancy for status: %w", err)
	}

	stage, err := s.GetPregnancyStage(ctx, horseID)
	if err != nil {
		fmt.Printf("WARN: Failed to get pregnancy stage for horse %d, returning partial status: %v\n", horseID, err)
		stage = "UNKNOWN"
	}

	if pregnancy.ConceptionDate == nil {
		return nil, fmt.Errorf("active pregnancy for horse ID %d lacks conception date", horseID)
	}

	return &models.PregnancyStatus{
		Stage:        stage,
		DaysPregnant: int(time.Since(*pregnancy.ConceptionDate).Hours() / 24),
		DueDate:      pregnancy.ExpectedDueDate(),
	}, nil
}

// GetPregnancyEvents retrieves all events for a pregnancy
func (s *PregnancyServiceImpl) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	return s.pregnancyRepo.GetEvents(ctx, horseID)
}

// GetGuidelines retrieves guidelines for a specific pregnancy stage
func (s *PregnancyServiceImpl) GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]models.Guideline, error) {
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
	return nil, fmt.Errorf("guidelines not found for pregnancy stage: %s", stage)
}

func (s *PregnancyServiceImpl) GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	allPregnancies, err := s.pregnancyRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pregnancies for user %s: %w", userID, err)
	}
	activePregnancies := []models.Pregnancy{}
	for _, p := range allPregnancies {
		if p.Status == models.PregnancyStatusActive {
			activePregnancies = append(activePregnancies, p)
		}
	}
	return activePregnancies, nil
}

func (s *PregnancyServiceImpl) UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	if err := s.pregnancyRepo.Update(ctx, pregnancy); err != nil {
		return fmt.Errorf("failed to update pregnancy: %w", err)
	}
	return nil
}

func (s *PregnancyServiceImpl) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	return s.pregnancyRepo.GetPreFoalingSigns(ctx, horseID)
}

func (s *PregnancyServiceImpl) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return s.pregnancyRepo.AddPreFoalingSign(ctx, sign)
}

func (s *PregnancyServiceImpl) GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error) {
	pregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "NOT_PREGNANT", nil
		}
		return "", fmt.Errorf("failed to get active pregnancy for stage calculation: %w", err)
	}

	if pregnancy.ConceptionDate == nil {
		return "", fmt.Errorf("active pregnancy for horse ID %d lacks conception date", horseID)
	}

	daysPregnant := int(time.Since(*pregnancy.ConceptionDate).Hours() / 24)

	oneThird := models.DefaultGestationDays / 3
	twoThirds := (models.DefaultGestationDays * 2) / 3
	overdueThreshold := models.DefaultGestationDays + 7

	switch {
	case daysPregnant < oneThird:
		return models.PregnancyStageEarly, nil
	case daysPregnant < twoThirds:
		return models.PregnancyStageMid, nil
	case daysPregnant <= overdueThreshold:
		return models.PregnancyStageLate, nil
	default:
		return models.PregnancyStageOverdue, nil
	}
}

func (s *PregnancyServiceImpl) EndPregnancy(ctx context.Context, horseID uint, status string, date time.Time) error {
	allowedStatuses := map[string]bool{
		models.PregnancyStatusComplete: true,
		models.PregnancyStatusLost:     true,
		models.PregnancyStatusAborted:  true,
	}
	if !allowedStatuses[status] {
		return fmt.Errorf("invalid end status: %s", status)
	}

	activePregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no active pregnancy found for horse ID %d to end", horseID)
		}
		return fmt.Errorf("failed to get active pregnancy to end: %w", err)
	}

	activePregnancy.Status = status
	activePregnancy.EndDate = &date

	if err := s.pregnancyRepo.Update(ctx, activePregnancy); err != nil {
		return fmt.Errorf("failed to update pregnancy status to %s: %w", status, err)
	}

	if err := s.pregnancyRepo.UpdatePregnancyStatus(ctx, horseID, false, nil); err != nil {
		fmt.Printf("WARN: Failed to update horse %d pregnancy status after ending tracking: %v\n", horseID, err)
	}

	return nil
}

// AddPregnancyEvent creates a new event for an active pregnancy of a horse.
func (s *PregnancyServiceImpl) AddPregnancyEvent(ctx context.Context, userID string, horseID uint, eventInput *models.PregnancyEventInputDTO) (*models.PregnancyEvent, error) {
	activePregnancy, err := s.pregnancyRepo.GetCurrentPregnancy(ctx, horseID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no active pregnancy found for horse ID %d to add event: %w", horseID, err)
		}
		return nil, fmt.Errorf("failed to retrieve active pregnancy for horse ID %d: %w", horseID, err)
	}

	newEvent := &models.PregnancyEvent{
		PregnancyID:  activePregnancy.ID,
		UserID:       userID,
		Type:         eventInput.Type,
		Description:  eventInput.Description,
		Date:         eventInput.Date,
	}

	err = s.pregnancyRepo.AddPregnancyEvent(ctx, newEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to create pregnancy event in repository: %w", err)
	}

	return newEvent, nil
}

func (s *PregnancyServiceImpl) GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	return s.pregnancyRepo.GetEvents(ctx, horseID)
}

func (s *PregnancyServiceImpl) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	return s.pregnancyRepo.GetPreFoalingChecklist(ctx, horseID)
}

func (s *PregnancyServiceImpl) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return s.pregnancyRepo.AddPreFoalingChecklistItem(ctx, item)
}

// GetPregnancyByID fetches a specific pregnancy by its ID.
func (s *PregnancyServiceImpl) GetPregnancyByID(ctx context.Context, pregnancyID uint) (*models.Pregnancy, error) {
	preg, err := s.pregnancyRepo.GetPregnancy(ctx, pregnancyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("pregnancy with ID %d not found: %w", pregnancyID, err)
		}
		return nil, fmt.Errorf("failed to get pregnancy by ID %d: %w", pregnancyID, err)
	}
	return preg, nil
}
