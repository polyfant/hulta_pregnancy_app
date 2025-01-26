package privacy

import (
	"context"
	"fmt"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetPrivacyPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error)
	UpdatePrivacyPreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error
	DeleteUserData(ctx context.Context, userID string, dataType string) error
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetPrivacyPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error) {
	return s.repo.GetPrivacyPreferences(ctx, userID)
}

func (s *Service) UpdatePrivacyPreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error {
	return s.repo.UpdatePrivacyPreferences(ctx, userID, prefs)
}

func (s *Service) DeleteUserData(ctx context.Context, userID string, dataType string) error {
	return s.repo.DeleteUserData(ctx, userID, dataType)
}

func (s *Service) DeleteExpiredData(ctx context.Context) error {
	// Get all users' privacy preferences
	prefs, err := s.repo.GetPrivacyPreferences(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get privacy preferences: %w", err)
	}

	if prefs.DataRetentionDays > 0 && prefs.AutoDeleteOldData {
		if err := s.repo.DeleteUserData(ctx, "", "expired"); err != nil {
			return fmt.Errorf("failed to delete expired data: %w", err)
		}
	}

	return nil
}
