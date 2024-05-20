package models

import "gorm.io/gorm"

type Pengajuan struct {
	gorm.Model
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	RumusanMasalah   string `json:"rumusan_masalah"`
	DosPem1          string `json:"dospem1"`
	DosPem2          string `json:"dospem2"`
	StatusAcc        string `json:"status_acc"`
}
