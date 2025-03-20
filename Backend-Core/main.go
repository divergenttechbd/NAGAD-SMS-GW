package main

import (
	"fmt"
	"log"
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

	// Initialize database
	db, err := utils.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models (only if ENABLE_AUTOMIGRATE is true)
	log.Println("AutoMigration Status: ", os.Getenv("ENABLE_AUTOMIGRATE"))

	if os.Getenv("ENABLE_AUTOMIGRATE") == "true" {
		log.Println("AutoMigrate is enabled. Running database migrations...")
		err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Campaign{},
			&models.SeederLog{}, &models.CampaignRecipient{}, &models.CampaignWorkflowProcessing{},
			&models.CampaignWorkflowProcessing{}, &models.CampaignWorkflow{}, &models.CampaignWorkflowUser{},
			&models.DND{}, &models.MNO{}, &models.MnoChannels{}, &models.MsgPriority{},
		)
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}

		// Seed initial data (only if not already seeded)
		// seeders.SeedData()
	} else {
		log.Println("AutoMigrate is disabled. Skipping database migrations.")
	}

	// Initialize Redis
	// redisClient := utils.InitRedis()

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	// router.Use(middleware.RateLimiter(redisClient))

	// Use the SetDBMiddleware to set the database connection in the context
	router.Use(middleware.SetDBMiddleware(db))

	// Swagger route (only if ENABLE_SWAGGER is true)
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		log.Println("Swagger is enabled. Serving API documentation at /swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		log.Println("Swagger is disabled. API documentation will not be served.")
	}

	/*---------- RabbitMQ ---------*/
	// Load RabbitMQ URLs from environment
	rabbitMQURLs := strings.Split(os.Getenv("RABBITMQ_URLS"), ",")

	// Initialize RabbitMQ
	rmq, err := rabbitmq.NewRabbitMQ(rabbitMQURLs)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	/* ------------- Declare Priority Queues ------------- *
	// Declare priority queues
	queues := []string{"general", "promotional", "transactional", "otp"}
	for _, queue := range queues {
		if err := rmq.DeclarePriorityQueue(queue); err != nil {
			log.Fatalf("Failed to declare queue %s: %v", queue, err)
		}
	}
	/* ------------- Declare Priority Queues ------------- */

	/*------------- Example Publish -------------*
	// Example: Publish messages with different priorities
	err = rmq.PublishWithPriority("general", []byte(`{"message": "This is a general message"}`), 1)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}

	err = rmq.PublishWithPriority("otp", []byte(`{"message": "This is an OTP message"}`), 4)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
	/*------------- Example Publish -------------*/

	/*------------- RabbitMQ -------------*/

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
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(fmt.Sprintf(":%s", port)))
}
