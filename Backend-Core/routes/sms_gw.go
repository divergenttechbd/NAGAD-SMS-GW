package routes

import (
	"myproject/config"
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// SetupSMSGatewayRoutes sets up the SMS Gateway routes
func SetupSMSGatewayRoutes(r *gin.RouterGroup, influxClient influxdb2.Client, cfg *config.Config) {
	// Initialize the SMS Gateway Controller
	smsController := controllers.NewSMSGatewayController(influxClient, cfg)

	smsRoutes := r.Group("/sms")
	smsRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		// smsRoutes.POST("/send", middleware.RBAC("send_sms"), smsController.ProcessSMS)
		smsRoutes.POST("/send", smsController.ProcessSMS)
	}
}
