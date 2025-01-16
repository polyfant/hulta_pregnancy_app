package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HealthService defines the interface for health-related operations
type HealthService interface {
	CreateRecord(ctx context.Context, record *models.HealthRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
	UpdateRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteRecord(ctx context.Context, id uint) error
}
