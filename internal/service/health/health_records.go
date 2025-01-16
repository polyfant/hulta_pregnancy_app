package health

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Service defines the interface for health record operations
type Service interface {
	CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error
	GetHealthRecord(ctx context.Context, id uint) (*models.HealthRecord, error)
	ListHealthRecords(ctx context.Context, horseID uint) ([]*models.HealthRecord, error)
	UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteHealthRecord(ctx context.Context, id uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new health service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Repository defines the interface for health record data operations
type Repository interface {
	CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error
	GetHealthRecord(ctx context.Context, id uint) (*models.HealthRecord, error)
	ListHealthRecords(ctx context.Context, horseID uint) ([]*models.HealthRecord, error)
	UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteHealthRecord(ctx context.Context, id uint) error
}

func (s *service) CreateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.repo.CreateHealthRecord(ctx, record)
}

func (s *service) GetHealthRecord(ctx context.Context, id uint) (*models.HealthRecord, error) {
	return s.repo.GetHealthRecord(ctx, id)
}

func (s *service) ListHealthRecords(ctx context.Context, horseID uint) ([]*models.HealthRecord, error) {
	return s.repo.ListHealthRecords(ctx, horseID)
}

func (s *service) UpdateHealthRecord(ctx context.Context, record *models.HealthRecord) error {
	return s.repo.UpdateHealthRecord(ctx, record)
}

func (s *service) DeleteHealthRecord(ctx context.Context, id uint) error {
	return s.repo.DeleteHealthRecord(ctx, id)
}
