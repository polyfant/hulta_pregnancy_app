package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

// HealthServiceImpl implements the HealthService interface
type HealthServiceImpl struct {
	repo repository.HealthRepository
}

// NewHealthService creates a new health service instance
func NewHealthService(repo repository.HealthRepository) HealthService {
	return &HealthServiceImpl{
		repo: repo,
	}
}

// GetRecords retrieves health records for a specific horse
func (s *HealthServiceImpl) GetRecords(ctx context.Context, horseID uint) ([]*models.HealthRecord, error) {
	records, err := s.repo.GetRecords(ctx, horseID)
	if err != nil {
		return nil, err
	}
	
	// Convert []models.HealthRecord to []*models.HealthRecord
	healthRecords := make([]*models.HealthRecord, len(records))
	for i := range records {
		healthRecords[i] = &records[i]
	}
	
	return healthRecords, nil
}

// CreateRecord adds a new health record
func (s *HealthServiceImpl) CreateRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.repo.CreateRecord(ctx, record)
}

// UpdateRecord modifies an existing health record
func (s *HealthServiceImpl) UpdateRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.repo.UpdateRecord(ctx, record)
}

// DeleteRecord removes a health record
func (s *HealthServiceImpl) DeleteRecord(ctx context.Context, id uint) error {
	return s.repo.DeleteRecord(ctx, id)
}
