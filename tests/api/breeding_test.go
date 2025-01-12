package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestBreedingService(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := setupTestContext(t)

	t.Run("AddBreedingRecord", func(t *testing.T) {
		record := &models.BreedingRecord{
			HorseID: 1,
			Date:    time.Now(),
			Status:  string(models.BreedingStatusActive),
		}

		mockBreedingRepo.On("CreateRecord", mock.Anything, record).
			Return(nil).Once()

		err := handler.BreedingService.CreateRecord(ctx, record)
		assert.NoError(t, err)

		mockBreedingRepo.AssertExpectations(t)
	})
} 