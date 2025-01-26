package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/api/types"
)



// SendSuccess sends a successful response
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(c *gin.Context, err error, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, types.APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}
