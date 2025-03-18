package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCampaignRecipients retrieves all campaign recipients
// @Summary Get all campaign recipients
// @Description Get all campaign recipients
// @Tags Campaign Recipients
// @Accept json
// @Produce json
// @Success 200 {array} models.CampaignRecipient
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-recipients [get]
func GetCampaignRecipients(c *gin.Context) {
	db := utils.GetDB()
	var recipients []models.CampaignRecipient

	if err := db.Preload("Campaign").Find(&recipients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign recipients"})
		return
	}

	c.JSON(http.StatusOK, recipients)
}

// CreateCampaignRecipient creates a new campaign recipient
// @Summary Create a new campaign recipient
// @Description Create a new campaign recipient with campaign ID, recipient ID, and status
// @Tags Campaign Recipients
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "Campaign recipient details"
// @Success 201 {object} models.CampaignRecipient
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-recipient [post]
func CreateCampaignRecipient(c *gin.Context) {
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

	recipientID, ok := input["recipient"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipient"})
		return
	}

	status, ok := input["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	recipient := models.CampaignRecipient{
		CampaignID: uint(campaignID),
		Recipient:  uint(recipientID),
		Status:     status,
	}

	db := utils.GetDB()

	if err := db.Create(&recipient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign recipient"})
		return
	}

	c.JSON(http.StatusCreated, recipient)
}

// UpdateCampaignRecipient updates an existing campaign recipient
// @Summary Update an existing campaign recipient
// @Description Update a campaign recipient by ID with optional fields: campaign ID, recipient ID, and status
// @Tags Campaign Recipients
// @Accept json
// @Produce json
// @Param id path string true "Campaign Recipient ID"
// @Param input body map[string]interface{} true "Campaign recipient details"
// @Success 200 {object} models.CampaignRecipient
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-recipient/{id} [put]
func UpdateCampaignRecipient(c *gin.Context) {
	var input map[string]interface{}

	recipientID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var recipient models.CampaignRecipient

	if err := db.First(&recipient, recipientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign recipient not found"})
		return
	}

	if campaignID, ok := input["campaign_id"].(float64); ok {
		recipient.CampaignID = uint(campaignID)
	}
	if recipientIDFloat, ok := input["recipient"].(float64); ok {
		recipient.Recipient = uint(recipientIDFloat)
	}
	if status, ok := input["status"].(string); ok {
		recipient.Status = status
	}

	if err := db.Save(&recipient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign recipient"})
		return
	}

	c.JSON(http.StatusOK, recipient)
}

// DeleteCampaignRecipient deletes an existing campaign recipient
// @Summary Delete an existing campaign recipient
// @Description Delete a campaign recipient by ID
// @Tags Campaign Recipients
// @Accept json
// @Produce json
// @Param id path string true "Campaign Recipient ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /campaign-recipient/{id} [delete]
func DeleteCampaignRecipient(c *gin.Context) {
	recipientID := c.Param("id")

	db := utils.GetDB()
	var recipient models.CampaignRecipient

	if err := db.First(&recipient, recipientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign recipient not found"})
		return
	}

	if err := db.Delete(&recipient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign recipient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign recipient deleted successfully"})
}

// GetCampaignRecipientDetails retrieves details of a campaign recipient by ID
// @Summary Get campaign recipient details
// @Description Get details of a campaign recipient by ID
// @Tags Campaign Recipients
// @Accept json
// @Produce json
// @Param id path string true "Campaign Recipient ID"
// @Success 200 {object} models.CampaignRecipient
// @Failure 404 {object} map[string]interface{}
// @Router /campaign-recipients/{id} [get]
func GetCampaignRecipientDetails(c *gin.Context) {
	recipientID := c.Param("id")

	db := utils.GetDB()
	var recipient models.CampaignRecipient

	if err := db.Preload("Campaign").First(&recipient, recipientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign recipient not found"})
		return
	}

	c.JSON(http.StatusOK, recipient)
}
