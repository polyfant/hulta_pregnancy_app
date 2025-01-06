package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

type healthRepository struct {
	db *gorm.DB
}

func NewHealthRepository(db *gorm.DB) repository.HealthRepository {
	return &healthRepository{db: db}
}

func (r *healthRepository) GetByHorseID(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	err := r.db.Where("horse_id = ?", horseID).Find(&records).Error
	return records, err
}

func (r *healthRepository) Create(ctx context.Context, record *models.HealthRecord) error {
	return r.db.Create(record).Error
}

func (r *healthRepository) Update(ctx context.Context, record *models.HealthRecord) error {
	return r.db.Save(record).Error
}

func (r *healthRepository) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(&models.HealthRecord{}, id).Error
} 