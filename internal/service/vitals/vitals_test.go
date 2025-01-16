package vitals

import (
	"context"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutil"
	"github.com/polyfant/hulta_pregnancy_app/internal/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVitalsService(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)

	repo := NewVitalsRepository(db)
	hub := websocket.NewHub()
	go hub.Run()

	service := NewService(repo, hub)

	t.Run("RecordVitalSigns", func(t *testing.T) {
		// Create a test pregnancy in late stage
		err := testutil.CreateTestPregnancy(db, 1, 15) // 15 days until due
		require.NoError(t, err)

		testCases := []struct {
			name          string
			vitals        *models.VitalSigns
			expectedAlert bool
		}{
			{
				name:          "Normal Vitals",
				vitals:        testutil.CreateTestVitalSigns(1),
				expectedAlert: false,
			},
			{
				name: "High Temperature",
				vitals: &models.VitalSigns{
					HorseID:     1,
					Temperature: 39.5, // High temperature
					HeartRate:   40,
					Respiration: 12,
					Hydration:   95,
					RecordedAt:  testutil.MockTimeNow(),
				},
				expectedAlert: true,
			},
			{
				name: "Low Heart Rate",
				vitals: &models.VitalSigns{
					HorseID:     1,
					Temperature: 38.0,
					HeartRate:   25, // Low heart rate
					Respiration: 12,
					Hydration:   95,
					RecordedAt:  testutil.MockTimeNow(),
				},
				expectedAlert: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// Record vital signs
				prediction, err := service.RecordVitalSigns(context.Background(), tc.vitals)
				require.NoError(t, err)
				assert.NotNil(t, prediction)

				// Check if alert was generated
				alerts, err := repo.GetAlerts(context.Background(), tc.vitals.HorseID, false)
				require.NoError(t, err)
				if tc.expectedAlert {
					assert.NotEmpty(t, alerts)
				} else {
					assert.Empty(t, alerts)
				}
			})
		}
	})

	t.Run("GetVitalSignsTrend", func(t *testing.T) {
		horseID := uint(2)
		
		// Create test data
		vitals := []models.VitalSigns{
			{
				HorseID:     horseID,
				Temperature: 38.0,
				HeartRate:   40,
				RecordedAt:  testutil.MockTimeNow().Add(-24 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.2,
				HeartRate:   42,
				RecordedAt:  testutil.MockTimeNow().Add(-12 * time.Hour),
			},
			{
				HorseID:     horseID,
				Temperature: 38.1,
				HeartRate:   41,
				RecordedAt:  testutil.MockTimeNow(),
			},
		}

		for _, v := range vitals {
			err := repo.SaveVitalSigns(context.Background(), &v)
			require.NoError(t, err)
		}

		// Test trend calculation
		trend, err := service.GetVitalSignsTrend(context.Background(), horseID, "TEMPERATURE", "24H")
		require.NoError(t, err)
		assert.NotNil(t, trend)
		assert.Len(t, trend.DataPoints, 3)
		assert.InDelta(t, 38.1, trend.Average, 0.1)
	})

	t.Run("CheckForAlerts", func(t *testing.T) {
		horseID := uint(3)

		testCases := []struct {
			name          string
			vitals        *models.VitalSigns
			expectedAlert bool
			alertType     string
		}{
			{
				name: "High Temperature Alert",
				vitals: &models.VitalSigns{
					HorseID:     horseID,
					Temperature: 39.5,
					HeartRate:   40,
					RecordedAt:  testutil.MockTimeNow(),
				},
				expectedAlert: true,
				alertType:     "HIGH_TEMPERATURE",
			},
			{
				name: "Low Hydration Alert",
				vitals: &models.VitalSigns{
					HorseID:     horseID,
					Temperature: 38.0,
					HeartRate:   40,
					Hydration:   85,
					RecordedAt:  testutil.MockTimeNow(),
				},
				expectedAlert: true,
				alertType:     "LOW_HYDRATION",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := service.checkForAlerts(context.Background(), tc.vitals)
				require.NoError(t, err)

				alerts, err := repo.GetAlerts(context.Background(), horseID, false)
				require.NoError(t, err)

				if tc.expectedAlert {
					assert.NotEmpty(t, alerts)
					assert.Equal(t, tc.alertType, alerts[0].AlertType)
				} else {
					assert.Empty(t, alerts)
				}
			})
		}
	})

	t.Run("Late Stage Detection", func(t *testing.T) {
		testCases := []struct {
			name           string
			daysUntilDue   int
			expectedResult bool
		}{
			{
				name:           "Late Stage - 15 days until due",
				daysUntilDue:   15,
				expectedResult: true,
			},
			{
				name:           "Not Late Stage - 45 days until due",
				daysUntilDue:   45,
				expectedResult: false,
			},
			{
				name:           "Very Late Stage - 5 days until due",
				daysUntilDue:   5,
				expectedResult: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				horseID := uint(100 + tc.daysUntilDue)
				err := testutil.CreateTestPregnancy(db, horseID, tc.daysUntilDue)
				require.NoError(t, err)

				result := service.isInLateStage(context.Background(), horseID)
				assert.Equal(t, tc.expectedResult, result)
			})
		}
	})
}
