package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
	"github.com/Seyditz/project-skripsi/utils"
	"github.com/gin-gonic/gin"
)

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get All Pengajuan
// @Description Get All Pengajuan
// @Param judul query string false "judul"
// @Produce application/json
// @Tags Pengajuan
// @Success 200 {object} []models.Pengajuan{}
// @Router /pengajuan/ [get]
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
	result := query.Find(&pengajuans).Order("created_at asc")

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all pengajuans", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuans})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Create Pengajuan
// @Description Create Pengajuan
// @Produce application/json
// @Param request body models.PengajuanCreateRequest true "Raw Request Body"
// @Tags Pengajuan
// @Success 200 {object} models.Pengajuan{}
// @Router /pengajuan [post]
func CreatePengajuan(c *gin.Context) {
	var input models.PengajuanCreateRequest
	var existingPengajuan models.Pengajuan
	var dospem1 models.Dosen
	var dospem2 models.Dosen

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
	} //a
	if input.Abstrak == "" {
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
	if result := database.DB.
		Where("mahasiswa_id = ?", input.MahasiswaId).
		Where("(status_acc = ? OR status_acc = ?)", "Pending", "Approved").
		Where("(status_acc_kaprodi = ? OR status_acc_kaprodi = ?)", "Pending", "Approved").First(&existingPengajuan); result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mahasiswa telah mengajukan Judul yang sudah pending/acc"})
		return
	}
	if result := database.DB.First(&dospem1, input.DosPem1Id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "DosPem 1 tidak ada"})
		return
	}
	if result := database.DB.First(&dospem2, input.DosPem2Id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "DosPem 2 tidak ada"})
		return
	}
	if len(dospem1.MahasiswaBimbinganId)+1 > dospem1.Kapasitas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kapasitas Dospem 1 sudah penuh, silahkan hubungi kaprodi jika ingin ditambahkan"})
		return
	}
	if len(dospem1.MahasiswaBimbinganId)+1 > dospem2.Kapasitas {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kapasitas Dospem 2 sudah penuh, silahkan hubungi kaprodi jika ingin ditambahkan"})
		return
	}

	pengajuan := models.Pengajuan{
		MahasiswaId:      input.MahasiswaId,
		Peminatan:        input.Peminatan,
		Judul:            input.Judul,
		TempatPenelitian: input.TempatPenelitian,
		Abstrak:          input.Abstrak,
		DosPem1Id:        input.DosPem1Id,
		DosPem2Id:        input.DosPem2Id,
		Metode:           input.Metode,
		StatusAcc:        "Pending",
		StatusAccKaprodi: "Pending",
		RejectedNote:     "",
		CreatedAt:        time.Now(),
	}

	// Create the pengajuan in the database
	if result := database.DB.Create(&pengajuan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Add Mahasiswa to Dosen Mahasiswa Bimbingan List
	dospem1UpdatedData := models.Dosen{
		MahasiswaBimbinganId: append(dospem1.MahasiswaBimbinganId, int64(input.MahasiswaId)),
	}
	dospem2UpdatedData := models.Dosen{
		MahasiswaBimbinganId: append(dospem2.MahasiswaBimbinganId, int64(input.MahasiswaId)),
	}

	if result := database.DB.Model(&models.Dosen{}).Where("id = ?", input.DosPem1Id).Updates(dospem1UpdatedData); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result := database.DB.Model(&models.Dosen{}).Where("id = ?", input.DosPem2Id).Updates(dospem2UpdatedData); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"result": &pengajuan})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
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
	var dospem1 models.Dosen
	var dospem2 models.Dosen

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the Pengajuan exists
	var existingPengajuan models.Pengajuan
	if result := database.DB.First(&existingPengajuan, pengajuanID); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengajuan not found"})
		return
	}

	// Validate required fields
	if input.MahasiswaId == 0 {
		input.MahasiswaId = existingPengajuan.MahasiswaId
	}
	if input.Peminatan == "" {
		input.Peminatan = existingPengajuan.Peminatan
	}
	if input.Judul == "" {
		input.Judul = existingPengajuan.Judul
	}
	if input.TempatPenelitian == "" {
		input.TempatPenelitian = existingPengajuan.TempatPenelitian
	}
	if input.Abstrak == "" {
		input.Abstrak = existingPengajuan.Abstrak
	}
	if input.DosPem1Id == 0 {
		input.DosPem1Id = existingPengajuan.DosPem1Id
	}
	if input.DosPem2Id == 0 {
		input.DosPem2Id = existingPengajuan.DosPem2Id
	}
	if input.StatusAcc == "" {
		input.StatusAcc = existingPengajuan.StatusAcc
	}
	if input.StatusAccKaprodi == "" {
		input.StatusAccKaprodi = existingPengajuan.StatusAccKaprodi
	}
	if result := database.DB.First(&dospem1, input.DosPem1Id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "DosPem 1 tidak ditemukan"})
		return
	}
	if result := database.DB.First(&dospem2, input.DosPem2Id); result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "DosPem 2 tidak ditemukan"})
		return
	}

	if input.StatusAcc == "Rejected" || input.StatusAccKaprodi == "Rejected" {
		input.StatusAcc = "Rejected"
		input.StatusAccKaprodi = "Rejected"
		if err := utils.RemoveMahasiswaBimbinganFromDosen(existingPengajuan.MahasiswaId, existingPengajuan.DosPem1Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := utils.RemoveMahasiswaBimbinganFromDosen(existingPengajuan.MahasiswaId, existingPengajuan.DosPem2Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if input.DosPem1Id != existingPengajuan.DosPem1Id && (input.StatusAcc != "Rejected" || input.StatusAccKaprodi != "Rejected") {
		if err := utils.RemoveMahasiswaBimbinganFromDosen(existingPengajuan.MahasiswaId, existingPengajuan.DosPem1Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := utils.AddMahasiswaBimbinganToDosen(existingPengajuan.MahasiswaId, input.DosPem1Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}
	if input.DosPem2Id != existingPengajuan.DosPem2Id && (input.StatusAcc != "Rejected" || input.StatusAccKaprodi != "Rejected") {
		if err := utils.RemoveMahasiswaBimbinganFromDosen(existingPengajuan.MahasiswaId, existingPengajuan.DosPem2Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := utils.AddMahasiswaBimbinganToDosen(existingPengajuan.MahasiswaId, input.DosPem2Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	pengajuan := models.Pengajuan{
		MahasiswaId:      input.MahasiswaId,
		Peminatan:        input.Peminatan,
		Judul:            input.Judul,
		TempatPenelitian: input.TempatPenelitian,
		Abstrak:          input.Abstrak,
		DosPem1Id:        input.DosPem1Id,
		DosPem2Id:        input.DosPem2Id,
		StatusAcc:        input.StatusAcc,
		StatusAccKaprodi: input.StatusAccKaprodi,
		RejectedNote:     input.RejectedNote,
		Metode:           input.Metode,
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
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
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

	if err := utils.RemoveMahasiswaBimbinganFromDosen(pengajuan.MahasiswaId, pengajuan.DosPem1Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.RemoveMahasiswaBimbinganFromDosen(pengajuan.MahasiswaId, pengajuan.DosPem2Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
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
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
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
	judul := c.Query("judul")

	// Convert Mahasiswa ID to an integer
	var id int
	if _, err := fmt.Sscanf(mahasiswaID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Mahasiswa ID"})
		return
	}

	// Find Pengajuan records by Mahasiswa ID
	var pengajuan []models.Pengajuan
	query := database.DB.Where("mahasiswa_id = ?", id)

	if judul != "" {
		query = query.Where("judul ILIKE ?", "%"+judul+"%")
	}

	result := query.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Preload("Mahasiswa", models.MahasiswaSafePreloadFunction).Find(&pengajuan)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all pengajuan by mahasiswa ID", "error": result.Error})
		return
	}

	c.JSON(200, gin.H{"result": pengajuan})
}

// CreateTags godoc
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Get Pengajuan By Dospem ID
// @Description Get Pengajuan By Dospem ID
// @Produce application/json
// @Param id path int true "Dospem ID"
// @Tags Pengajuan
// @Success 200 {object} models.Pengajuan{}
// @Router /pengajuan/dospem/{id} [get]
func GetPengajuanByDosPem1Id(c *gin.Context) {
	// Get the DosPem1 ID from the URL parameters
	dospem1ID := c.Param("id")
	judul := c.Query("judul")

	// Convert DosPem1 ID to an integer
	var id int
	if _, err := fmt.Sscanf(dospem1ID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DosPem1 ID"})
		return
	}

	// Find Pengajuan records by DosPem1 ID
	var pengajuan []models.Pengajuan
	query := database.DB.Where("dos_pem1_id = ?", id)

	if judul != "" {
		query = query.Where("judul ILIKE ?", "%"+judul+"%")
	}

	result := query.Preload("DosPem1", models.DosenSafePreloadFunction).Preload("DosPem2", models.DosenSafePreloadFunction).Preload("Mahasiswa", models.MahasiswaSafePreloadFunction).Find(&pengajuan)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get all pengajuan by DosPem1 ID", "error": result.Error})
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
// @param Authorization header string true "example : Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc2MDk3NDQsImlzcyI6IkJTRC1MSU5LIn0.DGqDz0YWO3RiqWUFOywVYkSOyImc3fDRtX9SvGpkINs"
// @Summary Similiarity Test Pengajuan
// @Description Similiarity Test Pengajuan
// @Produce application/json
// @Param request body models.SimilarityRequest true "Raw Request Body"
// @Tags Pengajuan
// @Param id query string false "id"
// @Success 200 {object} interface{}
// @Router /pengajuan/similarity-test [post]
func SimilartityTest(c *gin.Context) {
	var pengajuans []models.Pengajuan
	database.DB.Preload("Mahasiswa").Preload("DosPem1").Preload("DosPem2").Find(&pengajuans)
	repoJudul, err := utils.GetRepoTitlesFromDatabase()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error dari repository UPNVJ, silahkan dicoba beberapa saat lagi."})
		return
	}

	pengajuanId := c.Query("id")

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
		if pengajuanId != "" {
			id, _ := strconv.Atoi(pengajuanId)
			if id == pengajuan.Id {
				continue
			}
		}
		if pengajuan.StatusAcc != "Rejected" {
			continue
		}
		similarity := similarityPercentage(strings.ToLower(judul), strings.ToLower(pengajuan.Judul))
		if similarity > 30.0 {
			c.JSON(http.StatusOK, gin.H{
				"message":    "Judul serupa ditemukan",
				"similar":    pengajuan.Judul,
				"similarity": similarity,
			})
			return
		}
	}

	for _, judulRepo := range repoJudul {
		similarity := similarityPercentage(strings.ToLower(judul), strings.ToLower(judulRepo))
		if similarity > 30.0 {
			c.JSON(http.StatusOK, gin.H{
				"message":    "Judul serupa ditemukan",
				"similar":    judulRepo,
				"similarity": similarity,
			})
			return
		} //a
	}

	peminatan := utils.ClassifyTitle(judul)

	c.JSON(http.StatusOK, gin.H{"message": "Tidak ditemukan judul yang memiliki kesamaan diatas 30%", "peminatan": peminatan})
}
