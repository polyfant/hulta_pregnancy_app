package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/health"
)

// HealthHandler handles HTTP requests for health records
type HealthHandler struct {
	service health.Service
}

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler(service health.Service) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

// HandleCreateHealthRecord handles the creation of a new health record
func (h *HealthHandler) HandleCreateHealthRecord(c *gin.Context) {
	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateHealthRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// HandleGetHealthRecord handles retrieving a specific health record
func (h *HealthHandler) HandleGetHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	record, err := h.service.GetHealthRecord(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// HandleListHealthRecords handles retrieving all health records for a horse
func (h *HealthHandler) HandleListHealthRecords(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("horse_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	records, err := h.service.ListHealthRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// HandleUpdateHealthRecord handles updating a health record
func (h *HealthHandler) HandleUpdateHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.ID = uint(id)
	if err := h.service.UpdateHealthRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// HandleDeleteHealthRecord handles deleting a health record
func (h *HealthHandler) HandleDeleteHealthRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	if err := h.service.DeleteHealthRecord(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
