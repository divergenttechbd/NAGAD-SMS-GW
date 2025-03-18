package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCampaignWorkflows retrieves all campaign workflows
// @Summary Get all campaign workflows
// @Description Get all campaign workflows
// @Tags Campaign Workflows
// @Accept json
// @Produce json
// @Success 200 {array} models.CampaignWorkflow
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflows [get]
func GetCampaignWorkflows(c *gin.Context) {
	db := utils.GetDB()
	var workflows []models.CampaignWorkflow

	if err := db.Preload("CampaignWorkflowUsers").Preload("CampaignWorkflowProcessings").Find(&workflows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign workflows"})
		return
	}

	c.JSON(http.StatusOK, workflows)
}

// CreateCampaignWorkflow creates a new campaign workflow
// @Summary Create a new campaign workflow
// @Description Create a new campaign workflow with name
// @Tags Campaign Workflows
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Campaign workflow details"
// @Success 201 {object} models.CampaignWorkflow
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow [post]
func CreateCampaignWorkflow(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name, ok := input["name"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid name"})
		return
	}

	workflow := models.CampaignWorkflow{
		Name: name,
	}

	db := utils.GetDB()

	if err := db.Create(&workflow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign workflow"})
		return
	}

	c.JSON(http.StatusCreated, workflow)
}

// UpdateCampaignWorkflow updates an existing campaign workflow
// @Summary Update an existing campaign workflow
// @Description Update a campaign workflow by ID with optional fields: name
// @Tags Campaign Workflows
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow ID"
// @Param input body map[string]interface{} true "Campaign workflow details"
// @Success 200 {object} models.CampaignWorkflow
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow/{id} [put]
func UpdateCampaignWorkflow(c *gin.Context) {
	var input map[string]interface{}

	workflowID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var workflow models.CampaignWorkflow

	if err := db.First(&workflow, workflowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow not found"})
		return
	}

	if name, ok := input["name"].(string); ok {
		workflow.Name = name
	}

	if err := db.Save(&workflow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign workflow"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// DeleteCampaignWorkflow deletes an existing campaign workflow
// @Summary Delete an existing campaign workflow
// @Description Delete a campaign workflow by ID
// @Tags Campaign Workflows
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow/{id} [delete]
func DeleteCampaignWorkflow(c *gin.Context) {
	workflowID := c.Param("id")

	db := utils.GetDB()
	var workflow models.CampaignWorkflow

	if err := db.First(&workflow, workflowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow not found"})
		return
	}

	if err := db.Delete(&workflow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign workflow deleted successfully"})
}

// GetCampaignWorkflowDetails retrieves details of a campaign workflow by ID
// @Summary Get campaign workflow details
// @Description Get details of a campaign workflow by ID
// @Tags Campaign Workflows
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow ID"
// @Success 200 {object} models.CampaignWorkflow
// @Failure 404 {object} map[string]interface{}
// @Router /campaign-workflows/{id} [get]
func GetCampaignWorkflowDetails(c *gin.Context) {
	workflowID := c.Param("id")

	db := utils.GetDB()
	var workflow models.CampaignWorkflow

	if err := db.Preload("CampaignWorkflowUsers").Preload("CampaignWorkflowProcessings").First(&workflow, workflowID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow not found"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

// CreateCampaignWorkflowUser creates a new user associated with a campaign workflow
// @Summary Create a new user associated with a campaign workflow
// @Description Create a new user associated with a campaign workflow
// @Tags Campaign Workflow Users
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Campaign workflow user details"
// @Success 201 {object} models.CampaignWorkflowUser
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-user [post]
func CreateCampaignWorkflowUser(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	workflowID, ok := input["workflow_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow_id"})
		return
	}

	userID, ok := input["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	status, ok := input["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	workflowUser := models.CampaignWorkflowUser{
		WorkflowID: uint(workflowID),
		UserID:     uint(userID),
		Status:     status,
	}

	db := utils.GetDB()

	if err := db.Create(&workflowUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign workflow user"})
		return
	}

	c.JSON(http.StatusCreated, workflowUser)
}

// UpdateCampaignWorkflowUser updates an existing user associated with a campaign workflow
// @Summary Update an existing user associated with a campaign workflow
// @Description Update a user associated with a campaign workflow by ID
// @Tags Campaign Workflow Users
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow User ID"
// @Param input body map[string]interface{} true "Campaign workflow user details"
// @Success 200 {object} models.CampaignWorkflowUser
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-user/{id} [put]
func UpdateCampaignWorkflowUser(c *gin.Context) {
	var input map[string]interface{}

	workflowUserID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var workflowUser models.CampaignWorkflowUser

	if err := db.First(&workflowUser, workflowUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow user not found"})
		return
	}

	if status, ok := input["status"].(string); ok {
		workflowUser.Status = status
	}

	if err := db.Save(&workflowUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign workflow user"})
		return
	}

	c.JSON(http.StatusOK, workflowUser)
}

// DeleteCampaignWorkflowUser deletes an existing user associated with a campaign workflow
// @Summary Delete an existing user associated with a campaign workflow
// @Description Delete a user associated with a campaign workflow by ID
// @Tags Campaign Workflow Users
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-user/{id} [delete]
func DeleteCampaignWorkflowUser(c *gin.Context) {
	workflowUserID := c.Param("id")

	db := utils.GetDB()
	var workflowUser models.CampaignWorkflowUser

	if err := db.First(&workflowUser, workflowUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow user not found"})
		return
	}

	if err := db.Delete(&workflowUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign workflow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign workflow user deleted successfully"})
}

// GetCampaignWorkflowUserDetails retrieves details of a user associated with a campaign workflow by ID
// @Summary Get campaign workflow user details
// @Description Get details of a user associated with a campaign workflow by ID
// @Tags Campaign Workflow Users
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow User ID"
// @Success 200 {object} models.CampaignWorkflowUser
// @Failure 404 {object} map[string]interface{}
// @Router /campaign-workflow-users/{id} [get]
func GetCampaignWorkflowUserDetails(c *gin.Context) {
	workflowUserID := c.Param("id")

	db := utils.GetDB()
	var workflowUser models.CampaignWorkflowUser

	if err := db.Preload("CampaignWorkflow").Preload("User").First(&workflowUser, workflowUserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow user not found"})
		return
	}

	c.JSON(http.StatusOK, workflowUser)
}

// CreateCampaignWorkflowProcessing creates a new processing association between a campaign and a workflow
// @Summary Create a new processing association between a campaign and a workflow
// @Description Create a new processing association between a campaign and a workflow
// @Tags Campaign Workflow Processings
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Campaign workflow processing details"
// @Success 201 {object} models.CampaignWorkflowProcessing
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-processing [post]
func CreateCampaignWorkflowProcessing(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaignID, ok := input["campaign_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign_id"})
		return
	}

	workflowID, ok := input["workflow_id"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow_id"})
		return
	}

	processing := models.CampaignWorkflowProcessing{
		CampaignID: uint(campaignID),
		WorkflowID: uint(workflowID),
	}

	db := utils.GetDB()

	if err := db.Create(&processing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign workflow processing"})
		return
	}

	c.JSON(http.StatusCreated, processing)
}

// UpdateCampaignWorkflowProcessing updates an existing processing association between a campaign and a workflow
// @Summary Update an existing processing association between a campaign and a workflow
// @Description Update a processing association between a campaign and a workflow by ID
// @Tags Campaign Workflow Processings
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow Processing ID"
// @Param input body map[string]interface{} true "Campaign workflow processing details"
// @Success 200 {object} models.CampaignWorkflowProcessing
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-processing/{id} [put]
func UpdateCampaignWorkflowProcessing(c *gin.Context) {
	var input map[string]interface{}

	processingID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var processing models.CampaignWorkflowProcessing

	if err := db.First(&processing, processingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow processing not found"})
		return
	}

	if campaignID, ok := input["campaign_id"].(float64); ok {
		processing.CampaignID = uint(campaignID)
	}
	if workflowID, ok := input["workflow_id"].(float64); ok {
		processing.WorkflowID = uint(workflowID)
	}

	if err := db.Save(&processing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign workflow processing"})
		return
	}

	c.JSON(http.StatusOK, processing)
}

// DeleteCampaignWorkflowProcessing deletes an existing processing association between a campaign and a workflow
// @Summary Delete an existing processing association between a campaign and a workflow
// @Description Delete a processing association between a campaign and a workflow by ID
// @Tags Campaign Workflow Processings
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow Processing ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-workflow-processing/{id} [delete]
func DeleteCampaignWorkflowProcessing(c *gin.Context) {
	processingID := c.Param("id")

	db := utils.GetDB()
	var processing models.CampaignWorkflowProcessing

	if err := db.First(&processing, processingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow processing not found"})
		return
	}

	if err := db.Delete(&processing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign workflow processing"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign workflow processing deleted successfully"})
}

// GetCampaignWorkflowProcessingDetails retrieves details of a processing association between a campaign and a workflow by ID
// @Summary Get campaign workflow processing details
// @Description Get details of a processing association between a campaign and a workflow by ID
// @Tags Campaign Workflow Processings
// @Accept json
// @Produce json
// @Param id path string true "Campaign Workflow Processing ID"
// @Success 200 {object} models.CampaignWorkflowProcessing
// @Failure 404 {object} map[string]interface{}
// @Router /campaign-workflow-processings/{id} [get]
func GetCampaignWorkflowProcessingDetails(c *gin.Context) {
	processingID := c.Param("id")

	db := utils.GetDB()
	var processing models.CampaignWorkflowProcessing

	if err := db.Preload("Campaign").Preload("CampaignWorkflow").First(&processing, processingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign workflow processing not found"})
		return
	}

	c.JSON(http.StatusOK, processing)
}
