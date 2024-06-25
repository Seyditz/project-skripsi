package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"

	"github.com/gin-gonic/gin"
)

// func AdminRoutes(r *gin.Engine) {
// 	adminRoutes := r.Group("/admin", middleware.AuthMiddleware())
// 	{
// 		adminRoutes.GET("/", utils.Authorize("admin"), controllers.GetAllAdmin)
// 		adminRoutes.GET("/:id", utils.Authorize("admin"), controllers.GetAdminbyId)
// 		adminRoutes.POST("/", utils.Authorize("admin"), controllers.CreateAdmin)
// 		adminRoutes.DELETE("/:id", utils.Authorize("admin"), controllers.DeleteAdmin)
// 		adminRoutes.PUT("/:id", utils.Authorize("admin"), controllers.UpdateAdmin)
// 	}
// }

func AdminRoutes(r *gin.Engine) {
	adminRoutes := r.Group("/admin")
	{
		adminRoutes.GET("/", controllers.GetAllAdmin)
		adminRoutes.GET("/:id", controllers.GetAdminbyId)
		adminRoutes.POST("/", controllers.CreateAdmin)
		adminRoutes.DELETE("/:id", controllers.DeleteAdmin)
		adminRoutes.PUT("/:id", controllers.UpdateAdmin)
	}
}
