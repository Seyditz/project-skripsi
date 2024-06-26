package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	middleware "github.com/Seyditz/project-skripsi/middlewares"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
)

func PengajuanRoutes(r *gin.Engine) {
	pengajuanRoutes := r.Group("/pengajuan", middleware.AuthMiddleware())
	pengajuanRoutes.GET("/", utils.Authorize("mahasiswa"), controllers.GetAllPengajuan)
	pengajuanRoutes.GET("/:id", utils.Authorize("mahasiswa"), controllers.GetPengajuanByID)
	pengajuanRoutes.GET("/mahasiswa/:id", utils.Authorize("mahasiswa"), controllers.GetPengajuanByMahasiswaID)
	pengajuanRoutes.GET("/dospem/:id", utils.Authorize("dosen"), controllers.GetPengajuanByDosPem1Id)
	pengajuanRoutes.POST("/", utils.Authorize("mahasiswa"), controllers.CreatePengajuan)
	pengajuanRoutes.POST("/similarity-test", utils.Authorize("mahasiswa"), controllers.SimilartityTest)
	pengajuanRoutes.DELETE("/:id", utils.Authorize("mahasiswa"), controllers.DeletePengajuan)
	pengajuanRoutes.PUT("/:id", utils.Authorize("mahasiswa"), controllers.UpdatePengajuan)
}
