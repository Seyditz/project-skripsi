package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Go Gin Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

func main() {
	//Load the .env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	//Connect to database
	database.Connect()

	r := gin.Default()

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})
	routes.AdminRoutes(r)
	routes.MahasiswaRoutes(r)
	routes.PengajuanRoutes(r)

	port := os.Getenv("PORT")

	r.Run(":" + port)

}
