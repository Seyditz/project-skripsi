package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func PengajuanRoutes(r *gin.Engine) {
	pengajuanRoutes := r.Group("/pengajuan")
	pengajuanRoutes.GET("/", controllers.GetAllPengajuan)
	pengajuanRoutes.GET("/:id", controllers.GetPengajuanByID)
	pengajuanRoutes.GET("/mahasiswa/:id", controllers.GetPengajuanByMahasiswaID)
	pengajuanRoutes.POST("/", controllers.CreatePengajuan)
	pengajuanRoutes.POST("/similarity-test", controllers.SimilartityTest)
	pengajuanRoutes.DELETE("/:id", controllers.DeletePengajuan)
	pengajuanRoutes.PUT("/:id", controllers.UpdatePengajuan)
}
