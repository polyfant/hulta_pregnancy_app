package notification

import (
	"context"
	"testing"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetCurrentWeather(ctx context.Context, lat, lon float64) (*weather.WeatherData, error) {
	args := m.Called(ctx, lat, lon)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weather.WeatherData), args.Error(1)
}

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
				Message:  "High temperature detected: 35.0°C",
				Priority: Medium,
			},
		},
		{
			name: "low temperature alert",
			user: &models.User{
				ID: "user123",
				WeatherSettings: models.WeatherSettings{
					NotificationsEnabled: true,
					DefaultLatitude:     35.0,
					DefaultLongitude:    -120.0,
				},
			},
			weatherData: &weather.WeatherData{
				Temperature: -5.0,
				Description: "Cold",
				Conditions:  []string{"cold"},
			},
			expectedNotif: &Notification{
				Type:     WeatherAlert,
				UserID:   "user123",
				Title:    "Low Temperature Alert",
				Message:  "Low temperature detected: -5.0°C",
				Priority: Medium,
			},
		},
		{
			name: "severe weather alert",
			user: &models.User{
				ID: "user123",
				WeatherSettings: models.WeatherSettings{
					NotificationsEnabled: true,
					DefaultLatitude:     35.0,
					DefaultLongitude:    -120.0,
				},
			},
			weatherData: &weather.WeatherData{
				Temperature: 25.0,
				Description: "Thunderstorm",
				Conditions:  []string{"thunderstorm"},
				WindSpeed:   25.0,
			},
			expectedNotif: &Notification{
				Type:     WeatherAlert,
				UserID:   "user123",
				Title:    "Severe Weather Alert",
				Message:  "Severe weather conditions detected: Thunderstorm",
				Priority: High,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNotifService := new(MockRepository)
			mockWeatherService := new(MockWeatherService)
			service := NewWeatherNotificationService(NewService(mockNotifService), mockWeatherService)

			if tt.weatherData != nil {
				mockWeatherService.On("GetCurrentWeather", mock.Anything, tt.user.WeatherSettings.DefaultLatitude, tt.user.WeatherSettings.DefaultLongitude).
					Return(tt.weatherData, nil)
			}

			if tt.expectedNotif != nil {
				mockNotifService.On("SaveNotification", mock.Anything, mock.MatchedBy(func(n *Notification) bool {
					return n.Type == tt.expectedNotif.Type &&
						n.UserID == tt.expectedNotif.UserID &&
						n.Title == tt.expectedNotif.Title &&
						n.Priority == tt.expectedNotif.Priority &&
						!n.CreatedAt.IsZero()
				})).Return(nil)
			}

			err := service.ProcessWeatherNotifications(context.Background(), tt.user)
			assert.NoError(t, err)

			mockWeatherService.AssertExpectations(t)
			mockNotifService.AssertExpectations(t)
		})
	}
}
