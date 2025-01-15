package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
)

// WeatherServiceInterface defines the interface for weather service
type WeatherServiceInterface interface {
	GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*weather.WeatherData, error)
}

type NotificationService interface {
	SendNotification(ctx context.Context, notification *Notification) error
}

type WeatherNotificationService struct {
	notificationService *Service
	weatherService     WeatherServiceInterface
}

func NewWeatherNotificationService(notificationService *Service, weatherService WeatherServiceInterface) *WeatherNotificationService {
	return &WeatherNotificationService{
		notificationService: notificationService,
		weatherService:     weatherService,
	}
}

func (s *WeatherNotificationService) ProcessWeatherNotifications(ctx context.Context, user *models.User) error {
	if !user.WeatherSettings.NotificationsEnabled {
		return nil
	}

	// Get current weather data
	weatherData, err := s.weatherService.GetCurrentWeather(ctx, 
		user.WeatherSettings.DefaultLatitude, 
		user.WeatherSettings.DefaultLongitude)
	if err != nil {
		return fmt.Errorf("failed to get weather data: %w", err)
	}

	// Create notification based on weather conditions
	notification := s.createWeatherNotification(user.ID, weatherData)
	if notification != nil {
		if err := s.notificationService.SendNotification(ctx, notification); err != nil {
			return fmt.Errorf("failed to send weather notification: %w", err)
		}
	}

	return nil
}

func (s *WeatherNotificationService) createWeatherNotification(userID string, weatherData *weather.WeatherData) *Notification {
	if weatherData == nil {
		return nil
	}

	var notification *Notification

	// Check for severe weather conditions
	if weatherData.HasSevereConditions() {
		notification = &Notification{
			Type:     WeatherAlert,
			UserID:   userID,
			Title:    "Severe Weather Alert",
			Message:  fmt.Sprintf("Severe weather conditions detected: %s", weatherData.Description),
			Priority: High,
		}
	} else if weatherData.Temperature > 30 {
		notification = &Notification{
			Type:     WeatherAlert,
			UserID:   userID,
			Title:    "High Temperature Alert",
			Message:  fmt.Sprintf("High temperature detected: %.1f°C", weatherData.Temperature),
			Priority: Medium,
		}
	} else if weatherData.Temperature < 0 {
		notification = &Notification{
			Type:     WeatherAlert,
			UserID:   userID,
			Title:    "Low Temperature Alert",
			Message:  fmt.Sprintf("Low temperature detected: %.1f°C", weatherData.Temperature),
			Priority: Medium,
		}
	}

	if notification != nil {
		notification.CreatedAt = time.Now()
	}

	return notification
}
