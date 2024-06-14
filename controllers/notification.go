package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
)

// CreateTags godoc
// @Summary Create Notification
// @Description Create Notifications
// @Produce application/json
// @Param request body models.MobileNotificationCreateRequest true "Raw Request Body"
// @Tags Notification
// @Success 200 {object} models.MobileNotification{}
// @Router /notification [post]
func CreateNotification(c *gin.Context) {
	var input models.MobileNotificationCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is required"})
		return
	}
	if input.DataPengajuanId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data_pengajuan_id is required"})
		return
	}

	mobileNotification := models.MobileNotification{
		Message:         input.Message,
		DataPengajuanId: input.DataPengajuanId,
	}

	// Create the mobileNotification in the database
	if result := database.DB.Create(&mobileNotification); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mobileNotification})
}

// CreateTags godoc
// @Summary Get Notification By ID
// @Description Get Notification By ID
// @Produce application/json
// @Param id path int true "Notification ID"
// @Tags Notification
// @Success 200 {object} models.MobileNotification{}
// @Router /notification/{id} [get]
func GetNotificationbyId(c *gin.Context) {
	// Get the notification ID from the URL parameters
	notificationID := c.Param("id")

	// Find the notification by ID
	var notification models.MobileNotification
	if result := database.DB.Model(&models.MobileNotification{}).First(&notification, notificationID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(200, gin.H{"notification": &notification})
}
