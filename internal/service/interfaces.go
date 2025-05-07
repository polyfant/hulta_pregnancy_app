package service

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HorseService defines the interface for horse-related operations
type HorseService interface {
	GetByID(ctx context.Context, id uint) (*models.Horse, error)
	Create(ctx context.Context, horse *models.Horse) error
	Update(ctx context.Context, horse *models.Horse) error
	Delete(ctx context.Context, id uint) error
	GetPregnant(ctx context.Context, userID string) ([]models.Horse, error)
	GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error)
	GetOffspring(ctx context.Context, horseID uint) ([]models.Horse, error)
	GetPregnantHorses(ctx context.Context, userID string) ([]models.Horse, error)
	ListByUserID(ctx context.Context, userID string) ([]models.Horse, error)
}

// HealthService defines the interface for health-related operations
type HealthService interface {
	CreateRecord(ctx context.Context, record *models.HealthRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
	UpdateRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteRecord(ctx context.Context, id uint) error
	GetRecordByID(ctx context.Context, id uint) (*models.HealthRecord, error)
}

// PregnancyService defines the interface for pregnancy-related operations
type PregnancyService interface {
	GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	GetPregnancies(ctx context.Context, userID string) ([]models.Pregnancy, error)
	StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error
	GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error)
	GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, userID string, horseID uint, eventInput *models.PregnancyEventInputDTO) (*models.PregnancyEvent, error)
	GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]models.Guideline, error)
	GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error)
	GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error)
	EndPregnancy(ctx context.Context, horseID uint, status string, date time.Time) error
	GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error)
	AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error
	UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error
	GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
	GetEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	GetPregnancyByID(ctx context.Context, pregnancyID uint) (*models.Pregnancy, error)
}

// UserService defines the interface for user-related operations
type UserService interface {
	GetByID(ctx context.Context, userID string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error)
	GetProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, user *models.User) error
}

// BreedingService defines the interface for breeding-related operations
type BreedingService interface {
	CreateRecord(ctx context.Context, record *models.BreedingRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.BreedingRecord, error)
	UpdateRecord(ctx context.Context, record *models.BreedingRecord) error
	DeleteRecord(ctx context.Context, id uint) error
	GetRecordByID(ctx context.Context, id uint) (*models.BreedingRecord, error)
}
