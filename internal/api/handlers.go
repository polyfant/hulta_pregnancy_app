package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/service"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/checklist"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/health"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/notification"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/pregnancy"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/privacy"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/vitals"
	"github.com/polyfant/hulta_pregnancy_app/internal/service/weather"
	"github.com/polyfant/hulta_pregnancy_app/internal/websocket"
	"gorm.io/gorm"
)

// Handlers contains all HTTP handlers and services
type Handlers struct {
	HorseService      service.HorseService
	UserService       service.UserService
	PregnancyService  service.PregnancyService
	HealthService     health.Service
	BreedingService   service.BreedingService
	WeatherService    weather.Service
	ChecklistService  checklist.Service
	PrivacyService    privacy.Service
	VitalsService     vitals.Service
	NotificationService notification.Service
	WebSocketHub      *websocket.Hub

	// Added back dependencies
	Cache    cache.Cache
	Database *gorm.DB

	// Specific handler instances
	PregnancyHandler *PregnancyHandler
	HealthHandler    *HealthHandler
	VitalsHandler    *VitalsHandler
	WeatherHandler   *WeatherHandler
	PrivacyHandler   *PrivacyHandler
}

// NewHandlers creates a new Handlers instance with all services and handlers
func NewHandlers(
	horseService service.HorseService,
	userService service.UserService,
	pregnancyService service.PregnancyService,
	healthService health.Service,
	breedingService service.BreedingService,
	weatherService weather.Service,
	checklistService checklist.Service,
	privacyService privacy.Service,
	vitalsService vitals.Service,
	notificationService notification.Service,
	webSocketHub *websocket.Hub,
	cache cache.Cache,
	database *gorm.DB,
) *Handlers {
	return &Handlers{
		HorseService:      horseService,
		UserService:       userService,
		PregnancyService:  pregnancyService,
		HealthService:     healthService,
		BreedingService:   breedingService,
		WeatherService:    weatherService,
		ChecklistService:  checklistService,
		PrivacyService:    privacyService,
		VitalsService:     vitalsService,
		NotificationService: notificationService,
		WebSocketHub:      webSocketHub,
		Cache:             cache,
		Database:          database,

		PregnancyHandler: NewPregnancyHandler(pregnancyService, weatherService, checklistService),
		HealthHandler:    NewHealthHandler(healthService),
		VitalsHandler:    NewVitalsHandler(vitalsService),
		WeatherHandler:   NewWeatherHandler(weatherService, notificationService),
		PrivacyHandler:   NewPrivacyHandler(privacyService),
	}
}

// ListHorses handles GET /horses
func (h *Handlers) ListHorses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horses, err := h.HorseService.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, horses)
}

// AddHorse handles POST /horses
func (h *Handlers) AddHorse(c *gin.Context) {
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
	if err := h.HorseService.Create(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, horse)
}

// GetHorse handles GET /horses/:id
func (h *Handlers) GetHorse(c *gin.Context) {
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

	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(id))
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
func (h *Handlers) UpdateHorse(c *gin.Context) {
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
	existingHorse, err := h.HorseService.GetByID(c.Request.Context(), horse.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if existingHorse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	if err := h.HorseService.Update(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, horse)
}

// DeleteHorse handles DELETE /horses/:id
func (h *Handlers) DeleteHorse(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	if err := h.HorseService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetHealthRecords handles GET /horses/:id/health
func (h *Handlers) GetHealthRecords(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	records, err := h.HealthService.GetRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// AddHealthRecord handles POST /horses/:id/health
func (h *Handlers) AddHealthRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.HealthService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// UpdateHealthRecord handles PUT /horses/:id/health/:recordId
func (h *Handlers) UpdateHealthRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.HealthService.UpdateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteHealthRecord handles DELETE /horses/:id/health/:recordId
func (h *Handlers) DeleteHealthRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.HealthService.DeleteRecord(c.Request.Context(), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserProfile handles GET /user/profile
func (h *Handlers) GetUserProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	profile, err := h.UserService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateUserProfile handles PUT /user/profile
func (h *Handlers) UpdateUserProfile(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.UserService.UpdateProfile(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetDashboardStats handles GET /dashboard
func (h *Handlers) GetDashboardStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	stats, err := h.UserService.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetPregnancyGuidelines handles GET /horses/:id/pregnancy/guidelines
func (h *Handlers) GetPregnancyGuidelines(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	stage := models.PregnancyStage(c.Query("stage"))
	guidelines, err := h.PregnancyService.GetGuidelines(c.Request.Context(), stage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, guidelines)
}

// StartPregnancyTracking handles POST /horses/:id/pregnancy/start
func (h *Handlers) StartPregnancyTracking(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.PregnancyService.StartTracking(c.Request.Context(), uint(horseID), start); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetEvents handles GET /horses/:id/pregnancy/events
func (h *Handlers) GetEvents(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	events, err := h.PregnancyService.GetEvents(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// AddPregnancyEvent handles POST /horses/:id/pregnancy/events
func (h *Handlers) AddPregnancyEvent(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.PregnancyService.AddPregnancyEvent(c.Request.Context(), &event); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetBreedingRecords handles GET /horses/:id/breeding
func (h *Handlers) GetBreedingRecords(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	records, err := h.BreedingService.GetRecords(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// AddBreedingRecord handles POST /horses/:id/breeding
func (h *Handlers) AddBreedingRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.BreedingService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// UpdateBreedingRecord handles PUT /horses/:id/breeding/:recordId
func (h *Handlers) UpdateBreedingRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.BreedingService.UpdateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// DeleteBreedingRecord handles DELETE /horses/:id/breeding/:recordId
func (h *Handlers) DeleteBreedingRecord(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
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

	if err := h.BreedingService.DeleteRecord(c.Request.Context(), uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPregnancy handles GET /horses/:id/pregnancy
func (h *Handlers) GetPregnancy(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	pregnancy, err := h.PregnancyService.GetPregnancy(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, pregnancy)
}

// GetPregnancyStatus handles GET /horses/:id/pregnancy/status
func (h *Handlers) GetPregnancyStatus(c *gin.Context) {
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
	horse, err := h.HorseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied"})
		return
	}

	status, err := h.PregnancyService.GetStatus(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

// GetHorseService returns the horse service
func (h *Handlers) GetHorseService() service.HorseService {
    return h.HorseService
}

// GetUserService returns the user service
func (h *Handlers) GetUserService() service.UserService {
    return h.UserService
}

// GetPregnancyService returns the pregnancy service
func (h *Handlers) GetPregnancyService() service.PregnancyService {
    return h.PregnancyService
}

// GetHealthService returns the health service
func (h *Handlers) GetHealthService() health.Service {
    return h.HealthService
}

// GetBreedingService returns the breeding service
func (h *Handlers) GetBreedingService() service.BreedingService {
    return h.BreedingService
}

// GetChecklistService returns the checklist service
func (h *Handlers) GetChecklistService() checklist.Service {
    return h.ChecklistService
}

// GetPrivacyService returns the privacy service
func (h *Handlers) GetPrivacyService() privacy.Service {
    return h.PrivacyService
}

// GetVitalsService returns the vitals service
func (h *Handlers) GetVitalsService() vitals.Service {
    return h.VitalsService
}

// GetNotificationService returns the notification service
func (h *Handlers) GetNotificationService() notification.Service {
    return h.NotificationService
}

// GetWeatherService returns the weather service
func (h *Handlers) GetWeatherService() weather.Service {
    return h.WeatherService
}

// GetWebSocketHandler returns the websocket handler
func (h *Handlers) GetWebSocketHandler() *websocket.Hub {
    return h.WebSocketHub
}
