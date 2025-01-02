package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/database"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
)

// Handler handles HTTP requests
type Handler struct {
	db              *database.SQLiteStore
	pregnancyService *pregnancy.Service
}

// NewHandler creates a new Handler instance
func NewHandler(db *database.SQLiteStore) *Handler {
	return &Handler{
		db:              db,
		pregnancyService: pregnancy.NewService(db),
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
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

func (h *Handler) GetHorse(c *gin.Context) {
	id := c.Param("id")
	horseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		log.Printf("Error getting horse %d: %v", horseID, err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("Horse not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, horse)
}

func (h *Handler) AddHorse(c *gin.Context) {
	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid horse data: %v", err)})
		return
	}

	err := h.db.AddHorse(&horse)
	if err != nil {
		log.Printf("Error adding horse: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to add horse: %v", err)})
		return
	}
	c.JSON(http.StatusCreated, horse)
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
		log.Printf("Error getting health records for horse %d: %v", horseID, err)
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

	record.HorseID = horseID
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
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("Invalid horse ID: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(horseID)
	if err != nil {
		log.Printf("Error fetching horse: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch horse details"})
		return
	}

	// If stage is provided in query params, use GetPregnancyGuidelinesByStage
	stageStr := c.Query("stage")
	if stageStr != "" {
		stage := models.PregnancyStage(stageStr)
		guidelines, err := h.pregnancyService.GetPregnancyGuidelinesByStage(stage)
		if err != nil {
			log.Printf("Error getting pregnancy guidelines by stage: %v", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, guidelines)
		return
	}

	// Otherwise, get guidelines based on horse's pregnancy status
	stage := h.pregnancyService.GetPregnancyStage(horse)
	guidelines, err := h.pregnancyService.GetPregnancyGuidelinesByStage(stage)
	if err != nil {
		log.Printf("Error getting pregnancy guidelines for horse: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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

	costs, err := h.db.GetBreedingCosts(horseID)
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
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Invalid breeding cost data: %v", err)})
		return
	}

	cost.HorseID = horseID
	cost.Date = time.Now()

	err = h.db.AddBreedingCost(&cost)
	if err != nil {
		log.Printf("Error adding breeding cost for horse %d: %v", horseID, err)
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
	stats, err := h.db.GetDashboardStats()
	if err != nil {
		log.Printf("Error getting dashboard stats: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to get dashboard stats: %v", err)})
		return
	}
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
	stage := h.pregnancyService.GetPregnancyStage(horse)

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
		
		for i, desc := range models.DefaultPreFoalingChecklist {
			defaultItems[i] = models.PreFoalingChecklistItem{
				HorseID:     horseID,
				Description: desc,
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

	item.HorseID = horseID
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

	item.ID = itemID
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
