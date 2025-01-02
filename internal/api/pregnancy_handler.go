package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
)

type PregnancyHandler struct {
	service *pregnancy.Service
	db      *database.PostgresDB
}

func NewPregnancyHandler(service *pregnancy.Service, db *database.PostgresDB) *PregnancyHandler {
	return &PregnancyHandler{
		service: service,
		db:      db,
	}
}

// Add the missing methods
func (h *PregnancyHandler) GetPregnancies(c *gin.Context) {
	pregnancies, err := h.service.GetPregnancies(c.Request.Context())
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, pregnancies)
}

func (h *PregnancyHandler) GetPregnancyStage(c *gin.Context) {
	id := c.Param("id")
	pregnancyID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid pregnancy ID"})
		return
	}

	stage, err := h.service.GetPregnancyStage(c.Request.Context(), pregnancyID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, gin.H{"stage": stage})
}

// Add other pregnancy-specific handlers
func (h *PregnancyHandler) GetPregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	stage, err := h.service.GetPregnancyStage(c.Request.Context(), horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"horse": horse,
		"stage": stage,
		"isPregnant": horse.IsPregnant,
	})
}

func (h *PregnancyHandler) StartPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var data struct {
		ConceptionDate time.Time `json:"conceptionDate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Validate conception date
	if data.ConceptionDate.After(time.Now()) {
		c.JSON(400, ErrorResponse{Error: "Conception date cannot be in the future"})
		return
	}

	// Update horse pregnancy status
	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	horse.IsPregnant = true
	horse.ConceptionDate = &data.ConceptionDate

	if err := h.db.UpdateHorse(&horse); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	// Create pregnancy record
	pregnancy := &models.Pregnancy{
		HorseID:   uint(horseID),
		UserID:    c.GetString("userID"),
		StartDate: data.ConceptionDate,
		Status:    models.PregnancyStatusActive,
	}

	if err := h.db.AddPregnancy(pregnancy); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(201, pregnancy)
}

func (h *PregnancyHandler) EndPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var data struct {
		Status string    `json:"status"`
		Date   time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Update horse pregnancy status
	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	horse.IsPregnant = false
	horse.ConceptionDate = nil

	if err := h.db.UpdateHorse(&horse); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	// Update pregnancy record
	pregnancy, err := h.db.GetPregnancy(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	pregnancy.Status = data.Status
	pregnancy.EndDate = &data.Date

	if err := h.db.UpdatePregnancy(&pregnancy); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, pregnancy)
}

func (h *PregnancyHandler) GetPregnancyEvents(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	events, err := h.db.GetPregnancyEvents(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(200, events)
}

func (h *PregnancyHandler) AddPregnancyEvent(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var event models.PregnancyEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(400, ErrorResponse{Error: err.Error()})
		return
	}

	pregnancy, err := h.db.GetPregnancy(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	event.PregnancyID = pregnancy.ID
	event.Date = time.Now()

	if err := h.db.AddPregnancyEvent(&event); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(201, event)
}

func (h *PregnancyHandler) GetFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Get existing checklist
	items, err := h.db.GetPreFoalingChecklist(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	// If no checklist exists, initialize with default items
	if len(items) == 0 {
		// Get horse to check pregnancy status and due date
		horse, err := h.db.GetHorse(horseID)
		if err != nil {
			c.JSON(500, ErrorResponse{Error: err.Error()})
			return
		}

		if !horse.IsPregnant || horse.ConceptionDate == nil {
			c.JSON(400, ErrorResponse{Error: "Horse is not pregnant"})
			return
		}

		// Calculate due date and set checklist deadlines accordingly
		dueDate := pregnancy.CalculateDueDate(*horse.ConceptionDate, pregnancy.DefaultGestationDays)
		
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
			if err := h.db.AddPreFoalingChecklistItem(&item); err != nil {
				c.JSON(500, ErrorResponse{Error: "Failed to initialize checklist"})
				return
			}
		}
	}

	// Return checklist sorted by due date and priority
	c.JSON(200, items)
}

func (h *PregnancyHandler) GetPostFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Get pregnancy to find foaling date
	pregnancy, err := h.db.GetPregnancy(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}

	if pregnancy.EndDate == nil {
		c.JSON(400, ErrorResponse{Error: "Foaling date not recorded"})
		return
	}

	foalingDate := *pregnancy.EndDate

	// Create post-foaling checklist with specific timelines
	items := []models.PreFoalingChecklistItem{
		// Immediate (within first hour)
		{
			HorseID:     uint(horseID),
			Description: "Check mare's vital signs",
			Priority:    models.PriorityHigh,
				DueDate:     foalingDate.Add(1 * time.Hour),
				Notes:       "Temperature: 37.5-38.5°C, Heart rate: 28-44 bpm, Respiration: 8-16 breaths/min",
		},
		{
			HorseID:     uint(horseID),
			Description: "Assess foal's initial vitals",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.Add(30 * time.Minute),
			Notes:       "Breathing normal, heart rate 60-120 bpm, strong suckle reflex",
		},
		// 2-3 hours post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Monitor placenta expulsion",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.Add(3 * time.Hour),
			Notes:       "Should be complete within 3 hours. Save placenta for vet inspection",
		},
		{
			HorseID:     uint(horseID),
			Description: "Check foal nursing",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.Add(2 * time.Hour),
			Notes:       "Ensure proper latching and colostrum intake. Monitor for 15-20 minutes every 2 hours",
		},
		// 6-12 hours post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Monitor foal's first defecation",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.Add(12 * time.Hour),
			Notes:       "Meconium passage should begin within 12 hours. Watch for straining",
		},
		{
			HorseID:     uint(horseID),
			Description: "Check foal's leg strength",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.Add(6 * time.Hour),
			Notes:       "Should stand strongly and walk without significant weakness",
		},
		// 24 hours post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Veterinary examination",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.AddDate(0, 0, 1),
			Notes:       "Complete health check for mare and foal, IgG test for foal",
		},
		{
			HorseID:     uint(horseID),
			Description: "Monitor mare's milk production",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.AddDate(0, 0, 1),
			Notes:       "Check udder health, milk flow, and foal's nursing frequency",
		},
		// 48-72 hours post-foaling
		{
			HorseID:     uint(horseID),
			Description: "Check mare's uterine discharge",
			Priority:    models.PriorityHigh,
			DueDate:     foalingDate.AddDate(0, 0, 2),
			Notes:       "Monitor color and odor. Should be dark red to brown, no foul smell",
		},
		{
			HorseID:     uint(horseID),
			Description: "Assess foal bonding",
			Priority:    models.PriorityMedium,
			DueDate:     foalingDate.AddDate(0, 0, 3),
			Notes:       "Observe mare-foal interactions, ensure proper maternal behavior",
		},
	}

	c.JSON(200, items)
}

func (h *PregnancyHandler) GetPregnancyGuidelines(c *gin.Context) {
	stageStr := c.Query("stage")
	stage := models.PregnancyStage(stageStr)
	guidelines := h.service.GetPregnancyGuidelinesByStage(stage)
	c.JSON(200, guidelines)
}

func (h *PregnancyHandler) CheckPreFoalingSigns(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	signs, err := h.db.GetPreFoalingSigns(horseID)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, signs)
}

func (h *PregnancyHandler) RecordPreFoalingSign(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var sign models.PreFoalingSign
	if err := c.ShouldBindJSON(&sign); err != nil {
		c.JSON(400, ErrorResponse{Error: err.Error()})
		return
	}

	sign.HorseID = uint(horseID)
	sign.Date = time.Now()

	if err := h.db.AddPreFoalingSign(&sign); err != nil {
		c.JSON(500, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(201, sign)
}

// Add this error type
type ErrorResponse struct {
	Error string `json:"error"`
}
