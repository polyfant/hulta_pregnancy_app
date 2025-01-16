package repository

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

// VitalsRepository handles vital signs data persistence
type VitalsRepository struct {
	db *gorm.DB
}

// NewVitalsRepository creates a new vitals repository
func NewVitalsRepository(db *gorm.DB) *VitalsRepository {
	return &VitalsRepository{db: db}
}

// SaveVitalSigns saves vital signs data
func (r *VitalsRepository) SaveVitalSigns(ctx context.Context, vitals *models.VitalSigns) error {
	return r.db.WithContext(ctx).Create(vitals).Error
}

// GetVitalSigns retrieves vital signs for a date range
func (r *VitalsRepository) GetVitalSigns(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSigns, error) {
	var vitals []*models.VitalSigns
	result := r.db.WithContext(ctx).
		Where("horse_id = ? AND recorded_at BETWEEN ? AND ?", horseID, from, to).
		Order("recorded_at DESC").
		Find(&vitals)
	return vitals, result.Error
}

// GetLatestVitalSigns retrieves the most recent vital signs for a horse
func (r *VitalsRepository) GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error) {
	var vitals models.VitalSigns
	result := r.db.WithContext(ctx).
		Where("horse_id = ?", horseID).
		Order("recorded_at DESC").
		First(&vitals)
	if result.Error != nil {
		return nil, result.Error
	}
	return &vitals, nil
}

// SaveAlert saves a vital signs alert
func (r *VitalsRepository) SaveAlert(ctx context.Context, alert *models.VitalSignsAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetAlert retrieves a specific alert
func (r *VitalsRepository) GetAlert(ctx context.Context, alertID uint) (*models.VitalSignsAlert, error) {
	var alert models.VitalSignsAlert
	result := r.db.WithContext(ctx).First(&alert, alertID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &alert, nil
}

// GetAlerts retrieves alerts for a horse
func (r *VitalsRepository) GetAlerts(ctx context.Context, horseID uint, includeAcknowledged bool) ([]*models.VitalSignsAlert, error) {
	var alerts []*models.VitalSignsAlert
	query := r.db.WithContext(ctx).Where("horse_id = ?", horseID)
	if !includeAcknowledged {
		query = query.Where("acknowledged_at IS NULL")
	}
	result := query.Order("created_at DESC").Find(&alerts)
	return alerts, result.Error
}

// AcknowledgeAlert marks an alert as acknowledged
func (r *VitalsRepository) AcknowledgeAlert(ctx context.Context, alertID uint) error {
	return r.db.WithContext(ctx).Model(&models.VitalSignsAlert{}).
		Where("id = ?", alertID).
		Update("acknowledged_at", time.Now()).
		Error
}

// SaveTrend saves a vital signs trend
func (r *VitalsRepository) SaveTrend(ctx context.Context, trend *models.VitalSignsTrend) error {
	return r.db.WithContext(ctx).Create(trend).Error
}

// GetTrends retrieves vital signs trends for a date range
func (r *VitalsRepository) GetTrends(ctx context.Context, horseID uint, from, to time.Time) ([]*models.VitalSignsTrend, error) {
	var trends []*models.VitalSignsTrend
	result := r.db.WithContext(ctx).
		Where("horse_id = ? AND start_time >= ? AND end_time <= ?", horseID, from, to).
		Order("start_time DESC").
		Find(&trends)
	return trends, result.Error
}
