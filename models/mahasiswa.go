package models

import (
	"time"

	"gorm.io/gorm"
)

var MahasiswaSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nim, prodi, angkatan, sks, image, agama, jenis_kelamin, no_telp")
}

type Mahasiswa struct {
	Id           int       `json:"id" gorm:"primary_key"`
	Name         string    `json:"name"`
	NIM          string    `json:"nim"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Prodi        string    `json:"prodi"`
	Angkatan     int       `json:"angkatan"`
	SKS          int       `json:"sks"`
	Image        string    `json:"image"`
	NoTelp       string    `json:"no_telp"`
	TempatLahir  string    `json:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Agama        string    `json:"agama"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MahasiswaDataResponse struct {
	Id           int       `json:"id" gorm:"primary_key"`
	Name         string    `json:"name"`
	NIM          string    `json:"nim"`
	Email        string    `json:"email"`
	Prodi        string    `json:"prodi"`
	Angkatan     int       `json:"angkatan"`
	SKS          int       `json:"sks"`
	Image        string    `json:"image"`
	NoTelp       string    `json:"no_telp"`
	TempatLahir  string    `json:"tempat_lahir"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string    `json:"jenis_kelamin"`
	Agama        string    `json:"agama"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type MahasiswaCreateRequest struct {
	Name         string `json:"name"`
	NIM          string `json:"nim"`
	Email        string `json:"email"`
	Prodi        string `json:"prodi"`
	Password     string `json:"password"`
	Angkatan     int    `json:"angkatan"`
	TempatLahir  string `json:"tempat_lahir"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_lahir"`
	SKS          int    `json:"sks"`
	JenisKelamin string `json:"jenis_kelamin"`
	Agama        string `json:"agama"`
}

type MahasiswaUpdateRequest struct {
	Name         string `json:"name"`
	NIM          string `json:"nim"`
	Email        string `json:"email"`
	Prodi        string `json:"prodi"`
	Password     string `json:"password"`
	Angkatan     int    `json:"angkatan"`
	NoTelp       string `json:"no_telp"`
	TempatLahir  string `json:"tempat_lahir"`
	TanggalLahir string `json:"tanggal_lahir"`
	SKS          int    `json:"sks"`
	JenisKelamin string `json:"jenis_kelamin"`
	Agama        string `json:"agama"`
}

type MahasiswaLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MahasiswaLoginResponse struct {
	Data  MahasiswaDataResponse `json:"data"`
	Roles []string              `json:"roles"`
	Token string                `json:"token"`
}
