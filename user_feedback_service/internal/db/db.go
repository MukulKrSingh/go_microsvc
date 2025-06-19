package db

import (
	"fmt"
	"log"
	"os"

	"github.com/user_feedback_service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection using GORM
func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Make sure connection is alive
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
}

// MigrateSchema creates or updates the database schema
func MigrateSchema() {
	log.Println("Migrating database schema...")

	// Auto migrate the schema
	err := DB.AutoMigrate(&models.User{}, &models.Feedback{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	log.Println("Database schema migration completed")
}

// SeedData adds some initial data to the database
func SeedData() {
	// Check if users already exist
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	// Only seed if no users exist
	if userCount == 0 {
		// Seed a default user
		defaultUser := models.User{
			Username: "feedbackuser",
			Email:    "feedback@example.com",
		}

		result := DB.Create(&defaultUser)
		if result.Error != nil {
			log.Fatalf("Failed to insert default user: %v", result.Error)
		}

		log.Println("Successfully seeded default user")
	}
}
