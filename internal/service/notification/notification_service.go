package notification

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Service defines the interface for notification operations
type Service interface {
	SendWeatherAlert(ctx context.Context, userID string, alert string) error
	SendVitalSignsAlert(ctx context.Context, userID string, alert *models.VitalSignsAlert) error
	SendPregnancyAlert(ctx context.Context, userID string, alert *models.PregnancyAlert) error
	GetUserNotifications(ctx context.Context, userID string) ([]*models.Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID uint) error
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new notification service instance
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Repository defines the interface for notification data operations
type Repository interface {
	SaveNotification(ctx context.Context, notification *models.Notification) error
	GetUserNotifications(ctx context.Context, userID string) ([]*models.Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID uint) error
}

func (s *service) SendWeatherAlert(ctx context.Context, userID string, alert string) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    "WEATHER",
		Message: alert,
	}
	return s.repo.SaveNotification(ctx, notification)
}

func (s *service) SendVitalSignsAlert(ctx context.Context, userID string, alert *models.VitalSignsAlert) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    "VITAL_SIGNS",
		Message: alert.Message,
	}
	return s.repo.SaveNotification(ctx, notification)
}

func (s *service) SendPregnancyAlert(ctx context.Context, userID string, alert *models.PregnancyAlert) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    "PREGNANCY",
		Message: alert.Message,
	}
	return s.repo.SaveNotification(ctx, notification)
}

func (s *service) GetUserNotifications(ctx context.Context, userID string) ([]*models.Notification, error) {
	return s.repo.GetUserNotifications(ctx, userID)
}

func (s *service) MarkNotificationRead(ctx context.Context, notificationID uint) error {
	return s.repo.MarkNotificationRead(ctx, notificationID)
}
