package breeding

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type BreedingService struct {
	breedingRepo repository.BreedingRepository
	horseRepo    repository.HorseRepository
}

func NewBreedingService(breedingRepo repository.BreedingRepository, horseRepo repository.HorseRepository) *BreedingService {
	return &BreedingService{
		breedingRepo: breedingRepo,
		horseRepo:    horseRepo,
	}
}

func (s *BreedingService) CalculatePregnancyStage(pregnancy *models.Pregnancy) models.PregnancyStage {
	if pregnancy.ConceptionDate == nil {
		return models.PregnancyStageEarlyGestation
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

func (s *BreedingService) GetBreedingRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	return s.breedingRepo.GetRecords(ctx, horseID)
}

func (s *BreedingService) AddBreedingRecord(ctx context.Context, record *models.BreedingRecord) error {
	// Validate the record
	if record.Date.After(time.Now()) {
		return errors.New("breeding date cannot be in the future")
	}

	// Add the record
	return s.breedingRepo.CreateRecord(ctx, record)
}

func (s *BreedingService) GetBreedingCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error) {
	return s.breedingRepo.GetCosts(ctx, horseID)
}

func (s *BreedingService) AddBreedingCost(ctx context.Context, cost *models.BreedingCost) error {
	return s.breedingRepo.Create(ctx, cost)
}

func (s *BreedingService) GetBreedingStatus(ctx context.Context, horseID uint) (*models.BreedingStatus, error) {
	records, err := s.breedingRepo.GetRecords(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding records: %w", err)
	}

	status := models.BreedingStatusActive
	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		switch lastRecord.Status {
		case string(models.BreedingStatusCompleted):
			status = models.BreedingStatusCompleted
		case string(models.BreedingStatusFailed):
			status = models.BreedingStatusFailed
		case string(models.BreedingStatusCancelled):
			status = models.BreedingStatusCancelled
		}
	}

	return &status, nil
}

func (s *BreedingService) UpdateBreedingStatus(ctx context.Context, horseID uint, status string) error {
	record := &models.BreedingRecord{
		HorseID: horseID,
		Date:    time.Now(),
		Status:  "ACTIVE",
	}

	if err := s.breedingRepo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create breeding record: %w", err)
	}

	return nil
}

func (s *BreedingService) ValidateBreedingEligibility(ctx context.Context, horseID uint) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return fmt.Errorf("failed to get horse: %w", err)
	}

	if horse.IsPregnant {
		return fmt.Errorf("horse is already pregnant")
	}

	return nil
}

func (s *BreedingService) ValidateBreedingStatus(status string) error {
	switch status {
	case "ACTIVE", "COMPLETED", "CANCELLED":
		return nil
	default:
		return fmt.Errorf("invalid breeding status: %s", status)
	}
}
