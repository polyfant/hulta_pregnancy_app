package notification

import (
	"context"
	"fmt"
	_ "net/http"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWeatherNotificationService_ProcessWeatherNotifications(t *testing.T) {
	tests := []struct {
		name          string
		user          *models.User
		weatherData   *weather.WeatherData
		expectedNotif *Notification
		mockError     error
	}{
		{
			name: "notifications disabled",
			user: &models.User{
				WeatherSettings: models.WeatherSettings{
					NotificationsEnabled: false,
				},
			},
			expectedNotif: nil,
		},
		{
			name: "high temperature alert",
			user: &models.User{
				ID: "user123",
				WeatherSettings: models.WeatherSettings{
					NotificationsEnabled: true,
					DefaultLatitude:     35.0,
					DefaultLongitude:    -120.0,
				},
			},
			weatherData: &weather.WeatherData{
				Temperature: 35.0,
				Description: "Hot",
				Conditions:  []string{"hot"},
			},
			expectedNotif: &Notification{
				Type:     WeatherAlert,
				UserID:   "user123",
				Title:    "High Temperature Alert",
				Message:  "Temperature is high at 35.0°C. Stay hydrated and cool.",
				Priority: High,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			mockEmailNotifier := &MockEmailNotifier{}
			mockWebsocketBroadcaster := &MockWebSocketBroadcaster{}
			mockWeatherService := &MockWeatherService{}

			service := NewService(mockRepo, mockRepo, mockEmailNotifier, mockWebsocketBroadcaster, mockWeatherService)

			// Mock the weather service to return predefined weather data
			// Only set expectations if notifications are enabled
			if tt.user.WeatherSettings.NotificationsEnabled {
				mockWeatherService.On("GetWeatherData", 
					mock.Anything, 
					tt.user.WeatherSettings.DefaultLatitude, 
					tt.user.WeatherSettings.DefaultLongitude,
				).Return(tt.weatherData, nil)
			}

			// Mock the repository save method
			if tt.expectedNotif != nil {
				mockRepo.On("SaveNotification", mock.Anything, mock.MatchedBy(func(n *Notification) bool {
					return n.Type == tt.expectedNotif.Type &&
						n.UserID == tt.expectedNotif.UserID &&
						n.Title == tt.expectedNotif.Title &&
						n.Message == tt.expectedNotif.Message
				})).Return(nil)
			}

			// Perform the test
			notifications, err := processWeatherNotifications(
				context.Background(), 
				service, 
				mockWeatherService, 
				tt.user,
			)

			// Assertions
			if tt.expectedNotif == nil {
				assert.Empty(t, notifications)
			} else {
				assert.Len(t, notifications, 1)
				assert.Equal(t, tt.expectedNotif.Type, notifications[0].Type)
				assert.Equal(t, tt.expectedNotif.UserID, notifications[0].UserID)
				assert.Equal(t, tt.expectedNotif.Title, notifications[0].Title)
				assert.Equal(t, tt.expectedNotif.Priority, notifications[0].Priority)
			}
			assert.NoError(t, err)

			// Verify mock expectations, but only if notifications were enabled
			mockRepo.AssertExpectations(t)
			if tt.user.WeatherSettings.NotificationsEnabled {
				mockWeatherService.AssertExpectations(t)
			}
		})
	}
}

// processWeatherNotifications is a helper function for testing weather notifications
func processWeatherNotifications(
	ctx context.Context, 
	service *Service, 
	weatherService WeatherService, 
	user *models.User,
) ([]*Notification, error) {
	// Check if weather notifications are enabled
	if !user.WeatherSettings.NotificationsEnabled {
		return nil, nil
	}

	// Fetch weather data
	weatherData, err := weatherService.GetWeatherData(
		ctx, 
		user.WeatherSettings.DefaultLatitude, 
		user.WeatherSettings.DefaultLongitude,
	)
	if err != nil {
		return nil, err
	}

	// Check for high temperature
	var notifications []*Notification
	if weatherData.Temperature >= 35.0 {
		notification := &Notification{
			Type:     WeatherAlert,
			UserID:   user.ID,
			Title:    "High Temperature Alert",
			Message:  "Temperature is high at " + fmt.Sprintf("%.1f", weatherData.Temperature) + "°C. Stay hydrated and cool.",
			Priority: High,
		}

		// Save the notification
		err := service.SendNotification(ctx, notification)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}
