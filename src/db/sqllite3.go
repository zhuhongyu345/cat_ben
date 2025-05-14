package db

import (
	"cat_ben/src/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var dbLite *gorm.DB

func InitDb() {
	dbPath := config.Config.DbPath
	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbLite = dbConn
}
