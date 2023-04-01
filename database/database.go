package database

import (
	"fmt"
	"log"
	"rakamin/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


type ConfigDB struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func (config *ConfigDB) ConfigDB() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().AutoMigrate(
		&models.User{},
		&models.Photo{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
