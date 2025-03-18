package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a connection to a RabbitMQ cluster.
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ connection.
func NewRabbitMQ(urls []string) (*RabbitMQ, error) {
	if len(urls) == 0 {
		return nil, errors.New("at least one RabbitMQ URL is required")
	}

	for _, url := range urls {
		conn, err := amqp.Dial(url)
		if err == nil {
			ch, err := conn.Channel()
			if err != nil {
				return nil, fmt.Errorf("failed to open channel: %v", err)
			}

			log.Printf("Connected to RabbitMQ at %s", url)
			return &RabbitMQ{conn: conn, channel: ch}, nil
		}
		log.Printf("Failed to connect to RabbitMQ at %s: %v", url, err)
	}

	return nil, errors.New("failed to connect to any RabbitMQ node")
}

// DeclarePriorityQueue declares a priority queue with the given name.
func (r *RabbitMQ) DeclarePriorityQueue(queueName string) error {
	_, err := r.channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-max-priority": 4, // Set maximum priority to 4
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare priority queue: %v", err)
	}

	log.Printf("Declared priority queue: %s", queueName)
	return nil
}

// GetMessage retrieves a single message from the specified queue.
func (r *RabbitMQ) GetMessage(queueName string) (*amqp.Delivery, error) {
	msg, ok, err := r.channel.Get(
		queueName, // queue
		false,     // auto-ack (set to false to manually acknowledge messages)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message from queue %s: %v", queueName, err)
	}

	if !ok {
		return nil, nil // No message available
	}

	return &msg, nil
}

// Close closes the RabbitMQ connection.
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("RabbitMQ connection closed")
}

// MessageHandler processes a message.
func MessageHandler(body []byte) error {
	var message map[string]interface{}
	if err := json.Unmarshal(body, &message); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}

	log.Printf("Processing message: %v", message)
	// Simulate message processing
	time.Sleep(1 * time.Second)
	log.Printf("Finished processing message: %v", message)
	return nil
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Access environment variables
	rabbitMQURLs := strings.Split(os.Getenv("RABBITMQ_URLS"), ",")
	if len(rabbitMQURLs) == 0 {
		log.Fatal("RABBITMQ_URLS environment variable is not set")
	}

	// Initialize RabbitMQ
	rmq, err := NewRabbitMQ(rabbitMQURLs)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Declare priority queues
	queues := []string{"otp", "transactional", "promotional", "general"} // Ordered by priority
	for _, queue := range queues {
		if err := rmq.DeclarePriorityQueue(queue); err != nil {
			log.Fatalf("Failed to declare queue %s: %v", queue, err)
		}
	}

	// Process messages priority-wise
	for {
		processed := false

		// Check queues in priority order
		for _, queue := range queues {
			msg, err := rmq.GetMessage(queue)
			if err != nil {
				log.Printf("Failed to get message from queue %s: %v", queue, err)
				continue
			}

			if msg != nil {
				log.Printf("Received message from queue %s: %s", queue, string(msg.Body))
				if err := MessageHandler(msg.Body); err != nil {
					log.Printf("Failed to process message: %v", err)
					msg.Nack(false, true) // Requeue the message on failure
				} else {
					msg.Ack(false) // Acknowledge the message on success
				}
				processed = true
				break // Process one message at a time
			}
		}

		// If no messages were processed, wait for a short time before checking again
		if !processed {
			time.Sleep(1 * time.Second)
		}
	}
}