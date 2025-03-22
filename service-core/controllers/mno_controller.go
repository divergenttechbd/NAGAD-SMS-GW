package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMNOs retrieves all MNOs
// @Summary Get all MNOs
// @Description Get all MNOs with their channels
// @Tags MNOs
// @Accept json
// @Produce json
// @Success 200 {array} models.MNO
// @Failure 500 {object} map[string]interface{}
// @Router /mnos [get]
func GetMNOs(c *gin.Context) {
	db := utils.GetDB()
	var mnos []models.MNO

	if err := db.Preload("Channels").Find(&mnos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch MNOs"})
		return
	}

	c.JSON(http.StatusOK, mnos)
}

// CreateMNO creates a new Mobile Network Operator
// @Summary Create a new Mobile Network Operator
// @Description Create a new MNO with name, prefix, and status
// @Tags MNO
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "MNO details"
// @Success 201 {object} models.MNO
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno [post]
func CreateMNO(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mnoName, ok := input["mno_name"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mno_name"})
		return
	}

	prefix, ok := input["prefix"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prefix"})
		return
	}

	status, ok := input["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	mno := models.MNO{
		MNO_Name: mnoName,
		Prefix:   prefix,
		Status:   status,
	}

	db := utils.GetDB()

	if err := db.Create(&mno).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create MNO"})
		return
	}

	c.JSON(http.StatusCreated, mno)
}

// UpdateMNO updates an existing Mobile Network Operator
// @Summary Update an existing Mobile Network Operator
// @Description Update an MNO by ID with optional fields: name, prefix, and status
// @Tags MNO
// @Accept json
// @Produce json
// @Param id path string true "MNO ID"
// @Param input body map[string]interface{} true "MNO details"
// @Success 200 {object} models.MNO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno/{id} [put]
func UpdateMNO(c *gin.Context) {
	var input map[string]interface{}

	mnoID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var mno models.MNO

	if err := db.First(&mno, mnoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MNO not found"})
		return
	}

	if mnoName, ok := input["mno_name"].(string); ok {
		mno.MNO_Name = mnoName
	}
	if prefix, ok := input["prefix"].(string); ok {
		mno.Prefix = prefix
	}
	if status, ok := input["status"].(string); ok {
		mno.Status = status
	}

	if err := db.Save(&mno).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MNO"})
		return
	}

	c.JSON(http.StatusOK, mno)
}

// DeleteMNO deletes an existing Mobile Network Operator
// @Summary Delete an existing Mobile Network Operator
// @Description Delete an MNO by ID
// @Tags MNO
// @Accept json
// @Produce json
// @Param id path string true "MNO ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno/{id} [delete]
func DeleteMNO(c *gin.Context) {
	mnoID := c.Param("id")

	db := utils.GetDB()
	var mno models.MNO

	if err := db.First(&mno, mnoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MNO not found"})
		return
	}

	if err := db.Delete(&mno).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete MNO"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MNO deleted successfully"})
}

// GetMNODetails retrieves details of an MNO by ID (including channels)
// @Summary Get MNO details
// @Description Get details of an MNO by ID including its channels
// @Tags MNOs
// @Accept json
// @Produce json
// @Param id path int true "MNO ID"
// @Success 200 {object} models.MNO
// @Failure 404 {object} map[string]interface{}
// @Router /mnos/{id} [get]
func GetMNODetails(c *gin.Context) {
	mnoID := c.Param("id")

	db := utils.GetDB()
	var mno models.MNO

	if err := db.Preload("Channels").First(&mno, mnoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MNO not found"})
		return
	}

	c.JSON(http.StatusOK, mno)
}

// CreateMNOChannel creates a new channel for an MNO
// @Summary Create a new MNO channel
// @Description Create a new channel for an MNO with the provided details
// @Tags MNO Channels
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "MNO channel details"
// @Success 201 {object} models.MnoChannels
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno-channels [post]
func CreateMNOChannel(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mnoID, ok := input["mno_id"].(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mno_id"})
		return
	}

	channelType, ok := input["channel_type"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel_type"})
		return
	}

	priority, ok := input["priority"].(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid priority"})
		return
	}

	tps, ok := input["tps"].(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tps"})
		return
	}

	status, ok := input["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	channel := models.MnoChannels{
		MNOID:       mnoID,
		ChannelType: channelType,
		Priority:    priority,
		TPS:         tps,
		Status:      status,
	}

	db := utils.GetDB()

	if err := db.Create(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create MNO channel"})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

// UpdateMNOChannel updates an existing MNO channel
// @Summary Update an existing MNO channel
// @Description Update an existing MNO channel with the provided details
// @Tags MNO Channels
// @Accept json
// @Produce json
// @Param id path int true "MNO Channel ID"
// @Param input body map[string]interface{} true "MNO channel details"
// @Success 200 {object} models.MnoChannels
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno-channels/{id} [put]
func UpdateMNOChannel(c *gin.Context) {
	var input map[string]interface{}

	channelID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var channel models.MnoChannels

	if err := db.First(&channel, channelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MNO channel not found"})
		return
	}

	if channelType, ok := input["channel_type"].(string); ok {
		channel.ChannelType = channelType
	}
	if priority, ok := input["priority"].(int); ok {
		channel.Priority = priority
	}
	if tps, ok := input["tps"].(int); ok {
		channel.TPS = tps
	}
	if status, ok := input["status"].(string); ok {
		channel.Status = status
	}

	if err := db.Save(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update MNO channel"})
		return
	}

	c.JSON(http.StatusOK, channel)
}

// DeleteMNOChannel deletes an MNO channel
// @Summary Delete an MNO channel
// @Description Delete an MNO channel by ID
// @Tags MNO Channels
// @Accept json
// @Produce json
// @Param id path int true "MNO Channel ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mno-channels/{id} [delete]
func DeleteMNOChannel(c *gin.Context) {
	channelID := c.Param("id")

	db := utils.GetDB()
	var channel models.MnoChannels

	if err := db.First(&channel, channelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MNO channel not found"})
		return
	}

	if err := db.Delete(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete MNO channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MNO channel deleted successfully"})
}
