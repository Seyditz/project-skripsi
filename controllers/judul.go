package controllers

import (
	"fmt"
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
)

// GetAllJudul retrieves all Judul records from the database
func GetAllJudul(c *gin.Context) {
	juduls := []models.Judul{}
	result := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Find(&juduls)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all juduls", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &juduls})
}

// CreateJudul creates a new Judul record in the database
func CreateJudul(c *gin.Context) {
	var input models.JudulCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if input.Peminatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "peminatan is required"})
		return
	}
	if input.Judul == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "judul is required"})
		return
	}
	if input.TempatPenelitian == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tempat penelitian is required"})
		return
	}
	if input.RumusanMasalah == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rumusan masalah is required"})
		return
	}
	if input.DosPem1Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem1 is required"})
		return
	}
	if input.DosPem2Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem2 is required"})
		return
	}

	judul := models.Judul{
		MahasiswaId:      input.MahasiswaId,
		Peminatan:        input.Peminatan,
		Judul:            input.Judul,
		TempatPenelitian: input.TempatPenelitian,
		RumusanMasalah:   input.RumusanMasalah,
		DosPem1Id:        input.DosPem1Id,
		DosPem2Id:        input.DosPem2Id,
	}

	// Create the judul in the database
	if result := database.DB.Create(&judul); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &judul})
}

// UpdateJudul updates an existing Judul record in the database
func UpdateJudul(c *gin.Context) {
	var judul models.Judul
	judulID := c.Param("id")

	if err := c.ShouldBindJSON(&judul); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the Judul exists
	var existingJudul models.Judul
	if result := database.DB.First(&existingJudul, judulID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Judul not found"})
		return
	}

	// Validate required fields
	if judul.Peminatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "peminatan is required"})
		return
	}
	if judul.Judul == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "judul is required"})
		return
	}
	if judul.TempatPenelitian == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tempat penelitian is required"})
		return
	}
	if judul.RumusanMasalah == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rumusan masalah is required"})
		return
	}
	if judul.DosPem1Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem1 is required"})
		return
	}
	if judul.DosPem2Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem2 is required"})
		return
	}

	// Update the Judul in the database
	if result := database.DB.Model(&existingJudul).Updates(judul); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &judul})
}

// DeleteJudul deletes an existing Judul record from the database
func DeleteJudul(c *gin.Context) {
	// Get the Judul ID from the URL parameters
	judulID := c.Param("id")

	// Find the Judul by ID
	var judul models.Judul
	if result := database.DB.First(&judul, judulID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Judul not found"})
		return
	}

	// Delete the Judul from the database
	if result := database.DB.Delete(&judul); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Judul deleted successfully"})
}

// GetJudulByID retrieves a Judul record by ID from the database
func GetJudulByID(c *gin.Context) {
	// Get the Judul ID from the URL parameters
	judulID := c.Param("id")

	// Find the Judul by ID
	var judul models.Judul
	if result := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).First(&judul, judulID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Judul not found"})
		return
	}

	c.JSON(200, gin.H{"result": &judul})
}

// GetJudulByMahasiswaID retrieves Judul records by Mahasiswa ID from the database
func GetJudulByMahasiswaID(c *gin.Context) {
	// Get the Mahasiswa ID from the URL parameters
	mahasiswaID := c.Param("id")

	// Convert Mahasiswa ID to an integer (if necessary)
	var id int
	if _, err := fmt.Sscanf(mahasiswaID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Mahasiswa ID"})
		return
	}

	// Find Judul records by Mahasiswa ID
	var juduls []models.Judul
	if result := database.DB.Where("mahasiswa_id = ?", id).Find(&juduls); result.RowsAffected == 0 {
		message := fmt.Sprintf("Judul not found for mahasiswa_id = %s", mahasiswaID)
		c.JSON(http.StatusNotFound, gin.H{"error": message})
		return
	}

	c.JSON(200, gin.H{"result": juduls})
}
