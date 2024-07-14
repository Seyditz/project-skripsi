package database

import (
	"github.com/Seyditz/project-skripsi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "postgres://koyeb-adm:hk9Dzl8csGRd@ep-bitter-snowflake-a2tym0ho.eu-central-1.pg.koyeb.app/koyebdb"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Admin{})
	db.AutoMigrate(&models.Dosen{})
	db.AutoMigrate(&models.Mahasiswa{})
	db.AutoMigrate(&models.Pengajuan{})
	db.AutoMigrate(&models.Judul{})
	DB = db
}

//aa
