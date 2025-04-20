package notification

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
	"github.com/stretchr/testify/mock"
)

// MockRepository implements Repository interface for testing
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

func (m *MockRepository) GetNotificationByID(ctx context.Context, id uint) (*Notification, error) {
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

// Implement UserRepository methods to satisfy the interface
func (m *MockRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DashboardStats), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockEmailNotifier implements EmailNotifier interface for testing
type MockEmailNotifier struct {
	mock.Mock
}

func (m *MockEmailNotifier) SendEmail(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}

// MockWebSocketBroadcaster implements WebSocketBroadcaster interface for testing
type MockWebSocketBroadcaster struct {
	mock.Mock
}

func (m *MockWebSocketBroadcaster) Broadcast(ctx context.Context, message interface{}) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// MockWeatherService implements WeatherService interface for testing
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetWeatherData(ctx context.Context, latitude, longitude float64) (*weather.WeatherData, error) {
	args := m.Called(ctx, latitude, longitude)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*weather.WeatherData), args.Error(1)
}
