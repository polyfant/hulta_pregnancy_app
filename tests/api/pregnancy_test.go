package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestPregnancyHandler(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := setupTestContext(t)

	t.Run("StartPregnancyTracking", func(t *testing.T) {
		horseID := uint(1)
		start := models.PregnancyStart{
			ConceptionDate: time.Now(),
		}

		mockPregnancyRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Pregnancy")).
			Return(nil).Once()
		mockHorseRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Horse")).
			Return(nil).Once()

		err := handler.PregnancyService.StartTracking(ctx, horseID, start)
		assert.NoError(t, err)

		mockPregnancyRepo.AssertExpectations(t)
		mockHorseRepo.AssertExpectations(t)
	})
} 