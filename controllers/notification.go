package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
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
	}

	// Create the mobileNotification in the database
	if result := database.DB.Create(&mobileNotification); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mobileNotification})
}
