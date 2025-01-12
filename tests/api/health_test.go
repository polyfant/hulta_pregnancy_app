package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestHealthHandler(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := setupTestContext(t)

	t.Run("AddHealthRecord", func(t *testing.T) {
		horseID := uint(1)
		record := &models.HealthRecord{
			HorseID:     horseID,
			Type:        models.HealthRecordTypeVaccination,
			Description: "Annual vaccination",
			Date:        time.Now(),
		}

		mockHealthRepo.On("CreateRecord", mock.Anything, record).
			Return(nil).Once()

		err := handler.HealthService.CreateRecord(ctx, record)
		assert.NoError(t, err)

		mockHealthRepo.AssertExpectations(t)
	})
} 