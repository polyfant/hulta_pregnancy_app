package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
)

type PregnancyHandler struct {
	service service.PregnancyService
}

func NewPregnancyHandler(service service.PregnancyService) *PregnancyHandler {
	return &PregnancyHandler{
		service: service,
	}
}

// Add the missing methods
func (h *PregnancyHandler) GetPregnancies(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	pregnancies, err := h.service.GetActive(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, pregnancies)
}

func (h *PregnancyHandler) GetPregnancyStage(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	stage, err := h.service.GetPregnancyStage(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stage": stage})
}

// Add other pregnancy-specific handlers
func (h *PregnancyHandler) GetPregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	status, err := h.service.GetStatus(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *PregnancyHandler) StartPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var start models.PregnancyStart
	if err := c.ShouldBindJSON(&start); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request data"})
		return
	}

	if err := h.service.StartTracking(c.Request.Context(), uint(horseID), start); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Pregnancy tracking started successfully"})
}

func (h *PregnancyHandler) EndPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var data struct {
		Status string    `json:"status"`
		Date   time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request data"})
		return
	}

	pregnancy, err := h.service.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	pregnancy.Status = data.Status
	pregnancy.EndDate = &data.Date

	if err := h.service.UpdatePregnancy(c.Request.Context(), pregnancy); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

func (h *PregnancyHandler) GetPregnancyEvents(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	events, err := h.service.GetPregnancyEvents(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *PregnancyHandler) AddPregnancyEvent(c *gin.Context) {
	var event models.PregnancyEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.AddPregnancyEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *PregnancyHandler) GetFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Get existing checklist
	items, err := h.service.GetPreFoalingChecklist(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	// If no checklist exists, initialize with default items
	if len(items) == 0 {
		// Get horse to check pregnancy status and due date
		pregnancy, err := h.service.GetPregnancy(c.Request.Context(), uint(horseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
			return
		}

		if pregnancy.ConceptionDate == nil {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Horse is not pregnant"})
			return
		}

		// Calculate due date and set checklist deadlines accordingly
		dueDate := pregnancy.ConceptionDate.Add(340 * 24 * time.Hour)
		
		// Create timeline-based checklist items
		items = []models.PreFoalingChecklistItem{
			// 60 days before
			{
				HorseID:     uint(horseID),
				Description: "Schedule pre-foaling veterinary exam",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -60),
				Notes:       "Check vaccination status, overall health, and pregnancy progress",
			},
			// 45 days before
			{
				HorseID:     uint(horseID),
				Description: "Begin mammary gland monitoring",
				Priority:    models.PriorityMedium,
				DueDate:     dueDate.AddDate(0, 0, -45),
				Notes:       "Document any changes in size or appearance",
			},
			// 30 days before
			{
				HorseID:     uint(horseID),
				Description: "Prepare foaling kit",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -30),
				Notes:       "Include: sterile gloves, iodine, clean towels, flashlight, watch, emergency contacts, tail wrap, umbilical clamp",
			},
			{
				HorseID:     uint(horseID),
				Description: "Set up foaling notification system",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -30),
				Notes:       "Test cameras, alarms, and ensure backup power supply",
			},
			// 21 days before
			{
				HorseID:     uint(horseID),
				Description: "Start intensive udder monitoring",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -21),
				Notes:       "Check twice daily: size, firmness, waxing. Document with photos",
			},
			{
				HorseID:     uint(horseID),
				Description: "Prepare foaling stall",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -21),
				Notes:       "Clean thoroughly, fresh bedding, ensure good lighting and ventilation",
			},
			// 14 days before
			{
				HorseID:     uint(horseID),
				Description: "Begin vulva monitoring",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -14),
				Notes:       "Check for relaxation and color changes",
			},
			{
				HorseID:     uint(horseID),
				Description: "Review emergency procedures",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -14),
				Notes:       "Update contact numbers, review red-flag symptoms, plan transport route to clinic",
			},
			// 7 days before
			{
				HorseID:     uint(horseID),
				Description: "Begin temperature monitoring",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -7),
				Notes:       "Monitor twice daily: normal 37.5-38.5°C. Drop of 1°C may indicate imminent foaling",
			},
			{
				HorseID:     uint(horseID),
				Description: "Monitor behavioral changes",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -7),
				Notes:       "Watch for restlessness, pawing, sweating, frequent urination",
			},
		}

		// Add each item to database
		for _, item := range items {
			if err := h.service.AddPreFoalingChecklistItem(c.Request.Context(), &item); err != nil {
				c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to initialize checklist"})
				return
			}
		}
	}

	// Return checklist sorted by due date and priority
	c.JSON(http.StatusOK, items)
}

func (h *PregnancyHandler) GetPostFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var data struct {
		FoalingDate time.Time `json:"foalingDate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Foaling date is required"})
		return
	}

	items := []models.PreFoalingChecklistItem{
		// Immediate post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Check placenta expulsion",
			Priority:    models.PriorityHigh,
			DueDate:     data.FoalingDate.Add(3 * time.Hour),
			Notes:       "Should be expelled within 3 hours. Save for vet inspection if needed",
		},
		{
			HorseID:     uint(horseID),
			Description: "Monitor foal nursing",
			Priority:    models.PriorityHigh,
			DueDate:     data.FoalingDate.Add(6 * time.Hour),
			Notes:       "Ensure successful first nursing within 6 hours",
		},
		// 24 hours post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Veterinary check-up",
			Priority:    models.PriorityHigh,
			DueDate:     data.FoalingDate.AddDate(0, 0, 1),
			Notes:       "Schedule vet visit for mare and foal health assessment",
		},
		// Continue with other checklist items...
	}

	c.JSON(http.StatusOK, items)
}

func (h *PregnancyHandler) GetPregnancyGuidelines(c *gin.Context) {
	stageStr := c.Query("stage")
	if stageStr == "" {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Stage parameter is required"})
		return
	}
	
	stage := models.PregnancyStage(stageStr)
	
	guidelines, err := h.service.GetGuidelines(c.Request.Context(), stage)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, guidelines)
}

func (h *PregnancyHandler) CheckPreFoalingSigns(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	signs, err := h.service.GetPreFoalingSigns(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, signs)
}

func (h *PregnancyHandler) RecordPreFoalingSign(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var sign models.PreFoalingSign
	if err := c.ShouldBindJSON(&sign); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	sign.HorseID = uint(horseID)
	sign.Date = time.Now()

	if err := h.service.AddPreFoalingSign(c.Request.Context(), &sign); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sign)
}

func (h *PregnancyHandler) UpdatePregnancy(c *gin.Context) {
	var pregnancy models.Pregnancy
	if err := c.ShouldBindJSON(&pregnancy); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdatePregnancy(c.Request.Context(), &pregnancy); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

func (h *PregnancyHandler) GetPregnancy(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	pregnancy, err := h.service.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}
