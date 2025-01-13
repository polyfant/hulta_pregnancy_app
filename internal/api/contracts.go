package api

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
)

// APIHandlerConfig defines the configuration for creating a new handler
type APIHandlerConfig struct {
	Database         *database.PostgresDB
	UserService      service.UserService
	HorseService     service.HorseService
	PregnancyService service.PregnancyService
	HealthService    service.HealthService
	BreedingService  service.BreedingService
	Cache            *cache.MemoryCache
	HorseRepo        repository.HorseRepository
	BreedingRepo     repository.BreedingRepository
}

// HealthService defines the interface for health-related operations
type HealthService interface {
	CreateRecord(ctx context.Context, record *models.HealthRecord) error
	GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
	UpdateRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteRecord(ctx context.Context, id uint) error
}

// PregnancyService defines the interface for pregnancy-related operations
type PregnancyService interface {
	GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error
	GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error)
	GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error
	GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]models.Guideline, error)
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
