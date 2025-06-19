package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/restaurant_ordering_service/internal/db"
	"github.com/restaurant_ordering_service/internal/kafka"
	"github.com/restaurant_ordering_service/internal/models"
)

// AuthHandler handles user authentication
func AuthHandler(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Find the user in the database
	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, password, email, address FROM users WHERE username = $1",
		loginRequest.Username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Address)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	// Check if the password is correct
	if user.Password != loginRequest.Password {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not generate token",
		})
		return
	}

	// Return the token
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Authentication successful",
		Data: models.LoginResponse{
			Token: tokenString,
		},
	})
}

// GetFoodItemsHandler returns a list of food items
func GetFoodItemsHandler(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, name, price, quantity FROM food_items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error retrieving food items",
		})
		return
	}
	defer rows.Close()

	var foodItems []models.FoodItem
	for rows.Next() {
		var item models.FoodItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Error scanning food items",
			})
			return
		}
		foodItems = append(foodItems, item)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Food items retrieved successfully",
		Data:    foodItems,
	})
}

// GetUserProfileHandler returns the profile of the authenticated user
func GetUserProfileHandler(c *gin.Context) {
	// Get the user ID from the token
	userID := c.MustGet("user_id").(int)

	// Get the user from the database
	var user models.User
	err := db.DB.QueryRow(
		"SELECT id, username, email, address FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Address)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    user,
	})
}

// PlaceOrderHandler handles placing an order
func PlaceOrderHandler(c *gin.Context) {
	var orderRequest models.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Get the user ID from the token
	userID := c.MustGet("user_id").(int)

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not start transaction",
		})
		return
	}

	// Calculate total price and check if items exist
	var totalPrice float64
	foodItems := make(map[int]string) // Map to store food item names for event publishing

	for _, item := range orderRequest.Items {
		var price float64
		var name string
		err := tx.QueryRow(
			"SELECT price, name FROM food_items WHERE id = $1",
			item.FoodItemID,
		).Scan(&price, &name)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, models.APIResponse{
					Success: false,
					Message: "Food item not found: " + strconv.Itoa(item.FoodItemID),
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Database error",
			})
			return
		}

		totalPrice += price * float64(item.Quantity)
		foodItems[item.FoodItemID] = name
	}

	// Create the order
	var orderID int
	err = tx.QueryRow(
		"INSERT INTO orders (user_id, total_price, status) VALUES ($1, $2, $3) RETURNING id",
		userID, totalPrice, "pending",
	).Scan(&orderID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not create order",
		})
		return
	}

	// Create order items
	for _, item := range orderRequest.Items {
		_, err := tx.Exec(
			"INSERT INTO order_items (order_id, food_item_id, quantity) VALUES ($1, $2, $3)",
			orderID, item.FoodItemID, item.Quantity,
		)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Could not create order items",
			})
			return
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not commit transaction",
		})
		return
	}

	// Construct the order object for the response
	order := models.Order{
		ID:         orderID,
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	}

	// Publish order event to Kafka
	orderItems := make([]models.OrderItem, 0, len(orderRequest.Items))
	for _, item := range orderRequest.Items {
		orderItems = append(orderItems, models.OrderItem{
			OrderID:    orderID,
			FoodItemID: item.FoodItemID,
			Quantity:   item.Quantity,
		})
	}
	order.OrderItems = orderItems

	go func() {
		if err := kafka.PublishOrderEvent(order, foodItems); err != nil {
			log.Printf("Failed to publish order event: %v", err)
		}
	}()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order placed successfully",
		Data: gin.H{
			"order_id":    orderID,
			"total_price": totalPrice,
		},
	})
}

// HandleTransactionHandler handles a transaction for an order
func HandleTransactionHandler(c *gin.Context) {
	var transactionRequest models.TransactionRequest
	if err := c.ShouldBindJSON(&transactionRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not start transaction",
		})
		return
	}

	// Get order details to verify ownership
	var userID int
	var status string
	var orderItems []struct {
		FoodItemID int
		Quantity   int
		Name       string
	}

	err = tx.QueryRow(
		"SELECT user_id, status FROM orders WHERE id = $1",
		transactionRequest.OrderID,
	).Scan(&userID, &status)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Message: "Order not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}

	// Ensure the order belongs to the authenticated user
	authenticatedUserID := c.MustGet("user_id").(int)
	if userID != authenticatedUserID {
		tx.Rollback()
		c.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "You do not have permission to process this order",
		})
		return
	}

	// Ensure the order is in 'pending' status
	if status != "pending" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Order is not in pending status",
		})
		return
	}

	// Get order items
	rows, err := tx.Query(
		"SELECT oi.food_item_id, oi.quantity, fi.name FROM order_items oi JOIN food_items fi ON oi.food_item_id = fi.id WHERE oi.order_id = $1",
		transactionRequest.OrderID,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not retrieve order items",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item struct {
			FoodItemID int
			Quantity   int
			Name       string
		}
		if err := rows.Scan(&item.FoodItemID, &item.Quantity, &item.Name); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Error scanning order items",
			})
			return
		}
		orderItems = append(orderItems, item)
	}

	// Update food item quantities
	foodItems := make(map[int]string) // For Kafka event
	for _, item := range orderItems {
		result, err := tx.Exec(
			"UPDATE food_items SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $1",
			item.Quantity, item.FoodItemID,
		)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Error updating food item quantities",
			})
			return
		}

		// Check if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Error checking affected rows",
			})
			return
		}
		if rowsAffected == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Not enough quantity for food item: " + item.Name,
			})
			return
		}

		foodItems[item.FoodItemID] = item.Name
	}

	// Update order status
	_, err = tx.Exec(
		"UPDATE orders SET status = $1 WHERE id = $2",
		"completed", transactionRequest.OrderID,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error updating order status",
		})
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Could not commit transaction",
		})
		return
	}

	// Get the full order details for Kafka
	var order models.Order
	order.ID = transactionRequest.OrderID
	order.UserID = userID
	order.Status = "completed"

	// Get total price
	err = db.DB.QueryRow(
		"SELECT total_price FROM orders WHERE id = $1",
		transactionRequest.OrderID,
	).Scan(&order.TotalPrice)
	if err != nil {
		log.Printf("Error getting order total price: %v", err)
	}

	// Convert order items for Kafka message
	var orderItemsForEvent []models.OrderItem
	for _, item := range orderItems {
		orderItemsForEvent = append(orderItemsForEvent, models.OrderItem{
			OrderID:    transactionRequest.OrderID,
			FoodItemID: item.FoodItemID,
			Quantity:   item.Quantity,
		})
	}
	order.OrderItems = orderItemsForEvent

	// Publish completed order event to Kafka asynchronously
	go func() {
		if err := kafka.PublishOrderEvent(order, foodItems); err != nil {
			log.Printf("Failed to publish order completion event: %v", err)
		}
	}()

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Transaction completed successfully",
	})
}
