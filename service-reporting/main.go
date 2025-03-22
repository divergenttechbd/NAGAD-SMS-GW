// main.go
package main

import (
	"reporting-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin
	router := gin.Default()

	// Set up reporting routes
	routes.SetupReportRoutes(router)

	// Start the server
	router.Run(":8082")
}