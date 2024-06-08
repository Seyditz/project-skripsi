package models

import (
	"time"
)

type Pengajuan struct {
	Id               int       `json:"id" gorm:"primary_key"`
	MahasiswaId      int       `json:"mahasiswa_id"`
	Mahasiswa        Mahasiswa `json:"mahasiswa" gorm:"foreignKey:MahasiswaId;references:Id"`
	Peminatan        string    `json:"peminatan"`
	Judul            string    `json:"judul"`
	TempatPenelitian string    `json:"tempat_penelitian"`
	RumusanMasalah   string    `json:"rumusan_masalah"`
	DosPem1Id        int       `json:"dospem1_id"`
	DosPem1          Dosen     `json:"dospem1" gorm:"foreignKey:DosPem1Id;references:Id"`
	DosPem2Id        int       `json:"dospem2_id"`
	DosPem2          Dosen     `json:"dospem2" gorm:"foreignKey:DosPem2Id;references:Id"`
	StatusAcc        string    `json:"status_acc"`
	RejectedNote     string    `json:"rejected_note"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PengajuanDataResponse struct {
	Id               int       `json:"id" gorm:"primary_key"`
	Mahasiswa        Mahasiswa `json:"mahasiswa" gorm:"foreignKey:MahasiswaId;references:Id"`
	Peminatan        string    `json:"peminatan"`
	Judul            string    `json:"judul"`
	TempatPenelitian string    `json:"tempat_penelitian"`
	RumusanMasalah   string    `json:"rumusan_masalah"`
	DosPem1          Dosen     `json:"dospem1" gorm:"foreignKey:DosPem1Id;references:Id"`
	DosPem2          Dosen     `json:"dospem2" gorm:"foreignKey:DosPem2Id;references:Id"`
	StatusAcc        string    `json:"status_acc"`
	RejectedNote     string    `json:"rejected_note"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PengajuanCreateRequest struct {
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	RumusanMasalah   string `json:"rumusan_masalah"`
	DosPem1Id        int    `json:"dospem1_id"`
	DosPem2Id        int    `json:"dospem2_id"`
	StatusAcc        string `json:"status_acc"`
	RejectedNote     string `json:"rejected_note"`
}

type PengajuanUpdateRequest struct {
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	RumusanMasalah   string `json:"rumusan_masalah"`
	DosPem1Id        int    `json:"dospem1_id"`
	DosPem2Id        int    `json:"dospem2_id"`
	StatusAcc        string `json:"status_acc"`
	RejectedNote     string `json:"rejected_note"`
}

type SimilarityRequest struct {
	Judul string `json:"judul"`
}
