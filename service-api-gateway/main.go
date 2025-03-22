// main.go
package main

import (
	"api-gateway/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Gin
	router := gin.Default()

	// Set up gateway routes
	routes.SetupGatewayRoutes(router)

	// Start the server
	router.Run(":8080")
}
