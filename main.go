package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Seyditz/project-skripsi/controllers"
	_ "github.com/Seyditz/project-skripsi/docs"
	"github.com/gin-contrib/gzip"
	"github.com/google/uuid"

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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With, Accept")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Writer.Header().Set("X-Request-Id", uuid.String())
		c.Next()
	}
}

func main() {
	//Load the .env file
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	//Connect to database
	database.Connect()

	r := gin.Default()

	// r.Use(cors.Default())

	// config := cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "X-Requested-With"},
	// 	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }

	// Apply the CORS middleware to the router
	r.Use(CORSMiddleware())
	r.Use(RequestIDMiddleware())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	log.Printf("Route: %+v/n", r)

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello"})
	})

	testRoute := r.Group("/test")
	testRoute.GET("/get-admin", controllers.GetAllAdmin)

	routes.AdminRoutes(r)
	routes.DosenRoutes(r)
	routes.MahasiswaRoutes(r)
	routes.PengajuanRoutes(r)
	routes.JudulRoutes(r)
	routes.AuthRoutes(r)
	routes.NotificationRoutes(r)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")

	r.Run(":" + port)

}
