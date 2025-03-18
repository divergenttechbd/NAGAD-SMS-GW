package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMsgPriorities retrieves all SMS priority configurations
// @Summary Get all SMS priority configurations
// @Description Get all SMS priority configurations
// @Tags SMS Priority
// @Accept json
// @Produce json
// @Success 200 {array} models.MsgPriority
// @Failure 500 {object} map[string]interface{}
// @Router /msg-priorities [get]
func GetMsgPriorities(c *gin.Context) {
	db := utils.GetDB()
	var priorities []models.MsgPriority

	if err := db.Find(&priorities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SMS priority configurations"})
		return
	}

	c.JSON(http.StatusOK, priorities)
}

// CreateMsgPriority creates a new SMS priority configuration
// @Summary Create a new SMS priority configuration
// @Description Create a new SMS priority configuration with message type, priority level, and description
// @Tags SMS Priority
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "SMS priority details"
// @Success 201 {object} models.MsgPriority
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /msg-priority [post]
func CreateMsgPriority(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageType, ok := input["message_type"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message_type"})
		return
	}

	priorityLevel, ok := input["priority_level"].(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority_level"})
		return
	}

	description, ok := input["description"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid description"})
		return
	}

	priority := models.MsgPriority{
		Message_Type:   messageType,
		Priority_Level: priorityLevel,
		Description:    description,
	}

	db := utils.GetDB()

	if err := db.Create(&priority).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SMS priority configuration"})
		return
	}

	c.JSON(http.StatusCreated, priority)
}

// UpdateMsgPriority updates an existing SMS priority configuration
// @Summary Update an existing SMS priority configuration
// @Description Update an SMS priority configuration by ID with optional fields: message type, priority level, and description
// @Tags SMS Priority
// @Accept json
// @Produce json
// @Param id path string true "SMS Priority ID"
// @Param input body map[string]interface{} true "SMS priority details"
// @Success 200 {object} models.MsgPriority
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /msg-priority/{id} [put]
func UpdateMsgPriority(c *gin.Context) {
	var input map[string]interface{}

	priorityID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var priority models.MsgPriority

	if err := db.First(&priority, priorityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMS priority configuration not found"})
		return
	}

	if messageType, ok := input["message_type"].(string); ok {
		priority.Message_Type = messageType
	}
	if priorityLevel, ok := input["priority_level"].(int); ok {
		priority.Priority_Level = priorityLevel
	}
	if description, ok := input["description"].(string); ok {
		priority.Description = description
	}

	if err := db.Save(&priority).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SMS priority configuration"})
		return
	}

	c.JSON(http.StatusOK, priority)
}

// DeleteMsgPriority deletes an existing SMS priority configuration
// @Summary Delete an existing SMS priority configuration
// @Description Delete an SMS priority configuration by ID
// @Tags SMS Priority
// @Accept json
// @Produce json
// @Param id path string true "SMS Priority ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /msg-priority/{id} [delete]
func DeleteMsgPriority(c *gin.Context) {
	priorityID := c.Param("id")

	db := utils.GetDB()
	var priority models.MsgPriority

	if err := db.First(&priority, priorityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMS priority configuration not found"})
		return
	}

	if err := db.Delete(&priority).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS priority configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMS priority configuration deleted successfully"})
}

// GetMsgPriorityDetails retrieves details of an SMS priority configuration by ID
// @Summary Get SMS priority details
// @Description Get details of an SMS priority configuration by ID
// @Tags SMS Priority
// @Accept json
// @Produce json
// @Param id path string true "SMS Priority ID"
// @Success 200 {object} models.MsgPriority
// @Failure 404 {object} map[string]interface{}
// @Router /msg-priorities/{id} [get]
func GetMsgPriorityDetails(c *gin.Context) {
	priorityID := c.Param("id")

	db := utils.GetDB()
	var priority models.MsgPriority

	if err := db.First(&priority, priorityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SMS priority configuration not found"})
		return
	}

	c.JSON(http.StatusOK, priority)
}
