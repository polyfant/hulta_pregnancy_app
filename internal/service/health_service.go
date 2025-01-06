package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HealthService struct {
    repo repository.HealthRepository
}

func NewHealthService(repo repository.HealthRepository) *HealthService {
    return &HealthService{repo: repo}
}

func (s *HealthService) GetHealthRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
    return s.repo.GetByHorseID(ctx, horseID)
}

func (s *HealthService) CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
    return s.repo.Create(ctx, record)
}

func (s *HealthService) UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
    return s.repo.Update(ctx, record)
}

func (s *HealthService) DeleteHealthRecord(ctx context.Context, id uint) error {
    return s.repo.Delete(ctx, id)
} 