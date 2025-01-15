package notification

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SaveNotification(ctx context.Context, notification *Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockRepository) GetNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Notification), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id uint) (*Notification, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Notification), args.Error(1)
}

func (m *MockRepository) MarkAsRead(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_SendNotification(t *testing.T) {
	tests := []struct {
		name          string
		notification  *Notification
		mockError     error
		expectedError bool
	}{
		{
			name: "successful send",
			notification: &Notification{
				Type:      SystemNotification,
				UserID:    "user123",
				Title:     "Test Notification",
				Message:   "This is a test",
				Priority:  Medium,
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "repository error",
			notification: &Notification{
				Type:    SystemNotification,
				UserID:  "user123",
				Title:   "Test Notification",
				Message: "This is a test",
			},
			mockError:     errors.New("repository error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo)

			mockRepo.On("SaveNotification", mock.Anything, mock.MatchedBy(func(n *Notification) bool {
				return n.Type == tt.notification.Type &&
					n.UserID == tt.notification.UserID &&
					n.Title == tt.notification.Title &&
					n.Message == tt.notification.Message &&
					!n.CreatedAt.IsZero()
			})).Return(tt.mockError)

			err := service.SendNotification(context.Background(), tt.notification)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetUserNotifications(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		limit          int
		mockNotifs     []*Notification
		mockError      error
		expectedError  bool
		expectedLength int
	}{
		{
			name:   "successful retrieval",
			userID: "user123",
			limit:  10,
			mockNotifs: []*Notification{
				{ID: 1, UserID: "user123"},
				{ID: 2, UserID: "user123"},
			},
			mockError:      nil,
			expectedError:  false,
			expectedLength: 2,
		},
		{
			name:           "repository error",
			userID:         "user123",
			limit:          10,
			mockNotifs:     nil,
			mockError:      errors.New("repository error"),
			expectedError:  true,
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo)

			mockRepo.On("GetNotifications", mock.Anything, tt.userID, tt.limit).Return(tt.mockNotifs, tt.mockError)

			result, err := service.GetUserNotifications(context.Background(), tt.userID, tt.limit)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLength)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		id             uint
		mockNotif      *Notification
		mockError      error
		expectedError  bool
		expectedResult *Notification
	}{
		{
			name: "successful retrieval",
			id:   1,
			mockNotif: &Notification{
				ID:      1,
				Type:    SystemNotification,
				UserID:  "user123",
				Title:   "Test Notification",
				Message: "This is a test",
			},
			mockError:      nil,
			expectedError:  false,
			expectedResult: &Notification{ID: 1},
		},
		{
			name:           "not found",
			id:            2,
			mockNotif:     nil,
			mockError:     errors.New("not found"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			service := NewService(mockRepo)

			mockRepo.On("GetByID", mock.Anything, tt.id).Return(tt.mockNotif, tt.mockError)

			result, err := service.GetByID(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockNotif, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CheckPregnancyMilestones(t *testing.T) {
	conceptionDate := time.Now().AddDate(0, 0, -80) // 80 days ago
	tests := []struct {
		name            string
		horse           *models.Horse
		expectedNotifs  int
		expectedMessage string
	}{
		{
			name: "no conception date",
			horse: &models.Horse{
				ID: 1,
			},
			expectedNotifs: 0,
		},
		{
			name: "first trimester check due",
			horse: &models.Horse{
				ID:            1,
				ConceptionDate: &conceptionDate,
			},
			expectedNotifs:  1,
			expectedMessage: "First trimester check due",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(new(MockRepository))
			notifications := service.CheckPregnancyMilestones(tt.horse)

			assert.Len(t, notifications, tt.expectedNotifs)
			if tt.expectedNotifs > 0 {
				assert.Equal(t, tt.expectedMessage, notifications[0].Message)
				assert.Equal(t, PregnancyMilestone, notifications[0].Type)
				assert.Equal(t, int64(tt.horse.ID), notifications[0].HorseID)
			}
		})
	}
}
