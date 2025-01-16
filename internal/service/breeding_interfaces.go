package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// BreedingService defines the interface for breeding-related operations
type BreedingService interface {
	CreateRecord(ctx context.Context, record *models.BreedingRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error)
	UpdateRecord(ctx context.Context, record *models.BreedingRecord) error
	DeleteRecord(ctx context.Context, id uint) error
}
