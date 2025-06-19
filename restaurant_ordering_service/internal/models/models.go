package models

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

// FoodItem represents a food item in the menu
type FoodItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`    // Always 10 INR as per requirements
	Quantity int     `json:"quantity"` // Starting with 1000 as per requirements
}

// Order represents a user's order
type Order struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"` // pending, completed, cancelled
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID         int `json:"id"`
	OrderID    int `json:"order_id"`
	FoodItemID int `json:"food_item_id"`
	Quantity   int `json:"quantity"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	Token string `json:"token"`
}

// OrderRequest represents a request to place an order
type OrderRequest struct {
	Items []OrderItemRequest `json:"items"`
}

// OrderItemRequest represents an item in an order request
type OrderItemRequest struct {
	FoodItemID int `json:"food_item_id"`
	Quantity   int `json:"quantity"`
}

// TransactionRequest represents a request to process a transaction
type TransactionRequest struct {
	OrderID int `json:"order_id"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OrderEvent represents an order event that will be sent to Kafka
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
