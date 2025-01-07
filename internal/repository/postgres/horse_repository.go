package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

func NewHorseRepository(db *gorm.DB) repository.HorseRepository {
	return &horseRepository{db: db}
}

type horseRepository struct {
	db *gorm.DB
}

func (r *horseRepository) Create(ctx context.Context, horse *models.Horse) error {
	return r.db.Create(horse).Error
}

func (r *horseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	var horse models.Horse
	if err := r.db.First(&horse, id).Error; err != nil {
		return nil, err
	}
	return &horse, nil
}

func (r *horseRepository) Update(ctx context.Context, horse *models.Horse) error {
	return r.db.Save(horse).Error
}

func (r *horseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(&models.Horse{}, id).Error
}

func (r *horseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	if err := r.db.Where("user_id = ?", userID).Find(&horses).Error; err != nil {
		return nil, err
	}
	return horses, nil
}

func (r *horseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	if err := r.db.Where("user_id = ? AND is_pregnant = ?", userID, true).Find(&horses).Error; err != nil {
		return nil, err
	}
	return horses, nil
} 