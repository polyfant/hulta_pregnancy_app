package repository

import (
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

// PregnancyRepository handles database operations for pregnancies
type PregnancyRepository struct {
	db *gorm.DB
}

// NewPregnancyRepository creates a new pregnancy repository
func NewPregnancyRepository(db *gorm.DB) *PregnancyRepository {
	return &PregnancyRepository{
		db: db,
	}
}

// Create adds a new pregnancy record
func (r *PregnancyRepository) Create(pregnancy *models.Pregnancy) error {
	return r.db.Create(pregnancy).Error
}

// GetActivePregnancyByHorseID finds the active pregnancy for a horse
func (r *PregnancyRepository) GetActivePregnancyByHorseID(horseID uint) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.Where("horse_id = ? AND status = ?", horseID, models.PregnancyStatusActive).First(&pregnancy).Error
	if err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

// Update modifies an existing pregnancy record
func (r *PregnancyRepository) Update(pregnancy *models.Pregnancy) error {
	return r.db.Save(pregnancy).Error
}

// GetPregnantHorses retrieves all currently pregnant horses
func (r *PregnancyRepository) GetPregnantHorses() ([]models.Pregnancy, error) {
	var pregnancies []models.Pregnancy
	err := r.db.Where("status = ?", models.PregnancyStatusActive).Find(&pregnancies).Error
	return pregnancies, err
}

// GetPregnancy retrieves a pregnancy by ID
func (r *PregnancyRepository) GetPregnancy(pregnancyID int64) (*models.Pregnancy, error) {
	var pregnancy models.Pregnancy
	err := r.db.First(&pregnancy, pregnancyID).Error
	if err != nil {
		return nil, err
	}
	return &pregnancy, nil
}

// GetPregnancies retrieves all pregnancies for a user
func (r *PregnancyRepository) GetPregnancies(userID string) ([]models.Pregnancy, error) {
	var pregnancies []models.Pregnancy
	err := r.db.Where("user_id = ?", userID).Find(&pregnancies).Error
	return pregnancies, err
}

// AddPregnancy adds a new pregnancy record
func (r *PregnancyRepository) AddPregnancy(pregnancy *models.Pregnancy) error {
	return r.db.Create(pregnancy).Error
}

// UpdatePregnancy updates an existing pregnancy record
func (r *PregnancyRepository) UpdatePregnancy(pregnancy *models.Pregnancy) error {
	return r.db.Save(pregnancy).Error
}

// AddPregnancyEvent adds a new pregnancy event
func (r *PregnancyRepository) AddPregnancyEvent(event *models.PregnancyEvent) error {
	return r.db.Create(event).Error
}
