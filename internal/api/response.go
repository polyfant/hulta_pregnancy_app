package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendSuccess sends a successful response
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(c *gin.Context, err error, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}
