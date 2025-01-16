package repository

import (
	"context"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"gorm.io/gorm"
)

// FeedbackRepository handles feature request data persistence
type FeedbackRepository struct {
	db *gorm.DB
}

// NewFeedbackRepository creates a new feedback repository
func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db: db}
}

// GetFeatureRequests retrieves all feature requests
func (r *FeedbackRepository) GetFeatureRequests(ctx context.Context) ([]models.FeatureRequest, error) {
	var features []models.FeatureRequest
	result := r.db.WithContext(ctx).Order("vote_count DESC").Find(&features)
	return features, result.Error
}

// GetFeatureRequest retrieves a specific feature request
func (r *FeedbackRepository) GetFeatureRequest(ctx context.Context, id uint) (*models.FeatureRequest, error) {
	var feature models.FeatureRequest
	result := r.db.WithContext(ctx).First(&feature, id)
	return &feature, result.Error
}

// SaveFeatureVote saves a user's vote on a feature
func (r *FeedbackRepository) SaveFeatureVote(ctx context.Context, vote *models.UserFeatureVote) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Save the vote
		if err := tx.Create(vote).Error; err != nil {
			return err
		}

		// Update the vote count
		return tx.Model(&models.FeatureRequest{}).
			Where("id = ?", vote.FeatureID).
			UpdateColumn("vote_count", gorm.Expr("vote_count + ?", 1)).
			Error
	})
}

// SaveSurveyResponse saves a feature survey response
func (r *FeedbackRepository) SaveSurveyResponse(ctx context.Context, response *models.FeatureSurveyResponse) error {
	response.SubmittedAt = time.Now()
	return r.db.WithContext(ctx).Create(response).Error
}

// GetTopFeatures retrieves the most voted features
func (r *FeedbackRepository) GetTopFeatures(ctx context.Context, limit int) ([]models.FeatureRequest, error) {
	var features []models.FeatureRequest
	result := r.db.WithContext(ctx).
		Where("status = ?", "PROPOSED").
		Order("vote_count DESC").
		Limit(limit).
		Find(&features)
	return features, result.Error
}

// GetUserVotes retrieves all votes by a user
func (r *FeedbackRepository) GetUserVotes(ctx context.Context, userID string) ([]models.UserFeatureVote, error) {
	var votes []models.UserFeatureVote
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&votes)
	return votes, result.Error
}

// InitializeFeatures adds initial feature requests if none exist
func (r *FeedbackRepository) InitializeFeatures(ctx context.Context) error {
	var count int64
	r.db.WithContext(ctx).Model(&models.FeatureRequest{}).Count(&count)
	
	if count == 0 {
		features := models.InitialFeatures()
		return r.db.WithContext(ctx).Create(&features).Error
	}
	
	return nil
}
