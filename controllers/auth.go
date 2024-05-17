package controllers

import (
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(c *gin.Context) {
	var admin models.Admin
	var input models.Admin

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

	token, err := utils.GenerateJWT(admin.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

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

	token, err := utils.GenerateJWT(mahasiswa.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "data": mahasiswaDataResponse})
}
