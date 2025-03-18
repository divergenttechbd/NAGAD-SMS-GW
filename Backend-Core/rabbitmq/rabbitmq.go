package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a connection to a RabbitMQ cluster with Quorum Queues.
type RabbitMQ struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	urls          []string // List of RabbitMQ node URLs
	currentURL    string
	mu            sync.Mutex
	reconnectChan chan bool
	closed        bool
}

// NewRabbitMQ creates a new RabbitMQ connection.
func NewRabbitMQ(urls []string) (*RabbitMQ, error) {
	if len(urls) == 0 {
		return nil, errors.New("at least one RabbitMQ URL is required")
	}

	rmq := &RabbitMQ{
		urls:          urls,
		reconnectChan: make(chan bool),
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	go rmq.reconnectListener()

	return rmq, nil
}

// connect establishes a connection to the RabbitMQ cluster.
func (r *RabbitMQ) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, url := range r.urls {
		conn, err := amqp.Dial(url)
		if err == nil {
			r.conn = conn
			r.currentURL = url

			ch, err := conn.Channel()
			if err != nil {
				return fmt.Errorf("failed to open channel: %v", err)
			}
			r.channel = ch

			log.Printf("Connected to RabbitMQ at %s", url)
			return nil
		}
		log.Printf("Failed to connect to RabbitMQ at %s: %v", url, err)
	}

	return errors.New("failed to connect to any RabbitMQ node")
}

// reconnectListener listens for connection errors and attempts to reconnect.
func (r *RabbitMQ) reconnectListener() {
	for {
		select {
		case <-r.reconnectChan:
			if r.closed {
				return
			}
			log.Println("Attempting to reconnect to RabbitMQ...")
			if err := r.connect(); err != nil {
				log.Printf("Reconnect failed: %v", err)
				time.Sleep(5 * time.Second) // Wait before retrying
				r.reconnectChan <- true
			} else {
				log.Println("Reconnected to RabbitMQ")
			}
		}
	}
}

// DeclarePriorityQueue declares a priority queue with the given name.
func (r *RabbitMQ) DeclarePriorityQueue(queueName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return errors.New("RabbitMQ channel is not open")
	}

	// Declare a priority queue with x-max-priority set to 4
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

// PublishWithPriority publishes a message to a priority queue with the given priority.
func (r *RabbitMQ) PublishWithPriority(queueName string, message []byte, priority uint8) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return errors.New("RabbitMQ channel is not open")
	}

	// Publish the message with the specified priority
	err := r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent, // Ensure message durability
			Priority:     priority,        // Set message priority
		},
	)

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		r.reconnectChan <- true
		return err
	}

	log.Printf("Published message to queue %s with priority %d", queueName, priority)
	return nil
}

// Close closes the RabbitMQ connection.
func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.closed = true

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

	log.Println("RabbitMQ connection closed")
}
