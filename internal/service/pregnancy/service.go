package pregnancy

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type Service struct {
	db *database.PostgresDB
}

type MonitoringSchedule struct {
	CheckFrequency    int  // hours between checks
	TemperatureCheck  bool
	BehaviorCheck     bool
	UdderCheck        bool
	VulvaCheck        bool
	Priority          string
}

func NewService(db *database.PostgresDB) *Service {
	return &Service{db: db}
}

func (s *Service) GetPregnancies(ctx context.Context) ([]models.Pregnancy, error) {
	userID := ctx.Value(middleware.UserIDKey).(string)
	return s.db.GetPregnancies(userID)
}

func (s *Service) GetPregnancyStage(ctx context.Context, pregnancyID int64) (models.PregnancyStage, error) {
	pregnancy, err := s.db.GetPregnancy(pregnancyID)
	if err != nil {
		return "", fmt.Errorf("failed to get pregnancy: %w", err)
	}
	return DeterminePregnancyStage(pregnancy.StartDate, time.Now()), nil
}

func (s *Service) StartPregnancyTracking(ctx context.Context, horseID int64, startDate time.Time) error {
	pregnancy := &models.Pregnancy{
		HorseID:   uint(horseID),
		StartDate: startDate,
		Status:    models.PregnancyStatusActive,
		UserID:    ctx.Value(middleware.UserIDKey).(string),
	}
	return s.db.AddPregnancy(pregnancy)
}

func (s *Service) EndPregnancyTracking(ctx context.Context, pregnancyID int64, outcome string) error {
	pregnancy, err := s.db.GetPregnancy(pregnancyID)
	if err != nil {
		return fmt.Errorf("failed to get pregnancy: %w", err)
	}

	pregnancy.Status = outcome
	pregnancy.EndDate = &time.Time{}
	return s.db.UpdatePregnancy(&pregnancy)
}

func (s *Service) GetPregnancyGuidelinesByStage(stage models.PregnancyStage) map[string]string {
	guidelines := make(map[string]string)
	
	switch stage {
	case models.EarlyGestation:
		guidelines["Feeding"] = "Maintain normal feeding routine, focus on quality forage"
		guidelines["Exercise"] = "Continue regular exercise program if mare is accustomed"
		guidelines["Monitoring"] = "Schedule initial pregnancy check, watch for signs of discomfort"
		
	case models.MidGestation:
		guidelines["Feeding"] = "Gradually increase feed quality, maintain body condition"
		guidelines["Exercise"] = "Moderate exercise, avoid strenuous activities"
		guidelines["Monitoring"] = "Regular vet checks, monitor weight gain"
		
	case models.LateGestation:
		guidelines["Feeding"] = "Increase feed quantity, ensure proper nutrition"
		guidelines["Exercise"] = "Light exercise only, daily turnout"
		guidelines["Monitoring"] = "Weekly vet checks, watch for pre-foaling signs"
		
	case models.Foaling:
		guidelines["Feeding"] = "Feed small meals frequently, maintain hydration"
		guidelines["Exercise"] = "Minimal exercise, supervised turnout only"
		guidelines["Monitoring"] = "24-hour monitoring, prepare for foaling"
	}
	
	return guidelines
}

func (s *Service) AddPregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	userID := ctx.Value(middleware.UserIDKey).(string)
	pregnancy.UserID = userID
	return s.db.AddPregnancy(pregnancy)
}

func (s *Service) UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	return s.db.UpdatePregnancy(pregnancy)
}

func (s *Service) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return s.db.AddPregnancyEvent(event)
}

// Add other service methods as needed
