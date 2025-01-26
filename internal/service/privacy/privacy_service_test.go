package privacy

import (
	"context"
	"errors"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetPrivacyPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PrivacyPreferences), args.Error(1)
}

func (m *MockRepository) UpdatePrivacyPreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error {
	args := m.Called(ctx, userID, prefs)
	return args.Error(0)
}

func (m *MockRepository) DeleteUserData(ctx context.Context, userID string, dataType string) error {
	args := m.Called(ctx, userID, dataType)
	return args.Error(0)
}

func TestPrivacyService_GetPrivacyPreferences(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		expected := &models.PrivacyPreferences{
			DataRetentionDays: 30,
			AutoDeleteOldData: true,
		}

		mockRepo.On("GetPrivacyPreferences", ctx, "user123").Return(expected, nil).Once()

		result, err := service.GetPrivacyPreferences(ctx, "user123")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetPrivacyPreferences", ctx, "user123").Return(nil, errors.New("db error")).Once()

		result, err := service.GetPrivacyPreferences(ctx, "user123")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPrivacyService_UpdatePrivacyPreferences(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			DataRetentionDays: 60,
			AutoDeleteOldData: true,
		}

		mockRepo.On("UpdatePrivacyPreferences", ctx, "user123", prefs).Return(nil).Once()

		err := service.UpdatePrivacyPreferences(ctx, "user123", prefs)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			DataRetentionDays: 60,
			AutoDeleteOldData: true,
		}

		mockRepo.On("UpdatePrivacyPreferences", ctx, "user123", prefs).Return(errors.New("db error")).Once()

		err := service.UpdatePrivacyPreferences(ctx, "user123", prefs)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPrivacyService_DeleteUserData(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("success", func(t *testing.T) {
		mockRepo.On("DeleteUserData", ctx, "user123", "health_data").Return(nil).Once()

		err := service.DeleteUserData(ctx, "user123", "health_data")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("DeleteUserData", ctx, "user123", "health_data").Return(errors.New("db error")).Once()

		err := service.DeleteUserData(ctx, "user123", "health_data")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPrivacyService_DeleteExpiredData(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	t.Run("success - data deleted", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			DataRetentionDays: 30,
			AutoDeleteOldData: true,
		}

		mockRepo.On("GetPrivacyPreferences", ctx, "").Return(prefs, nil).Once()
		mockRepo.On("DeleteUserData", ctx, "", "expired").Return(nil).Once()

		err := service.DeleteExpiredData(ctx)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("success - auto delete disabled", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			DataRetentionDays: 30,
			AutoDeleteOldData: false,
		}

		mockRepo.On("GetPrivacyPreferences", ctx, "").Return(prefs, nil).Once()

		err := service.DeleteExpiredData(ctx)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error getting preferences", func(t *testing.T) {
		mockRepo.On("GetPrivacyPreferences", ctx, "").Return(nil, errors.New("db error")).Once()

		err := service.DeleteExpiredData(ctx)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error deleting data", func(t *testing.T) {
		prefs := &models.PrivacyPreferences{
			DataRetentionDays: 30,
			AutoDeleteOldData: true,
		}

		mockRepo.On("GetPrivacyPreferences", ctx, "").Return(prefs, nil).Once()
		mockRepo.On("DeleteUserData", ctx, "", "expired").Return(errors.New("db error")).Once()

		err := service.DeleteExpiredData(ctx)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
