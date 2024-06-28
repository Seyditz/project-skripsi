package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var DosenSafePreloadFunction = func(db *gorm.DB) *gorm.DB {
	return db.Select("id, name, email, nidn, kapasitas, mahasiswa_bimbingan_id, prodi, jabatan, kepakaran")
}

type Dosen struct {
	Id                   int           `json:"id" gorm:"primary_key"`
	Name                 string        `json:"name"`
	Nidn                 string        `json:"nidn"`
	Email                string        `json:"email"`
	Password             string        `json:"password"`
	Prodi                string        `json:"prodi"`
	Jabatan              string        `json:"jabatan"`
	Kepakaran            string        `json:"kepakaran"`
	Kapasitas            int           `json:"kapasitas"`
	MahasiswaBimbinganId pq.Int64Array `json:"mahasiswa_bimbingan_id" gorm:"type:integer[]"`
	Image                string        `json:"image"`
	Gelar                string        `json:"gelar"`
	JenjangAkademik      string        `json:"jenjang_akademik"`
	NoTelp               string        `json:"no_telp"`
	TanggalLahir         time.Time     `json:"tanggal_lahir"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type DosenDataResponse struct {
	Id                   int           `json:"id" gorm:"primary_key"`
	Name                 string        `json:"name"`
	Nidn                 string        `json:"nidn"`
	Email                string        `json:"email"`
	Prodi                string        `json:"prodi"`
	Jabatan              string        `json:"jabatan"`
	Kepakaran            string        `json:"kepakaran"`
	Kapasitas            int           `json:"kapasitas"`
	MahasiswaBimbinganId pq.Int64Array `json:"mahasiswa_bimbingan_id"`
	Image                string        `json:"image"`
	Gelar                string        `json:"gelar"`
	NoTelp               string        `json:"no_telp"`
	JenjangAkademik      string        `json:"jenjang_akademik"`
	TanggalLahir         time.Time     `json:"tanggal_lahir"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type DosenCreateRequest struct {
	Name                 string        `json:"name"`
	Nidn                 string        `json:"nidn"`
	Email                string        `json:"email"`
	Password             string        `json:"password"`
	Prodi                string        `json:"prodi"`
	Jabatan              string        `json:"jabatan"`
	Kepakaran            string        `json:"kepakaran"`
	Kapasitas            int           `json:"kapasitas"`
	Gelar                string        `json:"gelar"`
	JenjangAkademik      string        `json:"jenjang_akademik"`
	NoTelp               string        `json:"no_telp"`
	TanggalLahir         time.Time     `json:"tanggal_lahir"`
	MahasiswaBimbinganId pq.Int64Array `json:"mahasiswa_bimbingan_id"`
}

type DosenUpdateRequest struct {
	Name                 string        `json:"name"`
	Nidn                 string        `json:"nidn"`
	Email                string        `json:"email"`
	Password             string        `json:"password"`
	Prodi                string        `json:"prodi"`
	Jabatan              string        `json:"jabatan"`
	Kepakaran            string        `json:"kepakaran"`
	Kapasitas            int           `json:"kapasitas"`
	Gelar                string        `json:"gelar"`
	NoTelp               string        `json:"no_telp"`
	JenjangAkademik      string        `json:"jenjang_akademik"`
	TanggalLahir         time.Time     `json:"tanggal_lahir"`
	MahasiswaBimbinganId pq.Int64Array `json:"mahasiswa_bimbingan_id"`
}

type DosenLoginRequest struct {
	Nidn     string `json:"nidn"`
	Password string `json:"password"`
}

type DosenMahasiswaBimbinganResponse struct {
	MahasiswaList []MahasiswaDataResponse `json:"mahasiswa_list"`
}
