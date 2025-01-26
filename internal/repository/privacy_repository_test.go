package repository

import (
	"context"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db := testutils.SetupTestDB(t)

	// Clean up tables
	tables := []interface{}{
		&models.PrivacyPreferences{},
		&models.WeatherData{},
		&models.HealthRecord{},
		&models.PregnancyEvent{},
		&models.PrivacyChangeLog{},
	}

	for _, table := range tables {
		err := db.Migrator().DropTable(table)
		require.NoError(t, err)
		err = db.AutoMigrate(table)
		require.NoError(t, err)
	}

	return db
}

func TestPrivacyRepository_GetPrivacyPreferences(t *testing.T) {
	// Get test database connection
	db := setupTestDB(t)
	repo := NewPrivacyRepository(db)
	ctx := context.Background()

	t.Run("get non-existent preferences returns default", func(t *testing.T) {
		prefs, err := repo.GetPrivacyPreferences(ctx, "nonexistent")
		require.NoError(t, err)
		require.NotNil(t, prefs)
		require.False(t, prefs.WeatherTrackingEnabled)
		require.False(t, prefs.ShareHealthData)
		require.Equal(t, 0, prefs.DataRetentionDays)
	})

	t.Run("get existing preferences", func(t *testing.T) {
		// Create test preferences
		testPrefs := &models.PrivacyPreferences{
			UserID:                "test123",
			WeatherTrackingEnabled: true,
			ShareHealthData:       true,
			DataRetentionDays:     90,
		}
		err := repo.UpdatePrivacyPreferences(ctx, testPrefs)
		require.NoError(t, err)

		// Get preferences
		prefs, err := repo.GetPrivacyPreferences(ctx, "test123")
		require.NoError(t, err)
		require.NotNil(t, prefs)
		require.True(t, prefs.WeatherTrackingEnabled)
		require.True(t, prefs.ShareHealthData)
		require.Equal(t, 90, prefs.DataRetentionDays)
	})
}

func TestPrivacyRepository_UpdatePrivacyPreferences(t *testing.T) {
	// Get test database connection
	db := setupTestDB(t)
	repo := NewPrivacyRepository(db)
	ctx := context.Background()

	t.Run("create new preferences", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			UserID:                "new_user",
			WeatherTrackingEnabled: true,
			ShareHealthData:       false,
			DataRetentionDays:     90,
		}

		err := repo.UpdatePrivacyPreferences(ctx, prefs)
		require.NoError(t, err)

		// Verify preferences were created
		saved, err := repo.GetPrivacyPreferences(ctx, "new_user")
		require.NoError(t, err)
		require.True(t, saved.WeatherTrackingEnabled)
		require.False(t, saved.ShareHealthData)
		require.Equal(t, 90, saved.DataRetentionDays)
	})

	t.Run("update existing preferences", func(t *testing.T) {
		// Create initial preferences
		initial := &models.PrivacyPreferences{
			UserID:                "update_test",
			WeatherTrackingEnabled: true,
			ShareHealthData:       true,
			DataRetentionDays:     90,
		}
		err := repo.UpdatePrivacyPreferences(ctx, initial)
		require.NoError(t, err)

		// Update preferences
		updated := &models.PrivacyPreferences{
			UserID:                "update_test",
			WeatherTrackingEnabled: false,
			ShareHealthData:       false,
			DataRetentionDays:     30,
		}
		err = repo.UpdatePrivacyPreferences(ctx, updated)
		require.NoError(t, err)

		// Verify preferences were updated
		saved, err := repo.GetPrivacyPreferences(ctx, "update_test")
		require.NoError(t, err)
		require.False(t, saved.WeatherTrackingEnabled)
		require.False(t, saved.ShareHealthData)
		require.Equal(t, 30, saved.DataRetentionDays)
	})
}

func TestPrivacyRepository_DeleteUserData(t *testing.T) {
	// Get test database connection
	db := setupTestDB(t)
	repo := NewPrivacyRepository(db)
	ctx := context.Background()

	t.Run("delete old weather data", func(t *testing.T) {
		// Create test data
		now := time.Now()
		oldData := &models.WeatherData{
			UserID:    "test123",
			Timestamp: now.AddDate(0, 0, -31), // 31 days old
		}
		newData := &models.WeatherData{
			UserID:    "test123",
			Timestamp: now.AddDate(0, 0, -29), // 29 days old
		}
		require.NoError(t, db.Create(oldData).Error)
		require.NoError(t, db.Create(newData).Error)

		// Delete old data
		err := repo.DeleteUserData(ctx, "test123", "weather")
		require.NoError(t, err)

		// Verify only old data was deleted
		var count int64
		db.Model(&models.WeatherData{}).Count(&count)
		require.Equal(t, int64(1), count)
	})

	t.Run("delete old health data", func(t *testing.T) {
		// Create test data
		healthData := &models.HealthRecord{
			UserID: "test123",
		}
		require.NoError(t, db.Create(healthData).Error)

		// Delete data
		err := repo.DeleteUserData(ctx, "test123", "health")
		require.NoError(t, err)

		// Verify data was deleted
		var count int64
		db.Model(&models.HealthRecord{}).Count(&count)
		require.Equal(t, int64(0), count)
	})
}

func TestPrivacyRepository_AuditLog(t *testing.T) {
	// Get test database connection
	db := setupTestDB(t)
	repo := NewPrivacyRepository(db)
	ctx := context.Background()

	t.Run("verify audit log creation", func(t *testing.T) {
		// Update preferences to trigger audit log
		prefs := &models.PrivacyPreferences{
			UserID:                "audit_test",
			WeatherTrackingEnabled: true,
			ShareHealthData:       true,
		}
		err := repo.UpdatePrivacyPreferences(ctx, prefs)
		require.NoError(t, err)

		// Verify audit log entries
		var logs []models.PrivacyChangeLog
		err = db.Where("user_id = ?", "audit_test").Find(&logs).Error
		require.NoError(t, err)
		require.Len(t, logs, 2) // Should have 2 entries (weather and health)
	})
}
