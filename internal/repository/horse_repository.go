package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type GormHorseRepository struct {
	db *gorm.DB
}

func NewHorseRepository(db *gorm.DB) *GormHorseRepository {
	return &GormHorseRepository{db: db}
}

func (r *GormHorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	return r.db.WithContext(ctx).Create(horse).Error
}

func (r *GormHorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	var horse models.Horse
	result := r.db.WithContext(ctx).First(&horse, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrHorseNotFound
		}
		return nil, result.Error
	}
	return &horse, nil
}

func (r *GormHorseRepository) Update(ctx context.Context, horse *models.Horse) error {
	horse.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(horse).Error
}

func (r *GormHorseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Horse{}, id).Error
}

func (r *GormHorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&horses)
	return horses, result.Error
}

func (r *GormHorseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND is_pregnant = ?", userID, true).
		Find(&horses)
	return horses, result.Error
}

// Custom errors
var (
	ErrHorseNotFound = errors.New("horse not found")
)
