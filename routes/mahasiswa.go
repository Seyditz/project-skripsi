package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func MahasiswaRoutes(r *gin.Engine) {
	mahasiswaRoutes := r.Group("/mahasiswa")
	mahasiswaRoutes.GET("/", controllers.GetAllMahasiswa)
	mahasiswaRoutes.GET("/:id", controllers.GetMahasiswaByID)
	mahasiswaRoutes.POST("/", controllers.CreateMahasiswa)
	mahasiswaRoutes.DELETE("/:id", controllers.DeleteMahasiswa)
	mahasiswaRoutes.PUT("/:id", controllers.UpdateMahasiswa)
	mahasiswaRoutes.POST("/login", controllers.MahasiswaLogin)
}
