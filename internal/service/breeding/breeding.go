package breeding

import (
	"context"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

type BreedingService struct {
	repo repository.BreedingRepository
}

var _ service.BreedingService = (*BreedingService)(nil)

func NewBreedingService(repo repository.BreedingRepository) service.BreedingService {
	return &BreedingService{repo: repo}
}

func (s *BreedingService) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	if err := s.repo.CreateRecord(ctx, record); err != nil {
		return fmt.Errorf("failed to create breeding record: %w", err)
	}
	return nil
}

func (s *BreedingService) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	records, err := s.repo.GetRecords(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding records: %w", err)
	}
	return records, nil
}

func (s *BreedingService) UpdateRecord(ctx context.Context, record *models.BreedingRecord) error {
	return s.repo.UpdateRecord(ctx, record)
}

func (s *BreedingService) DeleteRecord(ctx context.Context, id uint) error {
	return s.repo.DeleteRecord(ctx, id)
}
