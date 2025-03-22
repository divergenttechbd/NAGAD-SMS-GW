// routes/gateway.go
package routes

import (
	"api-gateway/middleware"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupGatewayRoutes(router *gin.Engine) {
	// Apply token verification middleware to all routes
	router.Use(middleware.TokenVerifyMiddleware())

	// Route to Core Service
	router.Any("/core/*path", func(c *gin.Context) {
		// Forward the request to the Core Service
		forwardRequest(c, os.Getenv("CORE_SERVICE_URL"))
	})

	// Route to Reporting Service
	router.Any("/reporting/*path", func(c *gin.Context) {
		// Forward the request to the Reporting Service
		forwardRequest(c, os.Getenv("REPORTING_SERVICE_URL"))
	})
}

// forwardRequest forwards the request to the specified backend service
func forwardRequest(c *gin.Context, backendURL string) {
	path := c.Param("path")

	// log.Printf("Forwarding request to %s", backendURL+path)
	// log.Printf("Forwarding request to %s", c.Request.Method)

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request based on the incoming request
	req, err := http.NewRequest(c.Request.Method, backendURL+path, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers from the incoming request to the new request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send the request to the backend service
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Return the response from the backend service to the client
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}
