package repository

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type WeatherRepository struct {
	db *gorm.DB
}

func NewWeatherRepository(db *gorm.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

func (r *WeatherRepository) SaveWeatherData(ctx context.Context, data *models.WeatherData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *WeatherRepository) SaveWeatherImpact(ctx context.Context, impact *models.WeatherImpact) error {
	return r.db.WithContext(ctx).Create(impact).Error
}

func (r *WeatherRepository) GetLatestWeatherData(ctx context.Context, locationID uint) (*models.WeatherData, error) {
	var data models.WeatherData
	err := r.db.WithContext(ctx).
		Where("location_id = ?", locationID).
		Order("timestamp DESC").
		First(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *WeatherRepository) GetWeatherHistory(ctx context.Context, locationID uint, start, end time.Time) ([]models.WeatherData, error) {
	var data []models.WeatherData
	err := r.db.WithContext(ctx).
		Where("location_id = ? AND timestamp BETWEEN ? AND ?", locationID, start, end).
		Order("timestamp ASC").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *WeatherRepository) GetWeatherImpacts(ctx context.Context, horseID uint, start, end time.Time) ([]models.WeatherImpact, error) {
	var impacts []models.WeatherImpact
	err := r.db.WithContext(ctx).
		Preload("WeatherData").
		Where("horse_id = ? AND created_at BETWEEN ? AND ?", horseID, start, end).
		Order("created_at DESC").
		Find(&impacts).Error
	if err != nil {
		return nil, err
	}
	return impacts, nil
}
