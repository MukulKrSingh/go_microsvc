package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Feedback represents user feedback for orders
type Feedback struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Rating    uint8          `json:"rating" gorm:"type:smallint;not null;check:rating <= 5"` // Rating from 1 to 5
	Comment   string         `json:"comment" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
}

// Order represents an order from the restaurant service (for reference only)
type Order struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

// FeedbackRequest represents a request to add feedback for an order
type FeedbackRequest struct {
	OrderID uint   `json:"order_id" binding:"required"`
	Rating  uint8  `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

// FeedbackUpdateRequest represents a request to update existing feedback
type FeedbackUpdateRequest struct {
	Rating  uint8  `json:"rating" binding:"min=1,max=5"`
	Comment string `json:"comment"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	Token string `json:"token"`
}

// OrderEvent represents an order event received from Kafka
type OrderEvent struct {
	OrderID    int     `json:"order_id"`
	UserID     int     `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	Items      []Item  `json:"items"`
	Timestamp  int64   `json:"timestamp"`
}

// Item represents an item in an order event
type Item struct {
	FoodItemID int    `json:"food_item_id"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
}
