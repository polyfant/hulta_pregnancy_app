// Package service provides centralized service interfaces for the horse tracking application
package service

import (
	"context"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// HealthService defines the interface for health-related operations
type HealthService interface {
	GetByID(ctx context.Context, id uint) (*models.HealthRecord, error)
	CreateRecord(ctx context.Context, record *models.HealthRecord) error
	UpdateRecord(ctx context.Context, record *models.HealthRecord) error
	DeleteRecord(ctx context.Context, id uint) error
	GetRecords(ctx context.Context, horseID uint) ([]models.HealthRecord, error)
}

// WeatherService defines the interface for weather-related operations
type WeatherService interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*models.WeatherData, error)
	GetLatestWeatherData(ctx context.Context) (*models.WeatherData, error)
	GetPregnancyWeatherAdvice(ctx context.Context, stage string) (*models.PregnancyWeatherAdvice, error)
	GetRecommendations(ctx context.Context) ([]*models.WeatherRecommendation, error)
	GetLatestImpact(ctx context.Context, horseID uint) (*models.WeatherImpact, error)
	CalculateWeatherImpact(ctx context.Context, horseID uint, weatherData *models.WeatherData) (*models.WeatherImpact, error)
}

// PregnancyService defines the interface for pregnancy management
type PregnancyService interface {
	GetPregnancy(ctx context.Context, id uint) (*models.Pregnancy, error)
	GetByHorseID(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	GetByUserID(ctx context.Context, userID string) ([]*models.Pregnancy, error)
	Create(ctx context.Context, pregnancy *models.Pregnancy) error
	Update(ctx context.Context, pregnancy *models.Pregnancy) error
	StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error
	GetEvents(ctx context.Context, horseID uint) ([]*models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error
	GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]string, error)
	GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error)
	Remove(ctx context.Context, id uint) error
	CalculateProgress(ctx context.Context, horseID uint) (float64, error)
	GetStageInfo(ctx context.Context, horseID uint) (*models.PregnancyStageInfo, error)
}

// VitalsService defines the interface for vital signs monitoring
type VitalsService interface {
	RecordVitalSigns(ctx context.Context, vitals *models.VitalSigns) error
	GetLatestVitalSigns(ctx context.Context, horseID uint) (*models.VitalSigns, error)
	GetVitalSignsTrend(ctx context.Context, horseID uint) ([]*models.VitalSigns, error)
	CheckForAlerts(ctx context.Context, vitals *models.VitalSigns) (*models.VitalSignsAlert, error)
	IsInLateStage(ctx context.Context, horseID uint) (bool, error)
}

// UserService defines the interface for user management
type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID string) error
	GetProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateProfile(ctx context.Context, user *models.User) error
	GetDashboardStats(ctx context.Context, userID string) (*models.DashboardStats, error)
}

// HorseService defines the interface for horse management
type HorseService interface {
	GetByID(ctx context.Context, id uint) (*models.Horse, error)
	Create(ctx context.Context, horse *models.Horse) error
	Update(ctx context.Context, horse *models.Horse) error
	Delete(ctx context.Context, id uint) error
	GetPregnantHorses(ctx context.Context, userID string) ([]*models.Horse, error)
	GetHorsesByUser(ctx context.Context, userID string) ([]*models.Horse, error)
	GetOffspring(ctx context.Context, horseID uint) ([]*models.Horse, error)
	GetFamilyTree(ctx context.Context, horseID uint) (*models.FamilyTree, error)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	SendWeatherAlert(ctx context.Context, userID string, alert string) error
	SendVitalSignsAlert(ctx context.Context, userID string, alert *models.VitalSignsAlert) error
	SendPregnancyAlert(ctx context.Context, userID string, alert *models.PregnancyAlert) error
	GetUserNotifications(ctx context.Context, userID string) ([]*models.Notification, error)
	MarkNotificationRead(ctx context.Context, notificationID uint) error
}

// BreedingService defines the interface for breeding management
type BreedingService interface {
	GetByID(ctx context.Context, id uint) (*models.BreedingRecord, error)
	CreateRecord(ctx context.Context, record *models.BreedingRecord) error
	UpdateRecord(ctx context.Context, record *models.BreedingRecord) error
	DeleteRecord(ctx context.Context, id uint) error
	GetRecords(ctx context.Context, horseID uint) ([]*models.BreedingRecord, error)
	GetCosts(ctx context.Context, horseID uint) ([]models.BreedingCost, error)
}

// PrivacyService defines the interface for privacy management
type PrivacyService interface {
	GetPrivacyPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error)
	UpdatePrivacyPreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error
	DeleteUserData(ctx context.Context, userID string, dataType string) error
	DeleteExpiredData(ctx context.Context) error
	CheckFeatureEnabled(ctx context.Context, userID string, feature string) (bool, error)
}

// RecommendationService defines the interface for product recommendations
type RecommendationService interface {
	GetRecommendations(ctx context.Context, horseID uint) ([]*models.Recommendation, error)
	GetPersonalizedRecommendations(ctx context.Context, userID string) ([]*models.Recommendation, error)
	GetSeasonalRecommendations(ctx context.Context) ([]*models.Recommendation, error)
	GetHealthBasedRecommendations(ctx context.Context, healthRecord *models.HealthRecord) ([]*models.Recommendation, error)
}

// ChecklistService defines the interface for checklist-related operations
type ChecklistService interface {
	GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]*models.PreFoalingChecklistItem, error)
	UpdateItemStatus(ctx context.Context, item *models.PreFoalingChecklistItem) error
	GetProgress(ctx context.Context, horseID uint) (*models.ChecklistProgress, error)
}
