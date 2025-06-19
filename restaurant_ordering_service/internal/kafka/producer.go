package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/restaurant_ordering_service/internal/models"
	"github.com/segmentio/kafka-go"
)

const (
	OrderTopic = "orders"
)

var Writer *kafka.Writer

// InitKafka initializes the Kafka producer
func InitKafka() {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	Writer = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokers),
		Topic:    OrderTopic,
		Balancer: &kafka.LeastBytes{},
	}

	log.Println("Kafka producer initialized successfully")
}

// CloseKafka closes the Kafka producer connection
func CloseKafka() {
	if err := Writer.Close(); err != nil {
		log.Printf("Error closing Kafka writer: %v", err)
	}
}

// PublishOrderEvent publishes an order event to the Kafka topic
func PublishOrderEvent(order models.Order, foodItems map[int]string) error {
	// Prepare order items for the event
	var items []models.Item
	for _, orderItem := range order.OrderItems {
		items = append(items, models.Item{
			FoodItemID: orderItem.FoodItemID,
			Name:       foodItems[orderItem.FoodItemID],
			Quantity:   orderItem.Quantity,
		})
	}

	// Create event
	event := models.OrderEvent{
		OrderID:    order.ID,
		UserID:     order.UserID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		Items:      items,
		Timestamp:  time.Now().Unix(),
	}

	// Serialize to JSON
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Publish message
	err = Writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(strconv.Itoa(event.OrderID)),
			Value: value,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Order event published to Kafka: OrderID=%d, Status=%s", event.OrderID, event.Status)
	return nil
}
