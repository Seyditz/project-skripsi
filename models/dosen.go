package models

import (
	"time"

	"gorm.io/gorm"
)

var DosenSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nidn, kapasitas, total_mahasiswa, prodi")
}

type Dosen struct {
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
	Image          string `json:"image"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type DosenDataResponse struct {
	Id             int    `json:"id" gorm:"primary_key"`
	Name           string `json:"name"`
	Nidn           string `json:"nidn"`
	Email          string `json:"email"`
	Prodi          string `json:"prodi"`
	Jabatan        string `json:"jabatan"`
	Kepakaran      string `json:"kepakaran"`
	Kapasitas      int    `json:"kapasitas"`
	TotalMahasiswa int    `json:"total_mahasiswa"`
	Image          string `json:"image"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type DosenCreateRequest struct {
	Name           string `json:"name"`
	Nidn           string `json:"nidn"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Prodi          string `json:"prodi"`
	Jabatan        string `json:"jabatan"`
	Kepakaran      string `json:"kepakaran"`
	Kapasitas      int    `json:"kapasitas"`
	TotalMahasiswa int    `json:"total_mahasiswa"`
	Image          string `json:"image"`
}

type DosenUpdateRequest struct {
	Name           string `json:"name"`
	Nidn           string `json:"nidn"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Prodi          string `json:"prodi"`
	Jabatan        string `json:"jabatan"`
	Kepakaran      string `json:"kepakaran"`
	Kapasitas      int    `json:"kapasitas"`
	TotalMahasiswa int    `json:"total_mahasiswa"`
	Image          string `json:"image"`
}

type DosenLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
