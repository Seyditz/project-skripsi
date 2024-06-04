package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "postgres://koyeb-adm:5Ki1EOJRaDhj@ep-floral-dream-a2znj19n.eu-central-1.pg.koyeb.app/koyebdb"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&models.Admin{})
	// db.AutoMigrate(&models.Dosen{})
	// db.AutoMigrate(&models.Mahasiswa{})
	// db.AutoMigrate(&models.Pengajuan{})
	DB = db
}
