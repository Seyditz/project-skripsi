package models

import (
	"time"

	"gorm.io/gorm"
)

type Judul struct {
	gorm.Model
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	Abstrak          string `json:"abstrak"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type JudulCreateRequest struct {
	gorm.Model
	Id               int    `json:"id" gorm:"primary_key"`
	MahasiswaId      int    `json:"mahasiswa_id"`
	Peminatan        string `json:"peminatan"`
	Judul            string `json:"judul"`
	TempatPenelitian string `json:"tempat_penelitian"`
	Abstrak          string `json:"abstrak"`
}
