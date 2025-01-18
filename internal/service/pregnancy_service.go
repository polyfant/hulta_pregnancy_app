package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
	"github.com/polyfant/hulta_pregnancy_app/internal/validation"
)

// PregnancyServiceImpl handles pregnancy-related business logic
type PregnancyServiceImpl struct {
	horseRepo     repository.HorseRepository
	pregnancyRepo repository.PregnancyRepository
	calculator    *pregnancy.Calculator
	validator     *validation.Validator
	eventLogger   EventLogger
}

// NewPregnancyService creates a new pregnancy service instance
func NewPregnancyService(horseRepo repository.HorseRepository, pregnancyRepo repository.PregnancyRepository, eventLogger EventLogger) PregnancyService {
	return &PregnancyServiceImpl{
		horseRepo:     horseRepo,
		pregnancyRepo: pregnancyRepo,
		calculator:    pregnancy.NewCalculator(),
		validator:     validation.NewValidator(),
		eventLogger:   eventLogger,
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
func (s *PregnancyServiceImpl) StartTracking(ctx context.Context, horseID uint, conceptionDate time.Time) (*models.Pregnancy, error) {
	// 1. Input Validation
	if horseID == 0 {
		return nil, fmt.Errorf("invalid horse ID: cannot be zero")
	}

	// Validate conception date
	if conceptionDate.IsZero() {
		return nil, fmt.Errorf("conception date cannot be zero")
	}

	if conceptionDate.After(time.Now()) {
		return nil, fmt.Errorf("conception date cannot be in the future")
	}

	// 2. Retrieve Horse
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve horse: %w", err)
	}

	// 3. Validate Horse Pregnancy Eligibility
	if horse == nil {
		return nil, fmt.Errorf("horse with ID %d not found", horseID)
	}

	if horse.Gender != models.GenderMare {
		return nil, fmt.Errorf("only mares can be pregnant, current horse is a %s", horse.Gender)
	}

	// 4. Check for Existing Active Pregnancy
	existingPregnancy, err := s.pregnancyRepo.GetActivePregnancyByHorseID(ctx, horseID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("error checking existing pregnancies: %w", err)
	}

	if existingPregnancy != nil {
		return nil, fmt.Errorf("horse is already pregnant with an active pregnancy")
	}

	// 5. Create Pregnancy Record
	pregnancy := &models.Pregnancy{
		HorseID:         horseID,
		ConceptionDate:  conceptionDate,
		Status:          models.PregnancyStatusActive,
		EstimatedFoalingDate: calculateEstimatedFoalingDate(conceptionDate),
	}

	// 6. Validate Pregnancy Model
	if err := validation.ValidateStruct(pregnancy); err != nil {
		return nil, fmt.Errorf("invalid pregnancy data: %w", err)
	}

	// 7. Save Pregnancy
	savedPregnancy, err := s.pregnancyRepo.Create(ctx, pregnancy)
	if err != nil {
		return nil, fmt.Errorf("failed to save pregnancy record: %w", err)
	}

	// 8. Update Horse Pregnancy Status
	horse.IsPregnant = true
	horse.ConceptionDate = &conceptionDate
	if err := s.horseRepo.Update(ctx, horse); err != nil {
		// Log warning, but don't fail the entire operation
		log.Printf("Warning: Could not update horse pregnancy status: %v", err)
	}

	// 9. Log Tracking Event
	s.eventLogger.LogEvent(ctx, "pregnancy_tracking_started", map[string]interface{}{
		"horse_id":         horseID,
		"conception_date":  conceptionDate,
		"foaling_estimate": savedPregnancy.EstimatedFoalingDate,
	})

	return savedPregnancy, nil
}

// Helper function to calculate estimated foaling date
func calculateEstimatedFoalingDate(conceptionDate time.Time) time.Time {
	// Typical horse gestation is approximately 340 days
	return conceptionDate.AddDate(0, 0, 340)
}

// Rest of the code remains the same
