package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMsgPriorityRoutes(r *gin.RouterGroup) {
	msgPriorityRoutes := r.Group("/msg_priority")
	msgPriorityRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		msgPriorityRoutes.GET("/", middleware.RBAC("view_msg_priority"), controllers.GetMsgPriorities)
		msgPriorityRoutes.POST("/", middleware.RBAC("create_msg_priority"), controllers.CreateMsgPriority)
		msgPriorityRoutes.PUT("/:id", middleware.RBAC("edit_msg_priority"), controllers.UpdateMsgPriority)
		msgPriorityRoutes.DELETE("/:id", middleware.RBAC("delete_msg_priority"), controllers.DeleteMsgPriority)
		msgPriorityRoutes.GET("/:id", middleware.RBAC("get_msg_priority_details"), controllers.GetMsgPriorityDetails)
	}
}
