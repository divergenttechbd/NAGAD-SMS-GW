package routes

import (
	"myproject/controllers"
	"myproject/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCampaignWorkflowRoutes(r *gin.RouterGroup) {
	campaignWorkflowRoutes := r.Group("/campaign_workflow")
	campaignWorkflowRoutes.Use(middleware.JWTAuth()) // Ensure authentication middleware is applied
	{
		// Apply RBAC middleware to each route with the required permission
		campaignWorkflowRoutes.GET("/", middleware.RBAC("view_campaign_workflow"), controllers.GetCampaignWorkflows)
		campaignWorkflowRoutes.POST("/", middleware.RBAC("create_campaign_workflow"), controllers.CreateCampaignWorkflow)
		campaignWorkflowRoutes.PUT("/:id", middleware.RBAC("edit_campaign_workflow"), controllers.UpdateCampaignWorkflow)
		campaignWorkflowRoutes.DELETE("/:id", middleware.RBAC("delete_campaign_workflow"), controllers.DeleteCampaignWorkflow)
		campaignWorkflowRoutes.GET("/:id", middleware.RBAC("get_campaign_workflow_details"), controllers.GetCampaignWorkflowDetails)
	}
}
