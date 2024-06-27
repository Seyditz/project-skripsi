package utils

import (
	"fmt"

	"github.com/Seyditz/project-skripsi/database"
	"github.com/Seyditz/project-skripsi/models"
)

func RemoveMahasiswaBimbinganFromDosen(mahasiswaId int, dospemId int) error {
	var dospem models.Dosen

	if result := database.DB.First(&dospem, dospemId); result.RowsAffected == 0 {
		return fmt.Errorf("dosen pembimbing dengan ID %d tidak ditemukan", dospemId)
	}

	dospemMahasiswaArray := RemoveInt64FromArray(dospem.MahasiswaBimbinganId, int64(mahasiswaId))

	dospemUpdatedData := models.Dosen{
		MahasiswaBimbinganId: dospemMahasiswaArray,
	}

	if result := database.DB.Model(&models.Dosen{}).Where("id = ?", dospemId).Updates(dospemUpdatedData); result.Error != nil {
		return fmt.Errorf("gagal menghapus mahasiswa bimbingan dalam dosen pembimbing dengan ID %d", dospemId)
	}

	return nil
}

func AddMahasiswaBimbinganToDosen(mahasiswaId int, dospemId int) error {
	var dospem models.Dosen

	if result := database.DB.First(&dospem, dospemId); result.RowsAffected == 0 {
		return fmt.Errorf("dosen pembimbing dengan ID %d tidak ditemukan", dospemId)
	}

	newDospem1UpdatedData := models.Dosen{
		MahasiswaBimbinganId: append(dospem.MahasiswaBimbinganId, int64(mahasiswaId)),
	}

	if result := database.DB.Model(&models.Dosen{}).Where("id = ?", dospem.Id).Updates(newDospem1UpdatedData); result.Error != nil {
		return fmt.Errorf("gagal mengupdate mahasiswa bimbingan dosen pembimbing dengan ID %d", dospemId)
	}

	return nil
}
