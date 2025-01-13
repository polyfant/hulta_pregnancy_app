package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestBreedingService(t *testing.T) {
	handler, _, _, _, _, mockBreedingRepo := setupTestHandler()
	ctx := setupTestContext(t)

	t.Run("AddBreedingRecord", func(t *testing.T) {
		record := &models.BreedingRecord{
			HorseID: 1,
			Date:    time.Now(),
			Status:  string(models.BreedingStatusActive),
		}

		mockBreedingRepo.On("CreateRecord", mock.Anything, record).
			Return(nil).Once()

		err := handler.GetBreedingService().CreateRecord(ctx, record)
		assert.NoError(t, err)

		mockBreedingRepo.AssertExpectations(t)
	})

	t.Run("GetBreedingRecords", func(t *testing.T) {
		horseID := 1

		mockBreedingRepo.On("GetRecords", mock.Anything, uint(horseID)).
			Return([]models.BreedingRecord{}, nil).Once()

		records, err := handler.GetBreedingService().GetRecords(ctx, uint(horseID))
		assert.NoError(t, err)
		assert.Empty(t, records)

		mockBreedingRepo.AssertExpectations(t)
	})
}