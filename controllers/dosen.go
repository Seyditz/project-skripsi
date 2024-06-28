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
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Dosen
// @Description Get All Dosens
// @Produce application/json
// @Tags Dosen
// @Param name query string false "name"
// @Success 200 {object} []models.DosenDataResponse{}
// @Router /dosen/ [get]
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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Create Dosen
// @Description Create Dosen
// @Produce application/json
// @Param request body models.DosenCreateRequest true "Raw Request Body"
// @Tags Dosen
// @Success 200 {object} models.Dosen{}
// @Router /dosen [post]
func CreateDosen(c *gin.Context) {
	var input models.DosenCreateRequest

	input.Name = c.PostForm("name")
	input.Nidn = c.PostForm("nidn")
	input.Email = c.PostForm("email")
	input.Password = c.PostForm("password")
	input.Prodi = c.PostForm("prodi")
	input.Jabatan = c.PostForm("jabatan")
	input.Kepakaran = c.PostForm("kepakaran")
	input.Kapasitas, _ = strconv.Atoi(c.PostForm("kapasitas"))
	input.Gelar = c.PostForm("gelar")
	input.Kepakaran = c.PostForm("kepakaran")

	input.TanggalLahir, _ = time.Parse("0000-00-00", c.PostForm("tanggal_lahir"))

	// input.Gelar = pq.StringArray(c.PostFormArray("gelar"))
	// input.JenjangAkademik = pq.StringArray(c.PostFormArray("jenjang_akademik"))

	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	var existingDosen models.Dosen
	if result := database.DB.Where("nidn = ?", input.Nidn).First(&existingDosen); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIDN already exists"})
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
	if input.Nidn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nidn is required"})
		return
	}
	if input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is required"})
		return
	}
	if input.Prodi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prodi is required"})
		return
	}
	if input.Jabatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jabatan is required"})
		return
	}
	if input.Kepakaran == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kepakaran is required"})
		return
	}
	if input.Kapasitas == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kapasitas is required"})
		return
	}
	// if input.MahasiswaBimbinganId == nil {
	// 	input.MahasiswaBimbinganId = utils.IntArray{}
	// }

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

	dosen := models.Dosen{
		Name:                 input.Name,
		Nidn:                 input.Nidn,
		Email:                input.Email,
		Password:             input.Password,
		Prodi:                input.Prodi,
		Jabatan:              input.Jabatan,
		Kepakaran:            input.Kepakaran,
		Kapasitas:            input.Kapasitas,
		MahasiswaBimbinganId: pq.Int64Array{},
		Gelar:                input.Gelar,
		JenjangAkademik:      input.JenjangAkademik,
		TanggalLahir:         input.TanggalLahir,
		CreatedAt:            time.Now(),
		Image:                imageUrl,
	}

	// Create the dosen in the database
	if result := database.DB.Create(&dosen); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error(), "res": result})
		return
	}

	c.JSON(200, gin.H{"result": &dosen})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Update Dosen
// @Description Update Dosen
// @Produce application/json
// @Param request body models.DosenUpdateRequest true "Raw Request Body"
// @Param id path int true "Dosen ID"
// @Tags Dosen
// @Success 200 {object} models.Dosen{}
// @Router /dosen/{id} [put]
func UpdateDosen(c *gin.Context) {
	var input models.Dosen
	dosenID, _ := strconv.Atoi(c.Param("id"))

	// Bind the JSON body to the dosen struct
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Check if the email field is provided and if it's different from the current dosen's email
	var existingDosen models.Dosen
	if result := database.DB.First(&existingDosen, dosenID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen not found"})
		return
	}

	// Check if email already exists for another dosen
	if existingDosen.Id != dosenID {
		var emailCheck models.Dosen
		if result := database.DB.Where("email = ?", input.Email).First(&emailCheck); result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	// Validate required fields
	if input.Name == "" {
		input.Name = existingDosen.Name
	}
	if input.Email == "" {
		input.Email = existingDosen.Email
	}
	if input.Nidn == "" {
		input.Nidn = existingDosen.Nidn
	}
	if input.Prodi == "" {
		input.Prodi = existingDosen.Prodi
	}
	if input.Jabatan == "" {
		input.Jabatan = existingDosen.Jabatan
	}
	if input.Kepakaran == "" {
		input.Kepakaran = existingDosen.Kepakaran
	}
	if input.Kapasitas == 0 {
		input.Kapasitas = existingDosen.Kapasitas
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
		input.Password = existingDosen.Password
	}

	// Handle image upload
	imageUrl := existingDosen.Image
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()

		// Delete the old image from Firebase Storage
		if existingDosen.Image != "" {
			err = utils.DeleteFileFromFirebase(context.Background(), "sijudul-610fb.appspot.com", existingDosen.Image)
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

	dosen := models.Dosen{
		Id:                   existingDosen.Id,
		Name:                 input.Name,
		Nidn:                 input.Nidn,
		Email:                input.Email,
		Password:             input.Password,
		Prodi:                input.Prodi,
		Jabatan:              input.Jabatan,
		Kepakaran:            input.Kepakaran,
		Kapasitas:            input.Kapasitas,
		MahasiswaBimbinganId: input.MahasiswaBimbinganId,
		UpdatedAt:            time.Now(),
		Image:                imageUrl,
	}

	// Update the dosen in the database
	if result := database.DB.Model(&existingDosen).Updates(dosen); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &dosen})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Delete Dosen
// @Description Delete Dosen
// @Produce application/json
// @Param id path int true "Dosen ID"
// @Tags Dosen
// @Success 200
// @Router /dosen/{id} [delete]
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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get Dosen By ID
// @Description Get Dosen By ID
// @Produce application/json
// @Param id path int true "Dosen ID"
// @Tags Dosen
// @Success 200 {object} models.DosenDataResponse{}
// @Router /dosen/{id} [get]
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

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Mahasiswa Bimbingan
// @Param id path int true "Dosen ID"
// @Description Get All Mahasiswa Bimbingan
// @Produce application/json
// @Tags Dosen
// @Success 200 {object} models.DosenMahasiswaBimbinganResponse{}
// @Router /dosen/mahasiswa-bimbingan/{id} [get]
func GetAllMahasiswaBimbingan(c *gin.Context) {
	var mahasiswaList []models.MahasiswaDataResponse
	// // Get the dosen ID from the URL parameters
	dosenID := c.Param("id")

	// Find the dosen by ID
	var dosen models.Dosen
	if res := database.DB.First(&dosen, dosenID); res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dosen not found"})
		return
	}

	// Convert pq.Int64Array to []int
	intArray := make([]int, len(dosen.MahasiswaBimbinganId))
	for i, v := range dosen.MahasiswaBimbinganId {
		intArray[i] = int(v)
	}

	result := database.DB.
		Model(&models.Mahasiswa{}).
		Where("id IN ?", intArray).
		Find(&mahasiswaList)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all mahasiswas", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"mahasiswa_list": &mahasiswaList})
}
