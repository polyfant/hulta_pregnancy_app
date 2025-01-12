package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestHorseService(t *testing.T) {
	handler := setupTestHandler(t)
	ctx := setupTestContext(t)

	t.Run("GetHorse", func(t *testing.T) {
		horseID := uint(1)
		expectedHorse := &models.Horse{
			ID:     horseID,
			Name:   "Test Horse",
			UserID: "test_user",
		}

		mockHorseRepo.On("GetByID", mock.Anything, horseID).
			Return(expectedHorse, nil).Once()

		horse, err := handler.HorseService.GetByID(ctx, horseID)
		assert.NoError(t, err)
		assert.Equal(t, expectedHorse, horse)

		mockHorseRepo.AssertExpectations(t)
	})
} 