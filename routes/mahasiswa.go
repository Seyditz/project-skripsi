package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

// func MahasiswaRoutes(r *gin.Engine) {
// 	mahasiswaRoutes := r.Group("/mahasiswa", middleware.AuthMiddleware())
// 	mahasiswaRoutes.GET("/", utils.Authorize("mahasiswa"), controllers.GetAllMahasiswa)
// 	mahasiswaRoutes.GET("/:id", utils.Authorize("mahasiswa"), controllers.GetMahasiswaByID)
// 	mahasiswaRoutes.POST("/", controllers.CreateMahasiswa)
// 	mahasiswaRoutes.DELETE("/:id", utils.Authorize("mahasiswa"), controllers.DeleteMahasiswa)
// 	mahasiswaRoutes.PUT("/:id", utils.Authorize("mahasiswa"), controllers.UpdateMahasiswa)
// }

func MahasiswaRoutes(r *gin.Engine) {
	mahasiswaRoutes := r.Group("/mahasiswa")
	mahasiswaRoutes.GET("/", controllers.GetAllMahasiswa)
	mahasiswaRoutes.GET("/:id", controllers.GetMahasiswaByID)
	mahasiswaRoutes.POST("/", controllers.CreateMahasiswa)
	mahasiswaRoutes.DELETE("/:id", controllers.DeleteMahasiswa)
	mahasiswaRoutes.PUT("/:id", controllers.UpdateMahasiswa)
}
