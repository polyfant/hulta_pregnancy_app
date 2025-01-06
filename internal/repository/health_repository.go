package repository

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type HealthRepository interface {
    GetByHorseID(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
    Create(ctx context.Context, record *models.HealthRecord) error
    Update(ctx context.Context, record *models.HealthRecord) error
    Delete(ctx context.Context, id uint) error
} 