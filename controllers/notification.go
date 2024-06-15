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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Notification
// @Description Get All Notifications
// @Produce application/json
// @Tags Notification
// @Success 200 {object} []models.MobileNotification{}
// @Router /notification [get]
func GetAllNotification(c *gin.Context) {
	notifications := []models.MobileNotification{}

	db := database.DB

	name := c.Query("name")
	email := c.Query("email")

	// Build the query conditionally based on the parameters
	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if email != "" {
		db = db.Where("email ILIKE ?", "%"+email+"%")
	}

	// result := db.Find(&notifications)
	result := db.Model(&[]models.MobileNotification{}).Find(&notifications)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all notifications", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &notifications})
}
