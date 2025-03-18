package controllers

import (
	"myproject/models"
	"myproject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDNDs retrieves all DND entries
// @Summary Get all DND entries
// @Description Get all DND entries
// @Tags DND
// @Accept json
// @Produce json
// @Success 200 {array} models.DND
// @Failure 500 {object} map[string]interface{}
// @Router /api/dnd [get]
func GetDNDs(c *gin.Context) {
	db := utils.GetDB()
	var dnds []models.DND

	if err := db.Find(&dnds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DND entries"})
		return
	}

	c.JSON(http.StatusOK, dnds)
}

// CreateDND creates a new DND entry
// @Summary Create a new DND entry
// @Description Create a new DND entry with phone number, reason, and status
// @Tags DND
// @Accept json
// @Produce json
// @Param input body map[string]interface{} true "DND details"
// @Success 201 {object} models.DND
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/dnd [post]
func CreateDND(c *gin.Context) {
	var input map[string]interface{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phoneNumber, ok := input["phone_number"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone_number"})
		return
	}

	reason, ok := input["reason"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reason"})
		return
	}

	status, ok := input["status"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	dnd := models.DND{
		Phone_Number: phoneNumber,
		Reason:       reason,
		Status:       status,
	}

	db := utils.GetDB()

	if err := db.Create(&dnd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DND entry"})
		return
	}

	c.JSON(http.StatusCreated, dnd)
}

// UpdateDND updates an existing DND entry
// @Summary Update an existing DND entry
// @Description Update a DND entry by ID with optional fields: phone number, reason, and status
// @Tags DND
// @Accept json
// @Produce json
// @Param id path string true "DND ID"
// @Param input body map[string]interface{} true "DND details"
// @Success 200 {object} models.DND
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/dnd/{id} [put]
func UpdateDND(c *gin.Context) {
	var input map[string]interface{}

	dndID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := utils.GetDB()
	var dnd models.DND

	if err := db.First(&dnd, dndID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DND entry not found"})
		return
	}

	if phoneNumber, ok := input["phone_number"].(string); ok {
		dnd.Phone_Number = phoneNumber
	}
	if reason, ok := input["reason"].(string); ok {
		dnd.Reason = reason
	}
	if status, ok := input["status"].(string); ok {
		dnd.Status = status
	}

	if err := db.Save(&dnd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DND entry"})
		return
	}

	c.JSON(http.StatusOK, dnd)
}

// DeleteDND deletes an existing DND entry
// @Summary Delete an existing DND entry
// @Description Delete a DND entry by ID
// @Tags DND
// @Accept json
// @Produce json
// @Param id path string true "DND ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/dnd/{id} [delete]
func DeleteDND(c *gin.Context) {
	dndID := c.Param("id")

	db := utils.GetDB()
	var dnd models.DND

	if err := db.First(&dnd, dndID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DND entry not found"})
		return
	}

	if err := db.Delete(&dnd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DND entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DND entry deleted successfully"})
}

// GetDNDDetails retrieves details of a DND entry by ID
// @Summary Get DND details
// @Description Get details of a DND entry by ID
// @Tags DND
// @Accept json
// @Produce json
// @Param id path string true "DND ID"
// @Success 200 {object} models.DND
// @Failure 404 {object} map[string]interface{}
// @Router /api/dnd/{id} [get]
func GetDNDDetails(c *gin.Context) {
	dndID := c.Param("id")

	db := utils.GetDB()
	var dnd models.DND

	if err := db.First(&dnd, dndID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DND entry not found"})
		return
	}

	c.JSON(http.StatusOK, dnd)
}
