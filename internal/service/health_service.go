package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type HealthService struct {
	healthRepo repository.HealthRepository
}

func NewHealthService(healthRepo repository.HealthRepository) *HealthService {
	return &HealthService{
		healthRepo: healthRepo,
	}
}

func (s *HealthService) GetHealthRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	return s.healthRepo.GetByHorseID(ctx, horseID)
}

func (s *HealthService) CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.healthRepo.Create(ctx, record)
}

func (s *HealthService) UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.healthRepo.Update(ctx, record)
}

func (s *HealthService) DeleteHealthRecord(ctx context.Context, id uint) error {
	return s.healthRepo.Delete(ctx, id)
}

func (s *HealthService) AddHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.healthRepo.Create(ctx, record)
} 