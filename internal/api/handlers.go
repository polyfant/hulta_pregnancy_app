package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
)

// Handler handles HTTP requests
type Handler struct {
	horseService         *service.HorseService
	userService          *service.UserService
	pregnancyService     *service.PregnancyService
	healthService        *service.HealthService
	cache                *cache.MemoryCache
	db                   *database.PostgresDB
	horseRepo           repository.HorseRepository
	breedingRepo        repository.BreedingRepository
}

// NewHandler creates a new Handler instance
func NewHandler(config HandlerConfig) *Handler {
	return &Handler{
		horseService:     config.HorseService,
		userService:      config.UserService,
		pregnancyService: config.PregnancyService,
		healthService:    config.HealthService,
		cache:           config.Cache,
		db:             config.Database,
		horseRepo:      config.HorseRepo,
		breedingRepo:   config.BreedingRepo,
	}
}

func (h *Handler) ListHorses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horses, err := h.horseService.ListHorsesByUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}
	SendSuccess(c, horses)
}

func (h *Handler) DeleteHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.horseService.GetHorse(uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	err = h.horseService.DeleteHorse(uint(horseID))
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetHealthAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership first
	horse, err := h.horseService.GetHorse(uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	records, err := h.healthService.GetHealthRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Health records not found: %v", err)})
		return
	}
	SendSuccess(c, records)
}

func (h *Handler) AddHealthRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusBadRequest)
		return
	}

	record.HorseID = uint(horseID)
	record.Date = time.Now()

	err = h.healthService.AddHealthRecord(c.Request.Context(), &record)
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}
	SendSuccess(c, record)
}

func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	stage := c.Query("stage")
	guidelines, err := h.pregnancyService.GetPregnancyGuidelinesByStage(models.PregnancyStage(stage))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SendSuccess(c, guidelines)
}

func (h *Handler) GetBreedingCosts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	costs, err := h.breedingRepo.GetCosts(c.Request.Context(), uint(horseID))
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusNotFound)
		return
	}
	SendSuccess(c, costs)
}

func (h *Handler) AddBreedingCost(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var cost models.BreedingCost
	if err := c.ShouldBindJSON(&cost); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid breeding cost data: %v", err)})
		return
	}

	cost.HorseID = uint(horseID)
	cost.Date = time.Now()

	err = h.breedingRepo.Create(c.Request.Context(), &cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to add breeding cost: %v", err)})
		return
	}
	SendSuccess(c, cost)
}

func (h *Handler) GetHorseOffspring(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	offspring, err := h.horseRepo.GetOffspring(c.Request.Context(), uint(horseID))
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, offspring)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check cache first
	cacheKey := fmt.Sprintf("dashboard_stats_%s", userID)
	if cachedStats, found := h.cache.Get(cacheKey); found {
		c.JSON(http.StatusOK, cachedStats)
		return
	}

	stats, err := h.userService.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	// Cache the stats for 5 minutes
	h.cache.Set(cacheKey, stats, 5*time.Minute)

	SendSuccess(c, stats)
}

func (h *Handler) GetHorseFamilyTree(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	tree, err := h.horseService.GetFamilyTree(c.Request.Context(), uint(horseID))
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, tree)
}

// GetPregnancyProgress returns the due date and progress for a pregnant horse
func (h *Handler) GetPregnancyProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horseIDStr := c.Param("id")
	horseID, err := strconv.ParseInt(horseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	if !horse.IsPregnant || horse.ConceptionDate == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No active pregnancy for this horse"})
		return
	}

	// Get custom gestation days if provided, otherwise use default
	gestationDays := c.Query("gestationDays")
	customDays := 340 // Default gestation days
	if gestationDays != "" {
		if days, err := strconv.Atoi(gestationDays); err == nil && days > 0 {
			customDays = days
		}
	}

	dueDate := horse.ConceptionDate.AddDate(0, 0, customDays)
	daysSinceConception := int(time.Since(*horse.ConceptionDate).Hours() / 24)
	
	progress := float64(daysSinceConception) / float64(customDays) * 100
	if progress > 100 {
		progress = 100
	}
	
	daysRemaining := customDays - daysSinceConception
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// Determine stage
	stage := "Unknown"
	progressPercentage := float64(daysSinceConception) / float64(customDays)
	
	switch {
	case progressPercentage <= 0.33:
		stage = "Early Gestation"
	case progressPercentage <= 0.66:
		stage = "Mid Gestation"
	default:
		stage = "Late Gestation"
	}

	response := PregnancyProgress{
		DueDate:       dueDate,
		Progress:      progress,
		DaysRemaining: daysRemaining,
		Stage:         stage,
	}

	SendSuccess(c, response)
}

// GetPreFoalingChecklist returns the pre-foaling checklist for a horse
func (h *Handler) GetPreFoalingChecklist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	items, err := h.pregnancyService.GetPreFoalingChecklist(c.Request.Context(), uint(horseID))
	if err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, items)
}

// AddPreFoalingChecklistItem adds a new item to the checklist
func (h *Handler) AddPreFoalingChecklistItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	var item models.PreFoalingChecklistItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	item.HorseID = uint(horseID)
	if err := h.pregnancyService.AddPreFoalingChecklistItem(c.Request.Context(), &item); err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, item)
}

// UpdatePreFoalingChecklistItem updates an existing checklist item
func (h *Handler) UpdatePreFoalingChecklistItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid item ID"})
		return
	}

	var item models.PreFoalingChecklistItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Verify ownership through horse
	horse, err := h.horseRepo.GetByID(c.Request.Context(), item.HorseID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	item.ID = uint(itemID)
	if err := h.pregnancyService.UpdatePreFoalingChecklistItem(c.Request.Context(), &item); err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	SendSuccess(c, item)
}

// DeletePreFoalingChecklistItem deletes a checklist item
func (h *Handler) DeletePreFoalingChecklistItem(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid item ID"})
		return
	}

	// Get the checklist item to verify ownership
	item, err := h.pregnancyService.GetPreFoalingChecklistItem(c.Request.Context(), uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Checklist item not found"})
		return
	}

	// Verify ownership through horse
	horse, err := h.horseRepo.GetByID(c.Request.Context(), item.HorseID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	if err := h.pregnancyService.DeletePreFoalingChecklistItem(c.Request.Context(), uint(itemID)); err != nil {
		log.Printf("Error: %v", err)
		SendError(c, err, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) InitializePreFoalingChecklist(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	if err := h.pregnancyService.InitializePreFoalingChecklist(c.Request.Context(), uint(horseID)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to initialize checklist: %v", err)})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) GetHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id := c.Param("id")
	horseID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}

	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, horse)
}

func (h *Handler) AddHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse data"})
		return
	}

	horse.UserID = userID
	if err := h.horseService.CreateHorse(&horse); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to create horse: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, horse)
}

func (h *Handler) UpdateHorsePregnancyStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var req struct {
		IsPregnant     bool       `json:"is_pregnant"`
		ConceptionDate *time.Time `json:"conception_date,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Verify ownership
	horse, err := h.horseRepo.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	err = h.pregnancyService.UpdatePregnancyStatus(c.Request.Context(), uint(horseID), req.IsPregnant, req.ConceptionDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pregnancy status updated successfully"})
}

func (h *Handler) GetPregnantHorses(c *gin.Context) {
	// Extract user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Retrieve pregnant horses
	horses, err := h.horseRepo.GetPregnantHorses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, horses)
}

func (h *Handler) GetPregnancyStage(c *gin.Context) {
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	stage, err := h.pregnancyService.GetPregnancyStage(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stage": stage})
}

func (h *Handler) GetPregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	pregnancy, err := h.pregnancyService.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}
