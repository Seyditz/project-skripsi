package database

import (
	"github.com/Seyditz/project-skripsi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "postgres://crotdan:0GBMVxbNz9Wl@ep-floral-dream-a2znj19n.eu-central-1.pg.koyeb.app/koyebdb"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Admin{})
	DB = db
}
