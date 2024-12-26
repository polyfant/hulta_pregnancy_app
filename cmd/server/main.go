package main

import (
	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/api"
	"github.com/polyfant/horse_tracking/internal/database"
	"log"
)

func main() {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Initialize SQLite store
	db, err := database.NewSQLiteStore("horse_tracking.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize handler with dependencies
	handler := api.NewHandler(db)

	// Setup router with the handler and store
	router := api.SetupRouter(handler, db)

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
