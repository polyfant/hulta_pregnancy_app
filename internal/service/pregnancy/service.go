package pregnancy

import (
	"context"
	"fmt"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type Service struct {
	pregnancyRepo repository.PregnancyRepository
}

type MonitoringSchedule struct {
	CheckFrequency    int  // hours between checks
	TemperatureCheck  bool
	BehaviorCheck     bool
	UdderCheck        bool
	VulvaCheck        bool
	Priority          string
}

func NewService(pregnancyRepo repository.PregnancyRepository) *Service {
	return &Service{
		pregnancyRepo: pregnancyRepo,
	}
}

func (s *Service) GetPregnancies(ctx context.Context) ([]models.Pregnancy, error) {
	userID := ctx.Value(middleware.UserIDKey).(string)
	return s.pregnancyRepo.GetPregnancies(userID)
}

func (s *Service) GetPregnancyStage(ctx context.Context, pregnancyID int64) (models.PregnancyStage, error) {
	pregnancy, err := s.pregnancyRepo.GetPregnancy(pregnancyID)
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
	return s.pregnancyRepo.AddPregnancy(pregnancy)
}

func (s *Service) EndPregnancyTracking(ctx context.Context, pregnancyID int64, outcome string) error {
	pregnancy, err := s.pregnancyRepo.GetPregnancy(pregnancyID)
	if err != nil {
		return fmt.Errorf("failed to get pregnancy: %w", err)
	}

	pregnancy.Status = outcome
	pregnancy.EndDate = &time.Time{}
	return s.pregnancyRepo.UpdatePregnancy(&pregnancy)
}

func (s *Service) GetPregnancyGuidelinesByStage(stage models.PregnancyStage) []models.PregnancyGuideline {
	var guidelines []models.PregnancyGuideline
	
	switch stage {
	case models.EarlyGestation:
		guidelines = []models.PregnancyGuideline{
			{Stage: stage, Category: "Feeding", Description: "Maintain normal feeding routine, focus on quality forage"},
			{Stage: stage, Category: "Exercise", Description: "Continue regular exercise program if mare is accustomed"},
			{Stage: stage, Category: "Monitoring", Description: "Schedule initial pregnancy check, watch for signs of discomfort"},
		}
		
	case models.MidGestation:
		guidelines = []models.PregnancyGuideline{
			{Stage: stage, Category: "Feeding", Description: "Gradually increase feed quality, maintain body condition"},
			{Stage: stage, Category: "Exercise", Description: "Moderate exercise, avoid strenuous activities"},
			{Stage: stage, Category: "Monitoring", Description: "Regular vet checks, monitor weight gain"},
		}
		
	case models.LateGestation:
		guidelines = []models.PregnancyGuideline{
			{Stage: stage, Category: "Feeding", Description: "Increase feed quantity, ensure proper nutrition"},
			{Stage: stage, Category: "Exercise", Description: "Light exercise only, daily turnout"},
			{Stage: stage, Category: "Monitoring", Description: "Weekly vet checks, watch for pre-foaling signs"},
		}
		
	case models.Foaling:
		guidelines = []models.PregnancyGuideline{
			{Stage: stage, Category: "Feeding", Description: "Feed small meals frequently, maintain hydration"},
			{Stage: stage, Category: "Exercise", Description: "Minimal exercise, supervised turnout only"},
			{Stage: stage, Category: "Monitoring", Description: "24-hour monitoring, prepare for foaling"},
		}
	}
	
	return guidelines
}

func (s *Service) AddPregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	userID := ctx.Value(middleware.UserIDKey).(string)
	pregnancy.UserID = userID
	return s.pregnancyRepo.AddPregnancy(pregnancy)
}

func (s *Service) UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error {
	return s.pregnancyRepo.UpdatePregnancy(pregnancy)
}

func (s *Service) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return s.pregnancyRepo.AddPregnancyEvent(event)
}

// Add other service methods as needed
