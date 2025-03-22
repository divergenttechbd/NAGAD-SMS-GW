package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"myproject/config"
	"myproject/rabbitmq"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// SMSRequest represents an incoming SMS API request
type SMSRequest struct {
	SMSText string `json:"sms_text" example:"Hello, this is a test message"`
	MSISDN  string `json:"msisdn" example:"01712345678"`
}

// MessagePayload defines the structure of the message to be sent
type MessagePayload struct {
	MNO    string `json:"mno"`
	MsgID  string `json:"msg_id"`
	MSISDN string `json:"msisdn"`
	Status string `json:"status"`
	Text   string `json:"text"`
	Type   string `json:"type"`
}

// SMSGatewayController handles SMS processing
type SMSGatewayController struct {
	InfluxClient influxdb2.Client
	Config       *config.Config
	RabbitMQ     *rabbitmq.RabbitMQ
}

// NewSMSGatewayController initializes an SMSGatewayController
func NewSMSGatewayController(client influxdb2.Client, cfg *config.Config, rmq *rabbitmq.RabbitMQ) *SMSGatewayController {
	return &SMSGatewayController{InfluxClient: client, Config: cfg, RabbitMQ: rmq}
}

// ProcessSMS receives an SMS request and processes it
// @Summary Send an SMS message
// @Description Receives an SMS text and MSISDN, determines the carrier, queues the message, and logs it in InfluxDB
// @Tags SMS Gateway
// @Accept json
// @Produce json
// @Param smsRequest body SMSRequest true "SMS request payload"
// @Success 200 {object} map[string]interface{} "SMS received and queued"
// @Failure 400 {object} map[string]string "Invalid request format or carrier prefix"
// @Failure 500 {object} map[string]string "Failed to write to InfluxDB"
// @Router /sms/send [post]
func (s *SMSGatewayController) ProcessSMS(c *gin.Context) {
	var smsReq SMSRequest
	if err := c.ShouldBindJSON(&smsReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if len(smsReq.MSISDN) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "MSISDN must be an 11-digit number"})
		return
	}

	// Determine MNO based on MSISDN prefix
	mno := getMNO(smsReq.MSISDN[:3])
	if mno == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid carrier prefix"})
		return
	}

	// Generate unique message ID
	msgID := generateMsgID()

	// Prepare message payload for RabbitMQ
	messageData := map[string]string{
		"msg_id": msgID,
		"msisdn": smsReq.MSISDN,
		"text":   smsReq.SMSText,
		"mno":    mno,
		"type":   "general", // Can be OTP, transactional, promotional, etc.
		"status": "queued",
	}

	messageJSON, err := json.Marshal(messageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize message data"})
		return
	}

	// Publish to RabbitMQ queue
	err = s.RabbitMQ.PublishWithPriority("general", messageJSON, 1) // Adjust priority if needed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message to RabbitMQ"})
		return
	}

	// Log the message in InfluxDB
	writeAPI := s.InfluxClient.WriteAPIBlocking(s.Config.InfluxDBOrg, s.Config.InfluxDBBucket)
	point := influxdb2.NewPoint("sms_delivery",
		map[string]string{
			"msg_id": msgID,
			"type":   "general",
			"mno":    mno,
			"msisdn": smsReq.MSISDN,
			"text":   smsReq.SMSText,
			"status": "queued",
		},
		map[string]interface{}{
			"retry_count":           0,
			"queue_time":            time.Now().UnixMilli(),
			"carrier_response_time": 0,
		},
		time.Now())

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to InfluxDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS received and queued", "msg_id": msgID})
}

// PublishMillionMessages publishes 1 million messages to the specified queue and logs them in InfluxDB
// @Summary Publish 1 million test SMS messages
// @Description Publishes 1 million messages to a specified RabbitMQ queue with priority and logs them in InfluxDB
// @Tags SMS Gateway
// @Produce json
// @Param queueName query string true "Queue name (e.g., general, otp)" default(general)
// @Param priority query uint8 true "Priority level (0-4)" default(1)
// @Success 200 {object} map[string]interface{} "Messages published successfully"
// @Failure 500 {object} map[string]string "Failed to publish messages"
// @Router /sms/publish-million [post]
func (s *SMSGatewayController) PublishMillionMessages(c *gin.Context) {
	const totalMessages = 1_000_000
	const batchSize = 10_000 // Batch size for progress logging

	// Get query parameters
	queueName := c.DefaultQuery("queueName", "general")
	priorityStr := c.DefaultQuery("priority", "1")
	var priority uint8
	fmt.Sscanf(priorityStr, "%d", &priority)
	if priority > 4 {
		priority = 4 // Cap at max priority
	}

	// Base message payload
	baseMsg := MessagePayload{
		MNO:    "Robi",
		MsgID:  "2025032102343877835", // Will be overridden for uniqueness
		MSISDN: "01814266295",
		Status: "queued",
		Text:   "",
		Type:   "general",
	}

	startTime := time.Now()
	log.Printf("Starting to publish %d messages to queue %s with priority %d", totalMessages, queueName, priority)

	for i := 0; i < totalMessages; i++ {
		// Create a unique MsgID for each message
		msg := baseMsg
		msg.MsgID = fmt.Sprintf("20250321%010d", i) // Unique ID: 202503210000000001 to 20250321000999999

		// Marshal to JSON for RabbitMQ
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to marshal message %d: %v", i, err)})
			return
		}

		// Publish to RabbitMQ
		err = s.RabbitMQ.PublishWithPriority(queueName, msgBytes, priority)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to publish message %d: %v", i, err)})
			return
		}

		// Log to InfluxDB
		writeAPI := s.InfluxClient.WriteAPIBlocking(s.Config.InfluxDBOrg, s.Config.InfluxDBBucket)
		point := influxdb2.NewPoint("sms_delivery",
			map[string]string{
				"msg_id": msg.MsgID,
				"type":   msg.Type,
				"mno":    msg.MNO,
				"msisdn": msg.MSISDN,
				"text":   msg.Text,
				"status": msg.Status,
			},
			map[string]interface{}{
				"retry_count":           0,
				"queue_time":            time.Now().UnixMilli(),
				"carrier_response_time": 0,
			},
			time.Now())

		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to write message %d to InfluxDB: %v", i, err)})
			return
		}

		// Log progress every batchSize messages
		if (i+1)%batchSize == 0 {
			elapsed := time.Since(startTime)
			log.Printf("Published %d of %d messages (%.2f%%) to queue and InfluxDB in %v", i+1, totalMessages, float64(i+1)/float64(totalMessages)*100, elapsed)
		}
	}

	// Calculate and log total time taken
	totalDuration := time.Since(startTime)
	throughput := float64(totalMessages) / totalDuration.Seconds()
	log.Printf("Completed: Published %d messages to queue %s and InfluxDB in %v (throughput: %.2f messages/second)", totalMessages, queueName, totalDuration, throughput)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Successfully published 1 million messages",
		"queue":      queueName,
		"priority":   priority,
		"duration":   totalDuration.String(),
		"throughput": fmt.Sprintf("%.2f msg/s", throughput),
	})
}

func (s *SMSGatewayController) GetRabbitMQStatistics(c *gin.Context) {
	stats, err := s.RabbitMQ.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve RabbitMQ statistics: %v", err)})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// getMNO determines the carrier based on MSISDN prefix
func getMNO(prefix string) string {
	switch prefix {
	case "018":
		return "Robi"
	case "017":
		return "GP"
	case "016":
		return "Airtel"
	default:
		return ""
	}
}

// generateMsgID generates a unique message ID
func generateMsgID() string {
	return time.Now().Format("20060102150405") + fmt.Sprint(rand.Intn(100000))
}
