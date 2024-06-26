package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm"
)

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Admin
// @Description Get All Admins
// @Produce application/json
// @Param name query string false "name"
// @Param email query string false "email"
// @Tags Admin
// @Success 200 {object} []models.AdminDataResponse{}
// @Router /admin/ [get]
func GetAllAdmin(c *gin.Context) {
	admins := []models.AdminDataResponse{}

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

	// result := db.Find(&admins)
	result := db.Model(&[]models.Admin{}).Find(&admins)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all admins", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &admins})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Create Admin
// @Description Create Admins
// @Produce application/json
// @Param request body models.AdminCreateRequest true "Raw Request Body"
// @Tags Admin
// @Success 200 {object} models.Admin{}
// @Router /admin [post]
func CreateAdmin(c *gin.Context) {
	var input models.AdminCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	var existingAdmin models.Admin
	if result := database.DB.Where("email = ?", input.Email).First(&existingAdmin); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}

	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}
	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}
	if input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	input.Password = string(hashedPassword)

	admin := models.Admin{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: time.Now(),
	}

	// Create the admin in the database
	if result := database.DB.Create(&admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &admin})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Update Admin
// @Description Update Admins
// @Produce application/json
// @Param request body models.AdminUpdateRequest true "Raw Request Body"
// @Param id path int true "Admin ID"
// @Tags Admin
// @Success 200 {object} models.Admin{}
// @Router /admin/{id} [put]
func UpdateAdmin(c *gin.Context) {
	var input models.AdminUpdateRequest
	adminID, _ := strconv.Atoi(c.Param("id"))

	// Bind the JSON body to the admin struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the email field is provided and if it's different from the current admin's email
	var existingAdmin models.Admin
	if result := database.DB.First(&existingAdmin, adminID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	// Check if email already exists for another admin
	if existingAdmin.Id != adminID {
		var emailCheck models.Admin
		if result := database.DB.Where("email = ?", input.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if input.Name == "" {
		input.Name = existingAdmin.Name
	}
	if input.Email == "" {
		input.Email = existingAdmin.Email
	}

	// If password is provided, encrypt it
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
			return
		}
		input.Password = string(hashedPassword)
	} else {
		// If password is not provided, keep the existing one
		input.Password = existingAdmin.Password
	}

	admin := models.Admin{
		Id:        existingAdmin.Id,
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		UpdatedAt: time.Now(),
	}

	// Update the admin in the database
	if result := database.DB.Model(&existingAdmin).Updates(admin); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &admin})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Delete Admin
// @Description Delete Admins
// @Produce application/json
// @Param id path int true "Admin ID"
// @Tags Admin
// @Success 200
// @Router /admin/{id} [delete]
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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get Admin By ID
// @Description Get Admin By ID
// @Produce application/json
// @Param id path int true "Admin ID"
// @Tags Admin
// @Success 200 {object} models.AdminDataResponse{}
// @Router /admin/{id} [get]
func GetAdminbyId(c *gin.Context) {
	// Get the admin ID from the URL parameters
	adminID := c.Param("id")

	// Find the admin by ID
	var admin models.AdminDataResponse
	if result := database.DB.Model(&models.Admin{}).First(&admin, adminID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Admin not found"})
		return
	}

	c.JSON(200, gin.H{"admin": &admin})
}
