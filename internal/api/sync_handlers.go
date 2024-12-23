package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/models"
)

// @Summary Sync user data
// @Description Synchronize local data with server
// @Tags sync
// @Accept json
// @Produce json
// @Param data body models.SyncData true "Sync data"
// @Success 200 {object} models.SyncData
// @Router /sync [post]
func (h *Handler) SyncData(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var syncData models.SyncData
	if err := c.ShouldBindJSON(&syncData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user owns this data
	if syncData.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Data ownership mismatch"})
		return
	}

	// Begin transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update horses
	for _, horse := range syncData.Horses {
		if err := h.db.AddHorse(&horse); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync horses"})
			return
		}
	}

	// Update health records
	for _, record := range syncData.Health {
		if err := h.db.AddHealthRecord(&record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync health records"})
			return
		}
	}

	// Update pregnancy events
	for _, event := range syncData.Events {
		if err := h.db.AddPregnancyEvent(&event); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sync pregnancy events"})
			return
		}
	}

	// Update last sync time
	if err := h.db.UpdateUserLastSync(userID, time.Now()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sync time"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, syncData)
}

// SyncStatus represents the synchronization status
type SyncStatus struct {
	LastSync    time.Time `json:"lastSync"`
	IsUpToDate  bool     `json:"isUpToDate"`
	PendingSync int      `json:"pendingSync"`
}

// @Summary Get sync status
// @Description Get the last sync time and status
// @Tags sync
// @Produce json
// @Success 200 {object} SyncStatus
// @Router /sync/status [get]
func (h *Handler) GetSyncStatus(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	lastSync, err := h.db.GetLastSyncTime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sync status"})
		return
	}

	pendingChanges, err := h.db.GetPendingSyncCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending changes"})
		return
	}

	status := SyncStatus{
		LastSync:    lastSync,
		IsUpToDate:  pendingChanges == 0,
		PendingSync: pendingChanges,
	}

	c.JSON(http.StatusOK, status)
}

// @Summary Restore user data
// @Description Restore all user data from server
// @Tags sync
// @Produce json
// @Success 200 {object} models.SyncData
// @Router /sync/restore [get]
func (h *Handler) RestoreData(c *gin.Context) {
	userID := getUserIDFromContext(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get all user data
	horses, err := h.db.GetUserHorses(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get horses"})
		return
	}

	health, err := h.db.GetUserHealthRecords(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get health records"})
		return
	}

	events, err := h.db.GetUserPregnancyEvents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pregnancy events"})
		return
	}

	syncData := models.SyncData{
		UserID:    userID,
		Timestamp: time.Now(),
		Horses:    horses,
		Health:    health,
		Events:    events,
	}

	c.JSON(http.StatusOK, syncData)
}
