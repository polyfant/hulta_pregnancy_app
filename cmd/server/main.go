// @title           Hulta Pregnancy App API
// @version         1.0
// @description     API for tracking horse pregnancies and related information
// @host            localhost:8080
// @BasePath        /

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/api"
	"github.com/polyfant/horse_tracking/internal/database"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create handler and setup routes
	store := database.NewSQLiteStore(db)
	handler := api.NewHandler(store)
	router := api.SetupRouter(handler)

	// Start server
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
