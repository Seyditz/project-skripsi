package controllers

import (
	"fmt"
	"net/http"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
)

// GetAllPengajuan retrieves all Pengajuan records from the database
func GetAllPengajuan(c *gin.Context) {
	pengajuans := []models.Pengajuan{}
	judul := c.Query("judul")

	// result := database.DB.Find(&pengajuans)
	query := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction)

	if judul != "" {
		query = query.Where("judul ILIKE ?", "%"+judul+"%")
	}

	result := query.Find(&pengajuans)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all pengajuans", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuans})
}

// CreatePengajuan creates a new Pengajuan record in the database
func CreatePengajuan(c *gin.Context) {
	var input models.PengajuanCreateRequest

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

	pengajuan := models.Pengajuan{
		MahasiswaId:      input.MahasiswaId,
		Peminatan:        input.Peminatan,
		Judul:            input.Judul,
		TempatPenelitian: input.TempatPenelitian,
		RumusanMasalah:   input.RumusanMasalah,
		DosPem1Id:        input.DosPem1Id,
		DosPem2Id:        input.DosPem2Id,
		StatusAcc:        input.StatusAcc,
		RejectedNote:     input.RejectedNote,
	}

	// Create the pengajuan in the database
	if result := database.DB.Create(&pengajuan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuan})
}

// UpdatePengajuan updates an existing Pengajuan record in the database
func UpdatePengajuan(c *gin.Context) {
	var pengajuan models.Pengajuan
	pengajuanID := c.Param("id")

	if err := c.ShouldBindJSON(&pengajuan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if the Pengajuan exists
	var existingPengajuan models.Pengajuan
	if result := database.DB.First(&existingPengajuan, pengajuanID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan not found"})
		return
	}

	// Validate required fields
	if pengajuan.Peminatan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "peminatan is required"})
		return
	}
	if pengajuan.Judul == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "judul is required"})
		return
	}
	if pengajuan.TempatPenelitian == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tempat penelitian is required"})
		return
	}
	if pengajuan.RumusanMasalah == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rumusan masalah is required"})
		return
	}
	if pengajuan.DosPem1Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem1 is required"})
		return
	}
	if pengajuan.DosPem2Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dospem2 is required"})
		return
	}

	// Update the Pengajuan in the database
	if result := database.DB.Model(&existingPengajuan).Updates(pengajuan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuan})
}

// DeletePengajuan deletes an existing Pengajuan record from the database
func DeletePengajuan(c *gin.Context) {
	// Get the Pengajuan ID from the URL parameters
	pengajuanID := c.Param("id")

	// Find the Pengajuan by ID
	var pengajuan models.Pengajuan
	if result := database.DB.First(&pengajuan, pengajuanID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan not found"})
		return
	}

	// Delete the Pengajuan from the database
	if result := database.DB.Delete(&pengajuan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Pengajuan deleted successfully"})
}

// GetPengajuanByID retrieves a Pengajuan record by ID from the database
func GetPengajuanByID(c *gin.Context) {
	// Get the Pengajuan ID from the URL parameters
	pengajuanID := c.Param("id")

	// Find the Pengajuan by ID
	var pengajuan models.Pengajuan
	if result := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).First(&pengajuan, pengajuanID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan not found"})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuan})
}

// Get Pengajuan By Mahasiswa ID
func GetPengajuanByMahasiswaID(c *gin.Context) {
	// Get the Mahasiswa ID from the URL parameters
	mahasiswaID := c.Param("id")

	// Convert Mahasiswa ID to an integer (if necessary)
	var id int
	if _, err := fmt.Sscanf(mahasiswaID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Mahasiswa ID"})
		return
	}

	// Find Pengajuan records by Mahasiswa ID
	var pengajuan []models.Pengajuan
	if result := database.DB.Where("mahasiswa_id = ?", id).Find(&pengajuan); result.RowsAffected == 0 {
		message := fmt.Sprintf("Pengajuan not found at mahasiswa_id = %s", mahasiswaID)
		c.JSON(http.StatusNotFound, gin.H{"error": message})
		return
	}

	c.JSON(200, gin.H{"result": pengajuan})
}

func SimiliartityTest(c *gin.Context) {

}
