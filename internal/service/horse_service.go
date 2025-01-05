package service

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HorseService struct {
	horseRepo     repository.HorseRepository
	expenseRepo   repository.ExpenseRepository
}

func NewHorseService(
	horseRepo repository.HorseRepository, 
	expenseRepo repository.ExpenseRepository,
) *HorseService {
	return &HorseService{
		horseRepo:   horseRepo,
		expenseRepo: expenseRepo,
	}
}

func (s *HorseService) CreateHorse(ctx context.Context, horse *models.Horse) error {
	// Validate horse data
	if err := s.validateHorseCreation(horse); err != nil {
		return err
	}

	// Set creation timestamps
	horse.CreatedAt = time.Now()
	horse.UpdatedAt = time.Now()

	return s.horseRepo.Create(ctx, horse)
}

func (s *HorseService) GetHorseWithDetails(ctx context.Context, horseID uint, userID string) (*models.HorseDetails, error) {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return nil, err
	}

	// Verify user ownership
	if horse.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to horse")
	}

	// Fetch related expenses
	expenses, _ := s.expenseRepo.GetByHorseID(ctx, horseID)

	return &models.HorseDetails{
		Horse:     *horse,
		Expenses:  expenses,
	}, nil
}

func (s *HorseService) UpdateHorsePregnancyStatus(ctx context.Context, horseID uint, userID string, isPregnant bool) error {
	horse, err := s.horseRepo.GetByID(ctx, horseID)
	if err != nil {
		return err
	}

	// Verify user ownership
	if horse.UserID != userID {
		return fmt.Errorf("unauthorized horse modification")
	}

	horse.IsPregnant = isPregnant
	horse.UpdatedAt = time.Now()

	return s.horseRepo.Update(ctx, horse)
}

func (s *HorseService) ListUserHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.horseRepo.ListByUser(ctx, userID)
}

func (s *HorseService) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	return s.horseRepo.GetPregnantHorses(ctx, userID)
}

func (s *HorseService) validateHorseCreation(horse *models.Horse) error {
	if horse.Name == "" {
		return fmt.Errorf("horse name is required")
	}

	if horse.BirthDate.IsZero() {
		return fmt.Errorf("birth date is required")
	}

	// Additional validation logic
	if horse.BirthDate.After(time.Now()) {
		return fmt.Errorf("birth date cannot be in the future")
	}

	return nil
}
