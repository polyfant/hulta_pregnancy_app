package notification

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	_ "github.com/polyfant/hulta_pregnancy_app/internal/models"
)

func TestService_SendNotification(t *testing.T) {
	tests := []struct {
		name          string
		notification  *Notification
		mockError     error
		expectedError bool
	}{
		{
			name: "successful notification",
			notification: &Notification{
				Type:     SystemNotification,
				UserID:   "user123",
				Title:    "Test Notification",
				Message:  "This is a test notification",
				Priority: High,
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "repository error",
			notification: &Notification{
				Type:     SystemNotification,
				UserID:   "user123",
				Title:    "Test Notification",
				Message:  "This is a test notification",
				Priority: High,
			},
			mockError:     errors.New("repository error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			mockEmailNotifier := &MockEmailNotifier{}
			mockWeatherService := &MockWeatherService{}
			service := NewService(mockRepo, mockRepo, mockEmailNotifier, nil, mockWeatherService)

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
			name:           "successful retrieval",
			userID:         "user123",
			limit:          10,
			mockNotifs:     []*Notification{{ID: 1}, {ID: 2}},
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
			mockRepo := &MockRepository{}
			mockEmailNotifier := &MockEmailNotifier{}
			mockWeatherService := &MockWeatherService{}
			service := NewService(mockRepo, mockRepo, mockEmailNotifier, nil, mockWeatherService)

			mockRepo.On("GetNotifications", mock.Anything, tt.userID, tt.limit).
				Return(tt.mockNotifs, tt.mockError)

			notifications, err := service.GetUserNotifications(context.Background(), tt.userID, tt.limit)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, notifications)
			} else {
				assert.NoError(t, err)
				assert.Len(t, notifications, tt.expectedLength)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		id            uint
		mockNotif     *Notification
		mockError     error
		expectedError bool
	}{
		{
			name:      "successful retrieval",
			id:        1,
			mockNotif: &Notification{ID: 1, UserID: "user123"},
			mockError: nil,
			expectedError: false,
		},
		{
			name:          "not found",
			id:            999,
			mockNotif:     nil,
			mockError:     errors.New("notification not found"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			mockEmailNotifier := &MockEmailNotifier{}
			mockWeatherService := &MockWeatherService{}
			service := NewService(mockRepo, mockRepo, mockEmailNotifier, nil, mockWeatherService)

			mockRepo.On("GetNotificationByID", mock.Anything, tt.id).
				Return(tt.mockNotif, tt.mockError)

			notification, err := service.GetByID(context.Background(), tt.id)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, notification)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, notification)
				assert.Equal(t, tt.id, uint(notification.ID))
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CheckPregnancyMilestones(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		expected error
	}{
		{
			name:     "successful milestone check",
			userID:   "user123",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			mockEmailNotifier := &MockEmailNotifier{}
			mockWeatherService := &MockWeatherService{}
			service := NewService(mockRepo, mockRepo, mockEmailNotifier, nil, mockWeatherService)

			err := service.CheckPregnancyMilestones(context.Background(), tt.userID)
			assert.Equal(t, tt.expected, err)
		})
	}
}
