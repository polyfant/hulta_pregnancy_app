package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

type GrowthHandler struct {
	growthService service.GrowthService
}

func NewGrowthHandler(growthService service.GrowthService) *GrowthHandler {
	return &GrowthHandler{
		growthService: growthService,
	}
}

// RecordGrowthMeasurement handles recording a new growth measurement for a foal
func (h *GrowthHandler) RecordGrowthMeasurement(c *gin.Context) {
	foalIDStr := c.Param("foalId")
	foalID, err := strconv.ParseUint(foalIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid foal ID"})
		return
	}

	var measurementData struct {
		Weight float64 `json:"weight" binding:"required,min=0"`
		Height float64 `json:"height" binding:"required,min=0"`
	}

	if err := c.ShouldBindJSON(&measurementData); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.growthService.RecordGrowthMeasurement(c.Request.Context(), uint(foalID), measurementData.Weight, measurementData.Height); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to record growth measurement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Growth measurement recorded successfully"})
}

// GetFoalGrowthData retrieves all growth data for a specific foal
func (h *GrowthHandler) GetFoalGrowthData(c *gin.Context) {
	foalIDStr := c.Param("foalId")
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
	foalIDStr := c.Param("foalId")
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
