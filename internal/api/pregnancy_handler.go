package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"gorm.io/gorm"
)

type PregnancyHandler struct {
	service      service.PregnancyService
	validate     *validator.Validate
	horseService service.HorseService
}

func NewPregnancyHandler(pregService service.PregnancyService, horseService service.HorseService, validate *validator.Validate) *PregnancyHandler {
	return &PregnancyHandler{
		service:      pregService,
		horseService: horseService,
		validate:     validate,
	}
}

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
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseIDParam := c.Param("id")
	horseID, err := strconv.ParseUint(horseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
		}
		return
	}

	var eventInput models.PregnancyEventInputDTO
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	if err := h.validate.Struct(&eventInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrorsForPregnancyHandler(err)})
		return
	}

	if eventInput.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Event date cannot be in the future."})
		return
	}

	createdEvent, err := h.service.AddPregnancyEvent(c.Request.Context(), userID, uint(horseID), &eventInput)
	if err != nil {
		if strings.Contains(err.Error(), "no active pregnancy found") {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to add pregnancy event: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
}

func (h *PregnancyHandler) GetFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	items, err := h.service.GetPreFoalingChecklist(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if len(items) == 0 {
		pregnancy, err := h.service.GetPregnancy(c.Request.Context(), uint(horseID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
			return
		}

		if pregnancy.ConceptionDate == nil {
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Horse is not pregnant"})
			return
		}

		dueDate := pregnancy.ConceptionDate.Add(340 * 24 * time.Hour)
		
		items = []models.PreFoalingChecklistItem{
			{
				HorseID:     uint(horseID),
				Description: "Schedule pre-foaling veterinary exam",
				Priority:    models.PriorityHigh,
				DueDate:     dueDate.AddDate(0, 0, -60),
				Notes:       "Check vaccination status, overall health, and pregnancy progress",
			},
			{
				HorseID:     uint(horseID),
				Description: "Begin mammary gland monitoring",
				Priority:    models.PriorityMedium,
				DueDate:     dueDate.AddDate(0, 0, -45),
				Notes:       "Document any changes in size or appearance",
			},
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

		for _, item := range items {
			if err := h.service.AddPreFoalingChecklistItem(c.Request.Context(), &item); err != nil {
				c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to initialize checklist"})
				return
			}
		}
	}

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
		{
			HorseID:     uint(horseID),
			Description: "Veterinary check-up",
			Priority:    models.PriorityHigh,
			DueDate:     data.FoalingDate.AddDate(0, 0, 1),
			Notes:       "Schedule vet visit for mare and foal health assessment",
		},
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

func formatValidationErrorsForPregnancyHandler(err error) map[string]string {
	errorsMap := make(map[string]string)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			errorsMap[fieldName] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldName, fieldError.Tag())
		}
	} else {
		errorsMap["error"] = "Validation error: " + err.Error()
	}
	return errorsMap
}
