package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}

	// Try to ping the database
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Printf("Could not establish connection to database after %d attempts", maxRetries)
		return
	}

	log.Println("Successfully connected to database")
}

// CreateTables creates the necessary tables in the database
func CreateTables() {
	// Create Users table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			address TEXT
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create FoodItems table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS food_items (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			price NUMERIC(10,2) NOT NULL,
			quantity INT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create food_items table: %v", err)
	}

	// Create Orders table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			total_price NUMERIC(10,2) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}

	// Create OrderItems table
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INT REFERENCES orders(id),
			food_item_id INT REFERENCES food_items(id),
			quantity INT NOT NULL
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create order_items table: %v", err)
	}

	log.Println("Successfully created tables")
}

// SeedData adds some initial data to the database
func SeedData() {
	// Check if food items already exist
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM food_items").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check food items count: %v", err)
	}

	// Only seed if no food items exist
	if count == 0 {
		// Seed some food items (all with price 10 INR and quantity 1000)
		foodItems := []struct {
			name string
		}{
			{"Butter Chicken"},
			{"Paneer Tikka"},
			{"Biryani"},
			{"Masala Dosa"},
			{"Chole Bhature"},
			{"Pav Bhaji"},
			{"Gulab Jamun"},
			{"Samosa"},
			{"Naan"},
			{"Tandoori Roti"},
		}

		for _, item := range foodItems {
			_, err := DB.Exec(
				"INSERT INTO food_items (name, price, quantity) VALUES ($1, $2, $3)",
				item.name, 10.0, 1000,
			)
			if err != nil {
				log.Fatalf("Failed to insert food item: %v", err)
			}
		}

		log.Println("Successfully seeded food items")
	}

	// Check if users already exist
	err = DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to check users count: %v", err)
	}

	// Only seed if no users exist
	if count == 0 {
		// Seed a default user
		_, err := DB.Exec(
			"INSERT INTO users (username, password, email, address) VALUES ($1, $2, $3, $4)",
			"testuser", "password123", "test@example.com", "123 Test Street, Test City",
		)
		if err != nil {
			log.Fatalf("Failed to insert default user: %v", err)
		}

		log.Println("Successfully seeded default user")
	}
}
