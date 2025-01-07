package postgres

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type PregnancyRepository struct {
	db *gorm.DB
}

func NewPregnancyRepository(db *gorm.DB) *PregnancyRepository {
	return &PregnancyRepository{db: db}
}

func (r *PregnancyRepository) Create(ctx context.Context, pregnancy *models.Pregnancy) error {
	return r.db.WithContext(ctx).Create(pregnancy).Error
}

func (r *PregnancyRepository) GetByID(ctx context.Context, id uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.WithContext(ctx).First(&pregnancy, id).Error
	if err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

func (r *PregnancyRepository) Update(ctx context.Context, pregnancy *models.Pregnancy) error {
	return r.db.WithContext(ctx).Save(pregnancy).Error
}

func (r *PregnancyRepository) GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).First(&pregnancy).Error
	if err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

func (r *PregnancyRepository) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PregnancyRepository) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	var items []models.PreFoalingChecklistItem
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&items).Error
	return items, err
}

func (r *PregnancyRepository) GetPreFoalingChecklistItem(ctx context.Context, itemID uint) (*models.PreFoalingChecklistItem, error) {
	var item models.PreFoalingChecklistItem
	err := r.db.WithContext(ctx).First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PregnancyRepository) DeletePreFoalingChecklistItem(ctx context.Context, itemID uint) error {
	return r.db.WithContext(ctx).Delete(&models.PreFoalingChecklistItem{}, itemID).Error
}

func (r *PregnancyRepository) GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error) {
	return r.GetByID(ctx, id)
}

func (r *PregnancyRepository) GetByUserID(ctx context.Context, userID string) ([]models.Pregnancy, error) {
	var pregnancies []models.Pregnancy
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&pregnancies).Error
	return pregnancies, err
}

func (r *PregnancyRepository) InitializePreFoalingChecklist(ctx context.Context, horseID uint) error {
	// Implementation depends on your business logic
	return nil
}

func (r *PregnancyRepository) GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.WithContext(ctx).Where("horse_id = ? AND status = ?", horseID, models.PregnancyStatusActive).First(&pregnancy).Error
	if err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

func (r *PregnancyRepository) UpdatePregnancyStatus(ctx context.Context, horseID uint, isPregnant bool, conceptionDate *time.Time) error {
	return r.db.WithContext(ctx).Model(&models.Horse{}).Where("id = ?", horseID).
		Updates(map[string]interface{}{
			"is_pregnant":     isPregnant,
			"conception_date": conceptionDate,
		}).Error
}

func (r *PregnancyRepository) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return r.db.WithContext(ctx).Create(sign).Error
}

func (r *PregnancyRepository) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *PregnancyRepository) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	var events []models.PregnancyEvent
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&events).Error
	return events, err
}

func (r *PregnancyRepository) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	var signs []models.PreFoalingSign
	err := r.db.WithContext(ctx).Where("horse_id = ?", horseID).Find(&signs).Error
	return signs, err
}

func (r *PregnancyRepository) UpdatePreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return r.db.WithContext(ctx).Save(item).Error
} 