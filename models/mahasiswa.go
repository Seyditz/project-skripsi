package models

import "gorm.io/gorm"

var MahasiswaSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nim, prodi, angkatan, sks, image")
}

type Mahasiswa struct {
	gorm.Model
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
	SKS      int    `json:"sks"`
	Image    string `json:"image"`
}

type MahasiswaDataResponse struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Email    string `json:"email"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
	SKS      int    `json:"sks"`
	Image    string `json:"image"`
}
