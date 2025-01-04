package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

func SendErrorResponse(c *gin.Context, err error, statusCode ...int) {
	status := http.StatusInternalServerError
	if len(statusCode) > 0 {
		status = statusCode[0]
	}

	c.JSON(status, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}
