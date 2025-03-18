package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCampaignRecipientRoutes(r *gin.RouterGroup) {
	campaignRecipientRoutes := r.Group("/campaign_recipient")
	campaignRecipientRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		campaignRecipientRoutes.GET("/", middleware.RBAC("view_campaign_recipient"), controllers.GetCampaignRecipients)
		campaignRecipientRoutes.POST("/", middleware.RBAC("create_campaign_recipient"), controllers.CreateCampaignRecipient)
		campaignRecipientRoutes.PUT("/:id", middleware.RBAC("edit_campaign_recipient"), controllers.UpdateCampaignRecipient)
		campaignRecipientRoutes.DELETE("/:id", middleware.RBAC("delete_campaign_recipient"), controllers.DeleteCampaignRecipient)
		campaignRecipientRoutes.GET("/:id", middleware.RBAC("get_campaign_recipient_details"), controllers.GetCampaignRecipientDetails)
	}
}
