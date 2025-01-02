package notification

import (
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type NotificationType string

const (
	HealthCheckDue    NotificationType = "HEALTH_CHECK_DUE"
	PregnancyMilestone NotificationType = "PREGNANCY_MILESTONE"
	VaccinationDue     NotificationType = "VACCINATION_DUE"
)

type Notification struct {
	ID        int64           `json:"id"`
	Type      NotificationType `json:"type"`
	HorseID   int64           `json:"horseId"`
	Message   string          `json:"message"`
	DueDate   time.Time       `json:"dueDate"`
	Completed bool           `json:"completed"`
}

type Service struct {
	notifications []Notification
}

func NewService() *Service {
	return &Service{
		notifications: make([]Notification, 0),
	}
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
				Type:    PregnancyMilestone,
				HorseID: int64(horse.ID),
				Message: message,
				DueDate: time.Now().AddDate(0, 0, days-daysPregnant),
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
			Type:    VaccinationDue,
			HorseID: int64(horse.ID),
			Message: "Annual vaccination due",
			DueDate: time.Now().AddDate(0, 0, 7),
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
	for i := range s.notifications {
		if s.notifications[i].ID == id {
			s.notifications[i].Completed = true
			logger.Info("Marked notification as complete", map[string]interface{}{
				"notificationID": id,
			})
			return nil
		}
	}
	return nil
}
