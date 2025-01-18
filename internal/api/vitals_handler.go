package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/validation"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/vitals"
)

// VitalsHandler handles HTTP requests for vital signs
type VitalsHandler struct {
	service vitals.Service
}

// NewVitalsHandler creates a new VitalsHandler instance
func NewVitalsHandler(service vitals.Service) *VitalsHandler {
	return &VitalsHandler{
		service: service,
	}
}

// HandleRecordVitalSigns handles the recording of new vital signs
func (h *VitalsHandler) HandleRecordVitalSigns(c *gin.Context) {
	var vitals models.VitalSigns
	if err := c.ShouldBindJSON(&vitals); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate vital signs
	if !h.validateVitalSigns(&vitals) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vital signs values"})
		return
	}

	// Set recorded time if not provided
	if vitals.RecordedAt.IsZero() {
		vitals.RecordedAt = time.Now()
	}

	if err := h.service.RecordVitalSigns(c.Request.Context(), &vitals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vitals)
}

// HandleGetVitalSigns handles retrieval of vital signs history
func (h *VitalsHandler) HandleGetVitalSigns(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	from, to, err := h.getTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vitals, err := h.service.GetVitalSigns(c.Request.Context(), uint(horseID), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vitals)
}

// HandleGetLatestVitalSigns handles retrieval of the most recent vital signs
func (h *VitalsHandler) HandleGetLatestVitalSigns(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	vitals, err := h.service.GetLatestVitalSigns(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vitals)
}

// HandleGetAlerts handles retrieval of vital signs alerts
func (h *VitalsHandler) HandleGetAlerts(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	includeAcknowledged := c.DefaultQuery("include_acknowledged", "false") == "true"

	alerts, err := h.service.GetAlerts(c.Request.Context(), uint(horseID), includeAcknowledged)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// HandleGetAlert handles retrieval of a specific alert
func (h *VitalsHandler) HandleGetAlert(c *gin.Context) {
	alertID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	alert, err := h.service.GetAlert(c.Request.Context(), uint(alertID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// HandleAcknowledgeAlert handles acknowledging an alert
func (h *VitalsHandler) HandleAcknowledgeAlert(c *gin.Context) {
	alertID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	if err := h.service.AcknowledgeAlert(c.Request.Context(), uint(alertID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// HandleGetTrends handles retrieval of vital signs trends
func (h *VitalsHandler) HandleGetTrends(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	from, to, err := h.getTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trends, err := h.service.GetTrends(c.Request.Context(), uint(horseID), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

// Helper functions

func (h *VitalsHandler) validateVitalSigns(vitals *models.VitalSigns) bool {
	// Use the generic struct validation
	if err := validator.New().Struct(vitals); err != nil {
		return false
	}

	// Additional custom validations
	if vitals.HorseID == 0 {
		return false
	}

	// Temperature validation
	if vitals.Temperature < 35.0 || vitals.Temperature > 42.0 {
		return false
	}

	// Heart rate validation (typical range for horses)
	if vitals.HeartRate < 20 || vitals.HeartRate > 80 {
		return false
	}

	// Respiratory rate validation
	if vitals.RespiratoryRate < 8 || vitals.RespiratoryRate > 40 {
		return false
	}

	return true
}

func (h *VitalsHandler) getTimeRange(c *gin.Context) (time.Time, time.Time, error) {
	fromStr := c.Query("from")
	toStr := c.DefaultQuery("to", time.Now().Format(time.RFC3339))

	var from, to time.Time
	var err error

	if fromStr == "" {
		// Default to 24 hours ago if not specified
		from = time.Now().Add(-24 * time.Hour)
	} else {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	to, err = time.Parse(time.RFC3339, toStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return from, to, nil
}
