package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/logger"
	"github.com/polyfant/horse_tracking/internal/models"
	"github.com/polyfant/horse_tracking/internal/service/breeding"
	"github.com/polyfant/horse_tracking/internal/service/health"
	"github.com/polyfant/horse_tracking/internal/service/pregnancy"
)

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error string `json:"error"`
}

func getUserIDFromContext(c *gin.Context) int64 {
	if id, exists := c.Get("userID"); exists {
		if userID, ok := id.(int64); ok {
			return userID
		}
	}
	return 0
}

type Handler struct {
	db              models.DataStore
	breedingService *breeding.Service
	healthService   *health.HealthService
	pregnancyService *pregnancy.Service
}

func NewHandler(db models.DataStore) *Handler {
	return &Handler{
		db:              db,
		breedingService: breeding.NewService(db),
		healthService:   health.NewHealthService(db),
		pregnancyService: pregnancy.NewService(db),
	}
}

// @Summary List all horses
// @Description Get a list of all horses
// @Tags horses
// @Produce json
// @Success 200 {array} models.Horse
// @Router /horses [get]
func (h *Handler) ListHorses(c *gin.Context) {
	horses, err := h.db.GetAllHorses()
	if err != nil {
		logger.Error(err, "Failed to get horses", nil)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get horses"})
		return
	}

	// Calculate additional information for pregnant horses
	for i := range horses {
		if horses[i].ConceptionDate != nil {
			calculator := pregnancy.NewPregnancyCalculator(*horses[i].ConceptionDate)
			dueDate := calculator.GetDueDate()
			stage := calculator.GetStage()
			horses[i].DueDate = &dueDate
			horses[i].PregnancyStage = string(stage)
			horses[i].PregnancyProgress = calculator.GetProgressPercentage()
		}
	}

	c.JSON(http.StatusOK, horses)
}

// @Summary Get horse details
// @Description Get details of a specific horse
// @Tags horses
// @Accept json
// @Produce json
// @Param id path int true "Horse ID"
// @Success 200 {object} models.Horse
// @Failure 404 {object} ErrorResponse
// @Router /horses/{id} [get]
func (h *Handler) GetHorse(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}

	// Calculate additional information if horse is pregnant
	if horse.ConceptionDate != nil {
		calculator := pregnancy.NewPregnancyCalculator(*horse.ConceptionDate)
		dueDate := calculator.GetDueDate()
		stage := calculator.GetStage()
		horse.DueDate = &dueDate
		horse.PregnancyStage = string(stage)
		horse.PregnancyProgress = calculator.GetProgressPercentage()
	}

	c.JSON(http.StatusOK, horse)
}

// @Summary Add new horse
// @Description Add a new horse to the database
// @Tags horses
// @Accept json
// @Produce json
// @Param horse body models.Horse true "Horse object"
// @Success 201 {object} models.Horse
// @Failure 400 {object} ErrorResponse
// @Router /horses [post]
func (h *Handler) AddHorse(c *gin.Context) {
	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse data"})
		return
	}

	if err := h.db.AddHorse(&horse); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add horse"})
		return
	}

	// Calculate additional information if horse is pregnant
	if horse.ConceptionDate != nil {
		calculator := pregnancy.NewPregnancyCalculator(*horse.ConceptionDate)
		dueDate := calculator.GetDueDate()
		stage := calculator.GetStage()
		horse.DueDate = &dueDate
		horse.PregnancyStage = string(stage)
		horse.PregnancyProgress = calculator.GetProgressPercentage()
	}

	c.JSON(http.StatusCreated, horse)
}

// @Summary Get health assessment
// @Description Get comprehensive health assessment for a specific horse
// @Tags health
// @Produce json
// @Param id path int true "Horse ID"
// @Success 200 {object} health.HealthSummary
// @Router /horses/{id}/health [get]
func (h *Handler) GetHealthAssessment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(id)
	if err != nil {
		logger.Error(err, "Failed to get horse", map[string]interface{}{"id": id})
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}

	assessment := h.healthService.GetHealthSummary(horse)
	c.JSON(http.StatusOK, assessment)
}

// @Summary Add health record
// @Description Add a new health record for a specific horse
// @Tags health
// @Accept json
// @Produce json
// @Param id path int true "Horse ID"
// @Param record body models.HealthRecord true "Health record"
// @Success 201 {object} models.HealthRecord
// @Router /horses/{id}/health-records [post]
func (h *Handler) AddHealthRecord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid health record data"})
		return
	}

	record.HorseID = id
	if err := h.db.AddHealthRecord(&record); err != nil {
		logger.Error(err, "Failed to add health record", map[string]interface{}{"horseID": id})
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add health record"})
		return
	}
	c.JSON(http.StatusCreated, record)
}

// @Summary Get pregnancy guidelines
// @Description Get detailed pregnancy guidelines for a specific horse
// @Tags pregnancy
// @Produce json
// @Param id path int true "Horse ID"
// @Success 200 {object} pregnancy.PregnancyGuidelines
// @Router /horses/{id}/pregnancy-guidelines [get]
func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.db.GetHorse(id)
	if err != nil {
		logger.Error(err, "Failed to get horse", map[string]interface{}{"id": id})
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Horse not found"})
		return
	}

	guidelines, err := h.pregnancyService.GetPregnancyGuidelines(horse)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, guidelines)
}

// @Summary Get breeding costs
// @Description Get all breeding costs for a specific horse
// @Tags breeding
// @Produce json
// @Param id path int true "Horse ID"
// @Success 200 {array} models.BreedingCost
// @Router /horses/{id}/breeding-costs [get]
func (h *Handler) GetBreedingCosts(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	costs, err := h.db.GetBreedingCosts(id)
	if err != nil {
		logger.Error(err, "Failed to get breeding costs", map[string]interface{}{"horseID": id})
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get breeding costs"})
		return
	}
	c.JSON(http.StatusOK, costs)
}

// @Summary Add breeding cost
// @Description Add a new breeding cost for a specific horse
// @Tags breeding
// @Accept json
// @Produce json
// @Param id path int true "Horse ID"
// @Param cost body models.BreedingCost true "Breeding cost"
// @Success 201 {object} models.BreedingCost
// @Router /horses/{id}/breeding-costs [post]
func (h *Handler) AddBreedingCost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	var cost models.BreedingCost
	if err := c.ShouldBindJSON(&cost); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid breeding cost data"})
		return
	}

	cost.HorseID = id
	if err := h.db.AddBreedingCost(&cost); err != nil {
		logger.Error(err, "Failed to add breeding cost", map[string]interface{}{"horseID": id})
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to add breeding cost"})
		return
	}
	c.JSON(http.StatusCreated, cost)
}
