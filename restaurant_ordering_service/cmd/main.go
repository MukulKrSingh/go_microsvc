package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/restaurant_ordering_service/internal/api"
	"github.com/restaurant_ordering_service/internal/db"
	"github.com/restaurant_ordering_service/internal/kafka"
	"github.com/restaurant_ordering_service/internal/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Set default values for environment variables if not set
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8080")
	}
	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "your-secret-key")
	}

	// Initialize database connection
	db.InitDB()

	// Create tables and seed data
	db.CreateTables()
	db.SeedData()

	// Initialize Kafka
	kafka.InitKafka()
	defer kafka.CloseKafka()

	// Set up Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/auth", api.AuthHandler)
	router.GET("/food-items", api.GetFoodItemsHandler)

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/profile", api.GetUserProfileHandler)
		authorized.POST("/orders", api.PlaceOrderHandler)
		authorized.POST("/transactions", api.HandleTransactionHandler)
	}

	// Start the server
	port := os.Getenv("PORT")
	log.Printf("Restaurant Ordering Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
