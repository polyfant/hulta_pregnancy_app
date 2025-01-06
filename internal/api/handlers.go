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
	db                 *database.PostgresDB
	horseRepo       repository.HorseRepository
	breedingRepo    repository.BreedingRepository
}

// NewHandler creates a new Handler instance
func NewHandler(config HandlerConfig) *Handler {
	return &Handler{
		horseService:         config.HorseService,
		userService:          config.UserService,
		pregnancyService:     config.PregnancyService,
		healthService:        config.HealthService,
		cache:               config.Cache,
		db:                 config.Database,
		horseRepo:       config.HorseRepo,
		breedingRepo:    config.BreedingRepo,
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
		log.Printf("Error fetching horses: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to fetch horses: %v", err)})
		return
	}
	c.JSON(http.StatusOK, horses)
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
		log.Printf("Error deleting horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to delete horse: %v", err)})
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
	c.JSON(http.StatusOK, records)
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
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid health record data: %v", err)})
		return
	}

	record.HorseID = uint(horseID)
	record.Date = time.Now()

	// TODO: Move this to a health repository
	// err = h.healthRepo.Create(c.Request.Context(), &record)
	if err != nil {
		log.Printf("Error adding health record for horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to add health record: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	stage := c.Query("stage")
	guidelines, err := h.pregnancyService.GetPregnancyGuidelinesByStage(models.PregnancyStage(stage))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, guidelines)
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

	// TODO: Move this to a breeding repository
	// costs, err := h.breedingRepo.GetCosts(c.Request.Context(), uint(horseID))
	if err != nil {
		log.Printf("Error getting breeding costs for horse %d: %v", horseID, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Breeding costs not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, costs)
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

	// TODO: Move this to a breeding repository
	// err = h.breedingRepo.Create(c.Request.Context(), &cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to add breeding cost: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, cost)
}

func (h *Handler) GetHorseOffspring(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	offspring, err := h.db.GetOffspring(uint(horseID))
	if err != nil {
		log.Printf("Error getting offspring for horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to get offspring: %v", err)})
		return
	}

	c.JSON(http.StatusOK, offspring)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")

	// Check cache first
	if cachedStats, found := h.cache.Get(fmt.Sprintf("dashboard_stats_%s", userID)); found {
		c.JSON(http.StatusOK, cachedStats)
		return
	}

	stats, err := h.db.GetDashboardStats(userID)
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve dashboard stats"})
		return
	}

	// Cache dashboard stats for 30 minutes
	h.cache.Set(fmt.Sprintf("dashboard_stats_%s", userID), stats, 30*time.Minute)

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetHorseFamilyTree(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	tree, err := h.db.GetFamilyTree(horseID)
	if err != nil {
		log.Printf("Error getting family tree for horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to get family tree: %v", err)})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// PregnancyProgress represents the pregnancy progress response
type PregnancyProgress struct {
	DueDate       time.Time `json:"dueDate"`
	Progress      float64   `json:"progress"`
	DaysRemaining int       `json:"daysRemaining"`
	Stage         string    `json:"stage"`
}

// GetPregnancyProgress returns the due date and progress for a pregnant horse
func (h *Handler) GetPregnancyProgress(c *gin.Context) {
	horseIDStr := c.Param("id")
	horseID, err := strconv.ParseInt(horseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Horse not found"})
		return
	}

	if horse.ConceptionDate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No active pregnancy for this horse"})
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

	dueDate := time.Now().AddDate(0, 0, customDays)
	progress := 0.0
	daysRemaining := 0
	stage := "Unknown"

	if horse.ConceptionDate != nil {
		daysSinceConception := int(time.Since(*horse.ConceptionDate).Hours() / 24)
		progress = float64(daysSinceConception) / float64(customDays) * 100
		if progress > 100 {
			progress = 100
		}
		
		daysRemaining = customDays - daysSinceConception
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		// Determine stage
		stageThresholds := map[string]float64{
			"Early Gestation":  0.33,
			"Mid Gestation":    0.66,
			"Late Gestation":   1.0,
		}
		
		progressPercentage := float64(daysSinceConception) / float64(customDays)
		
		for stageName, threshold := range stageThresholds {
			if progressPercentage <= threshold {
				stage = stageName
				break
			}
		}
	}

	response := PregnancyProgress{
		DueDate:       dueDate,
		Progress:      progress,
		DaysRemaining: daysRemaining,
		Stage:         stage,
	}

	c.JSON(http.StatusOK, response)
}

// GetPreFoalingChecklist returns the pre-foaling checklist for a horse
func (h *Handler) GetPreFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	items, err := h.db.GetPreFoalingChecklist(uint(horseID))
	if err != nil {
		log.Printf("Error fetching pre-foaling checklist: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch checklist"})
		return
	}

	// If no items exist yet, return the default checklist
	if len(items) == 0 {
		defaultItems := make([]models.PreFoalingChecklistItem, len(models.DefaultPreFoalingChecklist))
		dueDate := time.Now().AddDate(0, 0, 7) // Default due date is 1 week from now
		
		for i, item := range models.DefaultPreFoalingChecklist {
			defaultItems[i] = models.PreFoalingChecklistItem{
				HorseID:     uint(horseID),
				Description: item.Description,
				IsCompleted: false,
				DueDate:     dueDate,
				Priority:    "MEDIUM",
			}
		}
		c.JSON(http.StatusOK, defaultItems)
		return
	}

	c.JSON(http.StatusOK, items)
}

// AddPreFoalingChecklistItem adds a new item to the checklist
func (h *Handler) AddPreFoalingChecklistItem(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var item models.PreFoalingChecklistItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	item.HorseID = uint(horseID)
	if err := h.db.AddPreFoalingChecklistItem(&item); err != nil {
		log.Printf("Error adding pre-foaling checklist item: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add checklist item"})
		return
	}

	c.JSON(http.StatusCreated, item)
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
	// TODO: Move this to a checklist repository
	if err := h.db.UpdatePreFoalingChecklistItem(&item); err != nil {
		log.Printf("Error updating pre-foaling checklist item: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update checklist item"})
		return
	}

	c.JSON(http.StatusOK, item)
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

	// TODO: Add ownership verification when we have the checklist repository
	// item, err := h.checklistRepo.GetByID(c.Request.Context(), uint(itemID))
	// if err != nil || item.Horse.UserID != userID {
	//     c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
	//     return
	// }

	// TODO: Move this to a checklist repository
	if err := h.db.DeletePreFoalingChecklistItem(itemID); err != nil {
		log.Printf("Error deleting pre-foaling checklist item: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete checklist item"})
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

	for _, item := range models.DefaultPreFoalingChecklist {
		checklistItem := models.PreFoalingChecklistItem{
			HorseID:     uint(horseID),
			Description: item.Description,
			Priority:    item.Priority,
			DueDate:     time.Now().AddDate(0, 0, 7), // Due in a week
		}
		// TODO: Move this to a checklist repository
		if err := h.db.AddPreFoalingChecklistItem(&checklistItem); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
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

	horse, err := h.horseService.GetHorse(uint(horseID))
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
	// Extract user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse horse ID from URL
	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid horse ID"})
		return
	}

	// Bind pregnancy status
	var pregnancyStatus struct {
		IsPregnant bool `json:"is_pregnant"`
	}
	if err := c.ShouldBindJSON(&pregnancyStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update pregnancy status
	if err := h.pregnancyService.UpdateHorsePregnancyStatus(
		c.Request.Context(), 
		uint(horseID), 
		userID, 
		pregnancyStatus.IsPregnant,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Horse pregnancy status updated"})
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

type HandlerConfig struct {
	Database         *database.PostgresDB
	UserService      *service.UserService
	HorseService     *service.HorseService
	PregnancyService *service.PregnancyService
	HealthService    *service.HealthService
	Cache            *cache.MemoryCache
	HorseRepo        repository.HorseRepository
	BreedingRepo     repository.BreedingRepository
}

func (h *Handler) handleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func (h *Handler) handleError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "error",
		"error":  err.Error(),
	})
}

func (h *Handler) handleResponse(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
