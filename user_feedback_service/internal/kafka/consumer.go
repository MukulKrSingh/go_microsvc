package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/user_feedback_service/internal/models"
	"gorm.io/gorm"
)

const (
	OrderTopic = "orders"
	GroupID    = "feedback-service-group"
)

var Reader *kafka.Reader
var DB *gorm.DB

// InitKafkaConsumer initializes the Kafka consumer
func InitKafkaConsumer(db *gorm.DB) {
	DB = db
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaBrokers},
		Topic:          OrderTopic,
		GroupID:        GroupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	log.Println("Kafka consumer initialized successfully")

	// Start consuming messages in a goroutine
	go consumeMessages()
}

// CloseKafkaConsumer closes the Kafka consumer connection
func CloseKafkaConsumer() {
	if Reader != nil {
		if err := Reader.Close(); err != nil {
			log.Printf("Error closing Kafka reader: %v", err)
		}
	}
}

// consumeMessages consumes messages from Kafka
func consumeMessages() {
	ctx := context.Background()
	for {
		message, err := Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Process the message
		var orderEvent models.OrderEvent
		if err := json.Unmarshal(message.Value, &orderEvent); err != nil {
			log.Printf("Error unmarshaling order event: %v", err)
			continue
		}

		// Process the order event, for example store in database for future feedback collection
		log.Printf("Received order event: OrderID=%d, Status=%s", orderEvent.OrderID, orderEvent.Status)

		// Here you could store the order information or perform other processing
		if orderEvent.Status == "completed" {
			// Check if user exists, create if not
			var user models.User
			if result := DB.FirstOrCreate(&user, models.User{ID: uint(orderEvent.UserID)}); result.Error != nil {
				log.Printf("Error creating user: %v", result.Error)
				continue
			}
		}

		// Commit the message offset
		if err := Reader.CommitMessages(ctx, message); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}
