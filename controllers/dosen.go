package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetAllDosens retrieves all dosens from the database
func GetAllDosens(c *gin.Context) {
	dosens := []models.Dosen{}
	db := database.DB

	name := c.Query("name")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}

	result := db.Find(&dosens)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all dosens", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &dosens})
}

// CreateDosen creates a new dosen in the database
func CreateDosen(c *gin.Context) {
	var dosen models.Dosen

	if err := c.ShouldBindJSON(&dosen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var existingDosen models.Dosen
	if result := database.DB.Where("email = ?", dosen.Email).First(&existingDosen); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	if dosen.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if dosen.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if dosen.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dosen.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	dosen.Password = string(hashedPassword)

	// Create the dosen in the database
	if result := database.DB.Create(&dosen); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &dosen})
}

// UpdateDosen updates an existing dosen in the database
func UpdateDosen(c *gin.Context) {
	var dosen models.Dosen

	// Bind the JSON body to the dosen struct
	if err := c.ShouldBindJSON(&dosen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the email field is provided and if it's different from the current dosen's email
	var existingDosen models.Dosen
	if result := database.DB.Where("email = ?", dosen.Email).First(&existingDosen); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen not found"})
		return
	}

	// Check if email already exists for another dosen
	if existingDosen.Id != dosen.Id {
		var emailCheck models.Dosen
		if result := database.DB.Where("email = ?", dosen.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if dosen.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if dosen.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// If password is provided, encrypt it
	if dosen.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dosen.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		dosen.Password = string(hashedPassword)
	} else {
		// If password is not provided, keep the existing one
		dosen.Password = existingDosen.Password
	}

	// Update the dosen in the database
	if result := database.DB.Model(&existingDosen).Updates(dosen); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &dosen})
}

// DeleteDosen deletes an existing dosen from the database
func DeleteDosen(c *gin.Context) {
	// Get the dosen ID from the URL parameters
	dosenID := c.Param("id")

	// Find the dosen by ID
	var dosen models.Dosen
	if result := database.DB.First(&dosen, dosenID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen not found"})
		return
	}

	// Delete the dosen from the database
	if result := database.DB.Delete(&dosen); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Dosen deleted successfully"})
}

// GetDosenByID retrieves a dosen by its ID from the database
func GetDosenByID(c *gin.Context) {
	// Get the dosen ID from the URL parameters
	dosenID := c.Param("id")

	// Find the dosen by ID
	var dosen models.Dosen
	if result := database.DB.First(&dosen, dosenID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen not found"})
		return
	}

	c.JSON(200, gin.H{"dosen": &dosen})
}
