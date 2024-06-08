package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/gin-gonic/gin"
)

// CreateTags godoc
// @Summary Get All Pengajuan
// @Description Get All Pengajuan
// @Produce application/json
// @Tags Pengajuan
// @Success 200 {object} []models.Pengajuan{}
// @Router /pengajuan [get]
func GetAllPengajuan(c *gin.Context) {
	// pengajuans := []models.PengajuanDataResponse{}
	pengajuans := []models.Pengajuan{}
	judul := c.Query("judul")

	// result := database.DB.Find(&pengajuans)
	query := database.DB.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Preload("Mahasiswa", models.MahasiswaSafePreloadFunction)

	if judul != "" {
		query = query.Where("judul ILIKE ?", "%"+judul+"%")
	}

	// result := query.Model(&models.Pengajuan{}).Find(&pengajuans)
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

// CreateTags godoc
// @Summary Update Pengajuan
// @Description Update Pengajuan
// @Produce application/json
// @Param request body models.PengajuanUpdateRequest true "Raw Request Body"
// @Param id path int true "Pengajuan ID"
// @Tags Pengajuan
// @Success 200 {object} models.Pengajuan{}
// @Router /pengajuan/{id} [put]
func UpdatePengajuan(c *gin.Context) {
	var input models.Pengajuan
	pengajuanID := c.Param("id")

	if err := c.ShouldBindJSON(&input); err != nil {
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
	if input.Peminatan == "" {
		input.Peminatan = existingPengajuan.Peminatan
	}
	if input.Judul == "" {
		input.Judul = existingPengajuan.Judul
	}
	if input.TempatPenelitian == "" {
		input.TempatPenelitian = existingPengajuan.TempatPenelitian
	}
	if input.RumusanMasalah == "" {
		input.RumusanMasalah = existingPengajuan.RumusanMasalah
	}
	if input.DosPem1Id == 0 {
		input.DosPem1Id = existingPengajuan.DosPem1Id
	}
	if input.DosPem2Id == 0 {
		input.DosPem2Id = existingPengajuan.DosPem2Id
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
		UpdatedAt:        time.Now(),
	}
	// Update the Pengajuan in the database
	if result := database.DB.Model(&existingPengajuan).Updates(&pengajuan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuan})
}

// CreateTags godoc
// @Summary Delete Pengajuan
// @Description Delete Pengajuan
// @Produce application/json
// @Param id path int true "Pengajuan ID"
// @Tags Pengajuan
// @Success 200
// @Router /pengajuan/{id} [delete]
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

// CreateTags godoc
// @Summary Get Pengajuan By ID
// @Description Get Pengajuan By ID
// @Produce application/json
// @Param id path int true "Pengajuan ID"
// @Tags Pengajuan
// @Success 200 {object} models.Pengajuan{}
// @Router /pengajuan/{id} [get]
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

// CreateTags godoc
// @Summary Get Pengajuan By Mahasiswa ID
// @Description Get Pengajuan By Mahasiswa ID
// @Produce application/json
// @Param id path int true "Mahasiswa ID"
// @Tags Pengajuan
// @Success 200 {object} models.Pengajuan{}
// @Router /pengajuan/mahasiswa/{id} [get]
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

// CreateTags godoc
// @Summary Similiarity Test Pengajuan
// @Description Similiarity Test Pengajuan
// @Produce application/json
// @Param request body models.SimilarityRequest true "Raw Request Body"
// @Tags Pengajuan
// @Success 200 {object} interface{}
// @Router /pengajuan/similarity-test [post]
func SimilartityTest(c *gin.Context) {
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
