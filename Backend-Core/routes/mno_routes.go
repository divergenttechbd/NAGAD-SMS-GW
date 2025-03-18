package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupMNORoutes(r *gin.RouterGroup) {
	Routes := r.Group("/mno")
	// userRoutes.Use(middleware.RBAC("manage_users")) // Ensure RBAC middleware is applied
	Routes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		Routes.GET("/", middleware.RBAC("view_mno"), controllers.GetMNOs)
		Routes.POST("/", middleware.RBAC("create_mno"), controllers.CreateMNO)
		Routes.PUT("/:id", middleware.RBAC("edit_mno"), controllers.UpdateMNO)
		Routes.DELETE("/:id", middleware.RBAC("delete_mno"), controllers.DeleteMNO)
		Routes.GET("/:id", middleware.RBAC("get_mno_details"), controllers.GetMNODetails)
	}
}
