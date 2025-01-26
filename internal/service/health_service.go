package service

import (
	"context"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HealthServiceImpl struct {
	repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) HealthService {
	return &HealthServiceImpl{repo: repo}
}

func (s *HealthServiceImpl) CreateRecord(ctx context.Context, record *models.HealthRecord) error {
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create health record: %w", err)
	}
	return nil
}

func (s *HealthServiceImpl) GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	records, err := s.repo.GetRecords(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health records: %w", err)
	}
	return records, nil
}

func (s *HealthServiceImpl) UpdateRecord(ctx context.Context, record *models.HealthRecord) error {
	if err := s.repo.UpdateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to update health record: %w", err)
	}
	return nil
}

func (s *HealthServiceImpl) DeleteRecord(ctx context.Context, id uint) error {
	if err := s.repo.DeleteRecord(ctx, id); err != nil {
		return fmt.Errorf("failed to delete health record: %w", err)
	}
	return nil
}