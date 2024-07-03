package controllers

import (
	"net/http"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
)

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
		UserId:          input.UserId,
		CreatedAt:       time.Now(),
	}

	// Create the mobileNotification in the database
	if result := database.DB.Create(&mobileNotification); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mobileNotification})
}

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

func GetAllNotification(c *gin.Context) {
	notifications := []models.MobileNotification{}

	db := database.DB

	claims := c.MustGet("claims").(*utils.Claims)

	// Use Preload to load the associated DataPengajuan
	result := db.
		Preload("DataPengajuan").
		Preload("DataPengajuan.Mahasiswa").
		Preload("DataPengajuan.DosPem1").
		Preload("DataPengajuan.DosPem2").
		Model(&[]models.MobileNotification{}).
		Where("user_id = ?", claims.UserId).
		Find(&notifications)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all notifications", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &notifications})
}

func GetRealAllNotification(c *gin.Context) {
	notifications := []models.MobileNotification{}

	db := database.DB

	// Use Preload to load the associated DataPengajuan
	result := db.Preload("DataPengajuan").Find(&notifications)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all notifications", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &notifications})
}
