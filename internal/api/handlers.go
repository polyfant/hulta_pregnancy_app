package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/repository"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"gorm.io/gorm"
)

// Handler handles HTTP requests
type Handler struct {
	horseService     service.HorseService
	userService      service.UserService
	pregnancyService service.PregnancyService
	healthService    service.HealthService
	breedingService  service.BreedingService
	weatherService   service.WeatherService
	pregnancyHandler *PregnancyHandler
	cache            cache.Cache
	db               *gorm.DB
	horseRepo        repository.HorseRepository
	breedingRepo     repository.BreedingRepository
}

// HandlerConfig defines the configuration for creating a new handler
type HandlerConfig struct {
	Database         *gorm.DB
	UserService      service.UserService
	HorseService     service.HorseService
	PregnancyService service.PregnancyService
	HealthService    service.HealthService
	BreedingService  service.BreedingService
	WeatherService   service.WeatherService
	Cache            cache.Cache
	HorseRepo        repository.HorseRepository
	BreedingRepo     repository.BreedingRepository
}

// NewHandler creates a new handler instance
func NewHandler(config HandlerConfig) *Handler {
	h := &Handler{
		horseService:     config.HorseService,
		userService:      config.UserService,
		pregnancyService: config.PregnancyService,
		healthService:    config.HealthService,
		breedingService:  config.BreedingService,
		weatherService:   config.WeatherService,
		cache:            config.Cache,
		db:               config.Database,
		horseRepo:        config.HorseRepo,
		breedingRepo:     config.BreedingRepo,
	}
	
	h.pregnancyHandler = NewPregnancyHandler(config.PregnancyService, config.WeatherService)
	return h
}

// Start starts the HTTP server
func (h *Handler) Start(port string) error {
	router := gin.Default()
	SetupRouter(router, h)
	return router.Run(":" + port)
}

// ListHorses handles GET /horses
func (h *Handler) ListHorses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horses, err := h.horseService.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, horses)
}

// AddHorse handles POST /horses
func (h *Handler) AddHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}
	
	horse.UserID = userID
	if err := h.horseService.Create(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, horse)
}

// GetHorse handles GET /horses/:id
func (h *Handler) GetHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	horse, err := h.horseService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	c.JSON(http.StatusOK, horse)
}

// UpdateHorse handles PUT /horses/:id
func (h *Handler) UpdateHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var horse models.Horse
	if err := c.ShouldBindJSON(&horse); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Verify ownership
	existingHorse, err := h.horseService.GetByID(c.Request.Context(), horse.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if existingHorse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	if err := h.horseService.Update(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, horse)
}

// DeleteHorse handles DELETE /horses/:id
func (h *Handler) DeleteHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	if err := h.horseService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHealthRecords handles GET /horses/:id/health
func (h *Handler) GetHealthRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	records, err := h.healthService.GetRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// AddHealthRecord handles POST /horses/:id/health
func (h *Handler) AddHealthRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	if err := h.healthService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// UpdateHealthRecord handles PUT /horses/:id/health/:recordId
func (h *Handler) UpdateHealthRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	if err := h.healthService.UpdateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteHealthRecord handles DELETE /horses/:id/health/:recordId
func (h *Handler) DeleteHealthRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid record ID"})
		return
	}

	if err := h.healthService.DeleteRecord(c.Request.Context(), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserProfile handles GET /user/profile
func (h *Handler) GetUserProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserProfile handles PUT /user/profile
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetDashboardStats handles GET /dashboard
func (h *Handler) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	stats, err := h.userService.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPregnancyGuidelines handles GET /horses/:id/pregnancy/guidelines
func (h *Handler) GetPregnancyGuidelines(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	stage := models.PregnancyStage(c.Query("stage"))
	guidelines, err := h.pregnancyService.GetGuidelines(c.Request.Context(), stage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, guidelines)
}

// StartPregnancyTracking handles POST /horses/:id/pregnancy/start
func (h *Handler) StartPregnancyTracking(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var start models.PregnancyStart
	if err := c.ShouldBindJSON(&start); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.pregnancyService.StartTracking(c.Request.Context(), uint(horseID), start); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetPregnancyEvents handles GET /horses/:id/pregnancy/events
func (h *Handler) GetPregnancyEvents(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	events, err := h.pregnancyService.GetPregnancyEvents(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// AddPregnancyEvent handles POST /horses/:id/pregnancy/events
func (h *Handler) AddPregnancyEvent(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var event models.PregnancyEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	event.PregnancyID = uint(horseID)
	event.UserID = userID

	if err := h.pregnancyService.AddPregnancyEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetBreedingRecords handles GET /horses/:id/breeding
func (h *Handler) GetBreedingRecords(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	records, err := h.breedingService.GetRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// AddBreedingRecord handles POST /horses/:id/breeding
func (h *Handler) AddBreedingRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	if err := h.breedingService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// UpdateBreedingRecord handles PUT /horses/:id/breeding/:recordId
func (h *Handler) UpdateBreedingRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	record.HorseID = uint(horseID)
	record.UserID = userID

	if err := h.breedingService.UpdateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteBreedingRecord handles DELETE /horses/:id/breeding/:recordId
func (h *Handler) DeleteBreedingRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid record ID"})
		return
	}

	if err := h.breedingService.DeleteRecord(c.Request.Context(), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPregnancy handles GET /horses/:id/pregnancy
func (h *Handler) GetPregnancy(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	pregnancy, err := h.pregnancyService.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

// GetPregnancyStatus handles GET /horses/:id/pregnancy/status
func (h *Handler) GetPregnancyStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	status, err := h.pregnancyService.GetStatus(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetHorseService returns the horse service
func (h *Handler) GetHorseService() service.HorseService {
    return h.horseService
}

// GetUserService returns the user service
func (h *Handler) GetUserService() service.UserService {
    return h.userService
}

// GetPregnancyService returns the pregnancy service
func (h *Handler) GetPregnancyService() service.PregnancyService {
    return h.pregnancyService
}

// GetHealthService returns the health service
func (h *Handler) GetHealthService() service.HealthService {
    return h.healthService
}

// GetBreedingService returns the breeding service
func (h *Handler) GetBreedingService() service.BreedingService {
    return h.breedingService
}
