package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func DosenRoutes(r *gin.Engine) {
	mahasiswaRoutes := r.Group("/dosen")
	mahasiswaRoutes.GET("/", controllers.GetAllDosens)
	mahasiswaRoutes.GET("/:id", controllers.GetDosenByID)
	mahasiswaRoutes.POST("/", controllers.CreateDosen)
	mahasiswaRoutes.DELETE("/:id", controllers.DeleteDosen)
	mahasiswaRoutes.PUT("/:id", controllers.UpdateDosen)
	// mahasiswaRoutes.POST("/login", controllers.Dosen)
}
