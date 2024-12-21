package breeding

import (
	"sort"
	"strings"
	"time"

	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/polyfant/horse_tracking/internal/logger"
)

const (
	EarlyPregnancyDays = 113
	MidPregnancyDays   = 226
	LatePregnancyDays  = 310
	OverdueDays        = 365
)

type BreedingService struct {
	db models.DataStore
}

type Service struct {
	db models.DataStore
}

func NewService(db models.DataStore) *Service {
	return &Service{db: db}
}

func NewBreedingService(db models.DataStore) *BreedingService {
	return &BreedingService{db: db}
}

func (s *BreedingService) CalculatePregnancySuccessRate(horses []models.Horse) float64 {
	if len(horses) == 0 {
		return 0.0
	}

	var pregnantCount, successfulBirths int
	for _, horse := range horses {
		if horse.ConceptionDate != nil {
			pregnantCount++
			// Check pregnancy events for successful birth
			events, err := s.db.GetPregnancyEvents(horse.ID)
			if err != nil {
				logger.Error(err, "Failed to get pregnancy events", map[string]interface{}{
					"horseID": horse.ID,
				})
				continue
			}
			for _, event := range events {
				if containsIgnoreCase(event.Description, "successful birth") {
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

func (s *BreedingService) GetUpcomingMilestones(horse models.Horse) []string {
	if horse.ConceptionDate == nil {
		return nil
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	var milestones []string

	// Define milestones with their day ranges
	milestonePeriods := []struct {
		start, end int
		message    string
	}{
		{80, 100, "First trimester check"},
		{130, 145, "Vaccination booster"},
		{145, 170, "Begin increasing feed"},
		{270, 300, "Check for udder development"},
		{310, 330, "Prepare foaling area"},
		{330, 365, "Monitor for signs of imminent foaling"},
	}

	for _, period := range milestonePeriods {
		if daysPregnant >= period.start && daysPregnant < period.end {
			milestones = append(milestones, period.message)
		}
	}

	return milestones
}

func (s *BreedingService) GetPregnancyStage(horse models.Horse) models.PregnancyStage {
	if horse.ConceptionDate == nil {
		return ""
	}

	daysPregnant := int(time.Since(*horse.ConceptionDate).Hours() / 24)

	switch {
	case daysPregnant < 114:
		return models.EarlyGestation
	case daysPregnant < 226:
		return models.MidGestation
	case daysPregnant < 310:
		return models.LateGestation
	default:
		return models.FinalGestation
	}
}

func (s *BreedingService) GetHighBreedingCosts(horses []models.Horse, threshold float64, startDate, endDate time.Time) []struct {
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
			if cost.Amount > threshold && !cost.Date.Before(startDate) && !cost.Date.After(endDate) {
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

	// Sort by amount descending and limit to top 5
	sortByAmount(highCosts)
	if len(highCosts) > 5 {
		highCosts = highCosts[:5]
	}

	return highCosts
}

func (s *BreedingService) GetPregnantHorses(horses []models.Horse) []models.Horse {
	var pregnant []models.Horse
	for _, horse := range horses {
		if horse.ConceptionDate != nil {
			pregnant = append(pregnant, horse)
		}
	}
	return pregnant
}

func containsIgnoreCase(s, substr string) bool {
	s, substr = strings.ToLower(s), strings.ToLower(substr)
	return strings.Contains(s, substr)
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
