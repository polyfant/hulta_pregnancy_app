package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type HealthRepository struct {
	db *gorm.DB
}

func NewHealthRepository(db *gorm.DB) *HealthRepository {
	return &HealthRepository{db: db}
}

func (r *HealthRepository) Create(ctx context.Context, record *models.HealthRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *HealthRepository) GetByHorseID(ctx context.Context, horseID uint) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&records).Error
	return records, err
}

func (r *HealthRepository) Update(ctx context.Context, record *models.HealthRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *HealthRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.HealthRecord{}, id).Error
}

func (r *HealthRepository) GetByID(ctx context.Context, id uint) (*models.HealthRecord, error) {
	var record models.HealthRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}