package checklist

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// Service handles checklist-related operations
type Service struct {
	repo Repository
}

// Repository defines the interface for checklist data operations
type Repository interface {
	GetChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error)
	SaveChecklist(ctx context.Context, items []models.PreFoalingChecklistItem) error
	UpdateItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
	GetSeasonalTemplate(ctx context.Context, season models.Season) (*models.SeasonalChecklistTemplate, error)
}

// NewService creates a new checklist service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// getCurrentSeason determines the current season based on the date
func getCurrentSeason(date time.Time) models.Season {
	month := date.Month()
	switch {
	case month >= 3 && month <= 5:
		return models.SeasonSpring
	case month >= 6 && month <= 8:
		return models.SeasonSummer
	case month >= 9 && month <= 11:
		return models.SeasonFall
	default:
		return models.SeasonWinter
	}
}

// GetSeasonSpecificItems returns checklist items specific to the current season
func (s *Service) GetSeasonSpecificItems(ctx context.Context, dueDate time.Time) ([]models.PreFoalingChecklistItem, error) {
	season := getCurrentSeason(dueDate)
	template, err := s.repo.GetSeasonalTemplate(ctx, season)
	if err != nil {
		return nil, err
	}
	return template.Items, nil
}

// GenerateChecklist creates a complete checklist including season-specific items
func (s *Service) GenerateChecklist(ctx context.Context, horseID uint, dueDate time.Time) ([]models.PreFoalingChecklistItem, error) {
	// Get base checklist items
	baseItems := getBaseChecklistItems(horseID, dueDate)
	
	// Get season-specific items
	seasonalItems, err := s.GetSeasonSpecificItems(ctx, dueDate)
	if err != nil {
		return nil, err
	}

	// Combine and return all items
	return append(baseItems, seasonalItems...), nil
}

// getBaseChecklistItems returns the standard checklist items
func getBaseChecklistItems(horseID uint, dueDate time.Time) []models.PreFoalingChecklistItem {
	return []models.PreFoalingChecklistItem{
		{
			HorseID:     horseID,
			Description: "Schedule pre-foaling veterinary exam",
			Priority:    models.PriorityHigh,
			DueDate:     dueDate.AddDate(0, 0, -60),
			Notes:       "Check vaccination status, overall health, and pregnancy progress",
			Category:    string(models.CategoryVeterinary),
			IsRequired:  true,
		},
		{
			HorseID:     horseID,
			Description: "Prepare foaling kit",
			Priority:    models.PriorityHigh,
			DueDate:     dueDate.AddDate(0, 0, -30),
			Notes:       "Include sterile supplies, emergency contacts, and monitoring equipment",
			Category:    string(models.CategorySupplies),
			IsRequired:  true,
		},
		// Add more base items...
	}
}

// GetProgress calculates the current progress of the checklist
func (s *Service) GetProgress(ctx context.Context, horseID uint) (*models.ChecklistProgress, error) {
	items, err := s.repo.GetChecklist(ctx, horseID)
	if err != nil {
		return nil, err
	}

	progress := &models.ChecklistProgress{}
	progress.TotalItems = len(items)
	
	for _, item := range items {
		if item.Completed {
			progress.CompletedItems++
			if item.IsRequired {
				progress.CompletedRequired++
			}
		}
		if item.IsRequired {
			progress.RequiredItems++
		}
	}

	if progress.TotalItems > 0 {
		progress.Progress = float64(progress.CompletedItems) / float64(progress.TotalItems) * 100
	}

	return progress, nil
}

// UpdateItemStatus updates the completion status of a checklist item
func (s *Service) UpdateItemStatus(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	item.UpdatedAt = time.Now()
	return s.repo.UpdateItem(ctx, item)
}
