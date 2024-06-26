package controllers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Mahasiswa
// @Param name query string false "name"
// @Param nim query string false "nim"
// @Description Get All Mahasiswa
// @Produce application/json
// @Tags Mahasiswa
// @Success 200 {object} []models.MahasiswaDataResponse{}
// @Router /mahasiswa [get]
func GetAllMahasiswa(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	mahasiswas := []models.MahasiswaDataResponse{}
	db := database.DB

	name := c.Query("name")
	nim := c.Query("nim")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if nim != "" {
		db = db.Where("nim ILIKE ?", "%"+nim+"%")
	}

	result := db.Model(&models.Mahasiswa{}).Find(&mahasiswas)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all mahasiswas", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswas})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Create Mahasiswa
// @Description Create Mahasiswa
// @Produce application/json
// @Param request body models.MahasiswaCreateRequest true "Raw Request Body"
// @Tags Mahasiswa
// @Success 200 {object} models.Mahasiswa{}
// @Router /mahasiswa [post]
func CreateMahasiswa(c *gin.Context) {
	var input models.MahasiswaCreateRequest

	input.Name = c.PostForm("name")
	input.NIM = c.PostForm("nim")
	input.Email = c.PostForm("email")
	input.Password = c.PostForm("password")
	input.Prodi = c.PostForm("prodi")
	input.Angkatan, _ = strconv.Atoi(c.PostForm("angkatan"))
	input.SKS, _ = strconv.Atoi(c.PostForm("sks"))

	var existingMahasiswa models.Mahasiswa
	if result := database.DB.Where("email = ?", input.Email).First(&existingMahasiswa); result.RowsAffected > 0 {
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
	if input.SKS <= 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKS minimum is 100"})
		return
	}

	// Encrypt the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt password"})
		return
	}
	input.Password = string(hashedPassword)

	// handle image upload
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))
	imageUrl, err := utils.UploadFileToFirebase(context.Background(), "sijudul-610fb.appspot.com", file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mahasiswa := models.Mahasiswa{
		Name:      input.Name,
		NIM:       input.NIM,
		Email:     input.Email,
		Prodi:     input.Prodi,
		Password:  input.Password,
		Angkatan:  input.Angkatan,
		SKS:       input.SKS,
		CreatedAt: time.Now(),
		Image:     imageUrl,
	}

	// Create the mahasiswa in the database
	if result := database.DB.Create(&mahasiswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswa})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Update Mahasiswa
// @Description Update Mahasiswa
// @Produce application/json
// @Param request body models.MahasiswaUpdateRequest true "Raw Request Body"
// @Param id path int true "Mahasiswa ID"
// @Tags Mahasiswa
// @Success 200 {object} models.Mahasiswa{}
// @Router /mahasiswa/{id} [put]
func UpdateMahasiswa(c *gin.Context) {
	var input models.MahasiswaUpdateRequest
	mahasiswaID, _ := strconv.Atoi(c.Param("id"))

	// Bind the form data to the mahasiswa struct
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Check if the email field is provided and if it's different from the current mahasiswa's email
	var existingMahasiswa models.Mahasiswa
	if result := database.DB.First(&existingMahasiswa, mahasiswaID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mahasiswa not found"})
		return
	}

	// Check if email already exists for another mahasiswa
	if existingMahasiswa.Id != mahasiswaID {
		var emailCheck models.Mahasiswa
		if result := database.DB.Where("email = ?", input.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if input.Name == "" {
		input.Name = existingMahasiswa.Name
	}
	if input.Email == "" {
		input.Email = existingMahasiswa.Email
	}
	if input.NIM == "" {
		input.NIM = existingMahasiswa.NIM
	}
	if input.Prodi == "" {
		input.Prodi = existingMahasiswa.Prodi
	}
	if input.Angkatan == 0 {
		input.Angkatan = existingMahasiswa.Angkatan
	}
	if input.SKS == 0 {
		input.SKS = existingMahasiswa.SKS
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
		input.Password = existingMahasiswa.Password
	}

	// Handle image upload
	imageUrl := existingMahasiswa.Image
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		// Delete the old image from Firebase Storage
		if existingMahasiswa.Image != "" {
			err = utils.DeleteFileFromFirebase(context.Background(), "sijudul-610fb.appspot.com", existingMahasiswa.Image)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		// Upload the new image to Firebase Storage
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))
		imageUrl, err = utils.UploadFileToFirebase(context.Background(), "sijudul-610fb.appspot.com", file, fileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	mahasiswa := models.Mahasiswa{
		Id:        existingMahasiswa.Id,
		Name:      input.Name,
		NIM:       input.NIM,
		Email:     input.Email,
		Prodi:     input.Prodi,
		Password:  input.Password,
		Angkatan:  input.Angkatan,
		SKS:       input.SKS,
		UpdatedAt: time.Now(),
		Image:     imageUrl,
	}

	// Update the mahasiswa in the database
	if result := database.DB.Model(&existingMahasiswa).Updates(mahasiswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &mahasiswa})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Delete Mahasiswa
// @Description Delete Mahasiswa
// @Produce application/json
// @Param id path int true "Mahasiswa ID"
// @Tags Mahasiswa
// @Success 200
// @Router /mahasiswa/{id} [delete]
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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get Mahasiswa By ID
// @Description Get Mahasiswa By ID
// @Produce application/json
// @Param id path int true "Mahasiswa ID"
// @Tags Mahasiswa
// @Success 200 {object} models.MahasiswaDataResponse{}
// @Router /mahasiswa/{id} [get]
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
