package db

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type Sto struct {
	ID    int64   `gorm:"column:id"`
	Name  string  `gorm:"column:name"`
	Type  int     `gorm:"column:type"`
	CHN   string  `gorm:"column:chn"`
	Yield float64 `gorm:"column:yield"`
	PE    float64 `gorm:"column:pe"`
	Price float64 `gorm:"column:price"`
	H52   float64 `gorm:"column:h52"`
	L52   float64 `gorm:"column:l52"`
	Hl    float64 `gorm:"column:hl"`
	Lb    float64 `gorm:"column:liangbi"`
	Sz    float64 `gorm:"column:shizhi"`
	Hsl   float64 `gorm:"column:huanshoulv"`
	Up    string  `gorm:"column:up"`
}

func CreateStos(stos []*Sto) error {
	if dbLite == nil {
		InitDb()
	}
	err := dbLite.Table("stock_basic").Create(stos).Error
	return err
}

func GetAllStockFromDB(hard string) ([]*Sto, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*Sto, 0)
	query := dbLite.Debug().Table("stock_basic").Where("1=1")
	if hard == "0" || hard == "" {
		query = query.Where("up!=?", time.Now().Format("2006-01-02"))
	}
	err := query.Find(&resp).Error
	return resp, err
}

func Search(name, hlLow, hlHigh, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType string) ([]*Sto, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*Sto, 0)
	debug := dbLite.Debug()
	db := debug.Table("stock_basic").Where("1=1")

	if name != "" {
		err := db.Where("name=?", strings.ToUpper(name)).Find(&resp).Error
		return resp, err
	}

	if tpe != "" {
		if float, err := strconv.ParseInt(tpe, 10, 64); err == nil {
			db = db.Where("type=?", float)
		}
	}
	if liangbi != "" {
		if float, err := strconv.ParseFloat(liangbi, 64); err == nil {
			db = db.Where("liangbi>?", float)
		}
	}
	if hlLow != "" {
		if float, err := strconv.ParseFloat(hlLow, 64); err == nil {
			db = db.Where("hl>=?", float)
		}
	}
	if hlHigh != "" {
		if float, err := strconv.ParseFloat(hlHigh, 64); err == nil {
			db = db.Where("hl<=?", float)
		}
	}
	if peHigh != "" {
		if float, err := strconv.ParseFloat(peHigh, 64); err == nil {
			db = db.Where("pe<=?", float)
		}
	}
	if peLow != "" {
		if float, err := strconv.ParseFloat(peLow, 64); err == nil {
			db = db.Where("pe>=?", float)
		}
	}
	if yield != "" {
		if float, err := strconv.ParseFloat(yield, 64); err == nil {
			db = db.Where("yield>=?", float)
		}
	}
	if priceLow != "" {
		if float, err := strconv.ParseFloat(priceLow, 64); err == nil {
			db = db.Where("price>=?", float)
		}
	}
	if priceHigh != "" {
		if float, err := strconv.ParseFloat(priceHigh, 64); err == nil {
			db = db.Where("price<=?", float)
		}
	}

	if sort != "" {
		if sortType == "desc" {
			db = db.Order(fmt.Sprintf("%s %s", sort, "desc"))
		} else {
			db = db.Order(fmt.Sprintf("%s %s", sort, "asc"))
		}
	}
	if size != "" {
		parseInt, _ := strconv.ParseInt(size, 10, 64)
		db = db.Limit(int(parseInt))
	}
	if skip != "" {
		parseInt, _ := strconv.ParseInt(skip, 10, 64)
		db = db.Offset(int(parseInt))
	}
	err := db.Find(&resp).Error
	return resp, err
}

func UpdateByID(id int64, pe, yield float64, chn string, price, h52, l52, hl, liangbi, shizhi, huanshoulv float64) error {
	if dbLite == nil {
		InitDb()
	}
	up := time.Now().Format("2006-01-02")
	update := map[string]interface{}{
		"pe":         math.Round(pe*10000) / 10000,
		"yield":      math.Round(yield*10000) / 10000,
		"chn":        chn,
		"price":      math.Round(price*10000) / 10000,
		"h52":        math.Round(h52*10000) / 10000,
		"l52":        math.Round(l52*10000) / 10000,
		"hl":         math.Round(hl*10000) / 10000,
		"liangbi":    math.Round(liangbi*10000) / 10000,
		"shizhi":     math.Round(shizhi*10000) / 10000,
		"huanshoulv": math.Round(huanshoulv*10000) / 10000,
		"up":         up,
	}
	err := dbLite.Table("stock_basic").Where("id = ?", id).Updates(update).Error
	return err

}
