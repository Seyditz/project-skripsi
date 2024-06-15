package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateTags godoc
// @Summary Admin Login
// @Description Admin Login
// @Produce application/json
// @Param request body models.AdminLoginRequest true "Raw Request Body"
// @Tags Auth
// @Success 200 {object} interface{}
// @Router /auth/admin/login [post]
func AdminLogin(c *gin.Context) {
	var admin models.Admin
	var input models.AdminLoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateJWT(admin.Email, admin.Id, []string{"admin", "dosen", "mahasiswa"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// CreateTags godoc
// @Summary Dosen Login
// @Description Dosen Login
// @Produce application/json
// @Param request body models.DosenLoginRequest true "Raw Request Body"
// @Tags Auth
// @Success 200 {object} interface{}
// @Router /auth/dosen/login [post]
func DosenLogin(c *gin.Context) {
	var dosen models.Dosen
	var input models.Dosen

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("nim = ?", input.Email).First(&dosen).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dosen.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Email or password"})
		return
	}

	dosenDataResponse := models.DosenDataResponse{
		Id:        dosen.Id,
		Name:      dosen.Name,
		Nidn:      dosen.Nidn,
		Email:     dosen.Email,
		Prodi:     dosen.Prodi,
		Kepakaran: dosen.Kepakaran,
		Jabatan:   dosen.Jabatan,
	}

	token, err := utils.GenerateJWT(dosen.Email, dosen.Id, []string{"dosen", "mahasiswa"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "data": dosenDataResponse})
}

// CreateTags godoc
// @Summary Mahasiswa Login
// @Description Mahasiswa Login
// @Produce application/json
// @Param request body models.MahasiswaLoginRequest true "Raw Request Body"
// @Tags Auth
// @Success 200 {object} interface{}
// @Router /auth/mahasiswa/login [post]
func MahasiswaLogin(c *gin.Context) {
	var mahasiswa models.Mahasiswa
	var input models.Mahasiswa

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Where("nim = ?", input.NIM).First(&mahasiswa).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid NIM or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(mahasiswa.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid NIM or password"})
		return
	}

	mahasiswaDataResponse := models.MahasiswaDataResponse{
		Id:       mahasiswa.Id,
		Name:     mahasiswa.Name,
		NIM:      mahasiswa.NIM,
		Email:    mahasiswa.Email,
		Prodi:    mahasiswa.Prodi,
		Angkatan: mahasiswa.Angkatan,
		SKS:      mahasiswa.SKS,
	}

	token, err := utils.GenerateJWT(mahasiswa.Email, mahasiswa.Id, []string{"mahasiswa"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "data": mahasiswaDataResponse})
}
