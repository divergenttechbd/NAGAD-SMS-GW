package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCampaignRoutes(r *gin.RouterGroup) {
	campaignRoutes := r.Group("/campaigns")
	campaignRoutes.Use(middleware.RBAC("manage_campaigns")) // Ensure RBAC middleware is applied
	{
		campaignRoutes.GET("/", controllers.GetCampaigns)
		campaignRoutes.POST("/", controllers.CreateCampaign)
		campaignRoutes.PUT("/:id", controllers.UpdateCampaign)
		campaignRoutes.DELETE("/:id", controllers.DeleteCampaign)
	}
}
