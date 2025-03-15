package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCampaigns retrieves all campaigns
func GetCampaigns(c *gin.Context) {
	db := utils.GetDB()
	var campaigns []models.Campaign

	if err := db.Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaigns"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

// CreateCampaign creates a new campaign
func CreateCampaign(c *gin.Context) {
	var input struct {
		Name   string `json:"name" binding:"required"`
		Status string `json:"status" binding:"required"`
		UserID uint   `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign := models.Campaign{
		Name:   input.Name,
		Status: input.Status,
		UserID: input.UserID,
	}

	db := utils.GetDB()
	if err := db.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create campaign"})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// UpdateCampaign updates an existing campaign
func UpdateCampaign(c *gin.Context) {
	var input struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	campaignID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var campaign models.Campaign
	if err := db.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	// Update fields if provided
	if input.Name != "" {
		campaign.Name = input.Name
	}
	if input.Status != "" {
		campaign.Status = input.Status
	}

	if err := db.Save(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update campaign"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}

// DeleteCampaign deletes a campaign
func DeleteCampaign(c *gin.Context) {
	campaignID := c.Param("id")

	db := utils.GetDB()
	var campaign models.Campaign
	if err := db.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	if err := db.Delete(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign deleted successfully"})
}
