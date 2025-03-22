package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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
	managementURL string // URL for RabbitMQ Management API (e.g., http://localhost:15672)
	username      string // Management API username
	password      string // Management API password
}

// QueueStats represents key statistics for a RabbitMQ queue
type QueueStats struct {
	Messages        int64   `json:"messages"`
	MessagesReady   int64   `json:"messages_ready"`
	MessagesUnacked int64   `json:"messages_unacknowledged"`
	PublishRate     float64 `json:"publish_rate,omitempty"`
	DeliverRate     float64 `json:"deliver_rate,omitempty"`
	AcknowledgeRate float64 `json:"acknowledge_rate,omitempty"`
	ConsumerCount   int     `json:"consumers"`
}

// NodeStats represents key statistics for a RabbitMQ node
type NodeStats struct {
	MemoryUsed    int64 `json:"mem_used"`
	FileDescUsed  int   `json:"fd_used"`
	FileDescTotal int   `json:"fd_total"`
	DiskFree      int64 `json:"disk_free"`
	Connections   int   `json:"connections"`
}

// Statistics aggregates queue and node statistics
type Statistics struct {
	Queues map[string]QueueStats `json:"queues"`
	Node   NodeStats             `json:"node"`
}

// NewRabbitMQ creates a new RabbitMQ connection with Management API access.
func NewRabbitMQ(urls []string, managementURL, username, password string) (*RabbitMQ, error) {
	if len(urls) == 0 {
		return nil, errors.New("at least one RabbitMQ URL is required")
	}
	if managementURL == "" {
		return nil, errors.New("management URL is required for statistics")
	}

	rmq := &RabbitMQ{
		urls:          urls,
		reconnectChan: make(chan bool),
		managementURL: managementURL,
		username:      username,
		password:      password,
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

	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-max-priority": 4,
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
	if r == nil {
		return fmt.Errorf("RabbitMQ instance is nil")
	}

	if r.channel == nil {
		return errors.New("RabbitMQ channel is not open")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         message,
			DeliveryMode: amqp.Persistent,
			Priority:     priority,
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

// GetStatistics retrieves key RabbitMQ statistics from the Management API.
func (r *RabbitMQ) GetStatistics() (Statistics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var stats Statistics
	stats.Queues = make(map[string]QueueStats)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/queues", r.managementURL), nil)
	if err != nil {
		return stats, fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(r.username, r.password)

	resp, err := client.Do(req)
	if err != nil {
		return stats, fmt.Errorf("failed to fetch queue stats: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return stats, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var queues []struct {
		Name            string `json:"name"`
		Messages        int64  `json:"messages"`
		MessagesReady   int64  `json:"messages_ready"`
		MessagesUnacked int64  `json:"messages_unacknowledged"`
		Consumers       int    `json:"consumers"`
		MessageStats    struct {
			PublishDetails struct {
				Rate float64 `json:"rate"`
			} `json:"publish_details"`
			DeliverDetails struct {
				Rate float64 `json:"rate"`
			} `json:"deliver_details"`
			AckDetails struct {
				Rate float64 `json:"rate"`
			} `json:"ack_details"`
		} `json:"message_stats"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&queues); err != nil {
		return stats, fmt.Errorf("failed to decode queue stats: %v", err)
	}

	for _, q := range queues {
		stats.Queues[q.Name] = QueueStats{
			Messages:        q.Messages,
			MessagesReady:   q.MessagesReady,
			MessagesUnacked: q.MessagesUnacked,
			PublishRate:     q.MessageStats.PublishDetails.Rate,     // Fixed typo here
			DeliverRate:     q.MessageStats.DeliverDetails.Rate,     // Fixed typo here
			AcknowledgeRate: q.MessageStats.AckDetails.Rate,         // Fixed typo here
			ConsumerCount:   q.Consumers,
		}
	}

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/nodes", r.managementURL), nil)
	if err != nil {
		return stats, fmt.Errorf("failed to create node request: %v", err)
	}
	req.SetBasicAuth(r.username, r.password)

	resp, err = client.Do(req)
	if err != nil {
		return stats, fmt.Errorf("failed to fetch node stats: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return stats, fmt.Errorf("unexpected node status code: %d", resp.StatusCode)
	}

	var nodes []NodeStats
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return stats, fmt.Errorf("failed to decode node stats: %v", err)
	}
	if len(nodes) > 0 {
		stats.Node = nodes[0] // Single node assumption
	}

	return stats, nil
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