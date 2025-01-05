package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type PregnancyService struct {
	horseRepo repository.HorseRepository
}

func NewPregnancyService(horseRepo repository.HorseRepository) *PregnancyService {
	return &PregnancyService{
		horseRepo: horseRepo,
	}
}

func (s *PregnancyService) StartPregnancy(ctx context.Context, horseID uint, userID string, conceptionDate time.Time) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return err
	}

	// Validate user ownership
	if horse.UserID != userID {
		return fmt.Errorf("unauthorized horse modification")
	}

	// Validate conception date
	if conceptionDate.After(time.Now()) {
		return fmt.Errorf("conception date cannot be in the future")
	}

	// Update horse pregnancy status
	horse.IsPregnant = true
	horse.ConceptionDate = conceptionDate
	horse.UpdatedAt = time.Now()

	return s.horseRepo.Update(ctx, horse)
}

func (s *PregnancyService) EndPregnancy(ctx context.Context, horseID uint, userID string) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return err
	}

	// Validate user ownership
	if horse.UserID != userID {
		return fmt.Errorf("unauthorized horse modification")
	}

	// Update horse pregnancy status
	horse.IsPregnant = false
	horse.ConceptionDate = time.Time{}
	horse.UpdatedAt = time.Now()

	return s.horseRepo.Update(ctx, horse)
}

func (s *PregnancyService) CalculateExpectedFoalingDate(conceptionDate time.Time) time.Time {
	// Average horse pregnancy is approximately 340 days
	return conceptionDate.AddDate(0, 0, 340)
}

func (s *PregnancyService) GetPregnancyStatus(ctx context.Context, horseID uint, userID string) (*models.PregnancyStatus, error) {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return nil, err
	}

	// Validate user ownership
	if horse.UserID != userID {
		return nil, fmt.Errorf("unauthorized horse access")
	}

	if !horse.IsPregnant {
		return &models.PregnancyStatus{
			IsPregnant: false,
		}, nil
	}

	expectedFoalingDate := s.CalculateExpectedFoalingDate(horse.ConceptionDate)
	daysPregnant := int(time.Since(horse.ConceptionDate).Hours() / 24)
	totalPregnancyDays := 340

	return &models.PregnancyStatus{
		IsPregnant:           true,
		ConceptionDate:       horse.ConceptionDate,
		ExpectedFoalingDate:  expectedFoalingDate,
		DaysPregnant:         daysPregnant,
		PregnancyPercentage:  calculatePregnancyPercentage(daysPregnant, totalPregnancyDays),
	}, nil
}

func calculatePregnancyPercentage(daysPregnant, totalPregnancyDays int) float64 {
	percentage := (float64(daysPregnant) / float64(totalPregnancyDays)) * 100
	if percentage > 100 {
		return 100
	}
	return percentage
}
