package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	middleware "github.com/Seyditz/project-skripsi/middlewares"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
)

func DosenRoutes(r *gin.Engine) {
	mahasiswaRoutes := r.Group("/dosen", middleware.AuthMiddleware())
	mahasiswaRoutes.GET("/", utils.Authorize("mahasiswa"), controllers.GetAllDosens)
	mahasiswaRoutes.GET("/:id", utils.Authorize("mahasiswa"), controllers.GetDosenByID)
	mahasiswaRoutes.POST("/", utils.Authorize("admin"), controllers.CreateDosen)
	mahasiswaRoutes.DELETE("/:id", utils.Authorize("admin"), controllers.DeleteDosen)
	mahasiswaRoutes.PUT("/:id", utils.Authorize("dosen"), controllers.UpdateDosen)
}
