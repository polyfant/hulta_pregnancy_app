package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

// FeedbackHandler handles feature request and feedback endpoints
type FeedbackHandler struct {
	repo *repository.FeedbackRepository
}

// NewFeedbackHandler creates a new feedback handler
func NewFeedbackHandler(repo *repository.FeedbackRepository) *FeedbackHandler {
	return &FeedbackHandler{repo: repo}
}

// GetFeatures handles retrieving all feature requests
func (h *FeedbackHandler) GetFeatures(c *gin.Context) {
	features, err := h.repo.GetFeatureRequests(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get features"})
		return
	}

	c.JSON(http.StatusOK, features)
}

// VoteForFeature handles recording a vote for a feature
func (h *FeedbackHandler) VoteForFeature(c *gin.Context) {
	featureID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid feature ID"})
		return
	}

	var vote models.UserFeatureVote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid vote data"})
		return
	}

	vote.FeatureID = uint(featureID)
	
	if err := h.repo.SaveFeatureVote(c.Request.Context(), &vote); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to save vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

// SubmitSurvey handles submitting a feature survey response
func (h *FeedbackHandler) SubmitSurvey(c *gin.Context) {
	var response models.FeatureSurveyResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid survey data"})
		return
	}

	if err := h.repo.SaveSurveyResponse(c.Request.Context(), &response); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to save survey"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Survey submitted successfully"})
}

// GetTopFeatures handles retrieving the most requested features
func (h *FeedbackHandler) GetTopFeatures(c *gin.Context) {
	limit := 5 // Default to top 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	features, err := h.repo.GetTopFeatures(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get top features"})
		return
	}

	c.JSON(http.StatusOK, features)
}

// GetUserVotes handles retrieving a user's feature votes
func (h *FeedbackHandler) GetUserVotes(c *gin.Context) {
	userID := c.GetString("user_id") // Assuming user ID is set in middleware
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	votes, err := h.repo.GetUserVotes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get user votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

// HandleFeatureRequest handles POST /feedback/feature
func (h *FeedbackHandler) HandleFeatureRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var request models.FeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid feature request data"})
		return
	}

	request.UserID = userID
	if err := h.repo.SaveFeatureRequest(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to save feature request"})
		return
	}

	c.JSON(http.StatusCreated, request)
}

// HandleListFeatures handles GET /feedback/features
func (h *FeedbackHandler) HandleListFeatures(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	features, err := h.repo.GetFeatureRequests(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get features"})
		return
	}

	c.JSON(http.StatusOK, features)
}

// HandleFeatureVote handles POST /feedback/features/:id/vote
func (h *FeedbackHandler) HandleFeatureVote(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	featureID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid feature ID"})
		return
	}

	var vote models.FeatureVote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid vote data"})
		return
	}

	vote.UserID = userID
	vote.FeatureID = uint(featureID)

	if err := h.repo.SaveFeatureVote(c.Request.Context(), &vote); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to save vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote recorded successfully"})
}

// HandleGetUserVotes handles GET /feedback/votes
func (h *FeedbackHandler) HandleGetUserVotes(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	votes, err := h.repo.GetUserVotes(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get user votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}
