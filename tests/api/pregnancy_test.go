package api_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestPregnancyService(t *testing.T) {
	handler, mockHorse, _, mockPregnancyRepo, _, _ := setupTestHandler()
	ctx := setupTestContext(t)

	t.Run("StartPregnancyTracking", func(t *testing.T) {
		horseID := uint(1)
		start := models.PregnancyStart{
			ConceptionDate: time.Now(),
		}

		mockPregnancyRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Pregnancy) bool {
			return p.HorseID == horseID && !p.ConceptionDate.IsZero()
		})).Return(nil).Once()

		mockHorse.On("GetByID", mock.Anything, horseID).
			Return(&models.Horse{ID: horseID}, nil).Once()

		mockHorse.On("Update", mock.Anything, mock.MatchedBy(func(h *models.Horse) bool {
			return h.ID == horseID && h.IsPregnant
		})).Return(nil).Once()

		mockPregnancyRepo.On("GetByHorseID", mock.Anything, horseID).
			Return(&models.Pregnancy{HorseID: horseID}, nil).Once()

		err := handler.GetPregnancyService().StartTracking(ctx, horseID, start)
		assert.NoError(t, err)

		pregnancy, err := handler.GetPregnancyService().GetPregnancy(ctx, horseID)
		assert.NoError(t, err)
		assert.Equal(t, horseID, pregnancy.HorseID)

		mockPregnancyRepo.AssertExpectations(t)
		mockHorse.AssertExpectations(t)
	})

	t.Run("AddAndGetPregnancyEvent", func(t *testing.T) {
		horseID := uint(1)
		event := &models.PregnancyEvent{
			PregnancyID: 1,
			UserID:      "user-1",
			Type:        "checkup",
			Description: "Vet checkup",
			Date:        time.Now(),
		}

		mockPregnancyRepo.On("AddPregnancyEvent", mock.Anything, event).Return(nil).Once()
		mockPregnancyRepo.On("GetEvents", mock.Anything, horseID).Return([]models.PregnancyEvent{*event}, nil).Once()

		err := handler.GetPregnancyService().AddPregnancyEvent(ctx, event)
		assert.NoError(t, err)

		events, err := handler.GetPregnancyService().GetPregnancyEvents(ctx, horseID)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Vet checkup", events[0].Description)
	})

	t.Run("AddAndGetPreFoalingSign", func(t *testing.T) {
		horseID := uint(1)
		sign := &models.PreFoalingSign{
			HorseID:     horseID,
			Description: "Waxing teats",
			Date:        time.Now(),
			Notes:       "Observed in the morning",
		}

		mockPregnancyRepo.On("AddPreFoaling", mock.Anything, sign).Return(nil).Once()
		mockPregnancyRepo.On("GetPreFoaling", mock.Anything, horseID).Return([]models.PreFoalingSign{*sign}, nil).Once()

		err := handler.GetPregnancyService().AddPreFoalingSign(ctx, sign)
		assert.NoError(t, err)

		signs, err := handler.GetPregnancyService().GetPreFoalingSigns(ctx, horseID)
		assert.NoError(t, err)
		assert.Len(t, signs, 1)
		assert.Equal(t, "Waxing teats", signs[0].Description)
	})

	t.Run("AddAndGetPreFoalingChecklistItem", func(t *testing.T) {
		horseID := uint(1)
		item := &models.PreFoalingChecklistItem{
			HorseID:     horseID,
			Description: "Prepare foaling kit",
			IsCompleted: false,
			DueDate:     time.Now().AddDate(0, 0, 7),
			Priority:    models.PriorityHigh,
		}

		mockPregnancyRepo.On("AddPreFoalingChecklistItem", mock.Anything, item).Return(nil).Once()
		mockPregnancyRepo.On("GetPreFoalingChecklist", mock.Anything, horseID).Return([]models.PreFoalingChecklistItem{*item}, nil).Once()

		err := handler.GetPregnancyService().AddPreFoalingChecklistItem(ctx, item)
		assert.NoError(t, err)

		items, err := handler.GetPregnancyService().GetPreFoalingChecklist(ctx, horseID)
		assert.NoError(t, err)
		assert.Len(t, items, 1)
		assert.Equal(t, "Prepare foaling kit", items[0].Description)
	})

	t.Run("GetGuidelinesForEachStage", func(t *testing.T) {
		stages := []models.PregnancyStage{
			models.PregnancyStageEarly,
			models.PregnancyStageMid,
			models.PregnancyStageLate,
		}
		for _, stage := range stages {
			guidelines, err := handler.GetPregnancyService().GetGuidelines(ctx, stage)
			assert.NoError(t, err)
			assert.NotEmpty(t, guidelines)
		}
	})

	t.Run("EndPregnancyAndUpdateStatus", func(t *testing.T) {
		horseID := uint(1)
		preg := &models.Pregnancy{
			HorseID: horseID,
			Status:  models.PregnancyStatusActive,
		}
		endDate := time.Now()
		mockPregnancyRepo.On("GetByHorseID", mock.Anything, horseID).Return(preg, nil).Once()
		mockPregnancyRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Pregnancy) bool {
			return p.Status == models.PregnancyStatusComplete && p.EndDate != nil
		})).Return(nil).Once()

		err := handler.GetPregnancyService().EndPregnancy(ctx, horseID, models.PregnancyStatusComplete, endDate)
		assert.NoError(t, err)
		assert.Equal(t, models.PregnancyStatusComplete, preg.Status)
		assert.NotNil(t, preg.EndDate)
	})
}