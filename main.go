package main

import (
	"database_app/handlers"
	"database_app/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Redis client
	utils.InitializeRedis("redis", "6379")

	// Create a Gin router
	router := gin.Default()

	// Load HTML templates
	router.LoadHTMLGlob("templates/*")

	// Serve frontend
	router.GET("/", handlers.ServeHTML)

	// API routes
	api := router.Group("/api")
	{
		api.POST("/connect", handlers.ConnectHandler) // Database connection handler
		api.GET("/jobs", handlers.GetJobs)            // Fetch jobs
	}

	// Start background job checker
	go handlers.StartJobChecker()

	// Run the server
	log.Println("Starting server on :8080")
	router.Run(":8080")
}
