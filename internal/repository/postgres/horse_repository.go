package postgres

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

type HorseRepository struct {
	db *gorm.DB
}

func NewHorseRepository(db *gorm.DB) *HorseRepository {
	return &HorseRepository{db: db}
}

func (r *HorseRepository) Create(ctx context.Context, horse *models.Horse) error {
	return r.db.WithContext(ctx).Create(horse).Error
}

func (r *HorseRepository) GetByID(ctx context.Context, id uint) (*models.Horse, error) {
	var horse models.Horse
	err := r.db.WithContext(ctx).First(&horse, id).Error
	if err != nil {
		return nil, err
	}
	return &horse, nil
}

func (r *HorseRepository) Update(ctx context.Context, horse *models.Horse) error {
	return r.db.WithContext(ctx).Save(horse).Error
}

func (r *HorseRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Horse{}, id).Error
}

func (r *HorseRepository) ListByUser(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&horses).Error
	return horses, err
}

func (r *HorseRepository) GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error) {
	var horses []models.Horse
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_pregnant = true", userID).Find(&horses).Error
	return horses, err
}

func (r *HorseRepository) GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error) {
	var horse models.Horse
	if err := r.db.WithContext(ctx).First(&horse, horseID).Error; err != nil {
		return nil, err
	}

	tree := &models.FamilyTree{
		Horse: &horse,
	}

	// Get mother if exists
	if horse.MotherId != nil {
		var mother models.Horse
		if err := r.db.WithContext(ctx).First(&mother, *horse.MotherId).Error; err == nil {
			tree.Mother = &mother
		}
	}

	// Get father if exists
	if horse.FatherId != nil {
		var father models.Horse
		if err := r.db.WithContext(ctx).First(&father, *horse.FatherId).Error; err == nil {
			tree.Father = &father
		}
	}

	// Get offspring
	var offspring []*models.Horse
	if err := r.db.WithContext(ctx).Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&offspring).Error; err != nil {
		return nil, err
	}
	tree.Offspring = offspring

	return tree, nil
}

func (r *HorseRepository) GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error) {
	var offspring []models.Horse
	err := r.db.WithContext(ctx).Where("mother_id = ? OR father_id = ?", horseID, horseID).Find(&offspring).Error
	return offspring, err
} 