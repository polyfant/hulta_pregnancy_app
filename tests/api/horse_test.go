package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestHorseService(t *testing.T) {
	handler, mockHorseRepo, _, _, _, _ := setupTestHandler()
	ctx := setupTestContext(t)

	t.Run("GetHorse", func(t *testing.T) {
		horseID := uint(1)

		mockHorseRepo.On("GetByID", mock.Anything, horseID).
			Return(&models.Horse{ID: horseID}, nil).Once()

		horse, err := handler.GetHorseService().GetByID(ctx, horseID)
		assert.NoError(t, err)
		assert.Equal(t, horseID, horse.ID)

		mockHorseRepo.AssertExpectations(t)
	})
} 