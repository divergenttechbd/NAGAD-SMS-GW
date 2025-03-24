package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myproject/docs" // Import the generated Swagger docs
	"myproject/rabbitmq"

	"myproject/config"
	"myproject/controllers"
	"myproject/middleware"
	"myproject/models"
	"myproject/routes"
	"myproject/utils"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// @title My Project API
// @version 1.0
// @description This is a sample server for My Project.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables
	config.LoadEnv()
	gin.SetMode(gin.DebugMode) // Ensure debug mode

	// Initialize utils (including configuration)
	utils.Init()

	// Open or create log files
	appLogFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open app log file:", err)
	}
	defer appLogFile.Close()

	errorLogFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}
	defer errorLogFile.Close()

	// Create separate loggers
	appLogger := log.New(appLogFile, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger := log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Redirect default log output to appLogger
	log.SetOutput(appLogFile)
	log.SetPrefix("APP: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Test error logging
	errorLogger.Println("This is a test error message") // Should appear in error.log

	// Initialize database
	db, err := utils.InitDB()
	if err != nil {
		errorLogger.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models (only if ENABLE_AUTOMIGRATE is true)
	appLogger.Println("AutoMigration Status: ", os.Getenv("ENABLE_AUTOMIGRATE"))

	if os.Getenv("ENABLE_AUTOMIGRATE") == "true" {
		appLogger.Println("AutoMigrate is enabled. Running database migrations...")
		err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Campaign{},
			&models.SeederLog{}, &models.CampaignRecipient{}, &models.CampaignWorkflowProcessing{},
			&models.CampaignWorkflowProcessing{}, &models.CampaignWorkflow{}, &models.CampaignWorkflowUser{},
			&models.DND{}, &models.MNO{}, &models.MnoChannels{}, &models.MsgPriority{},
		)
		if err != nil {
			errorLogger.Fatal("Failed to migrate database:", err)
		}
	} else {
		appLogger.Println("AutoMigrate is disabled. Skipping database migrations.")
	}

	// Initialize Redis
	redisClient := utils.InitRedis()

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(redisClient))
	router.Use(ErrorLogger(errorLogger)) // Add error logging middleware with errorLogger

	// Use the SetDBMiddleware to set the database connection in the context
	router.Use(middleware.SetDBMiddleware(db))

	// Swagger route (only if ENABLE_SWAGGER is true)
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		appLogger.Println("Swagger is enabled. Serving API documentation at /swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		appLogger.Println("Swagger is disabled. API documentation will not be served.")
	}

	// Initialize RabbitMQ
	rabbitMQURLs := strings.Split(os.Getenv("RABBITMQ_URLS"), ",")
	rabbitMQmanagementURL := os.Getenv("RABBITMQ_MANAGEMENT_URL")
	rabbitMQusername := os.Getenv("RABBITMQ_USER")
	rabbitMQpassword := os.Getenv("RABBITMQ_PASSWORD")

	rmq, err := rabbitmq.NewRabbitMQ(rabbitMQURLs, rabbitMQmanagementURL, rabbitMQusername, rabbitMQpassword)
	if err != nil {
		errorLogger.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Load Configuration
	cfg := config.GetConfig()

	// Initialize InfluxDB client
	influxClient := influxdb2.NewClient(cfg.InfluxDBURL, cfg.InfluxDBToken)
	defer influxClient.Close()

	// Routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", controllers.Login)
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/verify-token", controllers.VerifyToken)
	}

	apiRoutes := router.Group("/api")
	apiRoutes.Use(middleware.JWTAuth())
	{
		routes.SetupUserRoutes(apiRoutes)
		routes.SetupCampaignRoutes(apiRoutes)
		routes.SetupMNORoutes(apiRoutes)
		routes.SetupDndRoutes(apiRoutes)
		routes.SetupMsgPriorityRoutes(apiRoutes)
		routes.SetupCampaignRecipientRoutes(apiRoutes)
		routes.SetupCampaignWorkflowRoutes(apiRoutes)
		routes.SetupSMSGatewayRoutes(apiRoutes, influxClient, cfg, rmq)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	appLogger.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(fmt.Sprintf(":%s", port)))
}

// ErrorLogger middleware to log errors and recover from panics
func ErrorLogger(errorLogger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errorLogger.Printf("Panic recovered: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"details": fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
