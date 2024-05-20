package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllAdmin(c *gin.Context) {
	admins := []models.Admin{}
	// database.DB.Find(&admins)

	db := database.DB

	name := c.Query("name")
	email := c.Query("email")

	// Build the query conditionally based on the parameters
	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if email != "" {
		db = db.Where("email ILIKE ?", "%"+email+"%")
	}

	result := db.Find(&admins)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all admins", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &admins})
}

func CreateAdmin(c *gin.Context) {
	var admin models.Admin

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var existingAdmin models.Admin
	if result := database.DB.Where("email = ?", admin.Email).First(&existingAdmin); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	if admin.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if admin.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if admin.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	admin.Password = string(hashedPassword)

	// Create the admin in the database
	if result := database.DB.Create(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &admin})
}

func UpdateAdmin(c *gin.Context) {
	var admin models.Admin

	// Bind the JSON body to the admin struct
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the email field is provided and if it's different from the current admin's email
	var existingAdmin models.Admin
	if result := database.DB.Where("email = ?", admin.Email).First(&existingAdmin); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Check if email already exists for another admin
	if existingAdmin.ID != admin.ID {
		var emailCheck models.Admin
		if result := database.DB.Where("email = ?", admin.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if admin.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if admin.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// If password is provided, encrypt it
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		admin.Password = string(hashedPassword)
	} else {
		// If password is not provided, keep the existing one
		admin.Password = existingAdmin.Password
	}

	// Update the admin in the database
	if result := database.DB.Model(&existingAdmin).Updates(admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &admin})
}

func DeleteAdmin(c *gin.Context) {
	// Get the admin ID from the URL parameters
	adminID := c.Param("id")

	// Find the admin by ID
	var admin models.Admin
	if result := database.DB.First(&admin, adminID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Delete the admin from the database
	if result := database.DB.Delete(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Admin deleted successfully"})
}

func GetAdminbyId(c *gin.Context) {
	// Get the admin ID from the URL parameters
	adminID := c.Param("id")

	// Find the admin by ID
	var admin models.Admin
	if result := database.DB.First(&admin, adminID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(200, gin.H{"admin": &admin})
}
