package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllMahasiswa(c *gin.Context) {
	mahasiswas := []models.Mahasiswa{}
	db := database.DB

	name := c.Query("name")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	result := db.Find(&mahasiswas)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all mahasiswas", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswas})
}

func CreateMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa

	if err := c.ShouldBindJSON(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var existingMahasiswa models.Mahasiswa
	if result := database.DB.Where("email = ?", mahasiswa.Email).First(&existingMahasiswa); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	if mahasiswa.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if mahasiswa.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if mahasiswa.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(mahasiswa.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	mahasiswa.Password = string(hashedPassword)

	// Create the mahasiswa in the database
	if result := database.DB.Create(&mahasiswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswa})
}

func UpdateMahasiswa(c *gin.Context) {
	var mahasiswa models.Mahasiswa

	// Bind the JSON body to the mahasiswa struct
	if err := c.ShouldBindJSON(&mahasiswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the email field is provided and if it's different from the current mahasiswa's email
	var existingMahasiswa models.Mahasiswa
	if result := database.DB.Where("email = ?", mahasiswa.Email).First(&existingMahasiswa); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa not found"})
		return
	}

	// Check if email already exists for another mahasiswa
	if existingMahasiswa.Id != mahasiswa.Id {
		var emailCheck models.Mahasiswa
		if result := database.DB.Where("email = ?", mahasiswa.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if mahasiswa.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if mahasiswa.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// If password is provided, encrypt it
	if mahasiswa.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(mahasiswa.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		mahasiswa.Password = string(hashedPassword)
	} else {
		// If password is not provided, keep the existing one
		mahasiswa.Password = existingMahasiswa.Password
	}

	// Update the mahasiswa in the database
	if result := database.DB.Model(&existingMahasiswa).Updates(mahasiswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswa})
}

func DeleteMahasiswa(c *gin.Context) {
	// Get the mahasiswa ID from the URL parameters
	mahasiswaID := c.Param("id")

	// Find the mahasiswa by ID
	var mahasiswa models.Mahasiswa
	if result := database.DB.First(&mahasiswa, mahasiswaID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa not found"})
		return
	}

	// Delete the mahasiswa from the database
	if result := database.DB.Delete(&mahasiswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Mahasiswa deleted successfully"})
}

func GetMahasiswaByID(c *gin.Context) {
	// Get the mahasiswa ID from the URL parameters
	mahasiswaID := c.Param("id")

	// Find the mahasiswa by ID
	var mahasiswa models.Mahasiswa
	if result := database.DB.First(&mahasiswa, mahasiswaID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa not found"})
		return
	}

	c.JSON(200, gin.H{"mahasiswa": &mahasiswa})
}
