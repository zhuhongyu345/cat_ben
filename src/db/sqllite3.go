package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dbPath = "/Users/bytedance/go/src/cat_ben/src/db/sqllite3.db"

var dbLite *gorm.DB

func InitDb() {
	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbLite = dbConn
}
