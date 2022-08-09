package db

import (
	"fmt"
	"math"
	"strconv"
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
}

func CreateStos(stos []*Sto) error {
	if dbLite == nil {
		InitDb()
	}
	err := dbLite.Table("stock_basic").Create(stos).Error
	return err
}

func GetAllStockFromDB() ([]*Sto, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*Sto, 0)
	err := dbLite.Table("stock_basic").Where("1=1 and (hl is null or hl = -1 or shizhi = -1)").Find(&resp).Error
	return resp, err
}

func Search(hlLow, hlHigh, pe, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType string) ([]*Sto, error) {
	if dbLite == nil {
		InitDb()
	}
	resp := make([]*Sto, 0)
	debug := dbLite.Debug()
	db := debug.Table("stock_basic").Where("1=1")
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
	if pe != "" {
		if float, err := strconv.ParseFloat(pe, 64); err == nil {
			db = db.Where("pe<=?", float)
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
	}
	err := dbLite.Table("stock_basic").Where("id = ?", id).Updates(update).Error
	return err

}
