package breeding

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/middleware"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Pregnancy stage thresholds in days
const (
	EarlyPregnancyDays = 113
	MidPregnancyDays   = 226
	LatePregnancyDays  = 310
	OverdueDays        = 365
)

// Service handles breeding-related operations
type Service struct {
	db models.DataStore
}

// NewService creates a new breeding service instance
func NewService(db models.DataStore) *Service {
	return &Service{db: db}
}

// CalculatePregnancySuccessRate calculates the success rate of pregnancies
func (s *Service) CalculatePregnancySuccessRate(horses []models.Horse) float64 {
	if len(horses) == 0 {
		return 0.0
	}

	var pregnantCount, successfulBirths int
	for _, horse := range horses {
		if horse.ConceptionDate != nil {
			pregnantCount++
			events, err := s.db.GetPregnancyEvents(int64(horse.ID))
			if err != nil {
				logger.Error(err, "Failed to get pregnancy events", map[string]interface{}{
					"horseID": horse.ID,
				})
				continue
			}

			for _, event := range events {
				if strings.EqualFold(event.Type, models.EventFoaling) && event.Description == "SUCCESSFUL" {
					successfulBirths++
					break
				}
			}
		}
	}

	if pregnantCount == 0 {
		return 0.0
	}
	return float64(successfulBirths) / float64(pregnantCount) * 100
}

// GetUpcomingMilestones returns a list of upcoming pregnancy milestones
func (s *Service) GetUpcomingMilestones(horse models.Horse) []string {
	if horse.ConceptionDate == nil {
		return nil
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	var milestones []string

	// Define milestones with their days and descriptions
	milestoneDefs := []struct {
		days int
		desc string
	}{
		{14, "First ultrasound check"},
		{30, "Second ultrasound check"},
		{45, "Gender determination possible"},
		{60, "Fetal movement check"},
		{90, "Vaccination review"},
		{EarlyPregnancyDays, "End of early pregnancy stage"},
		{MidPregnancyDays, "End of mid pregnancy stage"},
		{LatePregnancyDays, "Prepare for foaling"},
	}

	for _, m := range milestoneDefs {
		if daysPregnant <= m.days {
			daysUntil := m.days - daysPregnant
			milestones = append(milestones, fmt.Sprintf("%s (in %d days)", m.desc, daysUntil))
		}
	}

	return milestones
}

// GetPregnancyStage determines the current stage of pregnancy
func (s *Service) GetPregnancyStage(horse models.Horse) models.PregnancyStage {
	if horse.ConceptionDate == nil {
		return ""
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)

	switch {
	case daysPregnant <= EarlyPregnancyDays:
		return models.EarlyGestation
	case daysPregnant <= MidPregnancyDays:
		return models.MidGestation
	case daysPregnant <= LatePregnancyDays:
		return models.LateGestation
	case daysPregnant <= OverdueDays:
		return models.PreFoaling
	default:
		return models.Foaling
	}
}

// GetHighBreedingCosts returns a list of high breeding costs
func (s *Service) GetHighBreedingCosts(horses []models.Horse, threshold float64, startDate, endDate time.Time) []struct {
	HorseName   string
	Description string
	Amount      float64
} {
	var highCosts []struct {
		HorseName   string
		Description string
		Amount      float64
	}

	for _, horse := range horses {
		costs, err := s.db.GetBreedingCosts(horse.ID)
		if err != nil {
			logger.Error(err, "Failed to get breeding costs", map[string]interface{}{
				"horseID": horse.ID,
			})
			continue
		}

		for _, cost := range costs {
			if cost.Amount >= threshold && cost.Date.After(startDate) && cost.Date.Before(endDate) {
				highCosts = append(highCosts, struct {
					HorseName   string
					Description string
					Amount      float64
				}{
					HorseName:   horse.Name,
					Description: cost.Description,
					Amount:      cost.Amount,
				})
			}
		}
	}

	sortByAmount(highCosts)
	return highCosts
}

// GetPregnantHorses returns a list of currently pregnant horses
func (s *Service) GetPregnantHorses(horses []models.Horse) []models.Horse {
	var pregnantHorses []models.Horse
	for _, horse := range horses {
		if horse.IsPregnant && horse.ConceptionDate != nil {
			pregnantHorses = append(pregnantHorses, horse)
		}
	}
	return pregnantHorses
}

// Helper functions
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func sortByAmount(costs []struct {
	HorseName   string
	Description string
	Amount      float64
}) {
	sort.Slice(costs, func(i, j int) bool {
		return costs[i].Amount > costs[j].Amount
	})
}

func (s *Service) ValidateBreeding(mare, stallion models.Horse) error {
	if mare.Gender != models.GenderMare {
		return fmt.Errorf("first horse must be a mare")
	}
	if stallion.Gender != models.GenderStallion {
		return fmt.Errorf("second horse must be a stallion")
	}
	if !mare.IsBreedingAge() {
		return fmt.Errorf("mare is not of breeding age")
	}
	if !stallion.IsBreedingAge() {
		return fmt.Errorf("stallion is not of breeding age")
	}
	return nil
}

func (s *Service) RecordBreeding(ctx context.Context, mareID, stallionID int64, date time.Time) error {
	stallionUint := uint(stallionID)
	record := &models.BreedingRecord{
		MareID:     uint(mareID),
		StallionID: &stallionUint,
		Date:       date,
		UserID:     ctx.Value(middleware.UserIDKey).(string),
	}
	return s.db.AddBreedingRecord(record)
}

func (s *Service) GetBreedingHistory(ctx context.Context, horseID int64) ([]models.BreedingRecord, error) {
	return s.db.GetBreedingRecords(horseID)
}
