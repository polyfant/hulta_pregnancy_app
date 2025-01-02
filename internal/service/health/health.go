package health

import (
	"sort"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/logger"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HealthSummary represents a comprehensive health summary for a horse
type HealthSummary struct {
	TotalRecords      int                  `json:"totalRecords"`
	LastCheckup       *time.Time           `json:"lastCheckup"`
	VaccinationStatus VaccinationStatus    `json:"vaccinationStatus"`
	RecentIssues      []models.HealthRecord `json:"recentIssues"`
}

type Service struct {
	db models.DataStore
}

func NewService(db models.DataStore) *Service {
	return &Service{db: db}
}

type HealthService struct {
	db models.DataStore
}

func NewHealthService(db models.DataStore) *HealthService {
	return &HealthService{db: db}
}

type VaccinationStatus struct {
	IsUpToDate bool
	LastDate   time.Time
	DueDate    time.Time
}
func (s *HealthService) GetVaccinationStatus(horse models.Horse) VaccinationStatus {
	records, err := s.db.GetHealthRecords(int64(horse.ID))
	if err != nil {
		logger.Error(err, "Failed to get health records", map[string]interface{}{
			"horseID": horse.ID,
		})
		return VaccinationStatus{IsUpToDate: false}
	}

	var lastVaccination time.Time
	for _, record := range records {
		if record.Type == "Vaccination" && record.Date.After(lastVaccination) {
			lastVaccination = record.Date
		}
	}

	if lastVaccination.IsZero() {
		return VaccinationStatus{
			IsUpToDate: false,
			DueDate:    time.Now(),
		}
	}

	// Vaccinations are due yearly
	dueDate := lastVaccination.AddDate(1, 0, 0)
	return VaccinationStatus{
		IsUpToDate: time.Now().Before(dueDate),
		LastDate:   lastVaccination,
		DueDate:    dueDate,
	}
}

func (s *HealthService) GetHealthSummary(horse models.Horse) HealthSummary {
	records, err := s.db.GetHealthRecords(int64(horse.ID))
	if err != nil {
		logger.Error(err, "Failed to get health records", map[string]interface{}{
			"horseID": horse.ID,
		})
		return HealthSummary{}
	}

	var lastCheckup *time.Time
	var recentIssues []models.HealthRecord
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	for _, record := range records {
		if record.Type == "Checkup" && (lastCheckup == nil || record.Date.After(*lastCheckup)) {
			lastCheckup = &record.Date
		}
		if record.Type == "Issue" && record.Date.After(threeMonthsAgo) {
			recentIssues = append(recentIssues, record)
		}
	}

	return HealthSummary{
		TotalRecords:      len(records),
		LastCheckup:       lastCheckup,
		VaccinationStatus: s.GetVaccinationStatus(horse),
		RecentIssues:      recentIssues,
	}
}

func (s *HealthService) GetUpcomingHealthChecks(horses []models.Horse) []struct {
	Horse    models.Horse
	DueDate  time.Time
	CheckType string
} {
	var upcoming []struct {
		Horse     models.Horse
		DueDate   time.Time
		CheckType string
	}

	for _, horse := range horses {
		// Check vaccinations
		vaccStatus := s.GetVaccinationStatus(horse)
		if !vaccStatus.IsUpToDate {
			upcoming = append(upcoming, struct {
				Horse     models.Horse
				DueDate   time.Time
				CheckType string
			}{
				Horse:     horse,
				DueDate:   vaccStatus.DueDate,
				CheckType: "Vaccination",
			})
		}

		// Check regular checkups
		records, err := s.db.GetHealthRecords(int64(horse.ID))
		if err != nil {
			logger.Error(err, "Failed to get health records", map[string]interface{}{
				"horseID": horse.ID,
			})
			continue
		}

		var lastCheckup time.Time
		for _, record := range records {
			if record.Type == "Checkup" && record.Date.After(lastCheckup) {
				lastCheckup = record.Date
			}
		}

		if lastCheckup.IsZero() || time.Since(lastCheckup) > 6*30*24*time.Hour {
			upcoming = append(upcoming, struct {
				Horse     models.Horse
				DueDate   time.Time
				CheckType string
			}{
				Horse:     horse,
				DueDate:   lastCheckup.AddDate(0, 6, 0),
				CheckType: "Regular Checkup",
			})
		}
	}

	// Sort by due date
	sortByDueDate(upcoming)
	return upcoming
}

func sortByDueDate(checks []struct {
	Horse     models.Horse
	DueDate   time.Time
	CheckType string
}) {
	sort.Slice(checks, func(i, j int) bool {
		return checks[i].DueDate.Before(checks[j].DueDate)
	})
}
