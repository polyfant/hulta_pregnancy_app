package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

type PrivacyService interface {
	GetUserPreferences(ctx context.Context, userID string) (*models.PrivacyPreferences, error)
	UpdatePreferences(ctx context.Context, userID string, prefs *models.PrivacyPreferences) error
	DeleteUserData(ctx context.Context, userID string, dataType string) error
	CheckFeatureEnabled(ctx context.Context, userID string, feature string) (bool, error)
}

type PrivacyHandler struct {
	service PrivacyService
}

func NewPrivacyHandler(service PrivacyService) *PrivacyHandler {
	return &PrivacyHandler{
		service: service,
	}
}

// GetPrivacySettings handles GET /privacy/settings
func (h *PrivacyHandler) GetPrivacySettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	prefs, err := h.service.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePrivacySettings handles PUT /privacy/settings
func (h *PrivacyHandler) UpdatePrivacySettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var updates models.PrivacyPreferences
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.UpdatePreferences(c.Request.Context(), userID, &updates); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// RequestDataDeletion handles POST /privacy/data-deletion
func (h *PrivacyHandler) RequestDataDeletion(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var req struct {
		DataTypes []string `json:"data_types"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request body"})
		return
	}

	for _, dataType := range req.DataTypes {
		if err := h.service.DeleteUserData(c.Request.Context(), userID, dataType); err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.Status(http.StatusOK)
}

// GetDataSharingStatus handles GET /privacy/sharing-status/:feature
func (h *PrivacyHandler) GetDataSharingStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	feature := c.Param("feature")
	enabled, err := h.service.CheckFeatureEnabled(c.Request.Context(), userID, feature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feature": feature,
		"enabled": enabled,
	})
}
