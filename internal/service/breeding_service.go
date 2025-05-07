package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

type BreedingServiceImpl struct {
	repo repository.BreedingRepository
	// other fields...
}

func NewBreedingService(repo repository.BreedingRepository /* other deps */) BreedingService {
	return &BreedingServiceImpl{repo: repo /* other init */}
}

// CreateRecord creates a new breeding record.
func (s *BreedingServiceImpl) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create breeding record: %w", err)
	}
	return nil
}

// GetRecords retrieves all breeding records for a specific horse.
func (s *BreedingServiceImpl) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	records, err := s.repo.GetRecords(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding records for horse ID %d: %w", horseID, err)
	}
	return records, nil
}

// UpdateRecord updates an existing breeding record.
func (s *BreedingServiceImpl) UpdateRecord(ctx context.Context, record *models.BreedingRecord) error {
	// It's good practice to ensure the record exists and the user has permission before updating.
	// This might involve a GetByID call first, or letting the repository handle it if it returns specific errors.
	if err := s.repo.UpdateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to update breeding record ID %d: %w", record.ID, err)
	}
	return nil
}

// DeleteRecord deletes a breeding record by its ID.
func (s *BreedingServiceImpl) DeleteRecord(ctx context.Context, id uint) error {
	// Similar to Update, check for existence/permissions if necessary.
	if err := s.repo.DeleteRecord(ctx, id); err != nil {
		return fmt.Errorf("failed to delete breeding record ID %d: %w", id, err)
	}
	return nil
}

// GetRecordByID retrieves a single breeding record by its ID.
func (s *BreedingServiceImpl) GetRecordByID(ctx context.Context, id uint) (*models.BreedingRecord, error) {
	record, err := s.repo.GetByID(ctx, id) // Method name GetByID is correct as per repository interface for BreedingRecord
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("breeding record with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get breeding record by ID %d: %w", id, err)
	}
	// TODO: Add authorization checks if necessary (e.g., can user associated with ctx view this record?)
	return record, nil
}

// ... (other existing methods of BreedingServiceImpl) ... 