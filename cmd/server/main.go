package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polyfant/horse_tracking/internal/database"
	"github.com/polyfant/horse_tracking/internal/api"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create store
	store := database.NewSQLiteStore(db)

	// Setup router
	r := gin.Default()

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Setup routes
	api.SetupRoutes(r, store)
	api.SetupPregnancyRoutes(r) // Add pregnancy tracking frontend routes

	// Start server
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
