package models

import "gorm.io/gorm"

type Judul struct {
	gorm.Model
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	RumusanMasalah   string `json:"rumusan_masalah"`
	DosPem1Id        int    `json:"dospem1_id"`
	DosPem1          Dosen  `json:"dospem1" gorm:"foreignKey:DosPem1Id;references:Id"`
	DosPem2Id        int    `json:"dospem2_id"`
	DosPem2          Dosen  `json:"dospem2" gorm:"foreignKey:DosPem2Id;references:Id"`
}

type JudulCreateRequest struct {
	gorm.Model
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	RumusanMasalah   string `json:"rumusan_masalah"`
	DosPem1Id        int    `json:"dospem1_id"`
	DosPem2Id        int    `json:"dospem2_id"`
}
