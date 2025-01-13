package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestPregnancyService(t *testing.T) {
	handler, _, _, _, mockPregnancyRepo, _ := setupTestHandler()
	ctx := setupTestContext(t)

	t.Run("StartPregnancyTracking", func(t *testing.T) {
		horseID := uint(1)
		start := models.PregnancyStart{
			ConceptionDate: time.Now(),
		}

		mockPregnancyRepo.On("StartTracking", mock.Anything, horseID, start).
			Return(nil).Once()
		mockPregnancyRepo.On("GetByHorseID", mock.Anything, horseID).
			Return(&models.Pregnancy{HorseID: horseID}, nil).Once()

		err := handler.GetPregnancyService().StartTracking(ctx, horseID, start)
		assert.NoError(t, err)

		pregnancy, err := handler.GetPregnancyService().GetPregnancy(ctx, horseID)
		assert.NoError(t, err)
		assert.Equal(t, horseID, pregnancy.HorseID)

		mockPregnancyRepo.AssertExpectations(t)
	})
} 