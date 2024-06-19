package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/Seyditz/project-skripsi/docs"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Sijudul API
// @version 1.0
// @description An API for Sijudul App using Gin

// @host projectskripsi-fvwdncsc.b4a.run
// @BasePath /
// @schemes https http
func main() {
	//Load the .env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	//Connect to database
	database.Connect()

	r := gin.Default()

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.Use(cors.Default())

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})
	routes.AdminRoutes(r)
	routes.DosenRoutes(r)
	routes.MahasiswaRoutes(r)
	routes.PengajuanRoutes(r)
	routes.JudulRoutes(r)
	routes.AuthRoutes(r)
	routes.NotificationRoutes(r)

	port := os.Getenv("PORT")

	r.Run(":" + port)

}
