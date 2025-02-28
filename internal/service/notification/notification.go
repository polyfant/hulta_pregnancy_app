package notification

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
	websocket "github.com/polyfant/hulta_pregnancy_app/internal/service/notification/websocket"
)

// Notification Types
type NotificationType string

const (
	SystemNotification   NotificationType = "SYSTEM"
	HealthCheckDue       NotificationType = "HEALTH_CHECK_DUE"
	PregnancyMilestone   NotificationType = "PREGNANCY_MILESTONE"
	VaccinationDue       NotificationType = "VACCINATION_DUE"
	WeatherAlert         NotificationType = "WEATHER_ALERT"
)

// Priority represents the importance level of a notification
type Priority string

const (
	High   Priority = "HIGH"
	Medium Priority = "MEDIUM"
	Low    Priority = "LOW"
)

// Notification struct defines the structure of a notification
type Notification struct {
	ID          int64           `json:"id"`
	Type        NotificationType `json:"type"`
	UserID      string          `json:"userId,omitempty"`
	HorseID     int64           `json:"horseId,omitempty"`
	Title       string          `json:"title,omitempty"`
	Message     string          `json:"message"`
	DueDate     time.Time       `json:"dueDate,omitempty"`
	Priority    Priority        `json:"priority"`
	Read        bool            `json:"read"`
	Completed   bool            `json:"completed"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// Repository interface defines the methods for interacting with the notification repository
type Repository interface {
	SaveNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error)
	GetNotificationByID(ctx context.Context, id uint) (*Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

// WeatherService defines the interface for retrieving weather data
type WeatherService interface {
	GetWeatherData(ctx context.Context, latitude, longitude float64) (*weather.WeatherData, error)
}

type Service struct {
	repo                 Repository
	userRepo             repository.UserRepository
	emailNotifier        EmailNotifier
	websocketBroadcaster websocket.WebSocketBroadcaster
	weatherService       WeatherService
}

func NewService(
	repo Repository, 
	userRepo repository.UserRepository, 
	emailNotifier EmailNotifier, 
	websocketBroadcaster websocket.WebSocketBroadcaster,
	weatherService WeatherService,
) *Service {
	return &Service{
		repo:                 repo,
		userRepo:             userRepo,
		emailNotifier:        emailNotifier,
		websocketBroadcaster: websocketBroadcaster,
		weatherService:       weatherService,
	}
}

// SendNotification creates and saves a new notification
func (s *Service) SendNotification(ctx context.Context, notification *Notification) error {
    // Set the CreatedAt timestamp to the current time
    notification.CreatedAt = time.Now()

    // Save the notification to the repository
    return s.repo.SaveNotification(ctx, notification)
}

// GetUserNotifications retrieves notifications for a specific user
func (s *Service) GetUserNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error) {
	return s.repo.GetNotifications(ctx, userID, limit)
}

// GetByID retrieves a specific notification by its ID
func (s *Service) GetByID(ctx context.Context, id uint) (*Notification, error) {
	return s.repo.GetNotificationByID(ctx, id)
}

// CheckPregnancyMilestones checks and potentially creates notifications for pregnancy milestones
func (s *Service) CheckPregnancyMilestones(ctx context.Context, userID string) error {
	// Placeholder implementation
	// In a real-world scenario, this would:
	// 1. Fetch user's pregnancy data
	// 2. Determine current milestone
	// 3. Create and save milestone notifications if needed
	return nil
}