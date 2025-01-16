package repository

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

// ChecklistRepository handles checklist data persistence
type ChecklistRepository struct {
	db *gorm.DB
}

// NewChecklistRepository creates a new checklist repository
func NewChecklistRepository(db *gorm.DB) *ChecklistRepository {
	return &ChecklistRepository{db: db}
}

// GetChecklist retrieves all checklist items for a horse
func (r *ChecklistRepository) GetChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	var items []models.PreFoalingChecklistItem
	result := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&items)
	return items, result.Error
}

// SaveChecklist saves multiple checklist items
func (r *ChecklistRepository) SaveChecklist(ctx context.Context, items []models.PreFoalingChecklistItem) error {
	return r.db.WithContext(ctx).Create(&items).Error
}

// UpdateItem updates a single checklist item
func (r *ChecklistRepository) UpdateItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// GetSeasonalTemplate retrieves season-specific checklist template
func (r *ChecklistRepository) GetSeasonalTemplate(ctx context.Context, season models.Season) (*models.SeasonalChecklistTemplate, error) {
	template := &models.SeasonalChecklistTemplate{
		Season: season,
		Items:  getSeasonalItems(season),
	}
	return template, nil
}

// getSeasonalItems returns checklist items specific to a season
func getSeasonalItems(season models.Season) []models.PreFoalingChecklistItem {
	now := time.Now()
	
	switch season {
	case models.SeasonWinter:
		return []models.PreFoalingChecklistItem{
			{
				Description: "Check stable insulation",
				Priority:    models.PriorityHigh,
				Notes:      "Ensure proper temperature regulation for winter foaling",
				Category:   string(models.CategoryEnvironment),
				IsRequired: true,
				Season:     models.SeasonWinter,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				Description: "Prepare heated water system",
				Priority:    models.PriorityHigh,
				Notes:      "Prevent water freezing and ensure constant access",
				Category:   string(models.CategoryEnvironment),
				IsRequired: true,
				Season:     models.SeasonWinter,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}
	case models.SeasonSummer:
		return []models.PreFoalingChecklistItem{
			{
				Description: "Install cooling systems",
				Priority:    models.PriorityHigh,
				Notes:      "Set up fans and ventilation for summer heat",
				Category:   string(models.CategoryEnvironment),
				IsRequired: true,
				Season:     models.SeasonSummer,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				Description: "Prepare electrolyte supplements",
				Priority:    models.PriorityMedium,
				Notes:      "Stock up for hot weather hydration",
				Category:   string(models.CategoryNutrition),
				IsRequired: false,
				Season:     models.SeasonSummer,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}
	case models.SeasonSpring:
		return []models.PreFoalingChecklistItem{
			{
				Description: "Check pasture fencing",
				Priority:    models.PriorityHigh,
				Notes:      "Ensure safe turnout areas for spring grass access",
				Category:   string(models.CategoryEnvironment),
				IsRequired: true,
				Season:     models.SeasonSpring,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				Description: "Monitor grass intake",
				Priority:    models.PriorityMedium,
				Notes:      "Prevent spring grass-related issues",
				Category:   string(models.CategoryNutrition),
				IsRequired: false,
				Season:     models.SeasonSpring,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}
	case models.SeasonFall:
		return []models.PreFoalingChecklistItem{
			{
				Description: "Check blanket condition",
				Priority:    models.PriorityMedium,
				Notes:      "Prepare for temperature fluctuations",
				Category:   string(models.CategorySupplies),
				IsRequired: false,
				Season:     models.SeasonFall,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			{
				Description: "Stock winter hay supply",
				Priority:    models.PriorityHigh,
				Notes:      "Ensure adequate nutrition for winter months",
				Category:   string(models.CategoryNutrition),
				IsRequired: true,
				Season:     models.SeasonFall,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		}
	default:
		return []models.PreFoalingChecklistItem{}
	}
}
