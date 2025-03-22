// routes/reports.go
package routes

import (
	"encoding/json"
	"net/http"
	"reporting-service/config"

	"github.com/gin-gonic/gin"
)

func SetupReportRoutes(router *gin.Engine) {
	reportGroup := router.Group("/api/reports")

	// Endpoint for generating reports
	reportGroup.GET("/generate", func(c *gin.Context) {
		// Fetch data from the Core Service via the API Gateway
		data, err := fetchDataFromCoreService(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from Core Service"})
			return
		}

		// Process the data and generate a report
		report := generateReport(data)

		// Return the report
		c.JSON(http.StatusOK, gin.H{"message": "Report generated successfully", "report": report})
	})

	// Endpoint for viewing reports
	reportGroup.GET("/view", func(c *gin.Context) {
		// Fetch data from the Core Service via the API Gateway
		data, err := fetchDataFromCoreService(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from Core Service"})
			return
		}

		// Process the data and generate a report
		report := generateReport(data)

		// Return the report
		c.JSON(http.StatusOK, gin.H{"message": "Report viewed successfully", "report": report})
	})
}

// fetchDataFromCoreService fetches data from the Core Service via the API Gateway
func fetchDataFromCoreService(c *gin.Context) (map[string]interface{}, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request to fetch data from the Core Service
	req, err := http.NewRequest("GET", config.APIGatewayURL+"/core/data", nil)
	if err != nil {
		return nil, err
	}

	// Copy headers from the incoming request to the new request
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Send the request to the Core Service via the API Gateway
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response body
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// generateReport processes the data and generates a report
func generateReport(data map[string]interface{}) map[string]interface{} {
	// Simulate report generation logic
	return map[string]interface{}{
		"summary": "Sample report summary",
		"data":    data,
	}
}
