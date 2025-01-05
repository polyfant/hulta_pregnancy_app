package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/pregnancy"
)

// Database interface defines the methods our handler needs
type Database interface {
	GetHorse(id int64) (models.Horse, error)
	GetAllHorses() ([]models.Horse, error)
	AddHorse(horse *models.Horse) error
	DB() *gorm.DB
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	AutoMigrate(dst ...interface{}) error
	GetHealthRecords(horseID int64) ([]models.HealthRecord, error)
	GetBreedingCosts(horseID uint) ([]models.BreedingCost, error)
	AddHealthRecord(record *models.HealthRecord) error
	AddBreedingCost(cost *models.BreedingCost) error
	GetOffspring(horseID int64) ([]models.Horse, error)
	GetDashboardStats(userID string) (*models.DashboardStats, error)
	GetFamilyTree(horseID int64) (*database.FamilyTree, error)
	GetPreFoalingChecklist(horseID int64) ([]models.PreFoalingChecklistItem, error)
	AddPreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error
	UpdatePreFoalingChecklistItem(item *models.PreFoalingChecklistItem) error
	DeletePreFoalingChecklistItem(itemID int64) error
	DeleteHorse(horseID int64) error
	GetPregnancyGuidelines(stage models.PregnancyStage) ([]models.PregnancyGuideline, error)
}

// Handler handles HTTP requests
type Handler struct {
	horseService     *service.HorseService
	expenseService   *service.ExpenseService
	userService      *service.UserService
	pregnancyService *pregnancy.Service
	cache            *cache.MemoryCache
	db               Database
	pregnancyRepo    repository.PregnancyRepository
}

// NewHandler creates a new Handler instance
func NewHandler(
	db Database, 
	horseRepo repository.HorseRepository,
	expenseRepo repository.ExpenseRepository,
	recurringExpenseRepo repository.RecurringExpenseRepository,
	userRepo repository.UserRepository,
	caches ...*cache.MemoryCache,
) *Handler {
	var memoryCache *cache.MemoryCache
	if len(caches) > 0 {
		memoryCache = caches[0]
	} else {
		memoryCache = cache.NewMemoryCache()
	}

	pregnancyRepo := repository.NewPregnancyRepository(db.DB())

	return &Handler{
		horseService: service.NewHorseService(horseRepo, expenseRepo),
		expenseService: service.NewExpenseService(expenseRepo, recurringExpenseRepo),
		userService: service.NewUserService(userRepo, horseRepo, expenseRepo),
		pregnancyService: pregnancy.NewService(pregnancyRepo),
		cache: memoryCache,
		db: db,
		pregnancyRepo: pregnancyRepo,
	}
}

func (h *Handler) ListHorses(c *gin.Context) {
	horses, err := h.db.GetAllHorses()
	if err != nil {
		log.Printf("Error fetching horses: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to fetch horses: %v", err)})
		return
	}
	c.JSON(http.StatusOK, horses)
}

func (h *Handler) DeleteHorse(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	err = h.db.DeleteHorse(horseID)
	if err != nil {
		log.Printf("Error deleting horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to delete horse: %v", err)})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetHealthAssessment(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	records, err := h.db.GetHealthRecords(horseID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Health records not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, records)
}

func (h *Handler) AddHealthRecord(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
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

	err = h.db.AddHealthRecord(&record)
	if err != nil {
		log.Printf("Error adding health record for horse %d: %v", horseID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to add health record: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, record)
}

func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	stageStr := c.Query("stage")
	stage := models.PregnancyStage(stageStr)
	guidelines := h.pregnancyService.GetPregnancyGuidelinesByStage(stage)
	
	if len(guidelines) == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("No guidelines found for stage: %s", stage)})
		return
	}
	
	c.JSON(http.StatusOK, guidelines)
}

func (h *Handler) GetBreedingCosts(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	costs, err := h.db.GetBreedingCosts(uint(horseID))
	if err != nil {
		log.Printf("Error getting breeding costs for horse %d: %v", horseID, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Breeding costs not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, costs)
}

func (h *Handler) AddBreedingCost(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var cost models.BreedingCost
	if err := c.ShouldBindJSON(&cost); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid breeding cost data: %v", err)})
		return
	}

	cost.HorseID = uint(horseID)
	cost.Date = time.Now()

	err = h.db.AddBreedingCost(&cost)
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

	offspring, err := h.db.GetOffspring(horseID)
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
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	stage, err := h.pregnancyService.GetPregnancyStage(c.Request.Context(), horseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		log.Printf("Error fetching horse: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch horse details"})
		return
	}

	if !horse.IsPregnant || horse.ConceptionDate == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Horse is not pregnant"})
		return
	}

	// Get custom gestation days if provided, otherwise use default
	gestationDays := c.Query("gestationDays")
	customDays := pregnancy.DefaultGestationDays
	if gestationDays != "" {
		if days, err := strconv.Atoi(gestationDays); err == nil && days > 0 {
			customDays = days
		}
	}

	dueDate := pregnancy.CalculateDueDate(*horse.ConceptionDate, customDays)
	progress, daysRemaining := pregnancy.CalculateGestationProgress(*horse.ConceptionDate, customDays)

	response := PregnancyProgress{
		DueDate:       dueDate,
		Progress:      progress,
		DaysRemaining: daysRemaining,
		Stage:         string(stage),
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

	items, err := h.db.GetPreFoalingChecklist(horseID)
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

	item.ID = uint(itemID)
	if err := h.db.UpdatePreFoalingChecklistItem(&item); err != nil {
		log.Printf("Error updating pre-foaling checklist item: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update checklist item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DeletePreFoalingChecklistItem deletes a checklist item
func (h *Handler) DeletePreFoalingChecklistItem(c *gin.Context) {
	itemID, err := strconv.ParseInt(c.Param("itemId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid item ID"})
		return
	}

	if err := h.db.DeletePreFoalingChecklistItem(itemID); err != nil {
		log.Printf("Error deleting pre-foaling checklist item: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete checklist item"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) InitializePreFoalingChecklist(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	for _, item := range models.DefaultPreFoalingChecklist {
		checklistItem := models.PreFoalingChecklistItem{
			HorseID:     uint(horseID),
			Description: item.Description,
			Priority:    item.Priority,
			DueDate:     time.Now().AddDate(0, 0, 7), // Due in a week
		}
		if err := h.db.AddPreFoalingChecklistItem(&checklistItem); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetHorse(c *gin.Context) {
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

	// Retrieve horse details
	horseDetails, err := h.horseService.GetHorseWithDetails(c.Request.Context(), uint(horseID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, horseDetails)
}

func (h *Handler) AddHorse(c *gin.Context) {
	// Extract user ID from context
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Bind horse data from request
	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID
	horse.UserID = userID

	// Create horse
	if err := h.horseService.CreateHorse(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if err := h.horseService.UpdateHorsePregnancyStatus(
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
	horses, err := h.horseService.GetPregnantHorses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, horses)
}
