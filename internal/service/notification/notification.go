package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type NotificationType string

const (
	SystemNotification  NotificationType = "SYSTEM"
	HealthCheckDue      NotificationType = "HEALTH_CHECK_DUE"
	PregnancyMilestone NotificationType = "PREGNANCY_MILESTONE"
	VaccinationDue     NotificationType = "VACCINATION_DUE"
	WeatherAlert       NotificationType = "WEATHER_ALERT"
)

// Priority represents the importance level of a notification
type Priority string

const (
	High   Priority = "HIGH"
	Medium Priority = "MEDIUM"
	Low    Priority = "LOW"
)

type Notification struct {
	ID        int64           `json:"id"`
	Type      NotificationType `json:"type"`
	UserID    string          `json:"userId,omitempty"`
	HorseID   int64           `json:"horseId,omitempty"`
	Title     string          `json:"title,omitempty"`
	Message   string          `json:"message"`
	DueDate   time.Time       `json:"dueDate,omitempty"`
	Priority  Priority        `json:"priority"`
	Read      bool            `json:"read"`
	Completed bool            `json:"completed"`
	CreatedAt time.Time       `json:"createdAt"`
}

type Repository interface {
	SaveNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error)
	GetByID(ctx context.Context, id uint) (*Notification, error)
	MarkAsRead(ctx context.Context, notificationID uint) error
	Delete(ctx context.Context, id uint) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SendNotification(ctx context.Context, notification *Notification) error {
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}
	
	if err := s.repo.SaveNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to save notification: %w", err)
	}
	return nil
}

func (s *Service) GetUserNotifications(ctx context.Context, userID string, limit int) ([]*Notification, error) {
	notifications, err := s.repo.GetNotifications(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifications, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*Notification, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return notification, nil
}

func (s *Service) MarkAsRead(ctx context.Context, id uint) error {
	if err := s.repo.MarkAsRead(ctx, id); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

func (s *Service) CheckPregnancyMilestones(horse *models.Horse) []Notification {
	if horse.ConceptionDate == nil {
		return nil
	}

	var notifications []Notification
	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)

	milestones := map[int]string{
		80:  "First trimester check due",
		130: "Vaccination booster due",
		145: "Begin increasing feed",
		270: "Check for udder development",
		310: "Prepare foaling area",
		330: "Monitor for signs of imminent foaling",
	}

	for days, message := range milestones {
		if daysPregnant >= days && daysPregnant <= days+7 {
			notification := Notification{
				Type:      PregnancyMilestone,
				HorseID:   int64(horse.ID),
				Title:     "Pregnancy Milestone",
				Message:   message,
				DueDate:   time.Now().AddDate(0, 0, days-daysPregnant),
				Priority:  Medium,
				CreatedAt: time.Now(),
			}
			notifications = append(notifications, notification)
			
			logger.Info("Created pregnancy milestone notification", map[string]interface{}{
				"horseID":  horse.ID,
				"message": message,
				"daysPregnant": daysPregnant,
			})
		}
	}

	return notifications
}

func (s *Service) CheckHealthRecords(horse *models.Horse, records []models.HealthRecord) []Notification {
	var notifications []Notification

	// Check last vaccination date
	var lastVaccination time.Time
	for _, record := range records {
		if record.Type == "Vaccination" && record.Date.After(lastVaccination) {
			lastVaccination = record.Date
		}
	}

	// If last vaccination was more than 11 months ago or never done
	if lastVaccination.IsZero() || time.Since(lastVaccination) > 11*30*24*time.Hour {
		notification := Notification{
			Type:      VaccinationDue,
			HorseID:   int64(horse.ID),
			Title:     "Vaccination Due",
			Message:   "Annual vaccination due",
			DueDate:   time.Now().AddDate(0, 0, 7),
			Priority:  High,
			CreatedAt: time.Now(),
		}
		notifications = append(notifications, notification)
		
		logger.Info("Created vaccination due notification", map[string]interface{}{
			"horseID":  horse.ID,
			"dueDate": notification.DueDate,
		})
	}

	return notifications
}

func (s *Service) MarkNotificationComplete(id int64) error {
	return fmt.Errorf("not implemented")
}
