package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine) {
	notificationRoutes := r.Group("/notification")
	notificationRoutes.GET("/:id", controllers.GetNotificationbyId)
	notificationRoutes.POST("/", controllers.CreateNotification)
}
