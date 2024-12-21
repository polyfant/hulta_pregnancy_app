package health

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

func (m *MockDB) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.HealthRecord), args.Error(1)
}

func TestGetVaccinationStatus(t *testing.T) {
	mockDB := new(MockDB)
	service := NewService(mockDB)

	t.Run("up to date vaccination", func(t *testing.T) {
		horse := models.Horse{ID: 1}
		records := []models.HealthRecord{
			{
				HorseID: 1,
				Type:    "Vaccination",
				Date:    time.Now().AddDate(0, -6, 0),
			},
		}
		mockDB.On("GetHealthRecords", int64(1)).Return(records, nil).Once()

		status := service.GetVaccinationStatus(horse)
		assert.True(t, status.IsUpToDate)
	})

	t.Run("outdated vaccination", func(t *testing.T) {
		horse := models.Horse{ID: 1}
		records := []models.HealthRecord{
			{
				HorseID: 1,
				Type:    "Vaccination",
				Date:    time.Now().AddDate(-2, 0, 0),
			},
		}
		mockDB.On("GetHealthRecords", int64(1)).Return(records, nil).Once()

		status := service.GetVaccinationStatus(horse)
		assert.False(t, status.IsUpToDate)
	})
}

func TestGetHealthSummary(t *testing.T) {
	mockDB := new(MockDB)
	service := NewService(mockDB)

	t.Run("complete health summary", func(t *testing.T) {
		horse := models.Horse{ID: 1}
		records := []models.HealthRecord{
			{
				HorseID: 1,
				Type:    "Checkup",
				Date:    time.Now().AddDate(0, -1, 0),
			},
			{
				HorseID: 1,
				Type:    "Issue",
				Date:    time.Now().AddDate(0, -2, 0),
			},
		}
		mockDB.On("GetHealthRecords", int64(1)).Return(records, nil).Once()

		summary := service.GetHealthSummary(horse)
		assert.Equal(t, 2, summary.TotalRecords)
		assert.NotNil(t, summary.LastCheckup)
	})
}

func TestGetNutritionPlan(t *testing.T) {
	service := NewService(nil)

	t.Run("maintenance nutrition plan", func(t *testing.T) {
		horse := models.Horse{
			ID:     1,
			Weight: 500,
		}

		plan := service.GetNutritionPlan(horse)
		assert.Equal(t, "Maintenance", plan.Guidelines.Category)
	})

	t.Run("late pregnancy nutrition plan", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -9, 0)
		horse := models.Horse{
			ID:             1,
			Weight:         500,
			ConceptionDate: &conceptionDate,
		}

		plan := service.GetNutritionPlan(horse)
		assert.Equal(t, "LatePregnancy", plan.Guidelines.Category)
	})
}
