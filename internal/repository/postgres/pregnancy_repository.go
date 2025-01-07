package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"gorm.io/gorm"
)

type pregnancyRepository struct {
	db *gorm.DB
}

func NewPregnancyRepository(db *gorm.DB) repository.PregnancyRepository {
	return &pregnancyRepository{
		db: db,
	}
}

func (r *pregnancyRepository) GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	if err := r.db.First(&pregnancy, id).Error; err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

func (r *pregnancyRepository) GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	if err := r.db.Where("horse_id = ?", horseID).First(&pregnancy).Error; err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

func (r *pregnancyRepository) Create(ctx context.Context, pregnancy *models.Pregnancy) error {
	return r.db.Create(pregnancy).Error
}

func (r *pregnancyRepository) Update(ctx context.Context, pregnancy *models.Pregnancy) error {
	return r.db.Save(pregnancy).Error
}

func (r *pregnancyRepository) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	var events []models.PregnancyEvent
	if err := r.db.Where("pregnancy_id IN (SELECT id FROM pregnancies WHERE horse_id = ?)", horseID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *pregnancyRepository) AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error {
	return r.db.Create(event).Error
}

func (r *pregnancyRepository) GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error) {
	var items []models.PreFoalingChecklistItem
	if err := r.db.Where("horse_id = ?", horseID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *pregnancyRepository) AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error {
	return r.db.Create(item).Error
}

func (r *pregnancyRepository) GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error) {
	var signs []models.PreFoalingSign
	if err := r.db.Where("horse_id = ?", horseID).Find(&signs).Error; err != nil {
		return nil, err
	}
	return signs, nil
}

func (r *pregnancyRepository) AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error {
	return r.db.Create(sign).Error
}

func (r *pregnancyRepository) GetCurrentPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.Where("horse_id = ? AND status = ?", horseID, models.PregnancyStatusActive).First(&pregnancy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get current pregnancy: %w", err)
	}
	return &pregnancy, nil
} 