package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	adminRoutes := r.Group("/admin")
	adminRoutes.GET("/", controllers.GetAllAdmin)
	adminRoutes.POST("/", controllers.CreateAdmin)
}
