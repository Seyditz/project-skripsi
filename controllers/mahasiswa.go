package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Mahasiswa
// @Description Get All Mahasiswa
// @Produce application/json
// @Tags Mahasiswa
// @Param name query string false "name"
// @Success 200 {object} []models.MahasiswaDataResponse{}
// @Router /mahasiswa [get]
func GetAllMahasiswa(c *gin.Context) {
	mahasiswas := []models.MahasiswaDataResponse{}
	db := database.DB

	name := c.Query("name")

	if name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
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

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

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

	// Handle the image file
	// file, err := c.FormFile("image")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
	// 	return
	// }

	// Create a unique filename using UUID
	// uniqueFileName := uuid.New().String() + filepath.Ext(file.Filename)
	// filePath := "./uploads/" + uniqueFileName

	// Save the file to the server
	// if err := c.SaveUploadedFile(file, filePath); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	mahasiswa := models.Mahasiswa{
		Name:     input.Name,
		NIM:      input.NIM,
		Email:    input.Email,
		Prodi:    input.Prodi,
		Password: input.Password,
		Angkatan: input.Angkatan,
		SKS:      input.SKS,
		Image:    "data:image/png;base64,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",
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

	// Bind the JSON body to the mahasiswa struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
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
	if input.Image == "" {
		input.Image = existingMahasiswa.Image
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

	mahasiswa := models.Mahasiswa{
		Name:      input.Name,
		NIM:       input.NIM,
		Email:     input.Email,
		Prodi:     input.Prodi,
		Password:  input.Password,
		Angkatan:  input.Angkatan,
		SKS:       input.SKS,
		Image:     input.Image,
		UpdatedAt: time.Now(),
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
