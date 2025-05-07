package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type BreedingRepository struct {
	db *gorm.DB
}

func NewBreedingRepository(db *gorm.DB) *BreedingRepository {
	return &BreedingRepository{db: db}
}

func (r *BreedingRepository) GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error) {
	var costs []models.BreedingCost
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&costs).Error
	return costs, err
}

func (r *BreedingRepository) Create(ctx context.Context, cost *models.BreedingCost) error {
	return r.db.WithContext(ctx).Create(cost).Error
}

func (r *BreedingRepository) GetByID(ctx context.Context, id uint) (*models.BreedingRecord, error) {
	var record models.BreedingRecord
	err := r.db.WithContext(ctx).First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *BreedingRepository) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	var records []models.BreedingRecord
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&records).Error
	return records, err
}

func (r *BreedingRepository) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *BreedingRepository) UpdateRecord(ctx context.Context, record *models.BreedingRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *BreedingRepository) DeleteRecord(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.BreedingRecord{}, id).Error
} 