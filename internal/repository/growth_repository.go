package repository

import (
	"context"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type GrowthRepository interface {
	CreateGrowthData(ctx context.Context, growthData *models.GrowthData) (*models.GrowthData, error)
	GetGrowthDataByFoalID(ctx context.Context, foalID uint) ([]models.GrowthData, error)
	UpdateGrowthData(ctx context.Context, growthData *models.GrowthData) error
	DeleteGrowthData(ctx context.Context, id uint) error
	CreateBodyCondition(ctx context.Context, data *models.BodyCondition) (*models.BodyCondition, error)
}

type PostgresGrowthRepository struct {
	db *gorm.DB
}

func NewGrowthRepository(db *gorm.DB) *PostgresGrowthRepository {
	return &PostgresGrowthRepository{db: db}
}

func (r *PostgresGrowthRepository) CreateGrowthData(ctx context.Context, growthData *models.GrowthData) (*models.GrowthData, error) {
	result := r.db.WithContext(ctx).Create(growthData)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create growth data: %w", result.Error)
	}
	return growthData, nil
}

func (r *PostgresGrowthRepository) GetGrowthDataByFoalID(ctx context.Context, foalID uint) ([]models.GrowthData, error) {
	var growthData []models.GrowthData
	result := r.db.WithContext(ctx).Where("foal_id = ?", foalID).Order("age ASC").Find(&growthData)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve growth data: %w", result.Error)
	}
	return growthData, nil
}

func (r *PostgresGrowthRepository) UpdateGrowthData(ctx context.Context, growthData *models.GrowthData) error {
	result := r.db.WithContext(ctx).Save(growthData)
	if result.Error != nil {
		return fmt.Errorf("failed to update growth data: %w", result.Error)
	}
	return nil
}

func (r *PostgresGrowthRepository) DeleteGrowthData(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.GrowthData{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete growth data: %w", result.Error)
	}
	return nil
}

func (r *PostgresGrowthRepository) CreateBodyCondition(ctx context.Context, data *models.BodyCondition) (*models.BodyCondition, error) {
	result := r.db.WithContext(ctx).Create(data)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create body condition record: %w", result.Error)
	}
	return data, nil
}
