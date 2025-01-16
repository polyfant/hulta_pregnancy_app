package repository

import (
	"context"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVitalsRepository(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	repo := NewVitalsRepository(db)
	ctx := context.Background()

	t.Run("SaveVitalSigns", func(t *testing.T) {
		testCases := []struct {
			name        string
			vitals      *models.VitalSigns
			expectError bool
		}{
			{
				name:        "Valid Vitals",
				vitals:      testutil.CreateTestVitalSigns(1),
				expectError: false,
			},
			{
				name: "Invalid Temperature Range",
				vitals: &models.VitalSigns{
					HorseID:     1,
					Temperature: 45.0, // Impossible temperature
					RecordedAt:  testutil.MockTimeNow(),
				},
				expectError: true,
			},
			{
				name: "Missing Required Fields",
				vitals: &models.VitalSigns{
					HorseID: 1,
					// Missing temperature and other fields
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := repo.SaveVitalSigns(ctx, tc.vitals)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("GetVitalSigns", func(t *testing.T) {
		// Create test data
		horseID := uint(2)
		now := testutil.MockTimeNow()
		testData := []*models.VitalSigns{
			{
				HorseID:     horseID,
				Temperature: 38.0,
				HeartRate:   40,
				RecordedAt:  now.Add(-48 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.2,
				HeartRate:   42,
				RecordedAt:  now.Add(-24 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.1,
				HeartRate:   41,
				RecordedAt:  now,
			},
		}

		for _, v := range testData {
			err := repo.SaveVitalSigns(ctx, v)
			require.NoError(t, err)
		}

		testCases := []struct {
			name          string
			from          time.Time
			to            time.Time
			expectedCount int
		}{
			{
				name:          "Full Range",
				from:          now.Add(-72 * time.Hour),
				to:            now.Add(time.Hour),
				expectedCount: 3,
			},
			{
				name:          "Last 24 Hours",
				from:          now.Add(-24 * time.Hour),
				to:            now.Add(time.Hour),
				expectedCount: 2,
			},
			{
				name:          "Future Range",
				from:          now.Add(24 * time.Hour),
				to:            now.Add(48 * time.Hour),
				expectedCount: 0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vitals, err := repo.GetVitalSigns(ctx, horseID, tc.from, tc.to)
				require.NoError(t, err)
				assert.Len(t, vitals, tc.expectedCount)
			})
		}
	})

	t.Run("GetLatestVitalSigns", func(t *testing.T) {
		horseID := uint(3)
		now := testutil.MockTimeNow()
		testData := []*models.VitalSigns{
			{
				HorseID:     horseID,
				Temperature: 38.0,
				HeartRate:   40,
				RecordedAt:  now.Add(-2 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.2,
				HeartRate:   42,
				RecordedAt:  now.Add(-1 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.1,
				HeartRate:   41,
				RecordedAt:  now,
			},
		}

		for _, v := range testData {
			err := repo.SaveVitalSigns(ctx, v)
			require.NoError(t, err)
		}

		latest, err := repo.GetLatestVitalSigns(ctx, horseID)
		require.NoError(t, err)
		assert.Equal(t, now, latest.RecordedAt)
		assert.Equal(t, float64(38.1), latest.Temperature)
		assert.Equal(t, uint(41), latest.HeartRate)
	})

	t.Run("Alerts", func(t *testing.T) {
		horseID := uint(4)
		now := testutil.MockTimeNow()

		// Create test alerts
		alerts := []*models.VitalSignsAlert{
			{
				HorseID:    horseID,
				Type:       "temperature_high",
				Message:    "Temperature above normal range",
				CreatedAt:  now.Add(-2 * time.Hour),
				Severity:   "warning",
				Parameter: "temperature",
				Value:     39.5,
			},
			{
				HorseID:    horseID,
				Type:       "heart_rate_low",
				Message:    "Heart rate below normal range",
				CreatedAt:  now.Add(-1 * time.Hour),
				Severity:   "critical",
				Parameter: "heart_rate",
				Value:     25,
			},
		}

		// Save alerts
		for _, alert := range alerts {
			err := repo.SaveAlert(ctx, alert)
			require.NoError(t, err)
		}

		// Test GetAlerts
		t.Run("GetAlerts", func(t *testing.T) {
			// Get unacknowledged alerts
			unacknowledged, err := repo.GetAlerts(ctx, horseID, false)
			require.NoError(t, err)
			assert.Len(t, unacknowledged, 2)

			// Acknowledge first alert
			err = repo.AcknowledgeAlert(ctx, unacknowledged[0].ID)
			require.NoError(t, err)

			// Get unacknowledged alerts again
			stillUnacknowledged, err := repo.GetAlerts(ctx, horseID, false)
			require.NoError(t, err)
			assert.Len(t, stillUnacknowledged, 1)

			// Get all alerts including acknowledged
			all, err := repo.GetAlerts(ctx, horseID, true)
			require.NoError(t, err)
			assert.Len(t, all, 2)
		})

		// Test GetAlert
		t.Run("GetAlert", func(t *testing.T) {
			alerts, err := repo.GetAlerts(ctx, horseID, true)
			require.NoError(t, err)
			require.NotEmpty(t, alerts)

			alert, err := repo.GetAlert(ctx, alerts[0].ID)
			require.NoError(t, err)
			assert.Equal(t, alerts[0].ID, alert.ID)
			assert.Equal(t, alerts[0].Type, alert.Type)
		})
	})

	t.Run("Trends", func(t *testing.T) {
		horseID := uint(5)
		now := testutil.MockTimeNow()

		// Create test trends
		trends := []*models.VitalSignsTrend{
			{
				HorseID:    horseID,
				Parameter:  "temperature",
				Trend:      "increasing",
				StartTime:  now.Add(-24 * time.Hour),
				EndTime:    now,
				StartValue: 38.0,
				EndValue:   38.5,
				Change:     0.5,
			},
			{
				HorseID:    horseID,
				Parameter:  "heart_rate",
				Trend:      "stable",
				StartTime:  now.Add(-12 * time.Hour),
				EndTime:    now,
				StartValue: 40,
				EndValue:   41,
				Change:     1,
			},
		}

		// Save trends
		for _, trend := range trends {
			err := repo.SaveTrend(ctx, trend)
			require.NoError(t, err)
		}

		// Test GetTrends
		t.Run("GetTrends", func(t *testing.T) {
			// Get all trends
			allTrends, err := repo.GetTrends(ctx, horseID, now.Add(-48*time.Hour), now.Add(time.Hour))
			require.NoError(t, err)
			assert.Len(t, allTrends, 2)

			// Get recent trends
			recentTrends, err := repo.GetTrends(ctx, horseID, now.Add(-12*time.Hour), now.Add(time.Hour))
			require.NoError(t, err)
			assert.Len(t, recentTrends, 1)
			assert.Equal(t, "heart_rate", recentTrends[0].Parameter)
		})
	})
})
