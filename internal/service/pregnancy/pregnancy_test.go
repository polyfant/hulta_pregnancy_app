package pregnancy

import (
	"database/sql"
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

func (m *MockDB) AddHealthRecord(record *models.HealthRecord) error {
	args := m.Called(record)
	return args.Error(0)
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockDB) GetAllHorses() ([]models.Horse, error) {
	args := m.Called()
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockDB) GetHorse(id int64) (models.Horse, error) {
	args := m.Called(id)
	return args.Get(0).(models.Horse), args.Error(1)
}

func (m *MockDB) AddHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDB) UpdateHorse(horse *models.Horse) error {
	args := m.Called(horse)
	return args.Error(0)
}

func (m *MockDB) DeleteHorse(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) GetBreedingCosts(horseID int64) ([]models.BreedingCost, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.BreedingCost), args.Error(1)
}

func (m *MockDB) AddBreedingCost(cost *models.BreedingCost) error {
	args := m.Called(cost)
	return args.Error(0)
}

func (m *MockDB) DeleteBreedingCost(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockDB) GetLastSyncTime(userID int64) (time.Time, error) {
	args := m.Called(userID)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) GetPendingSyncCount(userID int64) (int, error) {
	args := m.Called(userID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockDB) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	args := m.Called(horseID)
	return args.Get(0).([]models.PregnancyEvent), args.Error(1)
}

func (m *MockDB) UpdatePregnancyEvent(event *models.PregnancyEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockDB) DeletePregnancyEvent(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDB) GetUserPregnancyEvents(userID int64) ([]models.PregnancyEvent, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.PregnancyEvent), args.Error(1)
}

func (m *MockDB) GetUserHorses(userID int64) ([]models.Horse, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Horse), args.Error(1)
}

func (m *MockDB) GetHorseByName(name string) (models.Horse, error) {
	args := m.Called(name)
	return args.Get(0).(models.Horse), args.Error(1)
}

func (m *MockDB) UpdateUserLastSync(userID int64, t time.Time) error {
	args := m.Called(userID, t)
	return args.Error(0)
}

func (m *MockDB) GetUserLastSync(userID int64) (time.Time, error) {
	args := m.Called(userID)
	return args.Get(0).(time.Time), args.Error(1)
}

func TestPregnancyService(t *testing.T) {
	mockDB := new(MockDB)
	tx := &sql.Tx{}
	mockDB.On("Begin").Return(tx, nil)
	service := NewService(mockDB)

	t.Run("add pregnancy event", func(t *testing.T) {
		event := &models.PregnancyEvent{
			HorseID:     1,
			EventType:   "Checkup",
			Notes:       "Regular checkup",
			EventDate:   time.Now(),
		}

		mockDB.On("AddPregnancyEvent", event).Return(nil)
		err := service.AddEvent(event)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("get pregnancy stage", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -2, 0) // 60 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		stage := service.GetPregnancyStage(horse)
		assert.Equal(t, models.EarlyGestation, stage)
	})

	t.Run("get pregnancy guidelines", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -5, 0) // 150 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		guidelines, err := service.GetPregnancyGuidelines(horse)
		assert.NoError(t, err)
		assert.NotNil(t, guidelines)
		assert.Equal(t, models.MidGestation, guidelines.Stage)
	})

	t.Run("check pre-foaling signs", func(t *testing.T) {
		conceptionDate := time.Now().AddDate(0, -11, 0) // 330 days
		horse := models.Horse{
			ConceptionDate: &conceptionDate,
		}

		signs := service.CheckPreFoalingSigns(horse)
		assert.NotNil(t, signs)
		assert.Greater(t, len(signs), 0)
	})

	t.Run("record pre-foaling sign", func(t *testing.T) {
		mockDB.On("AddPregnancyEvent", mock.AnythingOfType("*models.PregnancyEvent")).Return(nil).Once()

		err := service.RecordPreFoalingSign(1, "Udder Development")
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
