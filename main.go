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

func main() {
	//Load the .env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	//Connect to database
	database.Connect()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})
	routes.AdminRoutes(r)
	routes.MahasiswaRoutes(r)

	port := os.Getenv("PORT")

	r.Run(":" + port)

}
