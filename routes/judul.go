package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func JudulRoutes(r *gin.Engine) {
	judulRoutes := r.Group("/judul")
	judulRoutes.GET("/", controllers.GetAllJudul)
	judulRoutes.GET("/:id", controllers.GetJudulByID)
	judulRoutes.GET("/mahasiswa/:id", controllers.GetJudulByMahasiswaID)
	judulRoutes.POST("/", controllers.CreateJudul)
	judulRoutes.DELETE("/:id", controllers.DeleteJudul)
	judulRoutes.PUT("/:id", controllers.UpdateJudul)
	judulRoutes.GET("/titles", controllers.FetchTitles)
}
