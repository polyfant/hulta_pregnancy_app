package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PrivacyRepository struct {
	db *gorm.DB
}

func NewPrivacyRepository(db *gorm.DB) *PrivacyRepository {
	return &PrivacyRepository{db: db}
}

func (r *PrivacyRepository) GetPrivacyPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error) {
	var prefs models.PrivacyPreferences
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&prefs).Error
	if err == gorm.ErrRecordNotFound {
		// Return default preferences if none exist
		return &models.PrivacyPreferences{
			UserID:    userID,
			UpdatedAt: time.Now(),
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func (r *PrivacyRepository) UpdatePrivacyPreferences(ctx context.Context, prefs *models.PrivacyPreferences) error {
	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get current preferences
	var current models.PrivacyPreferences
	err := tx.Where("user_id = ?", prefs.UserID).First(&current).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return err
	}

	// Update preferences
	err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(prefs).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Only log changes if we found existing preferences
	if err != gorm.ErrRecordNotFound {
		// Log weather tracking changes
		if current.WeatherTrackingEnabled != prefs.WeatherTrackingEnabled {
			log := &models.PrivacyChangeLog{
				UserID:    prefs.UserID,
				Field:     "weather_tracking_enabled",
				OldValue:  current.WeatherTrackingEnabled,
				NewValue:  prefs.WeatherTrackingEnabled,
				Timestamp: time.Now(),
			}
			if err := tx.Create(log).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		// Log health data sharing changes
		if current.ShareHealthData != prefs.ShareHealthData {
			log := &models.PrivacyChangeLog{
				UserID:    prefs.UserID,
				Field:     "share_health_data",
				OldValue:  current.ShareHealthData,
				NewValue:  prefs.ShareHealthData,
				Timestamp: time.Now(),
			}
			if err := tx.Create(log).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

func (r *PrivacyRepository) DeleteUserData(ctx context.Context, userID string, dataType string) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	switch dataType {
	case "weather":
		// Delete weather data older than 30 days
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
		if err := tx.Where("user_id = ? AND timestamp < ?", userID, thirtyDaysAgo).Delete(&models.WeatherData{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	case "health":
		// Delete all health records for the user
		if err := tx.Where("user_id = ?", userID).Delete(&models.HealthRecord{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	default:
		tx.Rollback()
		return fmt.Errorf("unsupported data type: %s", dataType)
	}

	return tx.Commit().Error
}
