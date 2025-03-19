package routes

// import (
// 	"myproject/controllers"
// 	"myproject/middleware"

// 	"github.com/gin-gonic/gin"
// )

// // SetupSMSGatewayRoutes defines routes for SMS Gateway
// func SetupSMSGatewayRoutes(r *gin.RouterGroup) {
// 	smsRoutes := r.Group("/sms")
// 	smsRoutes.Use(middleware.JWTAuth()) // Apply authentication middleware
// 	{
// 		// Apply RBAC middleware to secure SMS Gateway routes
// 		smsRoutes.POST("/send", middleware.RBAC("send_sms"), controllers.ProcessSMS)
// 	}
// }
