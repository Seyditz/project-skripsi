package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine) {
	adminRoutes := r.Group("/admin")
	adminRoutes.GET("/", controllers.GetAllAdmin)
	adminRoutes.GET("/:id", controllers.GetAdminbyId)
	adminRoutes.POST("/", controllers.CreateAdmin)
	adminRoutes.DELETE("/:id", controllers.DeleteAdmin)
	adminRoutes.PUT("/:id", controllers.UpdateAdmin)
	adminRoutes.POST("/login", controllers.AdminLogin)
}
