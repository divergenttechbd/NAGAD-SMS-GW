package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDndRoutes(r *gin.RouterGroup) {
	dndRoutes := r.Group("/dnd")
	dndRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		dndRoutes.GET("/", middleware.RBAC("view_dnd"), controllers.GetDNDs)
		dndRoutes.POST("/", middleware.RBAC("create_dnd"), controllers.CreateDND)
		dndRoutes.PUT("/:id", middleware.RBAC("edit_dnd"), controllers.UpdateDND)
		dndRoutes.DELETE("/:id", middleware.RBAC("delete_dnd"), controllers.DeleteDND)
		dndRoutes.GET("/:id", middleware.RBAC("get_dnd_details"), controllers.GetDNDDetails)
	}
}
