package database

import (
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	db *gorm.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&models.Horse{},
		&models.HealthRecord{},
		&models.PregnancyEvent{},
		&models.PreFoalingSign{},
		&models.PreFoalingChecklistItem{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &PostgresDB{db: db}, nil
}

// Horse methods
func (p *PostgresDB) GetHorse(id int64) (models.Horse, error) {
	var horse models.Horse
	if err := p.db.First(&horse, id).Error; err != nil {
		return horse, fmt.Errorf("failed to get horse: %w", err)
	}
	return horse, nil
}

func (p *PostgresDB) GetAllHorses() ([]models.Horse, error) {
	var horses []models.Horse
	if err := p.db.Find(&horses).Error; err != nil {
		return nil, fmt.Errorf("failed to get horses: %w", err)
	}
	return horses, nil
}

func (p *PostgresDB) AddHorse(horse *models.Horse) error {
	if err := p.db.Create(horse).Error; err != nil {
		return fmt.Errorf("failed to add horse: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdateHorse(horse *models.Horse) error {
	if err := p.db.Save(horse).Error; err != nil {
		return fmt.Errorf("failed to update horse: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeleteHorse(id int64) error {
	if err := p.db.Delete(&models.Horse{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete horse: %w", err)
	}
	return nil
}

// Health record methods
func (p *PostgresDB) GetHealthRecords(horseID int64) ([]models.HealthRecord, error) {
	var records []models.HealthRecord
	if err := p.db.Where("horse_id = ?", horseID).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get health records: %w", err)
	}
	return records, nil
}

func (p *PostgresDB) AddHealthRecord(record *models.HealthRecord) error {
	if err := p.db.Create(record).Error; err != nil {
		return fmt.Errorf("failed to add health record: %w", err)
	}
	return nil
}

// Pregnancy methods
func (p *PostgresDB) GetPregnancyEvents(horseID int64) ([]models.PregnancyEvent, error) {
	var events []models.PregnancyEvent
	if err := p.db.Where("horse_id = ?", horseID).Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to get pregnancy events: %w", err)
	}
	return events, nil
}

func (p *PostgresDB) AddPregnancyEvent(event *models.PregnancyEvent) error {
	if err := p.db.Create(event).Error; err != nil {
		return fmt.Errorf("failed to add pregnancy event: %w", err)
	}
	return nil
}

// Pre-foaling methods
func (p *PostgresDB) GetPreFoalingSigns(horseID int64) ([]models.PreFoalingSign, error) {
	var signs []models.PreFoalingSign
	if err := p.db.Where("horse_id = ?", horseID).Find(&signs).Error; err != nil {
		return nil, fmt.Errorf("failed to get pre-foaling signs: %w", err)
	}
	return signs, nil
}

func (p *PostgresDB) AddPreFoalingSign(sign *models.PreFoalingSign) error {
	if err := p.db.Create(sign).Error; err != nil {
		return fmt.Errorf("failed to add pre-foaling sign: %w", err)
	}
	return nil
}

// Pre-foaling checklist methods
func (p *PostgresDB) GetPreFoalingChecklist(horseID int64) ([]models.PreFoalingChecklistItem, error) {
	var items []models.PreFoalingChecklistItem
	if err := p.db.Where("horse_id = ?", horseID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get pre-foaling checklist: %w", err)
	}
	return items, nil
}

func (p *PostgresDB) AddPreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	if err := p.db.Create(item).Error; err != nil {
		return fmt.Errorf("failed to add pre-foaling checklist item: %w", err)
	}
	return nil
}

func (p *PostgresDB) UpdatePreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error {
	if err := p.db.Save(item).Error; err != nil {
		return fmt.Errorf("failed to update pre-foaling checklist item: %w", err)
	}
	return nil
}

func (p *PostgresDB) DeletePreFoalingChecklistItem(id int64) error {
	if err := p.db.Delete(&models.PreFoalingChecklistItem{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete pre-foaling checklist item: %w", err)
	}
	return nil
} 