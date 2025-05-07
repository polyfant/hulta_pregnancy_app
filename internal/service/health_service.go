package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

type HealthServiceImpl struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) HealthService {
	return &HealthServiceImpl{repo: repo}
}

func (s *HealthServiceImpl) CreateRecord(ctx context.Context, record *models.HealthRecord) error {
	if err := s.repo.Create(ctx, record); err != nil {
		return fmt.Errorf("failed to create health record: %w", err)
	}
	return nil
}

func (s *HealthServiceImpl) GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	records, err := s.repo.GetByHorseID(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health records: %w", err)
	}
	return records, nil
}

func (s *HealthServiceImpl) UpdateRecord(ctx context.Context, record *models.HealthRecord) error {
	if err := s.repo.Update(ctx, record); err != nil {
		return fmt.Errorf("failed to update health record: %w", err)
	}
	return nil
}

func (s *HealthServiceImpl) DeleteRecord(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete health record: %w", err)
	}
	return nil
}

func (s *HealthServiceImpl) GetRecordByID(ctx context.Context, id uint) (*models.HealthRecord, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("health record with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get health record by ID %d: %w", id, err)
	}
	return record, nil
}