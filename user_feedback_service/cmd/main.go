package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/user_feedback_service/internal/api"
	"github.com/user_feedback_service/internal/db"
	"github.com/user_feedback_service/internal/kafka"
	"github.com/user_feedback_service/internal/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set default values for environment variables if not set
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8081")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-secret-key")
	}

	// Initialize database connection
	db.InitDB()

	// Create tables and seed data
	db.MigrateSchema()
	db.SeedData()

	// Initialize Kafka (pass the database connection for consumer use)
	kafka.InitKafkaConsumer(db.DB)
	defer kafka.CloseKafkaConsumer()

	// Set up Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/auth", api.AuthHandler)

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Feedback endpoints
		authorized.GET("/feedback", api.GetUserFeedbackHandler)
		authorized.POST("/feedback", api.CreateFeedbackHandler)
		authorized.PUT("/feedback/:id", api.UpdateFeedbackHandler)
		authorized.DELETE("/feedback/:id", api.DeleteFeedbackHandler)

		// Get feedback stats
		authorized.GET("/feedback/stats", api.GetFeedbackStatsHandler)
	}

	// Start the server
	port := os.Getenv("PORT")
	log.Printf("Feedback Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
