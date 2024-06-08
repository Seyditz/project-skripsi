package models

import (
	"time"

	"gorm.io/gorm"
)

var MahasiswaSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nim, prodi, angkatan, sks, image")
}

type Mahasiswa struct {
	Id        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	NIM       string `json:"nim"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Prodi     string `json:"prodi"`
	Angkatan  int    `json:"angkatan"`
	SKS       int    `json:"sks"`
	Image     string `json:"image"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MahasiswaDataResponse struct {
	Id        int    `json:"id" gorm:"primary_key"`
	Name      string `json:"name"`
	NIM       string `json:"nim"`
	Email     string `json:"email"`
	Prodi     string `json:"prodi"`
	Angkatan  int    `json:"angkatan"`
	SKS       int    `json:"sks"`
	Image     string `json:"image"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MahasiswaCreateRequest struct {
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Email    string `json:"email"`
	Prodi    string `json:"prodi"`
	Password string `json:"password"`
	Angkatan int    `json:"angkatan"`
	SKS      int    `json:"sks"`
	Image    string `json:"image"`
}

type MahasiswaUpdateRequest struct {
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Email    string `json:"email"`
	Prodi    string `json:"prodi"`
	Password string `json:"password"`
	Angkatan int    `json:"angkatan"`
	SKS      int    `json:"sks"`
	Image    string `json:"image"`
}
