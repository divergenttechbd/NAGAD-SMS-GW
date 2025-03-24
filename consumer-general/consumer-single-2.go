package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"os/signal"
// 	"path/filepath"
// 	"strings"
// 	"sync"
// 	"syscall"
// 	"time"

// 	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
// 	"github.com/influxdata/influxdb-client-go/v2/api"
// 	"github.com/joho/godotenv"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"golang.org/x/time/rate"
// )

// // Constants
// const (
// 	QueueName         = "general"
// 	SMSRateLimit      = 500
// 	SMSRateInterval   = time.Second
// 	MaxWorkers        = 100
// 	InfluxBatchSize   = 1000
// 	InfluxFlushPeriod = 5 * time.Second
// 	ReconnectDelay    = 5 * time.Second
// 	ConnectionTimeout = 30 * time.Second
// 	DefaultMessageTTL = 86400 * 1000 // 24 hours in milliseconds
// 	LogDir            = "logs"
// )

// var (
// 	infoLogger  *log.Logger
// 	errorLogger *log.Logger
// )

// // Initialize logging
// func initLogging() error {
// 	// Create logs directory if it doesn't exist
// 	if err := os.MkdirAll(LogDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create logs directory: %v", err)
// 	}

// 	// Create date-wise log file
// 	currentDate := time.Now().Format("02012006") // DDMMYYYY format
// 	logFileName := filepath.Join(LogDir, "consumer"+currentDate+".log")

// 	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return fmt.Errorf("failed to open log file: %v", err)
// 	}

// 	// Create multi-writer to log to both file and stdout
// 	multiWriter := io.MultiWriter(os.Stdout, logFile)

// 	infoLogger = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	errorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

// 	return nil
// }

// // RabbitMQ implementation (same as before with logger updates)
// type RabbitMQ struct {
// 	urls     []string
// 	conn     *amqp.Connection
// 	channel  *amqp.Channel
// 	mu       sync.Mutex
// 	quitChan chan struct{}
// }

// func NewRabbitMQ(urls []string) (*RabbitMQ, error) {
// 	if len(urls) == 0 {
// 		return nil, errors.New("at least one RabbitMQ URL is required")
// 	}

// 	rmq := &RabbitMQ{
// 		urls:     urls,
// 		quitChan: make(chan struct{}),
// 	}

// 	if err := rmq.connect(); err != nil {
// 		return nil, err
// 	}

// 	go rmq.connectionWatcher()
// 	return rmq, nil
// }

// func (r *RabbitMQ) connect() error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	for _, url := range r.urls {
// 		conn, err := amqp.DialConfig(url, amqp.Config{
// 			Dial:      amqp.DefaultDial(ConnectionTimeout),
// 			Heartbeat: 10 * time.Second,
// 		})
// 		if err == nil {
// 			ch, err := conn.Channel()
// 			if err != nil {
// 				conn.Close()
// 				continue
// 			}

// 			if err := ch.Qos(100, 0, false); err != nil {
// 				ch.Close()
// 				conn.Close()
// 				continue
// 			}

// 			r.conn = conn
// 			r.channel = ch
// 			infoLogger.Printf("Connected to RabbitMQ at %s", obfuscateURL(url))
// 			return nil
// 		}
// 	}

// 	return errors.New("failed to connect to any RabbitMQ node")
// }

// func (r *RabbitMQ) connectionWatcher() {
// 	notifyClose := make(chan *amqp.Error)
// 	r.channel.NotifyClose(notifyClose)

// 	select {
// 	case <-notifyClose:
// 		errorLogger.Println("RabbitMQ connection lost, attempting to reconnect...")
// 		for {
// 			select {
// 			case <-r.quitChan:
// 				return
// 			default:
// 				if err := r.connect(); err == nil {
// 					infoLogger.Println("Successfully reconnected to RabbitMQ")
// 					go r.connectionWatcher()
// 					return
// 				}
// 				time.Sleep(ReconnectDelay)
// 			}
// 		}
// 	case <-r.quitChan:
// 		return
// 	}
// }

// func obfuscateURL(url string) string {
// 	parts := strings.Split(url, "@")
// 	if len(parts) > 1 {
// 		return "***@" + parts[1]
// 	}
// 	return url
// }

// func (r *RabbitMQ) DeclareQueue(queueName string) error {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	_, err := r.channel.QueueDeclare(
// 		queueName,
// 		true,
// 		false,
// 		false,
// 		false,
// 		amqp.Table{
// 			"x-message-ttl": DefaultMessageTTL,
// 		},
// 	)
// 	return err
// }

// func (r *RabbitMQ) ConsumeMessages(queueName string, handler func([]byte) error) error {
// 	msgs, err := r.channel.Consume(
// 		queueName,
// 		"",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	var wg sync.WaitGroup
// 	workerChan := make(chan amqp.Delivery, MaxWorkers*2)

// 	for i := 0; i < MaxWorkers; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for msg := range workerChan {
// 				if err := handler(msg.Body); err != nil {
// 					errorLogger.Printf("Processing failed: %v", err)
// 					msg.Nack(false, true)
// 				} else {
// 					msg.Ack(false)
// 				}
// 			}
// 		}()
// 	}

// 	go func() {
// 		for msg := range msgs {
// 			select {
// 			case workerChan <- msg:
// 			default:
// 				infoLogger.Println("Worker pool overloaded, slowing down consumption")
// 				time.Sleep(100 * time.Millisecond)
// 				workerChan <- msg
// 			}
// 		}
// 		close(workerChan)
// 	}()

// 	wg.Wait()
// 	return nil
// }

// func (r *RabbitMQ) Close() {
// 	close(r.quitChan)
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	if r.channel != nil {
// 		r.channel.Close()
// 	}
// 	if r.conn != nil {
// 		r.conn.Close()
// 	}
// 	infoLogger.Println("RabbitMQ connection closed")
// }

// // InfluxDB implementation with proper WriteAPI type
// type InfluxDB struct {
// 	client   influxdb2.Client
// 	writeAPI api.WriteAPI
// 	org      string
// 	bucket   string
// }

// func NewInfluxDB(url, token, org, bucket string) *InfluxDB {
// 	client := influxdb2.NewClient(url, token)
// 	writeAPI := client.WriteAPI(org, bucket)

// 	errorsCh := writeAPI.Errors()
// 	go func() {
// 		for err := range errorsCh {
// 			errorLogger.Printf("InfluxDB write error: %v", err)
// 		}
// 	}()

// 	return &InfluxDB{
// 		client:   client,
// 		writeAPI: writeAPI,
// 		org:      org,
// 		bucket:   bucket,
// 	}
// }

// func (db *InfluxDB) UpdateSMSStatus(msgID, status string) {
// 	point := influxdb2.NewPoint(
// 		"sms_delivery",
// 		map[string]string{"msg_id": msgID},
// 		map[string]interface{}{
// 			"status":    status,
// 			"timestamp": time.Now().UnixNano(),
// 		},
// 		time.Now(),
// 	)
// 	db.writeAPI.WritePoint(point)
// }

// func (db *InfluxDB) Close() {
// 	db.writeAPI.Flush()
// 	db.client.Close()
// }

// // SMSProcessor implementation
// type SMSProcessor struct {
// 	db           *InfluxDB
// 	limiter      *rate.Limiter
// 	successCnt   uint64
// 	failureCnt   uint64
// 	rateLimitCnt uint64
// 	mu           sync.Mutex
// }

// func NewSMSProcessor(db *InfluxDB) *SMSProcessor {
// 	return &SMSProcessor{
// 		db:      db,
// 		limiter: rate.NewLimiter(rate.Every(SMSRateInterval/SMSRateLimit), SMSRateLimit),
// 	}
// }

// func (p *SMSProcessor) ProcessMessage(body []byte) error {
// 	var message struct {
// 		MsgID string `json:"msg_id"`
// 	}

// 	if err := json.Unmarshal(body, &message); err != nil {
// 		return fmt.Errorf("failed to unmarshal message: %v", err)
// 	}

// 	if message.MsgID == "" {
// 		return errors.New("empty msg_id")
// 	}

// 	if !p.limiter.Allow() {
// 		p.mu.Lock()
// 		p.rateLimitCnt++
// 		p.mu.Unlock()
// 		p.db.UpdateSMSStatus(message.MsgID, "rate_limited")
// 		return nil
// 	}

// 	processingTime := time.Duration(50+rand.Intn(150)) * time.Millisecond
// 	time.Sleep(processingTime)

// 	if rand.Float32() < 0.95 {
// 		p.mu.Lock()
// 		p.successCnt++
// 		p.mu.Unlock()
// 		p.db.UpdateSMSStatus(message.MsgID, "delivered")
// 	} else {
// 		p.mu.Lock()
// 		p.failureCnt++
// 		p.mu.Unlock()
// 		p.db.UpdateSMSStatus(message.MsgID, "failed")
// 	}

// 	return nil
// }

// func (p *SMSProcessor) Stats() (success, failure, rateLimited uint64) {
// 	p.mu.Lock()
// 	defer p.mu.Unlock()
// 	return p.successCnt, p.failureCnt, p.rateLimitCnt
// }

// func main() {
// 	// Initialize logging
// 	if err := initLogging(); err != nil {
// 		log.Fatalf("Failed to initialize logging: %v", err)
// 	}

// 	// Load environment variables
// 	if err := godotenv.Load(); err != nil {
// 		infoLogger.Printf("Warning: Error loading .env file: %v", err)
// 	}

// 	// Initialize components
// 	rmq, err := NewRabbitMQ(strings.Split(os.Getenv("RABBITMQ_URLS"), ","))
// 	if err != nil {
// 		errorLogger.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}
// 	defer rmq.Close()

// 	influxDB := NewInfluxDB(
// 		os.Getenv("INFLUXDB_URL"),
// 		os.Getenv("INFLUXDB_TOKEN"),
// 		os.Getenv("INFLUXDB_ORG"),
// 		os.Getenv("INFLUXDB_BUCKET"),
// 	)
// 	defer influxDB.Close()

// 	processor := NewSMSProcessor(influxDB)

// 	// Set up graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 	// Start periodic stats logging
// 	go func() {
// 		ticker := time.NewTicker(10 * time.Second)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-ticker.C:
// 				success, failure, rateLimited := processor.Stats()
// 				infoLogger.Printf("Stats - Success: %d, Failure: %d, RateLimited: %d",
// 					success, failure, rateLimited)
// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()

// 	// Start consuming messages
// 	infoLogger.Printf("Starting consumer with %d workers...", MaxWorkers)
// 	if err := rmq.ConsumeMessages(QueueName, processor.ProcessMessage); err != nil {
// 		errorLogger.Fatalf("Failed to consume messages: %v", err)
// 	}

// 	<-quit
// 	infoLogger.Println("Shutting down gracefully...")
// }
