package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("/users")
	// userRoutes.Use(middleware.RBAC("manage_users")) // Ensure RBAC middleware is applied
	userRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		userRoutes.GET("/", middleware.RBAC("view_user"), controllers.GetUsers)
		userRoutes.POST("/", middleware.RBAC("create_user"), controllers.CreateUser)
		userRoutes.PUT("/:id", middleware.RBAC("edit_use"), controllers.UpdateUser)
		userRoutes.DELETE("/:id", middleware.RBAC("delete_user"), controllers.DeleteUser)
	}
}
