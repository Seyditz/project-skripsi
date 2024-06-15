package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	middleware "github.com/Seyditz/project-skripsi/middlewares"
	"github.com/Seyditz/project-skripsi/utils"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine) {
	notificationRoutes := r.Group("/notification", middleware.AuthMiddleware())
	notificationRoutes.GET("/:id", utils.Authorize("mahasiswa"), controllers.GetNotificationbyId)
	notificationRoutes.GET("/", utils.Authorize("mahasiswa"), controllers.GetAllNotification)
	notificationRoutes.POST("/", utils.Authorize("mahasiswa"), controllers.CreateNotification)
}
