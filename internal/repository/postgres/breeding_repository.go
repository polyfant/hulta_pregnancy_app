package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
)

type breedingRepository struct {
	db *gorm.DB
}

func (r *breedingRepository) GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error) {
	var records []models.BreedingRecord
	if err := r.db.Where("horse_id = ?", horseID).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get breeding records: %w", err)
	}
	return records, nil
}

func (r *breedingRepository) GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error) {
	var costs []models.BreedingCost
	if err := r.db.Where("horse_id = ?", horseID).Find(&costs).Error; err != nil {
		return nil, fmt.Errorf("failed to get breeding costs: %w", err)
	}
	return costs, nil
}

func (r *breedingRepository) CreateCost(ctx context.Context, cost *models.BreedingCost) error {
	if err := r.db.Create(cost).Error; err != nil {
		return fmt.Errorf("failed to create breeding cost: %w", err)
	}
	return nil
}

func (r *breedingRepository) UpdateCost(ctx context.Context, cost *models.BreedingCost) error {
	if err := r.db.Save(cost).Error; err != nil {
		return fmt.Errorf("failed to update breeding cost: %w", err)
	}
	return nil
}

func (r *breedingRepository) DeleteCost(ctx context.Context, id uint) error {
	if err := r.db.Delete(&models.BreedingCost{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete breeding cost: %w", err)
	}
	return nil
}

func (r *breedingRepository) GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error) {
	var events []models.PregnancyEvent
	if err := r.db.Where("horse_id = ?", horseID).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get pregnancy events: %w", err)
	}
	return events, nil
}

func (r *breedingRepository) CreateRecord(ctx context.Context, record *models.BreedingRecord) error {
	if err := r.db.Create(record).Error; err != nil {
		return fmt.Errorf("failed to create breeding record: %w", err)
	}
	return nil
}

func (r *breedingRepository) GetBreedingStatus(ctx context.Context, horseID uint) (*models.BreedingStatus, error) {
	records, err := r.GetRecords(ctx, horseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get breeding records: %w", err)
	}

	status := &models.BreedingStatus{
		HorseID: horseID,
		LastBreedingDate: nil,
		IsBreeding: false,
	}

	if len(records) > 0 {
		lastRecord := records[len(records)-1]
		status.LastBreedingDate = &lastRecord.Date
		status.IsBreeding = lastRecord.Status == "ACTIVE"
	}

	return status, nil
}

func NewBreedingRepository(db *gorm.DB) repository.BreedingRepository {
	return &breedingRepository{
		db: db,
	}
} 