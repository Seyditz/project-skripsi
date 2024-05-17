package models

import "gorm.io/gorm"

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
}

type MahasiswaDataResponse struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	NIM      string `json:"nim"`
	Email    string `json:"email"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
	SKS      int    `json:"sks"`
}
