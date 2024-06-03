package models

import "gorm.io/gorm"

var DosenSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nidn, kapasitas, total_mahasiswa, prodi")
}

type Dosen struct {
	gorm.Model
	Id             int    `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	Nidn           string `json:"nidn"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Prodi          string `json:"prodi"`
	Jabatan        string `json:"jabatan"`
	Kepakaran      string `json:"kepakaran"`
	Kapasitas      int    `json:"kapasitas"`
	TotalMahasiswa int    `json:"total_mahasiswa"`
}
