package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/hulta_pregnancy_app/internal/models"
	"github.com/polyfant/hulta_pregnancy_app/internal/types"
)

type WeatherService interface {
	GetCurrentWeather(ctx context.Context, latitude, longitude float64) (*models.WeatherData, error)
	GetLatestImpact(ctx context.Context, horseID uint) (*models.WeatherImpact, error)
	GetWeatherHistory(ctx context.Context, horseID uint, start, end string) ([]*models.WeatherData, error)
}

type NotificationService interface {
	SendWeatherAlert(ctx context.Context, userID string, alert string) error
}

type WeatherHandler struct {
	weatherService WeatherService
	notifyService NotificationService
}

func NewWeatherHandler(ws WeatherService, ns NotificationService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: ws,
		notifyService: ns,
	}
}

// GetCurrentWeather handles GET /weather/current
func (h *WeatherHandler) GetCurrentWeather(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "unauthorized"})
		return
	}

	latitude, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid latitude"})
		return
	}

	longitude, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid longitude"})
		return
	}

	weather, err := h.weatherService.GetCurrentWeather(c.Request.Context(), latitude, longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "failed to get weather data"})
		return
	}

	c.JSON(http.StatusOK, weather)
}

// GetWeatherImpact handles GET /weather/impact/:horseId
func (h *WeatherHandler) GetWeatherImpact(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "unauthorized"})
		return
	}

	horseID, err := strconv.ParseUint(c.Param("horseId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "invalid horse ID"})
		return
	}

	impact, err := h.weatherService.GetLatestImpact(c.Request.Context(), uint(horseID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "failed to get weather impact"})
		return
	}

	c.JSON(http.StatusOK, impact)
}
