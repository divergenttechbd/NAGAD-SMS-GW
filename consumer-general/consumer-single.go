package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueName         = "general"
	MaxWorkers        = 200
	PrefetchCount     = 500
	RedisLockTTL      = 30 * time.Second
	RedisRateWindow   = 1 * time.Second
	InfluxBatchSize   = 5000
	HeartbeatInterval = 5 * time.Second
	ReconnectDelay    = 5 * time.Second
)

// MNO-specific TPS limits (slightly lower than provided to be safe)
var mnoTPSLimits = map[string]int{
	"Robi":     1, // Provided: 2, Using: 1
	"Gp":       2, // Provided: 3, Using: 2
	"Airtel":   4, // Provided: 5, Using: 4
	"Bl":       5, // Provided: 6, Using: 5
	"Teletalk": 1, // Provided: 1, Using: 1 (no lower option, but safe)
	"Tpgw":     1, // Provided: 1, Using: 1 (no lower option, but safe)
}

var (
	ctx               = context.Background()
	influxFlushPeriod = 100 * time.Millisecond
)

type SafeConsumer struct {
	instanceID   string
	rabbitConn   *amqp.Connection
	rabbitChan   *amqp.Channel
	influxAPI    api.WriteAPI
	redisClient  *redis.Client
	successCount uint64
	failureCount uint64
	rateLimited  uint64
	rabbitURLs   []string
}

func NewSafeConsumer() (*SafeConsumer, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning loading .env: %v", err)
	}

	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		return nil, errors.New("INSTANCE_ID must be set")
	}

	rabbitURLs := strings.Split(os.Getenv("RABBITMQ_URLS"), ",")
	if len(rabbitURLs) == 0 || rabbitURLs[0] == "" {
		return nil, errors.New("RABBITMQ_URLS must be set with at least one URL")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS_URL must be set")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	influxURL := os.Getenv("INFLUXDB_URL")
	influxToken := os.Getenv("INFLUXDB_TOKEN")
	influxOrg := os.Getenv("INFLUXDB_ORG")
	influxBucket := os.Getenv("INFLUXDB_BUCKET")
	if influxURL == "" || influxToken == "" || influxOrg == "" || influxBucket == "" {
		return nil, errors.New("INFLUXDB_URL, INFLUXDB_TOKEN, INFLUXDB_ORG, and INFLUXDB_BUCKET must be set")
	}

	influxClient := influxdb2.NewClientWithOptions(
		influxURL,
		influxToken,
		influxdb2.DefaultOptions().
			SetBatchSize(InfluxBatchSize).
			SetFlushInterval(uint(influxFlushPeriod.Milliseconds())),
	)
	writeAPI := influxClient.WriteAPI(influxOrg, influxBucket)

	conn, err := connectRabbitMQ(rabbitURLs)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ connection failed: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Channel creation failed: %v", err)
	}

	if err := ch.Qos(PrefetchCount, 0, false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("Qos failed: %v", err)
	}

	return &SafeConsumer{
		instanceID:  instanceID,
		rabbitConn:  conn,
		rabbitChan:  ch,
		influxAPI:   writeAPI,
		redisClient: redisClient,
		rabbitURLs:  rabbitURLs,
	}, nil
}

func connectRabbitMQ(urls []string) (*amqp.Connection, error) {
	for _, url := range urls {
		conn, err := amqp.DialConfig(url, amqp.Config{
			Heartbeat: HeartbeatInterval,
		})
		if err == nil {
			log.Printf("Connected to RabbitMQ at %s", url)
			return conn, nil
		}
		log.Printf("Failed to connect to %s: %v", url, err)
	}
	return nil, errors.New("failed to connect to any RabbitMQ node")
}

func (c *SafeConsumer) Close() {
	if c.rabbitChan != nil {
		c.rabbitChan.Close()
	}
	if c.rabbitConn != nil {
		c.rabbitConn.Close()
	}
	if c.redisClient != nil {
		c.redisClient.Close()
	}
	c.influxAPI.Flush()
}

func (c *SafeConsumer) acquireMessageLock(messageID string) (bool, error) {
	return c.redisClient.SetNX(ctx, "lock:"+messageID, c.instanceID, RedisLockTTL).Result()
}

func (c *SafeConsumer) checkMNORateLimit(mno string) (bool, error) {
	tpsLimit, ok := mnoTPSLimits[mno]
	if !ok {
		return true, fmt.Errorf("unknown MNO: %s", mno)
	}

	window := time.Now().Truncate(RedisRateWindow).Unix()
	key := fmt.Sprintf("rate:%s:%d", mno, window)

	count, err := c.redisClient.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		c.redisClient.Expire(ctx, key, RedisRateWindow)
	}

	return count > int64(tpsLimit), nil
}

func (c *SafeConsumer) submitToMNOAPI(mno, msgID string) error {
	// Simulate API call (replace with actual HTTP call to MNO SMS API)
	log.Printf("Submitting to %s SMS API for msg_id: %s", mno, msgID)
	time.Sleep(50 * time.Millisecond) // Simulate network delay
	return nil
}

func (c *SafeConsumer) ProcessMessage(msg amqp.Delivery) {
	locked, err := c.acquireMessageLock(msg.MessageId)
	if err != nil || !locked {
		msg.Nack(false, true)
		return
	}
	defer c.redisClient.Del(ctx, "lock:"+msg.MessageId)

	var message struct {
		MsgID string `json:"msg_id"`
		MNO   string `json:"mno"`
	}
	if err := json.Unmarshal(msg.Body, &message); err != nil {
		atomic.AddUint64(&c.failureCount, 1)
		msg.Nack(false, true)
		return
	}

	if message.MNO == "" {
		log.Printf("Message %s has no MNO specified", message.MsgID)
		atomic.AddUint64(&c.failureCount, 1)
		msg.Nack(false, true)
		return
	}

	limited, err := c.checkMNORateLimit(message.MNO)
	if err != nil {
		log.Printf("Rate limit check failed for %s: %v", message.MNO, err)
		atomic.AddUint64(&c.failureCount, 1)
		msg.Nack(false, true)
		return
	}
	if limited {
		atomic.AddUint64(&c.rateLimited, 1)
		msg.Ack(false)
		return
	}

	processingTime := time.Duration(50+rand.Intn(100)) * time.Millisecond
	time.Sleep(processingTime)

	// Submit to MNO SMS API and set status
	var status string
	err = c.submitToMNOAPI(message.MNO, message.MsgID)
	if err != nil {
		status = "failed"
		atomic.AddUint64(&c.failureCount, 1)
		log.Printf("Failed to submit to %s API: %v", message.MNO, err)
	} else {
		status = "delivered"
		atomic.AddUint64(&c.successCount, 1)
	}

	point := influxdb2.NewPoint(
		"final_sms_delivery",
		map[string]string{
			"msg_id":   message.MsgID,
			"instance": c.instanceID,
			"mno":      message.MNO,
			"status":   status,
		},
		map[string]interface{}{
			"processing_time_ms": processingTime.Milliseconds(),
		},
		time.Now(),
	)
	c.influxAPI.WritePoint(point)

	msg.Ack(false)
}

func (c *SafeConsumer) Run() error {
	go func() {
		for err := range c.influxAPI.Errors() {
			log.Printf("InfluxDB write error: %v", err)
		}
	}()

	msgs, err := c.rabbitChan.Consume(
		QueueName,
		"sms-consumer-"+c.instanceID,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Printf("[%s] Stats - Success: %d, Failure: %d, RateLimited: %d",
				c.instanceID,
				atomic.LoadUint64(&c.successCount),
				atomic.LoadUint64(&c.failureCount),
				atomic.LoadUint64(&c.rateLimited),
			)
		}
	}()

	sem := make(chan struct{}, MaxWorkers)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c.rabbitConn.NotifyClose(make(chan *amqp.Error))
		log.Printf("RabbitMQ connection closed, attempting to reconnect...")
		for {
			conn, err := connectRabbitMQ(c.rabbitURLs)
			if err == nil {
				c.rabbitConn = conn
				ch, err := conn.Channel()
				if err == nil {
					c.rabbitChan = ch
					if err := ch.Qos(PrefetchCount, 0, false); err == nil {
						msgs, err = ch.Consume(QueueName, "sms-consumer-"+c.instanceID, false, false, false, false, nil)
						if err == nil {
							log.Printf("Reconnected and resumed consuming")
							return
						}
					}
				}
			}
			log.Printf("Reconnect failed: %v, retrying in %v", err, ReconnectDelay)
			time.Sleep(ReconnectDelay)
		}
	}()

	for {
		select {
		case msg := <-msgs:
			sem <- struct{}{}
			go func(m amqp.Delivery) {
				defer func() { <-sem }()
				c.ProcessMessage(m)
			}(msg)
		case <-sig:
			log.Printf("[%s] Shutting down gracefully...", c.instanceID)
			return nil
		}
	}
}

func main() {
	consumer, err := NewSafeConsumer()
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.Run(); err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
}
