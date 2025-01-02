package api

import (
	"net/http"
	"strconv"
	"time"

	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
)

type PregnancyResponse struct {
	Horse         HorseInfo                    `json:"horse"`
	Progress      float64                      `json:"progress"`
	CurrentDay    int                          `json:"currentDay"`
	RemainingDays int                         `json:"remainingDays"`
	Stage         string                      `json:"stage"`
	Schedule      pregnancy.MonitoringSchedule `json:"schedule"`
}

type HorseInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PregnancyHandler struct {
	service *pregnancy.Service
	store   models.DataStore
}

func NewPregnancyHandler(service *pregnancy.Service, store models.DataStore) *PregnancyHandler {
	return &PregnancyHandler{
		service: service,
		store:   store,
	}
}

func (h *PregnancyHandler) SetupPregnancyRoutes(router *gin.Engine) {
	// Serve static files
	router.Static("/static", filepath.Join("frontend", "static"))
	
	// Serve templates
	router.LoadHTMLGlob(filepath.Join("frontend", "templates", "*"))

	// Routes
	router.GET("/pregnancy/:id", h.handlePregnancyPage)
	router.GET("/api/pregnancy/status/:id", h.handlePregnancyStatus)
	router.POST("/api/pregnancy/start/:id", h.StartPregnancyTracking)
	router.POST("/api/pregnancy/end/:id", h.EndPregnancyTracking)
	router.GET("/api/pregnancy/status/:id", h.GetPregnancyStatus)
	router.POST("/api/pregnancy/event/:id", h.AddPregnancyEvent)
	router.GET("/api/pregnancy/events/:id", h.GetPregnancyEvents)
	router.GET("/api/pregnancy/guidelines/:id", h.GetPregnancyGuidelines)
	router.GET("/api/pregnancy/pre-foaling-signs/:id", h.CheckPreFoalingSigns)
	router.POST("/api/pregnancy/pre-foaling-sign/:id", h.RecordPreFoalingSign)
	router.GET("/api/pregnancy/foaling-checklist", h.GetFoalingChecklist)
	router.GET("/api/pregnancy/post-foaling-checklist", h.GetPostFoalingChecklist)
}

func (h *PregnancyHandler) handlePregnancyPage(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid horse ID",
		})
		return
	}

	// In a real app, get this from the database
	horse := HorseInfo{
		ID:   int(horseID),
		Name: "Bella",
	}

	calculator := pregnancy.NewPregnancyCalculator(getConceptionDate())
	
	response := PregnancyResponse{
		Horse:         horse,
		Progress:      calculator.GetProgressPercentage(),
		CurrentDay:    calculator.GetCurrentDay(),
		RemainingDays: calculator.GetRemainingDays(),
		Stage:         string(calculator.GetStage()),
		Schedule:      calculator.GetMonitoringSchedule(),
	}

	c.HTML(http.StatusOK, "pregnancy.html", response)
}

func (h *PregnancyHandler) handlePregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid horse ID",
		})
		return
	}

	horse := HorseInfo{
		ID:   int(horseID),
		Name: "Bella",
	}

	calculator := pregnancy.NewPregnancyCalculator(getConceptionDate())
	
	response := PregnancyResponse{
		Horse:         horse,
		Progress:      calculator.GetProgressPercentage(),
		CurrentDay:    calculator.GetCurrentDay(),
		RemainingDays: calculator.GetRemainingDays(),
		Stage:         string(calculator.GetStage()),
		Schedule:      calculator.GetMonitoringSchedule(),
	}

	c.JSON(http.StatusOK, response)
}

// Temporary function - replace with database call
func getConceptionDate() time.Time {
	// For testing, set conception date to 150 days ago
	return time.Now().AddDate(0, 0, -150)
}

func (h *PregnancyHandler) StartPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	var req struct {
		ConceptionDate string `json:"conceptionDate"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	conceptionDate, err := time.Parse("2006-01-02", req.ConceptionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conception date format"})
		return
	}

	if err := h.service.StartPregnancyTracking(horseID, conceptionDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PregnancyHandler) EndPregnancyTracking(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	var req struct {
		Outcome string `json:"outcome"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.EndPregnancyTracking(horseID, req.Outcome); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *PregnancyHandler) GetPregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	horse, err := h.store.GetHorse(horseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := h.service.GetPregnancyStatus(horse.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (h *PregnancyHandler) AddPregnancyEvent(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	var req struct {
		EventType   string `json:"eventType"`
		Description string `json:"description"`
		Notes       string `json:"notes"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.AddPregnancyEvent(horseID, req.EventType, req.Description, req.Notes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *PregnancyHandler) GetPregnancyEvents(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	events, err := h.service.GetPregnancyEvents(horseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *PregnancyHandler) GetPregnancyGuidelines(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	horse, err := h.store.GetHorse(horseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	guidelines, err := h.service.GetPregnancyGuidelines(horse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, guidelines)
}

func (h *PregnancyHandler) CheckPreFoalingSigns(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	horse, err := h.store.GetHorse(horseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signs := h.service.CheckPreFoalingSigns(horse)
	c.JSON(http.StatusOK, signs)
}

func (h *PregnancyHandler) RecordPreFoalingSign(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	var req struct {
		SignName string `json:"signName"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.RecordPreFoalingSign(horseID, req.SignName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *PregnancyHandler) GetFoalingChecklist(c *gin.Context) {
	checklist := h.service.GetFoalingChecklist()
	c.JSON(http.StatusOK, checklist)
}

func (h *PregnancyHandler) GetPostFoalingChecklist(c *gin.Context) {
	checklist := h.service.GetPostFoalingChecklist()
	c.JSON(http.StatusOK, checklist)
}
