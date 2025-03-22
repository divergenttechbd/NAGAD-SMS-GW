package routes

import (
	"myproject/config"
	"myproject/controllers"
	"myproject/middleware"
	"myproject/rabbitmq"

	"github.com/gin-gonic/gin"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// SetupSMSGatewayRoutes sets up the SMS Gateway routes
func SetupSMSGatewayRoutes(r *gin.RouterGroup, influxClient influxdb2.Client, cfg *config.Config, rmq *rabbitmq.RabbitMQ) {
	// Initialize the SMS Gateway Controller
	// smsController := controllers.NewSMSGatewayController(influxClient, cfg)
	smsController := controllers.NewSMSGatewayController(influxClient, cfg, rmq)

	smsRoutes := r.Group("/sms")
	smsRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		// smsRoutes.POST("/send", middleware.RBAC("send_sms"), smsController.ProcessSMS)
		smsRoutes.POST("/send", smsController.ProcessSMS)
		smsRoutes.GET("/test-million-msg", smsController.PublishMillionMessages)
		smsRoutes.GET("/rabbitmq-stats", smsController.GetRabbitMQStatistics)
	}
}
