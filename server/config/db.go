package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error

	db, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_CONNECTION")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	d, _ := db.DB()
	d.Close()
}
