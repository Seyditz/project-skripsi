package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	middleware "github.com/Seyditz/project-skripsi/middlewares"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
)

func MahasiswaRoutes(r *gin.Engine) {
	mahasiswaRoutes := r.Group("/mahasiswa", middleware.AuthMiddleware())
	mahasiswaRoutes.GET("/", utils.Authorize("mahasiswa"), controllers.GetAllMahasiswa)
	mahasiswaRoutes.GET("/:id", utils.Authorize("mahasiswa"), controllers.GetMahasiswaByID)
	mahasiswaRoutes.POST("/", utils.Authorize("mahasiswa"), controllers.CreateMahasiswa)
	mahasiswaRoutes.DELETE("/:id", utils.Authorize("mahasiswa"), controllers.DeleteMahasiswa)
	mahasiswaRoutes.PUT("/:id", utils.Authorize("mahasiswa"), controllers.UpdateMahasiswa)
}
