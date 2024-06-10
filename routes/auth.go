package routes

import (
	"github.com/Seyditz/project-skripsi/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	authRoutes.POST("/admin/login", controllers.AdminLogin)
	authRoutes.POST("/mahasiswa/login", controllers.MahasiswaLogin)
	authRoutes.POST("/dosen/login", controllers.DosenLogin)
}
