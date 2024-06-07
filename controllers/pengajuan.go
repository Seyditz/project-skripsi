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
	query := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Preload("Mahasiswa", models.MahasiswaSafePreloadFunction)

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
	if result := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Preload("Mahasiswa", models.MahasiswaSafePreloadFunction).First(&pengajuan, pengajuanID); result.RowsAffected == 0 {
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

func levenshteinDistance(s, t string) int {
	m := len(s)
	n := len(t)
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			cost := 0
			if s[i-1] != t[j-1] {
				cost = 1
			}
			d[i][j] = min(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[m][n]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	} else if b < c {
		return b
	}
	return c
}

func similarityPercentage(s, t string) float64 {
	distance := levenshteinDistance(s, t)
	maxLength := len(s)
	if len(t) > maxLength {
		maxLength = len(t)
	}
	return 100.0 - (float64(distance) / float64(maxLength) * 100.0)
}

func SimiliartityTest(c *gin.Context) {
	var pengajuans []models.Pengajuan
	database.DB.Preload("Mahasiswa").Preload("DosPem1").Preload("DosPem2").Find(&pengajuans)

	var req models.SimilarityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	judul := req.Judul
	if judul == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "judul is required"})
		return
	}

	for _, pengajuan := range pengajuans {
		similarity := similarityPercentage(judul, pengajuan.Judul)
		if similarity > 60.0 {
			c.JSON(http.StatusOK, gin.H{
				"message":    "Similar title found",
				"similar":    pengajuan.Judul,
				"similarity": similarity,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "No similar title found with similarity > 60%"})
}
