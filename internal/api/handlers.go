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
	"github.com/polyfant/hulta_pregnancy_app/internal/service/breeding"
)

// Handler handles HTTP requests
type Handler struct {
	horseService         *service.HorseService
	userService          *service.UserService
	pregnancyService     *service.PregnancyService
	healthService        *service.HealthService
	breedingService      *breeding.BreedingService
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
		breedingService:  config.BreedingService,
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

	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
	// Return array directly
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

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Validate record type
	if record.Type == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Health record type is required"})
		return
	}

	// Validate date is not in future
	if record.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Health record date cannot be in the future"})
		return
	}

	// Verify ownership
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	err = h.healthService.AddHealthRecord(c.Request.Context(), &record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Return 201 Created for successful creation
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	// Check if specific horse ID is provided
	if horseID := c.Param("id"); horseID != "" {
		id, err := strconv.ParseInt(horseID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
			return
		}

		// Verify horse exists
		horse, err := h.horseService.GetHorse(c.Request.Context(), uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
			return
		}
		if horse.UserID != c.GetString("user_id") {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
			return
		}
	}

	stage := c.Query("stage")
	guidelines, err := h.pregnancyService.GetPregnancyGuidelinesByStage(models.PregnancyStage(stage))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Return array directly instead of using SendSuccess
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}
	// Return array directly
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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

	// Return 201 Created instead of 200 OK
	c.JSON(http.StatusCreated, cost)
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, []models.Horse{})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Return array directly
	c.JSON(http.StatusOK, offspring)
}

func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get total horses
	horses, err := h.horseService.ListHorsesByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get horses"})
		return
	}

	// Get pregnant horses
	pregnantHorses, err := h.horseRepo.GetPregnantHorses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get pregnant horses"})
		return
	}

	// Calculate stats
	stats := gin.H{
		"total_horses":     len(horses),
		"pregnant_horses":  len(pregnantHorses),
		"upcoming_births": 0, // TODO: Calculate from due dates
		"recent_activity": 0, // TODO: Get from activity log
	}

	c.JSON(http.StatusOK, stats)
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Get family tree
	tree, err := h.horseService.GetFamilyTree(c.Request.Context(), uint(horseID))
	if err != nil {
		if err.Error() == "record not found" {
			// Return empty tree structure instead of 404
			c.JSON(http.StatusOK, gin.H{
				"horse": horse,
				"parents": gin.H{
					"mother": nil,
					"father": nil,
				},
				"offspring": []models.Horse{},
				"siblings":  []models.Horse{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), item.HorseID)
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), item.HorseID)
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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

	// Validate required fields
	if horse.Name == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Horse name is required"})
		return
	}

	// Validate gender
	if horse.Gender != models.GenderMare && horse.Gender != models.GenderStallion {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid gender. Must be MARE or STALLION"})
		return
	}

	// Validate birth date not in future
	if !horse.BirthDate.IsZero() && horse.BirthDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Birth date cannot be in the future"})
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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

	// Get horse first to verify it exists
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		// Return 404 when horse not found
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}

	// If horse is not pregnant, return empty status with 200
	if !horse.IsPregnant {
		c.JSON(http.StatusOK, gin.H{
			"is_pregnant": false,
			"stage": "NONE",
		})
		return
	}

	pregnancy, err := h.pregnancyService.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, gin.H{
				"is_pregnant": false,
				"stage": "NONE",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

func (h *Handler) AddBreedingRecord(c *gin.Context) {
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

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request data"})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	if err := h.breedingService.AddBreedingRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetBreedingRecords(c *gin.Context) {
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

	// Verify horse exists and ownership
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	records, err := h.breedingService.GetBreedingRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, []models.BreedingRecord{})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *Handler) GetHealthRecords(c *gin.Context) {
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
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
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
		// Return empty array instead of 404 when no records found
		if err.Error() == "record not found" {
			c.JSON(http.StatusOK, []models.HealthRecord{})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to get health records: %v", err)})
		return
	}

	c.JSON(http.StatusOK, records)
}

func (h *Handler) StartPregnancyTracking(c *gin.Context) {
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

	var data struct {
		ConceptionDate time.Time `json:"conception_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Conception date is required"})
		return
	}

	// Validate conception date
	if data.ConceptionDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Conception date cannot be in the future"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	// Start pregnancy tracking
	err = h.pregnancyService.StartPregnancy(c.Request.Context(), uint(horseID), userID, data.ConceptionDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// Return 201 Created instead of 200 OK
	c.JSON(http.StatusCreated, gin.H{"message": "Pregnancy tracking started successfully"})
}

func (h *Handler) RecordPreFoalingSign(c *gin.Context) {
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

	var sign models.PreFoalingSign
	if err := c.ShouldBindJSON(&sign); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request data"})
		return
	}

	// Validate required fields
	if sign.Description == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Description is required"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	sign.HorseID = uint(horseID)
	sign.Date = time.Now()

	if err := h.pregnancyService.AddPreFoalingSign(c.Request.Context(), &sign); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sign)
}

func (h *Handler) EndPregnancyTracking(c *gin.Context) {
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

	var data struct {
		Outcome     string    `json:"outcome" binding:"required"`
		FoalingDate time.Time `json:"foalingDate" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Outcome and foaling date are required"})
		return
	}

	// Validate end date
	if data.FoalingDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Foaling date cannot be in the future"})
		return
	}

	// Map outcome to status
	status := ""
	switch data.Outcome {
	case "FOALED":
		status = string(models.PregnancyStatusCompleted)
	case "ABORTED":
		status = string(models.PregnancyStatusAborted)
	default:
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid outcome. Must be FOALED or ABORTED"})
		return
	}

	// Verify horse ownership and pregnancy
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	if !horse.IsPregnant {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Horse is not pregnant"})
		return
	}

	// End pregnancy with foaling date
	err = h.pregnancyService.EndPregnancy(c.Request.Context(), uint(horseID), status, data.FoalingDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pregnancy tracking ended successfully"})
}

func (h *Handler) GetPregnancyEvents(c *gin.Context) {
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

	// Verify horse exists and ownership
	horse, err := h.horseService.GetHorse(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	events, err := h.pregnancyService.GetPregnancyEvents(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}
