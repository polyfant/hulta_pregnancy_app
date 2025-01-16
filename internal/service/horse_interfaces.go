package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HorseService defines the interface for horse-related operations
type HorseService interface {
	GetByID(ctx context.Context, id uint) (*models.Horse, error)
	Create(ctx context.Context, horse *models.Horse) error
	Update(ctx context.Context, horse *models.Horse) error
	Delete(ctx context.Context, id uint) error
	GetPregnant(ctx context.Context, userID string) ([]models.Horse, error)
	GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error)
	GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error)
	GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error)
	ListByUserID(ctx context.Context, userID string) ([]models.Horse, error)
}
