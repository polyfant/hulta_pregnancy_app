package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// These are generic mock types that mirror the interfaces in the notification package
// without directly importing it to avoid circular dependencies

type Notification struct {
	ID          int64
	Type        string
	UserID      string
	HorseID     int64
	Title       string
	Message     string
	DueDate     time.Time
	Priority    string
	Read        bool
	Completed   bool
	CreatedAt   time.Time
}

type Repository interface {
	mock.Mock
	SaveNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error)
	GetByID(ctx context.Context, id uint) (*Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	Delete(ctx context.Context, id uint) error
}

type EmailNotifier interface {
	mock.Mock
	SendEmail(ctx context.Context, to, subject, body string) error
}

type WebSocketBroadcaster interface {
	mock.Mock
	Broadcast(ctx context.Context, message interface{}) error
}

type WeatherService interface {
	mock.Mock
	GetWeatherData(ctx context.Context, latitude, longitude float64) (*WeatherData, error)
}

type WeatherData struct {
	Temperature float64
	Description string
	Conditions  []string
}
