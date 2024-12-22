package api

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/service/pregnancy"
)

type PregnancyResponse struct {
	Horse         HorseInfo                    `json:"horse"`
	Progress      float64                      `json:"progress"`
	CurrentDay    int                          `json:"currentDay"`
	RemainingDays int                         `json:"remainingDays"`
	Stage         string                      `json:"stage"`
	Schedule      pregnancy.MonitoringSchedule `json:"schedule"`
}

type HorseInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func SetupPregnancyRoutes(router *gin.Engine) {
	// Serve static files
	router.Static("/static", filepath.Join("frontend", "static"))
	
	// Serve templates
	router.LoadHTMLGlob(filepath.Join("frontend", "templates", "*"))

	// Routes
	router.GET("/pregnancy/:id", handlePregnancyPage)
	router.GET("/api/pregnancy/status/:id", handlePregnancyStatus)
}

func handlePregnancyPage(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid horse ID",
		})
		return
	}

	// In a real app, get this from the database
	horse := HorseInfo{
		ID:   int(horseID),
		Name: "Bella",
	}

	calculator := pregnancy.NewPregnancyCalculator(getConceptionDate())
	
	response := PregnancyResponse{
		Horse:         horse,
		Progress:      calculator.GetProgressPercentage(),
		CurrentDay:    calculator.GetCurrentDay(),
		RemainingDays: calculator.GetRemainingDays(),
		Stage:         string(calculator.GetStage()),
		Schedule:      calculator.GetMonitoringSchedule(),
	}

	c.HTML(http.StatusOK, "pregnancy.html", response)
}

func handlePregnancyStatus(c *gin.Context) {
	horseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid horse ID",
		})
		return
	}

	horse := HorseInfo{
		ID:   int(horseID),
		Name: "Bella",
	}

	calculator := pregnancy.NewPregnancyCalculator(getConceptionDate())
	
	response := PregnancyResponse{
		Horse:         horse,
		Progress:      calculator.GetProgressPercentage(),
		CurrentDay:    calculator.GetCurrentDay(),
		RemainingDays: calculator.GetRemainingDays(),
		Stage:         string(calculator.GetStage()),
		Schedule:      calculator.GetMonitoringSchedule(),
	}

	c.JSON(http.StatusOK, response)
}

// Temporary function - replace with database call
func getConceptionDate() time.Time {
	// For testing, set conception date to 150 days ago
	return time.Now().AddDate(0, 0, -150)
}
