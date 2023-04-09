package configs

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func ConnectDB() *gorm.DB {
	dsn := EnvMySQL()
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	return db

}

// Client instance
var DB *gorm.DB = ConnectDB()
