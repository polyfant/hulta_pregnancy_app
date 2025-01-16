package repository

import (
	"context"
	"testing"
	"time"

	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
)

func TestFeedbackRepository(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)
	repo := NewFeedbackRepository(db)

	t.Run("InitializeFeatures", func(t *testing.T) {
		// Test initial feature creation
		err := repo.InitializeFeatures(context.Background())
		require.NoError(t, err)

		// Verify features were created
		features, err := repo.GetFeatureRequests(context.Background())
		require.NoError(t, err)
		assert.NotEmpty(t, features)

		// Test idempotency (calling again shouldn't create duplicates)
		err = repo.InitializeFeatures(context.Background())
		require.NoError(t, err)

		featuresAfter, err := repo.GetFeatureRequests(context.Background())
		require.NoError(t, err)
		assert.Equal(t, len(features), len(featuresAfter))
	})

	t.Run("GetFeatureRequests", func(t *testing.T) {
		ctx := context.Background()

		// Create test features with different vote counts
		testFeatures := []struct {
			title     string
			voteCount int
		}{
			{"Most Voted", 10},
			{"Second Most", 5},
			{"Least Voted", 2},
		}

		for _, tf := range testFeatures {
			feature, err := testutil.CreateTestFeatureRequest(db, tf.title)
			require.NoError(t, err)

			// Update vote count
			err = db.Model(&models.FeatureRequest{}).
				Where("id = ?", feature.ID).
				Update("vote_count", tf.voteCount).Error
			require.NoError(t, err)
		}

		// Test retrieval
		features, err := repo.GetFeatureRequests(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, features)
		
		// Verify ordering (should be by vote count DESC)
		assert.Equal(t, "Most Voted", features[0].Title)
		assert.Equal(t, "Least Voted", features[len(features)-1].Title)
	})

	t.Run("GetFeatureRequest", func(t *testing.T) {
		tests := []struct {
			name      string
			id        uint
			wantErr   bool
			wantTitle string
		}{
			{
				name:      "Existing feature",
				id:        1,
				wantErr:   false,
				wantTitle: "Test Feature 1",
			},
			{
				name:      "Non-existent feature",
				id:        999,
				wantErr:   true,
				wantTitle: "",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				feature, err := repo.GetFeatureRequest(nil, tc.id)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.wantTitle, feature.Title)
				}
			})
		}
	})

	t.Run("SaveFeatureVote", func(t *testing.T) {
		tests := []struct {
			name      string
			featureID uint
			userID    string
			wantErr   bool
		}{
			{
				name:      "Valid vote",
				featureID: 1,
				userID:    "user123",
				wantErr:   false,
			},
			{
				name:      "Invalid feature ID",
				featureID: 999,
				userID:    "user123",
				wantErr:   true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				vote := &models.UserFeatureVote{
					FeatureID: tc.featureID,
					UserID:    tc.userID,
				}

				err := repo.SaveFeatureVote(nil, vote)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					// Verify vote was saved
					var votes []models.UserFeatureVote
					err = db.Where("feature_id = ? AND user_id = ?", tc.featureID, tc.userID).Find(&votes).Error
					assert.NoError(t, err)
					assert.Len(t, votes, 1)
				}
			})
		}
	})

	t.Run("SaveSurveyResponse", func(t *testing.T) {
		ctx := context.Background()

		testCases := []struct {
			name        string
			response    *models.FeatureSurveyResponse
			expectError bool
		}{
			{
				name: "Valid Response",
				response: &models.FeatureSurveyResponse{
					UserID:     "test_user",
					FeatureID:  1,
					WouldUse:   true,
					WouldPay:   true,
					PricePoint: ptr(9.99),
					UseCase:    "Daily monitoring",
					Feedback:   "Great feature!",
				},
				expectError: false,
			},
			{
				name: "Invalid Price Point",
				response: &models.FeatureSurveyResponse{
					UserID:     "test_user",
					FeatureID:  1,
					WouldUse:   true,
					WouldPay:   true,
					PricePoint: ptr(-9.99), // Negative price
					UseCase:    "Daily monitoring",
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				err := repo.SaveSurveyResponse(ctx, tc.response)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("GetTopFeatures", func(t *testing.T) {
		ctx := context.Background()

		// Create features with different vote counts
		voteData := []struct {
			title string
			votes int
		}{
			{"Top Feature", 20},
			{"Second Feature", 15},
			{"Third Feature", 10},
			{"Fourth Feature", 5},
			{"Last Feature", 1},
		}

		for _, vd := range voteData {
			feature, err := testutil.CreateTestFeatureRequest(db, vd.title)
			require.NoError(t, err)

			err = db.Model(&models.FeatureRequest{}).
				Where("id = ?", feature.ID).
				Update("vote_count", vd.votes).Error
			require.NoError(t, err)
		}

		testCases := []struct {
			name          string
			limit         int
			expectedCount int
		}{
			{
				name:          "Get Top 3",
				limit:         3,
				expectedCount: 3,
			},
			{
				name:          "Get All",
				limit:         10,
				expectedCount: 5,
			},
			{
				name:          "Zero Limit",
				limit:         0,
				expectedCount: 0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				features, err := repo.GetTopFeatures(ctx, tc.limit)
				require.NoError(t, err)
				assert.Len(t, features, tc.expectedCount)

				if tc.expectedCount > 0 {
					// Verify ordering
					assert.Equal(t, "Top Feature", features[0].Title)
				}
			})
		}
	})

	t.Run("GetUserVotes", func(t *testing.T) {
		userID := "test_user"
		
		// Create test votes
		testVotes := []models.UserFeatureVote{
			{
				UserID:    userID,
				FeatureID: 1,
			},
			{
				UserID:    userID,
				FeatureID: 2,
			},
		}

		for _, vote := range testVotes {
			err := repo.SaveFeatureVote(nil, &vote)
			require.NoError(t, err)
		}

		// Get user votes
		votes, err := repo.GetUserVotes(nil, userID)
		require.NoError(t, err)
		assert.Len(t, votes, 2)
	})

	t.Run("Concurrent Voting", func(t *testing.T) {
		// Create test feature
		feature := &models.FeatureRequest{
			Title:       "Concurrent Test Feature",
			Description: "Test feature for concurrent voting",
			Status:      "PLANNED",
		}
		err := db.Create(feature).Error
		require.NoError(t, err)

		// Test concurrent voting
		numVoters := 5
		done := make(chan bool)

		for i := 0; i < numVoters; i++ {
			go func(userID string) {
				vote := &models.UserFeatureVote{
					UserID:    userID,
					FeatureID: feature.ID,
				}
				err := repo.SaveFeatureVote(nil, vote)
				require.NoError(t, err)
				done <- true
			}("user_" + strconv.Itoa(i))
		}

		// Wait for all votes
		for i := 0; i < numVoters; i++ {
			<-done
		}

		// Verify vote count
		var updatedFeature models.FeatureRequest
		err = db.First(&updatedFeature, feature.ID).Error
		require.NoError(t, err)
		assert.Equal(t, numVoters, updatedFeature.VoteCount)
	})

	t.Run("GetTopFeatures", func(t *testing.T) {
		ctx := context.Background()

		// Create features with different vote counts
		voteData := []struct {
			title string
			votes int
		}{
			{"Top Feature", 20},
			{"Second Feature", 15},
			{"Third Feature", 10},
			{"Fourth Feature", 5},
			{"Last Feature", 1},
		}

		for _, vd := range voteData {
			feature, err := testutil.CreateTestFeatureRequest(db, vd.title)
			require.NoError(t, err)

			err = db.Model(&models.FeatureRequest{}).
				Where("id = ?", feature.ID).
				Update("vote_count", vd.votes).Error
			require.NoError(t, err)
		}

		testCases := []struct {
			name          string
			limit         int
			expectedCount int
		}{
			{
				name:          "Get Top 3",
				limit:         3,
				expectedCount: 3,
			},
			{
				name:          "Get All",
				limit:         10,
				expectedCount: 5,
			},
			{
				name:          "Zero Limit",
				limit:         0,
				expectedCount: 0,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				features, err := repo.GetTopFeatures(ctx, tc.limit)
				require.NoError(t, err)
				assert.Len(t, features, tc.expectedCount)

				if tc.expectedCount > 0 {
					// Verify ordering
					assert.Equal(t, "Top Feature", features[0].Title)
				}
			})
		}
	})

	t.Run("GetUserVotes", func(t *testing.T) {
		userID := "test_user"
		
		// Create test votes
		testVotes := []models.UserFeatureVote{
			{
				UserID:    userID,
				FeatureID: 1,
			},
			{
				UserID:    userID,
				FeatureID: 2,
			},
		}

		for _, vote := range testVotes {
			err := repo.SaveFeatureVote(nil, &vote)
			require.NoError(t, err)
		}

		// Get user votes
		votes, err := repo.GetUserVotes(nil, userID)
		require.NoError(t, err)
		assert.Len(t, votes, 2)
	})

	t.Run("Concurrent Voting", func(t *testing.T) {
		ctx := context.Background()
		feature, err := testutil.CreateTestFeatureRequest(db, "Concurrent Test Feature")
		require.NoError(t, err)

		// Simulate concurrent votes
		for i := 0; i < 10; i++ {
			go func(i int) {
				vote := &models.UserFeatureVote{
					UserID:    "user_" + strconv.Itoa(i),
					FeatureID: feature.ID,
					Priority:  "MUST_HAVE",
					VotedAt:   testutil.MockTimeNow(),
				}
				_ = repo.SaveFeatureVote(ctx, vote)
			}(i)
		}

		// Give time for goroutines to complete
		time.Sleep(100 * time.Millisecond)

		// Verify final vote count
		var updatedFeature models.FeatureRequest
		err = db.First(&updatedFeature, feature.ID).Error
		require.NoError(t, err)
		assert.Greater(t, updatedFeature.VoteCount, 0)
	})
}

// Helper function to create pointer to float64
func ptr(v float64) *float64 {
	return &v
}
