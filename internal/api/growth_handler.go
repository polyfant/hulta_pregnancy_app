package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

type GrowthHandler struct {
	growthService service.GrowthService
	validate      *validator.Validate
}

func NewGrowthHandler(growthService service.GrowthService, validate *validator.Validate) *GrowthHandler {
	return &GrowthHandler{
		growthService: growthService,
		validate:      validate,
	}
}

// RecordGrowthMeasurement handles recording a new growth measurement for a foal
func (h *GrowthHandler) RecordGrowthMeasurement(c *gin.Context) {
	foalIDStr := c.Param("id")
	foalID, err := strconv.ParseUint(foalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid foal ID in path"})
		return
	}

	var growthInput models.GrowthData
	if err := c.ShouldBindJSON(&growthInput); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	growthInput.FoalID = uint(foalID)

	if err := h.validate.Struct(&growthInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	if growthInput.MeasurementDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Measurement date cannot be in the future."})
		return
	}

	createdGrowthData, err := h.growthService.RecordGrowthMeasurement(c.Request.Context(), &growthInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to record growth measurement: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Growth measurement recorded successfully", "data": createdGrowthData})
}

// GetFoalGrowthData retrieves all growth data for a specific foal
func (h *GrowthHandler) GetFoalGrowthData(c *gin.Context) {
	foalIDStr := c.Param("id")
	foalID, err := strconv.ParseUint(foalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid foal ID"})
		return
	}

	growthData, err := h.growthService.GetFoalGrowthData(c.Request.Context(), uint(foalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve growth data"})
		return
	}

	c.JSON(http.StatusOK, growthData)
}

// AnalyzeGrowthTrends provides an analysis of a foal's growth trends
func (h *GrowthHandler) AnalyzeGrowthTrends(c *gin.Context) {
	foalIDStr := c.Param("id")
	foalID, err := strconv.ParseUint(foalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid foal ID"})
		return
	}

	analysis, err := h.growthService.AnalyzeGrowthTrends(c.Request.Context(), uint(foalID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to analyze growth trends"})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// RecordBodyCondition handles recording a new body condition score for a foal
func (h *GrowthHandler) RecordBodyCondition(c *gin.Context) {
	foalIDStr := c.Param("id")
	foalID, err := strconv.ParseUint(foalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid foal ID in path"})
		return
	}

	var bodyConditionInput models.BodyCondition
	if err := c.ShouldBindJSON(&bodyConditionInput); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	bodyConditionInput.FoalID = uint(foalID)

	if err := h.validate.Struct(&bodyConditionInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	if bodyConditionInput.LastUpdated.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Last updated date cannot be in the future."})
		return
	}

	createdBodyCondition, err := h.growthService.RecordBodyCondition(c.Request.Context(), &bodyConditionInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to record body condition: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Body condition record created successfully", "data": createdBodyCondition})
}
