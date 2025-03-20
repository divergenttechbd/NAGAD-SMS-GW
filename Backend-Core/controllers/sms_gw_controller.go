package controllers

import (
	"context"
	"fmt"
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

	// Simulate queuing the message (Replace this with actual queue logic)
	// go func() {
	// 	time.Sleep(2 * time.Second) // Simulating processing delay
	// }()

	// Publish to RabbitMQ
	err := s.RabbitMQ.PublishWithPriority("general", []byte(`{"message": "This is a general message"}`), 1) // Adjust priority if needed
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	// Log the message in InfluxDB
	writeAPI := s.InfluxClient.WriteAPIBlocking(s.Config.InfluxDBOrg, s.Config.InfluxDBBucket)
	point := influxdb2.NewPoint("sms_delivery",
		map[string]string{
			"msg_id": msgID,
			"type":   "general", // otp, transactional, promotional etc
			"mno":    mno,
			"msisdn": smsReq.MSISDN,
			"text":   smsReq.SMSText,
			"status": "pending",
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
