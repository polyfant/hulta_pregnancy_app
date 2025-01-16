package service

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

// PregnancyService defines the interface for pregnancy-related operations
type PregnancyService interface {
	GetPregnancy(ctx context.Context, horseID uint) (*models.Pregnancy, error)
	StartTracking(ctx context.Context, horseID uint, start models.PregnancyStart) error
	GetStatus(ctx context.Context, horseID uint) (*models.PregnancyStatus, error)
	GetPregnancyEvents(ctx context.Context, horseID uint) ([]models.PregnancyEvent, error)
	AddPregnancyEvent(ctx context.Context, event *models.PregnancyEvent) error
	GetGuidelines(ctx context.Context, stage models.PregnancyStage) ([]models.Guideline, error)
	GetActive(ctx context.Context, userID string) ([]models.Pregnancy, error)
	GetPregnancyStage(ctx context.Context, horseID uint) (models.PregnancyStage, error)
	EndPregnancy(ctx context.Context, horseID uint, status string, date time.Time) error
	GetPreFoalingSigns(ctx context.Context, horseID uint) ([]models.PreFoalingSign, error)
	AddPreFoalingSign(ctx context.Context, sign *models.PreFoalingSign) error
	UpdatePregnancy(ctx context.Context, pregnancy *models.Pregnancy) error
	GetPreFoalingChecklist(ctx context.Context, horseID uint) ([]models.PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(ctx context.Context, item *models.PreFoalingChecklistItem) error
}
