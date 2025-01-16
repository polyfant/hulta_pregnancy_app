package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedbackHandler(t *testing.T) {
	// Setup
	db, err := testutil.SetupTestDB()
	require.NoError(t, err)

	repo := repository.NewFeedbackRepository(db)
	handler := NewFeedbackHandler(repo)

	// Initialize test features
	err = repo.InitializeFeatures(nil)
	require.NoError(t, err)

	t.Run("GetFeatures", func(t *testing.T) {
		c, w := testutil.CreateTestContext()

		handler.GetFeatures(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var features []models.FeatureRequest
		err := testutil.ParseResponse(w, &features)
		require.NoError(t, err)
		assert.NotEmpty(t, features)
	})

	t.Run("VoteForFeature", func(t *testing.T) {
		testCases := []struct {
			name           string
			featureID      uint
			vote          models.UserFeatureVote
			expectedStatus int
		}{
			{
				name:      "Valid Vote",
				featureID: 1,
				vote: models.UserFeatureVote{
					UserID:   "test_user",
					Priority: "MUST_HAVE",
					Comment:  "Great feature!",
				},
				expectedStatus: http.StatusOK,
			},
			{
				name:      "Invalid Priority",
				featureID: 1,
				vote: models.UserFeatureVote{
					UserID:   "test_user",
					Priority: "INVALID",
				},
				expectedStatus: http.StatusBadRequest,
			},
			{
				name:      "Missing User ID",
				featureID: 1,
				vote: models.UserFeatureVote{
					Priority: "MUST_HAVE",
				},
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				c.Params = gin.Params{{Key: "id", Value: "1"}}

				body, err := json.Marshal(tc.vote)
				require.NoError(t, err)
				c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))

				handler.VoteForFeature(c)

				assert.Equal(t, tc.expectedStatus, w.Code)

				if tc.expectedStatus == http.StatusOK {
					// Verify vote was recorded
					votes, err := repo.GetUserVotes(nil, tc.vote.UserID)
					require.NoError(t, err)
					assert.NotEmpty(t, votes)
				}
			})
		}
	})

	t.Run("SubmitSurvey", func(t *testing.T) {
		testCases := []struct {
			name           string
			survey        models.FeatureSurveyResponse
			expectedStatus int
		}{
			{
				name: "Valid Survey",
				survey: models.FeatureSurveyResponse{
					UserID:     "test_user",
					FeatureID:  1,
					WouldUse:   true,
					WouldPay:   true,
					PricePoint: ptr(9.99),
					UseCase:    "Daily monitoring",
					Feedback:   "Looks great!",
				},
				expectedStatus: http.StatusOK,
			},
			{
				name: "Missing Required Fields",
				survey: models.FeatureSurveyResponse{
					UserID: "test_user",
				},
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()

				body, err := json.Marshal(tc.survey)
				require.NoError(t, err)
				c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(body))

				handler.SubmitSurvey(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
			})
		}
	})

	t.Run("GetTopFeatures", func(t *testing.T) {
		// Create features with different vote counts
		features := []struct {
			title string
			votes int
		}{
			{"Most Popular", 10},
			{"Second Popular", 5},
			{"Least Popular", 2},
		}

		for _, f := range features {
			feature, err := testutil.CreateTestFeatureRequest(db, f.title)
			require.NoError(t, err)

			// Add votes
			for i := 0; i < f.votes; i++ {
				vote := &models.UserFeatureVote{
					UserID:    "user_" + string(i),
					FeatureID: feature.ID,
					Priority:  "MUST_HAVE",
				}
				err := repo.SaveFeatureVote(nil, vote)
				require.NoError(t, err)
			}
		}

		testCases := []struct {
			name          string
			limit         string
			expectedSize  int
			expectedFirst string
		}{
			{
				name:          "Default Limit",
				limit:         "",
				expectedSize:  5,
				expectedFirst: "Most Popular",
			},
			{
				name:          "Custom Limit",
				limit:         "2",
				expectedSize:  2,
				expectedFirst: "Most Popular",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				if tc.limit != "" {
					c.Request, _ = http.NewRequest("GET", "/?limit="+tc.limit, nil)
				} else {
					c.Request, _ = http.NewRequest("GET", "/", nil)
				}

				handler.GetTopFeatures(c)

				assert.Equal(t, http.StatusOK, w.Code)

				var features []models.FeatureRequest
				err := testutil.ParseResponse(w, &features)
				require.NoError(t, err)

				if tc.expectedSize > 0 {
					assert.Len(t, features, tc.expectedSize)
				}
				if tc.expectedFirst != "" {
					assert.Equal(t, tc.expectedFirst, features[0].Title)
				}
			})
		}
	})

	t.Run("GetUserVotes", func(t *testing.T) {
		userID := "test_user_votes"

		// Create some votes for the user
		votes := []models.UserFeatureVote{
			{
				UserID:    userID,
				FeatureID: 1,
				Priority:  "MUST_HAVE",
			},
			{
				UserID:    userID,
				FeatureID: 2,
				Priority:  "NICE_TO_HAVE",
			},
		}

		for _, vote := range votes {
			err := repo.SaveFeatureVote(nil, &vote)
			require.NoError(t, err)
		}

		c, w := testutil.CreateTestContext()
		c.Set("user_id", userID)

		handler.GetUserVotes(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var userVotes []models.UserFeatureVote
		err := testutil.ParseResponse(w, &userVotes)
		require.NoError(t, err)
		assert.Len(t, userVotes, 2)
	})

	t.Run("HandleFeatureRequest", func(t *testing.T) {
		testCases := []struct {
			name           string
			userID        string
			request       models.FeatureRequest
			expectedStatus int
		}{
			{
				name:   "Valid Request",
				userID: "test-user",
				request: models.FeatureRequest{
					Title:       "New Feature",
					Description: "This is a new feature request",
					Priority:    "high",
				},
				expectedStatus: http.StatusCreated,
			},
			{
				name:   "Missing User ID",
				userID: "",
				request: models.FeatureRequest{
					Title:       "New Feature",
					Description: "This is a new feature request",
				},
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				c.Set("user_id", tc.userID)

				body, err := json.Marshal(tc.request)
				require.NoError(t, err)
				c.Request = httptest.NewRequest(http.MethodPost, "/feedback/feature", bytes.NewReader(body))

				handler.HandleFeatureRequest(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
				if tc.expectedStatus == http.StatusCreated {
					var response models.FeatureRequest
					err := testutil.ParseResponse(w, &response)
					require.NoError(t, err)
					assert.Equal(t, tc.request.Title, response.Title)
					assert.Equal(t, tc.request.Description, response.Description)
					assert.Equal(t, tc.userID, response.UserID)
				}
			})
		}
	})

	t.Run("HandleListFeatures", func(t *testing.T) {
		testCases := []struct {
			name           string
			userID        string
			expectedStatus int
		}{
			{
				name:           "Valid Request",
				userID:         "test-user",
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Missing User ID",
				userID:         "",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				c.Set("user_id", tc.userID)
				c.Request = httptest.NewRequest(http.MethodGet, "/feedback/features", nil)

				handler.HandleListFeatures(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
				if tc.expectedStatus == http.StatusOK {
					var features []models.FeatureRequest
					err := testutil.ParseResponse(w, &features)
					require.NoError(t, err)
					assert.NotEmpty(t, features)
				}
			})
		}
	})

	t.Run("HandleFeatureVote", func(t *testing.T) {
		testCases := []struct {
			name           string
			userID        string
			featureID     string
			vote          models.FeatureVote
			expectedStatus int
		}{
			{
				name:       "Valid Vote",
				userID:     "test-user",
				featureID:  "1",
				vote:      models.FeatureVote{Vote: 1},
				expectedStatus: http.StatusOK,
			},
			{
				name:       "Missing User ID",
				userID:     "",
				featureID:  "1",
				vote:      models.FeatureVote{Vote: 1},
				expectedStatus: http.StatusUnauthorized,
			},
			{
				name:       "Invalid Feature ID",
				userID:     "test-user",
				featureID:  "invalid",
				vote:      models.FeatureVote{Vote: 1},
				expectedStatus: http.StatusBadRequest,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				c.Set("user_id", tc.userID)
				c.AddParam("id", tc.featureID)

				body, err := json.Marshal(tc.vote)
				require.NoError(t, err)
				c.Request = httptest.NewRequest(http.MethodPost, "/feedback/features/"+tc.featureID+"/vote", bytes.NewReader(body))

				handler.HandleFeatureVote(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
				if tc.expectedStatus == http.StatusOK {
					var response gin.H
					err := testutil.ParseResponse(w, &response)
					require.NoError(t, err)
					assert.Equal(t, "Vote recorded successfully", response["message"])
				}
			})
		}
	})

	t.Run("HandleGetUserVotes", func(t *testing.T) {
		testCases := []struct {
			name           string
			userID        string
			expectedStatus int
		}{
			{
				name:           "Valid Request",
				userID:         "test-user",
				expectedStatus: http.StatusOK,
			},
			{
				name:           "Missing User ID",
				userID:         "",
				expectedStatus: http.StatusUnauthorized,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				c, w := testutil.CreateTestContext()
				c.Set("user_id", tc.userID)
				c.Request = httptest.NewRequest(http.MethodGet, "/feedback/votes", nil)

				handler.HandleGetUserVotes(c)

				assert.Equal(t, tc.expectedStatus, w.Code)
				if tc.expectedStatus == http.StatusOK {
					var votes []models.FeatureVote
					err := testutil.ParseResponse(w, &votes)
					require.NoError(t, err)
					// Votes might be empty if user hasn't voted yet
				}
			})
		}
	})
}

// Helper function to create pointer to float64
func ptr(v float64) *float64 {
	return &v
}
