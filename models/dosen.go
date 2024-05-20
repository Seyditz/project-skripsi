package models

import "gorm.io/gorm"

type Dosen struct {
	gorm.Model
	Id             int    `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	NIDN           string `json:"nidn"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Prodi          string `json:"prodi"`
	Kapasitas      int    `json:"kapasitas"`
	TotalMahasiswa int    `json:"total_mahasiswa"`
}
