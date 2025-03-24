package main

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"
// 	"time"

// 	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
// 	"github.com/joho/godotenv"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// // RabbitMQ represents a connection to a RabbitMQ cluster.
// type RabbitMQ struct {
// 	conn    *amqp.Connection
// 	channel *amqp.Channel
// }

// // NewRabbitMQ creates a new RabbitMQ connection.
// func NewRabbitMQ(urls []string) (*RabbitMQ, error) {
// 	if len(urls) == 0 {
// 		return nil, errors.New("at least one RabbitMQ URL is required")
// 	}

// 	for _, url := range urls {
// 		conn, err := amqp.Dial(url)
// 		if err == nil {
// 			ch, err := conn.Channel()
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to open channel: %v", err)
// 			}

// 			log.Printf("Connected to RabbitMQ at %s", url)
// 			return &RabbitMQ{conn: conn, channel: ch}, nil
// 		}
// 		log.Printf("Failed to connect to RabbitMQ at %s: %v", url, err)
// 	}

// 	return nil, errors.New("failed to connect to any RabbitMQ node")
// }

// // DeclareQueue ensures the general queue exists.
// func (r *RabbitMQ) DeclareQueue(queueName string) error {
// 	_, err := r.channel.QueueDeclare(
// 		queueName, // queue name
// 		true,      // durable
// 		false,     // auto-delete
// 		false,     // exclusive
// 		false,     // no-wait
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to declare queue: %v", err)
// 	}

// 	log.Printf("Declared queue: %s", queueName)
// 	return nil
// }

// // ConsumeMessages starts consuming messages from the general queue.
// func (r *RabbitMQ) ConsumeMessages(queueName string, handler func([]byte) error) error {
// 	msgs, err := r.channel.Consume(
// 		queueName,
// 		"",
// 		false, // manual acknowledgment
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("failed to start consuming messages: %v", err)
// 	}

// 	log.Printf("Listening for messages on queue: %s", queueName)
// 	for msg := range msgs {
// 		log.Printf("Received message: %s", msg.Body)

// 		err := handler(msg.Body)
// 		if err != nil {
// 			log.Printf("Processing failed: %v", err)
// 			msg.Nack(false, true) // Requeue message on failure
// 		} else {
// 			msg.Ack(false) // Acknowledge message on success
// 		}
// 	}
// 	return nil
// }

// // Close closes the RabbitMQ connection.
// func (r *RabbitMQ) Close() {
// 	if r.channel != nil {
// 		r.channel.Close()
// 	}
// 	if r.conn != nil {
// 		r.conn.Close()
// 	}
// 	log.Println("RabbitMQ connection closed")
// }

// // InfluxDB represents the database connection.
// type InfluxDB struct {
// 	client influxdb2.Client
// 	org    string
// 	bucket string
// }

// // NewInfluxDB initializes the InfluxDB client.
// func NewInfluxDB(url, token, org, bucket string) *InfluxDB {
// 	client := influxdb2.NewClient(url, token)
// 	return &InfluxDB{client: client, org: org, bucket: bucket}
// }

// // UpdateSMSStatus updates the SMS status in InfluxDB.
// func (db *InfluxDB) UpdateSMSStatus(msgID, status string) error {
// 	writeAPI := db.client.WriteAPIBlocking(db.org, db.bucket)

// 	point := influxdb2.NewPoint("sms_delivery",
// 		map[string]string{
// 			"msg_id": msgID,
// 		},
// 		map[string]interface{}{
// 			"status": status,
// 		},
// 		time.Now(),
// 	)

// 	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
// 		return fmt.Errorf("failed to update InfluxDB: %v", err)
// 	}

// 	log.Printf("Updated msg_id %s to status %s in InfluxDB", msgID, status)
// 	return nil
// }

// // MessageHandler processes messages and updates InfluxDB accordingly.
// func MessageHandler(db *InfluxDB) func([]byte) error {
// 	return func(body []byte) error {
// 		var message map[string]interface{}
// 		if err := json.Unmarshal(body, &message); err != nil {
// 			return fmt.Errorf("failed to unmarshal message: %v", err)
// 		}

// 		msgID, ok := message["msg_id"].(string)
// 		if !ok || msgID == "" {
// 			return errors.New("invalid msg_id in message")
// 		}

// 		log.Printf("Processing SMS message: %v", message)

// 		// Simulate SMS sending
// 		time.Sleep(1 * time.Second) // Simulating delay

// 		// Simulate SMS delivery success or failure
// 		if success := simulateSMSDelivery(); success {
// 			db.UpdateSMSStatus(msgID, "successful")
// 		} else {
// 			db.UpdateSMSStatus(msgID, "failed")
// 		}

// 		return nil
// 	}
// }

// // Simulate SMS delivery logic (Replace with actual API call)
// func simulateSMSDelivery() bool {
// 	return time.Now().Unix()%2 == 0 // Random success/failure for testing
// }

// func main() {
// 	// Load environment variables from .env file
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	// Access environment variables
// 	rabbitMQURLs := strings.Split(os.Getenv("RABBITMQ_URLS"), ",")
// 	if len(rabbitMQURLs) == 0 {
// 		log.Fatal("RABBITMQ_URLS environment variable is not set")
// 	}

// 	influxDBURL := os.Getenv("INFLUXDB_URL")
// 	influxDBToken := os.Getenv("INFLUXDB_TOKEN")
// 	influxDBOrg := os.Getenv("INFLUXDB_ORG")
// 	influxDBBucket := os.Getenv("INFLUXDB_BUCKET")

// 	// Initialize RabbitMQ
// 	rmq, err := NewRabbitMQ(rabbitMQURLs)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}
// 	defer rmq.Close()

// 	// Declare the general queue
// 	queueName := "general"
// 	// if err := rmq.DeclareQueue(queueName); err != nil {
// 	// 	log.Fatalf("Failed to declare queue %s: %v", queueName, err)
// 	// }

// 	// Initialize InfluxDB
// 	influxDB := NewInfluxDB(influxDBURL, influxDBToken, influxDBOrg, influxDBBucket)
// 	defer influxDB.client.Close()

// 	// Start consuming messages
// 	if err := rmq.ConsumeMessages(queueName, MessageHandler(influxDB)); err != nil {
// 		log.Fatalf("Failed to consume messages: %v", err)
// 	}
// }
