package pregnancy

import (
	"testing"
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestGetPregnancyStage(t *testing.T) {
	service := NewService(nil)

	t.Run("early gestation", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -2, 0) // 60 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		stage := service.GetPregnancyStage(horse)
		assert.Equal(t, models.EarlyGestation, stage)
	})

	t.Run("mid gestation", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -5, 0) // 150 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		stage := service.GetPregnancyStage(horse)
		assert.Equal(t, models.MidGestation, stage)
	})

	t.Run("late gestation", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -8, 0) // 240 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		stage := service.GetPregnancyStage(horse)
		assert.Equal(t, models.LateGestation, stage)
	})

	t.Run("final gestation", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -11, 0) // 330 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		stage := service.GetPregnancyStage(horse)
		assert.Equal(t, models.FinalGestation, stage)
	})
}

func TestGetPregnancyGuidelines(t *testing.T) {
	service := NewService(nil)

	t.Run("get guidelines for pregnant horse", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -5, 0)
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		guidelines, err := service.GetPregnancyGuidelines(horse)
		assert.NoError(t, err)
		assert.NotNil(t, guidelines)
		assert.Equal(t, models.MidGestation, guidelines.Stage)
	})

	t.Run("get guidelines for non-pregnant horse", func(t *testing.T) {
		horse := models.Horse{}

		guidelines, err := service.GetPregnancyGuidelines(horse)
		assert.Error(t, err)
		assert.Nil(t, guidelines)
	})
}

func TestCheckPreFoalingSigns(t *testing.T) {
	service := NewService(nil)

	t.Run("check signs for late pregnancy", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -11, 0)
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		signs := service.CheckPreFoalingSigns(horse)
		assert.NotNil(t, signs)
		assert.Greater(t, len(signs), 0)
	})

	t.Run("check signs for early pregnancy", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -2, 0)
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		signs := service.CheckPreFoalingSigns(horse)
		assert.Nil(t, signs)
	})
}

func TestRecordPreFoalingSign(t *testing.T) {
	mockDB := new(MockDB)
	service := NewService(mockDB)

	t.Run("record pre-foaling sign", func(t *testing.T) {
		mockDB.On("AddPregnancyEvent", mock.AnythingOfType("*models.PregnancyEvent")).Return(nil).Once()

		err := service.RecordPreFoalingSign(1, "Udder Development")
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
