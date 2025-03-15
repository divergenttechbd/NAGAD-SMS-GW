package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "myproject/docs" // Import the generated Swagger docs

	"myproject/config"
	"myproject/controllers"
	"myproject/middleware"
	"myproject/models"
	"myproject/routes"
	"myproject/utils"
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

	// Initialize utils (including configuration)
	utils.Init()

	// Initialize database
	db, err := utils.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models (only if ENABLE_AUTOMIGRATE is true)
	if os.Getenv("ENABLE_AUTOMIGRATE") == "true" {
		log.Println("AutoMigrate is enabled. Running database migrations...")
		err = db.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Campaign{},
			&models.SeederLog{})
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}

		// Seed initial data (only if not already seeded)
		// seeders.SeedData()
	} else {
		log.Println("AutoMigrate is disabled. Skipping database migrations.")
	}

	// Initialize Redis
	redisClient := utils.InitRedis()

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter(redisClient))

	// Use the SetDBMiddleware to set the database connection in the context
	router.Use(middleware.SetDBMiddleware(db))

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
	}

	// Swagger route (only if ENABLE_SWAGGER is true)
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		log.Println("Swagger is enabled. Serving API documentation at /swagger/index.html")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	} else {
		log.Println("Swagger is disabled. API documentation will not be served.")
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(fmt.Sprintf(":%s", port)))
}
