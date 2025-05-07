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
		userID := "user-1" 
		horseID := uint(1)
		pregnancyID := uint(1)
		currentTime := time.Now()

		eventInput := &models.PregnancyEventInputDTO{
			Type:        "checkup",
			Description: "Vet checkup",
			Date:        currentTime,
		}

		expectedRepoEvent := &models.PregnancyEvent{
			PregnancyID: pregnancyID,
			UserID:      userID,
			Type:        eventInput.Type,
			Description: eventInput.Description,
			Date:        eventInput.Date,
		}

		mockPregnancyRepo.On("GetCurrentPregnancy", mock.Anything, horseID).
			Return(&models.Pregnancy{ID: pregnancyID, HorseID: horseID}, nil).Once()

		mockPregnancyRepo.On("AddPregnancyEvent", mock.Anything, mock.MatchedBy(func(e *models.PregnancyEvent) bool {
			return e.PregnancyID == expectedRepoEvent.PregnancyID && 
				   e.UserID == expectedRepoEvent.UserID && 
				   e.Type == expectedRepoEvent.Type && 
				   e.Description == expectedRepoEvent.Description && 
				   e.Date.Equal(expectedRepoEvent.Date)
		})).Return(nil).Once()
		
		mockPregnancyRepo.On("GetEvents", mock.Anything, horseID).Return([]models.PregnancyEvent{*expectedRepoEvent}, nil).Once()

		createdEvent, err := handler.GetPregnancyService().AddPregnancyEvent(ctx, userID, horseID, eventInput)
		assert.NoError(t, err)
		assert.NotNil(t, createdEvent)
		assert.Equal(t, expectedRepoEvent.PregnancyID, createdEvent.PregnancyID)
		assert.Equal(t, expectedRepoEvent.UserID, createdEvent.UserID)
		assert.Equal(t, expectedRepoEvent.Type, createdEvent.Type)
		assert.Equal(t, expectedRepoEvent.Description, createdEvent.Description)

		events, err := handler.GetPregnancyService().GetPregnancyEvents(ctx, horseID)
		assert.NoError(t, err)
		assert.Len(t, events, 1)
		assert.Equal(t, "Vet checkup", events[0].Description)

		mockPregnancyRepo.AssertExpectations(t)
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