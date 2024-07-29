package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const Local = true

// const dbPath = "C:/sqllite3.db"
var dbPath = "D:/workplace/cat_ben/src/db/sqllite3.db"

var dbLite *gorm.DB

func InitDb() {
	if !Local {
		dbPath = "C:/sqllite3.db"
	}
	dbConn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbLite = dbConn
}
