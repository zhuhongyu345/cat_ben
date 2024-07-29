package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// const dbPath = "C:/sqllite3.db"
const dbPath = "D:/workplace/cat_ben/src/db/sqllite3.db"

var dbLite *gorm.DB

func InitDb() {
	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbLite = dbConn
}
