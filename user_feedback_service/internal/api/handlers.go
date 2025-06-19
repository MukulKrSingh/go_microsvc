package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/user_feedback_service/internal/db"
	"github.com/user_feedback_service/internal/models"
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
	result := db.DB.Where("username = ?", loginRequest.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Invalid username or password",
		})
		return
	}

	// In a real-world scenario, we would check the password hash here
	// For this example we're omitting actual password checking

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

// GetUserFeedbackHandler returns all feedback from the authenticated user
func GetUserFeedbackHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var feedbacks []models.Feedback
	result := db.DB.Where("user_id = ?", userID).Find(&feedbacks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error fetching feedback",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Feedback retrieved successfully",
		Data:    feedbacks,
	})
}

// CreateFeedbackHandler creates new feedback for an order
func CreateFeedbackHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var feedbackRequest models.FeedbackRequest
	if err := c.ShouldBindJSON(&feedbackRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Check if the user has already provided feedback for this order
	var existingFeedback models.Feedback
	result := db.DB.Where("order_id = ? AND user_id = ?", feedbackRequest.OrderID, userID).First(&existingFeedback)
	if result.Error == nil {
		c.JSON(http.StatusConflict, models.APIResponse{
			Success: false,
			Message: "You have already provided feedback for this order",
		})
		return
	}

	// Create the feedback
	feedback := models.Feedback{
		OrderID: feedbackRequest.OrderID,
		UserID:  userID,
		Rating:  feedbackRequest.Rating,
		Comment: feedbackRequest.Comment,
	}

	result = db.DB.Create(&feedback)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error creating feedback: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Feedback submitted successfully",
		Data:    feedback,
	})
}

// UpdateFeedbackHandler updates existing feedback
func UpdateFeedbackHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	feedbackID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid feedback ID",
		})
		return
	}

	var updateRequest models.FeedbackUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Find the feedback to ensure it belongs to the user
	var feedback models.Feedback
	result := db.DB.Where("id = ? AND user_id = ?", feedbackID, userID).First(&feedback)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Feedback not found or not authorized to update",
		})
		return
	}

	// Update the feedback
	updates := map[string]interface{}{}
	if updateRequest.Rating > 0 {
		updates["rating"] = updateRequest.Rating
	}
	if updateRequest.Comment != "" {
		updates["comment"] = updateRequest.Comment
	}

	result = db.DB.Model(&feedback).Updates(updates)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error updating feedback: " + result.Error.Error(),
		})
		return
	}

	// Fetch the updated feedback
	db.DB.First(&feedback, feedbackID)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Feedback updated successfully",
		Data:    feedback,
	})
}

// DeleteFeedbackHandler deletes feedback
func DeleteFeedbackHandler(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	feedbackID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid feedback ID",
		})
		return
	}

	// Find the feedback to ensure it belongs to the user
	var feedback models.Feedback
	result := db.DB.Where("id = ? AND user_id = ?", feedbackID, userID).First(&feedback)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Feedback not found or not authorized to delete",
		})
		return
	}

	// Delete the feedback
	result = db.DB.Delete(&feedback)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Error deleting feedback: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Feedback deleted successfully",
	})
}

// GetFeedbackStatsHandler returns feedback statistics
func GetFeedbackStatsHandler(c *gin.Context) {
	var stats struct {
		TotalFeedback int64   `json:"total_feedback"`
		AverageRating float64 `json:"average_rating"`
		RatingCounts  []struct {
			Rating int   `json:"rating"`
			Count  int64 `json:"count"`
		} `json:"rating_counts"`
	}

	// Get total feedback count
	db.DB.Model(&models.Feedback{}).Count(&stats.TotalFeedback)

	// Get average rating
	db.DB.Model(&models.Feedback{}).Select("COALESCE(AVG(rating), 0) as average_rating").Row().Scan(&stats.AverageRating)

	// Get count by rating
	stats.RatingCounts = make([]struct {
		Rating int   `json:"rating"`
		Count  int64 `json:"count"`
	}, 5)

	for i := 1; i <= 5; i++ {
		var count int64
		db.DB.Model(&models.Feedback{}).Where("rating = ?", i).Count(&count)
		stats.RatingCounts[i-1] = struct {
			Rating int   `json:"rating"`
			Count  int64 `json:"count"`
		}{
			Rating: i,
			Count:  count,
		}
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Feedback statistics retrieved successfully",
		Data:    stats,
	})
}
