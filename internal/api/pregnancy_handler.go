package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/checklist"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
)

type PregnancyHandler struct {
	service        service.PregnancyService
	weatherService service.WeatherService
	checklistSvc   *checklist.Service
}

func NewPregnancyHandler(service service.PregnancyService, weatherService service.WeatherService, checklistSvc *checklist.Service) *PregnancyHandler {
	return &PregnancyHandler{
		service:        service,
		weatherService: weatherService,
		checklistSvc:   checklistSvc,
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

	// Get pregnancy to check status and due date
	pregnancy, err := h.service.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if pregnancy.ConceptionDate == nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Horse is not pregnant"})
		return
	}

	// Calculate due date
	dueDate := pregnancy.ConceptionDate.Add(340 * 24 * time.Hour)
	
	// Generate checklist with season-specific items
	items, err := h.checklistSvc.GenerateChecklist(c.Request.Context(), uint(horseID), dueDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to generate checklist"})
		return
	}

	// Get progress information
	progress, err := h.checklistSvc.GetProgress(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get checklist progress"})
		return
	}

	response := gin.H{
		"items":    items,
		"progress": progress,
		"due_date": dueDate,
	}

	c.JSON(http.StatusOK, response)
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

func (h *PregnancyHandler) GetPregnancyWeatherAdvice(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Get pregnancy stage first
	stage, err := h.service.GetPregnancyStage(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get pregnancy stage"})
		return
	}

	// Get weather advice based on stage
	advice, err := h.weatherService.GetPregnancyWeatherAdvice(c.Request.Context(), string(stage))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get weather advice"})
		return
	}

	c.JSON(http.StatusOK, advice)
}

func (h *PregnancyHandler) HandleCreatePregnancy(c *gin.Context) {
	var pregnancy models.Pregnancy
	if err := c.ShouldBindJSON(&pregnancy); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &pregnancy); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pregnancy)
}

func (h *PregnancyHandler) HandleGetPregnancy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid pregnancy ID"})
		return
	}

	pregnancy, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

func (h *PregnancyHandler) HandleUpdatePregnancy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid pregnancy ID"})
		return
	}

	var pregnancy models.Pregnancy
	if err := c.ShouldBindJSON(&pregnancy); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request body"})
		return
	}

	pregnancy.ID = uint(id)
	if err := h.service.Update(c.Request.Context(), &pregnancy); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

func (h *PregnancyHandler) HandleDeletePregnancy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid pregnancy ID"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PregnancyHandler) HandleListPregnancies(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	pregnancies, err := h.service.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancies)
}
