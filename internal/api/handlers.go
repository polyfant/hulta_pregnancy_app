package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
	"github.com/polyfant/hulta_pregnancy_app/internal/cache"
	"github.com/polyfant/hulta_pregnancy_app/internal/config"
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
	growthService    service.GrowthService
	expenseService   service.ExpenseService
	cache            cache.Cache
	db               *gorm.DB
	horseRepo        repository.HorseRepository
	breedingRepo     repository.BreedingRepository
	growthRepo       repository.GrowthRepository
	config           HandlerConfig
	growthHandler    *GrowthHandler
	validate         *validator.Validate
}

// HandlerConfig defines the configuration for creating a new handler
type HandlerConfig struct {
	Database         *gorm.DB
	UserService      service.UserService
	HorseService     service.HorseService
	PregnancyService service.PregnancyService
	HealthService    service.HealthService
	BreedingService  service.BreedingService
	GrowthService    service.GrowthService
	ExpenseService   service.ExpenseService
	Cache            cache.Cache
	HorseRepo        repository.HorseRepository
	BreedingRepo     repository.BreedingRepository
	GrowthRepo       repository.GrowthRepository
	Auth0            config.Auth0Config
}

// NewHandler creates a new handler instance
func NewHandler(config HandlerConfig) *Handler {
	// Create growth handler
	validate := validator.New()
	growthHandler := NewGrowthHandler(config.GrowthService, validate)

	return &Handler{
		horseService:     config.HorseService,
		userService:      config.UserService,
		pregnancyService: config.PregnancyService,
		healthService:    config.HealthService,
		breedingService:  config.BreedingService,
		growthService:    config.GrowthService,
		expenseService:   config.ExpenseService,
		cache:            config.Cache,
		db:               config.Database,
		horseRepo:        config.HorseRepo,
		breedingRepo:     config.BreedingRepo,
		growthRepo:       config.GrowthRepo,
		config:           config,
		growthHandler:    growthHandler,
		validate:         validate,
	}
}

// Start starts the HTTP server
func (h *Handler) Start(port string) error {
	router := gin.Default()
	router = SetupRouter(router, h)
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
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}
	
	horse.UserID = userID

	if err := h.validate.Struct(&horse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	if !horse.ValidateGender() {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid gender specified."})
		return
	}
	if !horse.ValidatePregnancy() {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid pregnancy data: If pregnant, gender must be MARE, conception date must be set, valid, and in the past."})
		return
	}

	if horse.BirthDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Birth date must be in the past."})
		return
	}

	if err := h.horseService.Create(c.Request.Context(), &horse); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create horse: " + err.Error()})
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

	horseIDParam := c.Param("id")
	id, err := strconv.ParseUint(horseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	var horseUpdates models.Horse
	if err := c.ShouldBindJSON(&horseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	existingHorse, err := h.horseService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
		}
		return
	}

	if existingHorse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied: You do not own this horse"})
		return
	}

	horseUpdates.ID = uint(id)
	horseUpdates.UserID = userID
	horseUpdates.CreatedAt = existingHorse.CreatedAt

	if err := h.validate.Struct(&horseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	if !horseUpdates.ValidateGender() {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid gender specified."})
		return
	}
	if !horseUpdates.ValidatePregnancy() {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid pregnancy data: If pregnant, gender must be MARE, conception date must be set, valid, and in the past."})
		return
	}

	if horseUpdates.BirthDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Birth date must be in the past."})
		return
	}

	if err := h.horseService.Update(c.Request.Context(), &horseUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update horse: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, horseUpdates)
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
	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID)) 
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
        }
		return
	}
	// Assumes GetByID checks ownership or it's handled by middleware/context

	var record models.HealthRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Set IDs before validation as they are required
	record.HorseID = uint(horseID)
	record.UserID = userID

	// Validate the struct
	if err := h.validate.Struct(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// Custom validation for Date
	if record.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Record date cannot be in the future."})
		return
	}

	if err := h.healthService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create health record: " + err.Error()})
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
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid record ID in path"})
		return
	}

	// Verify horse ownership first 
	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
        }
        return
    }

	// Fetch existing record to verify it belongs to the horse and exists
	// Use the request context (c.Request.Context()) for service calls
	existingRecord, err := h.healthService.GetRecordByID(c.Request.Context(), uint(recordID))
	if err != nil {
	    if errors.Is(err, gorm.ErrRecordNotFound) {
	        c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Health record not found"})
	    } else {
	        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve health record: " + err.Error()})
	    }
	    return
	}

    if existingRecord.HorseID != uint(horseID) {
        c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Health record does not belong to the specified horse"})
        return
    }

	var recordUpdates models.HealthRecord
	if err := c.ShouldBindJSON(&recordUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	recordUpdates.ID = uint(recordID)
	recordUpdates.HorseID = uint(horseID)
	recordUpdates.UserID = userID 
	recordUpdates.CreatedAt = existingRecord.CreatedAt

	if err := h.validate.Struct(&recordUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	if recordUpdates.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Record date cannot be in the future."})
		return
	}

	if err := h.healthService.UpdateRecord(c.Request.Context(), &recordUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update health record: " + err.Error()})
		return
	}

    updatedRecord, err := h.healthService.GetRecordByID(c.Request.Context(), uint(recordID)) 
    if err != nil {
        fmt.Printf("WARN: Failed to fetch updated health record %d after update: %v\n", recordID, err)
        c.JSON(http.StatusOK, recordUpdates) 
        return
    }
	c.JSON(http.StatusOK, updatedRecord) // Use updatedRecord here
}

// DeleteHealthRecord handles DELETE /horses/:id/health/:recordId
// Needs review for variable declaration fixing (potential := vs = issue)
func (h *Handler) DeleteHealthRecord(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("id"), 10, 64) // Declare err here
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID"})
		return
	}

	// Verify horse ownership first
	// Use existing err variable with =
	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
        }
		return
	}

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 64) // Declare err again for this scope if needed, or use =
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid record ID"})
		return
	}

    // Fetch the record to ensure it belongs to the horse before deleting
    recordToDelete, err := h.healthService.GetRecordByID(c.Request.Context(), uint(recordID)) // Use GetRecordByID
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Health record not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve health record for deletion: " + err.Error()})
        }
        return
    }

    if recordToDelete.HorseID != uint(horseID) {
        c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Health record does not belong to the specified horse"})
        return
    }
    // User ownership verified via horse ownership

	if err = h.healthService.DeleteRecord(c.Request.Context(), uint(recordID)); err != nil { // Use = for err
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to delete health record: " + err.Error()})
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
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var userUpdates models.User
	if err := c.ShouldBindJSON(&userUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Ensure the UserID for the update is the authenticated user's ID
	// and not one from the request body, for security.
	userUpdates.ID = userID 

	// Validate the user struct. 
	// Note: User struct itself is validated, which includes WeatherSettings if it has `validate:"dive"` or similar.
	// If WeatherSettings can be entirely omitted, its fields are validated only if present due to `omitempty` on them.
	if err := h.validate.Struct(&userUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// If there are any model-level custom validations for User (like ValidateGender for Horse), they would go here.
	// For example: if !userUpdates.ValidateSomethingCustom() { ... }

	// Fetch existing user to preserve fields not typically updated by this endpoint (e.g., CreatedAt, IsActive, etc.)
	// The service layer should ideally handle this, by only updating allowed fields.
	// For now, we pass `userUpdates` which has the ID set correctly.
	// The service.UpdateProfile should be designed to only update specific fields from `userUpdates`.
	// (e.g., Email, WeatherSettings, but not IsActive, CreatedAt, HashedPassword etc.)

	if err := h.userService.UpdateProfile(c.Request.Context(), &userUpdates); err != nil {
		// Handle specific errors from the service, e.g., email already taken by another user (if applicable)
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update user profile: " + err.Error()})
		return
	}

	// Fetch the updated profile to return the full, clean data, 
	// as userUpdates might only contain the fields sent in the request.
	// Alternatively, the service.UpdateProfile could return the full updated user.
	updatedProfile, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		// Log this error, but the update itself was successful, so we might still return OK or the partial data.
		// For simplicity now, we return the data from userUpdates if fresh fetch fails, but log the error.
		fmt.Printf("Error fetching updated profile after update (user %s): %v\n", userID, err) // Basic logging
		c.JSON(http.StatusOK, userUpdates) // Return the input data if fresh fetch fails
		return
	}

	c.JSON(http.StatusOK, updatedProfile)
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

	// Verify horse ownership - it's good practice to ensure the user owns the horse
	// for which they are starting pregnancy tracking.
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
		}
		return
	}
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "Access denied to this horse"})
		return
	}

	var startRequest models.PregnancyStart
	if err := c.ShouldBindJSON(&startRequest); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Validate the PregnancyStart DTO
	if err := h.validate.Struct(&startRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// Custom validation for ConceptionDate being in the past
	if startRequest.ConceptionDate.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Conception date must be in the past."})
		return
	}

	if err := h.pregnancyService.StartTracking(c.Request.Context(), uint(horseID), startRequest); err != nil {
		// Handle specific service errors, e.g., horse already pregnant, horse not a mare
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to start pregnancy tracking: " + err.Error()})
		return
	}

	c.Status(http.StatusCreated) // Or return the created/updated pregnancy record
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

	horseIDParam := c.Param("id")
	horseID, err := strconv.ParseUint(horseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	// Verify horse ownership implicitly via horseService.GetByID if it enforces user ID
	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
		}
		return
	}

	var eventInput models.PregnancyEventInputDTO
	if err := c.ShouldBindJSON(&eventInput); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Validate the input DTO using tags
	if err := h.validate.Struct(&eventInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// Custom validation for the event date
	if eventInput.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Event date cannot be in the future."})
		return
	}

	// Call the updated service method with the NEW signature
	createdEvent, err := h.pregnancyService.AddPregnancyEvent(c.Request.Context(), userID, uint(horseID), &eventInput)
	if err != nil {
		// Check for specific errors returned by the service
		if strings.Contains(err.Error(), "no active pregnancy found") { // Basic check, better to use custom error types
			c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to add pregnancy event: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdEvent)
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
	_, err = h.horseService.GetByID(c.Request.Context(), uint(horseID)) 
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
        }
		return
	}

	var record models.BreedingRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Set IDs before validation
	record.HorseID = uint(horseID)
	record.UserID = userID

	// Validate the struct
	if err := h.validate.Struct(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// Custom validation for Date
	if record.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Record date cannot be in the future."})
		return
	}
	// TODO: Add validation for Status enum if defined (this should be handled by `max=50` and specific values if any)

	if err := h.breedingService.CreateRecord(c.Request.Context(), &record); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create breeding record: " + err.Error()})
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
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	recordID, err := strconv.ParseUint(c.Param("recordId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid record ID in path"})
		return
	}

	// Verify horse ownership (user owns the horse)
	horse, err := h.horseService.GetByID(c.Request.Context(), uint(horseID))
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Horse not found"})
        } else {
            c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve horse: " + err.Error()})
        }
        return
    }
	if horse.UserID != userID {
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "User does not own this horse"})
		return
	}

	// Fetch existing breeding record by recordID
	existingRecord, err := h.breedingService.GetRecordByID(c.Request.Context(), uint(recordID))
	if err != nil {
	    if errors.Is(err, gorm.ErrRecordNotFound) {
	        c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Breeding record not found"})
	    } else {
	        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to retrieve breeding record: " + err.Error()})
	    }
	    return
	}

    // Verify the existing record belongs to the specified horse and user
    if existingRecord.HorseID != uint(horseID) {
        c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Breeding record does not belong to the specified horse"})
        return
    }
	if existingRecord.UserID != userID { // Double check, though horse ownership check should cover this
		c.JSON(http.StatusForbidden, types.ErrorResponse{Error: "User does not own this breeding record"})
		return
	}

	var recordUpdates models.BreedingRecord
	if err := c.ShouldBindJSON(&recordUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Set immutable fields and preserve CreatedAt
	recordUpdates.ID = uint(recordID)
	recordUpdates.HorseID = uint(horseID)
	recordUpdates.UserID = userID
	recordUpdates.CreatedAt = existingRecord.CreatedAt // Preserve original creation timestamp

	// Validate the updated struct
	if err := h.validate.Struct(&recordUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	// Custom validation for Date
	if recordUpdates.Date.After(time.Now()) {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Record date cannot be in the future."})
		return
	}
	// TODO: Add validation for Status enum if defined (this should be handled by `max=50` and specific values if any)

	if err := h.breedingService.UpdateRecord(c.Request.Context(), &recordUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update breeding record: " + err.Error()})
		return
	}

    // Fetch and return the full updated record for consistency
    updatedRecord, err := h.breedingService.GetRecordByID(c.Request.Context(), uint(recordID))
    if err != nil {
        // Log the error, but still return the data we have if the primary update succeeded
        fmt.Printf("WARN: Failed to fetch updated breeding record %d after update: %v\n", recordID, err)
        c.JSON(http.StatusOK, recordUpdates) // Fallback to returning the input if fetch fails
        return
    }
	c.JSON(http.StatusOK, updatedRecord)
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

// GetGrowthService returns the growth service
func (h *Handler) GetGrowthService() service.GrowthService {
    return h.growthService
}

// GetExpenseService returns the expense service
func (h *Handler) GetExpenseService() service.ExpenseService {
	return h.expenseService
}

// AddExpense handles POST /expenses
// TODO: Consider if HorseID should be validated to ensure it belongs to the user.
// The service layer's GetHorseExpenses does this, but perhaps a check here too for creating an expense against a horse.
// For now, assuming UserID on the expense is the primary check for direct expenses not linked to a horse,
// or that HorseID linkage implies user ownership (checked if listing by horse).
func (h *Handler) AddExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	expense.UserID = userID // Ensure UserID is set from authenticated user

	if err := h.validate.Struct(&expense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}
	
	// Additional business logic validation if any (e.g. Date must be in the past - already handled by validator tag 'past')

	createdExpense, err := h.expenseService.CreateExpense(c.Request.Context(), &expense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdExpense)
}

// ListExpensesByHorse handles GET /horses/:horseId/expenses
func (h *Handler) ListExpensesByHorse(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	horseIDParam := c.Param("horseId")
	horseID, err := strconv.ParseUint(horseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid horse ID in path"})
		return
	}

	expenses, err := h.expenseService.GetHorseExpenses(c.Request.Context(), userID, uint(horseID))
	if err != nil {
		// Service layer handles logging and specific error types (e.g., permission denied)
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if strings.Contains(err.Error(), "not found"){ // Or a more specific error type if service returns one
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to list expenses for horse: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetExpense handles GET /expenses/:expenseId
func (h *Handler) GetExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId")
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid expense ID in path"})
		return
	}

	expense, err := h.expenseService.GetExpenseByID(c.Request.Context(), userID, uint(expenseID))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") { // More robust check for not found
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Expense not found or access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateExpense handles PUT /expenses/:expenseId
func (h *Handler) UpdateExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId")
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid expense ID in path"})
		return
	}

	var expenseUpdates models.Expense
	if err := c.ShouldBindJSON(&expenseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	// Ensure UserID and ID are set for the service layer, which performs ownership check
	expenseUpdates.UserID = userID
	expenseUpdates.ID = uint(expenseID)


	if err := h.validate.Struct(&expenseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	updatedExpense, err := h.expenseService.UpdateExpense(c.Request.Context(), &expenseUpdates)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Expense not found or access denied for update"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedExpense)
}

// DeleteExpense handles DELETE /expenses/:expenseId
func (h *Handler) DeleteExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId")
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid expense ID in path"})
		return
	}

	err = h.expenseService.DeleteExpense(c.Request.Context(), userID, uint(expenseID))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Expense not found or access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to delete expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

// AddRecurringExpense handles POST /recurring-expenses
func (h *Handler) AddRecurringExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	var recurringExpense models.RecurringExpense
	if err := c.ShouldBindJSON(&recurringExpense); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	recurringExpense.UserID = userID // Ensure UserID is set

	if err := h.validate.Struct(&recurringExpense); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}
	
	// Ensure StartDate is not in the future (handled by 'past' tag if applicable or add custom validation)
	// The model validator includes `validate:"required,past"` for StartDate if that's defined.

	createdRecExpense, err := h.expenseService.CreateRecurringExpense(c.Request.Context(), &recurringExpense)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create recurring expense: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdRecExpense)
}

// ListRecurringExpenses handles GET /recurring-expenses
func (h *Handler) ListRecurringExpenses(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenses, err := h.expenseService.GetUserRecurringExpenses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to list recurring expenses: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

// GetRecurringExpense handles GET /recurring-expenses/:expenseId
func (h *Handler) GetRecurringExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId") // Parameter name should be consistent
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid recurring expense ID in path"})
		return
	}

	expense, err := h.expenseService.GetRecurringExpenseByID(c.Request.Context(), userID, uint(expenseID))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Recurring expense not found or access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get recurring expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, expense)
}

// UpdateRecurringExpense handles PUT /recurring-expenses/:expenseId
func (h *Handler) UpdateRecurringExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId") // Parameter name should be consistent
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid recurring expense ID in path"})
		return
	}

	var expenseUpdates models.RecurringExpense
	if err := c.ShouldBindJSON(&expenseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request payload: " + err.Error()})
		return
	}

	expenseUpdates.UserID = userID
	expenseUpdates.ID = uint(expenseID)

	if err := h.validate.Struct(&expenseUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationErrors(err)})
		return
	}

	updatedExpense, err := h.expenseService.UpdateRecurringExpense(c.Request.Context(), &expenseUpdates)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Recurring expense not found or access denied for update"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update recurring expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedExpense)
}

// DeleteRecurringExpense handles DELETE /recurring-expenses/:expenseId
func (h *Handler) DeleteRecurringExpense(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	expenseIDParam := c.Param("expenseId") // Parameter name consistency
	expenseID, err := strconv.ParseUint(expenseIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid recurring expense ID in path"})
		return
	}

	err = h.expenseService.DeleteRecurringExpense(c.Request.Context(), userID, uint(expenseID))
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			c.JSON(http.StatusForbidden, types.ErrorResponse{Error: err.Error()})
		} else if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Recurring expense not found or access denied"})
		} else {
			c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to delete recurring expense: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recurring expense deleted successfully"})
}

// GetUserTotalExpensesHandler handles GET /expenses/total
func (h *Handler) GetUserTotalExpensesHandler(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
		return
	}

	total, err := h.expenseService.GetUserTotalExpenses(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get total expenses: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_expenses": total})
}

// GetExpensesByTypeHandler handles GET /expenses/type/:type
func (h *Handler) GetExpensesByTypeHandler(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "User not authenticated"})
        return
    }

    expenseType := c.Param("type")
    if expenseType == "" {
        c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Expense type parameter is required"})
        return
    }

    expenses, err := h.expenseService.GetExpensesByType(c.Request.Context(), userID, expenseType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to get expenses by type: " + err.Error()})
        return
    }

    c.JSON(http.StatusOK, expenses)
}

// Helper function to format validation errors
func formatValidationErrors(err error) map[string]string {
	errorsMap := make(map[string]string)
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			errorsMap[fieldName] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldName, fieldError.Tag())
		}
	} else {
		errorsMap["error"] = "Validation error: " + err.Error()
	}
	return errorsMap
}
