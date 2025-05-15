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
	Mic   string  `gorm:"column:mic"`
	Type  int     `gorm:"column:type"`
	CHN   string  `gorm:"column:chn"`
	Yield float64 `gorm:"column:yield"`
	PE    float64 `gorm:"column:pe"`
	PEF   float64 `gorm:"column:pef"`
	Price float64 `gorm:"column:price"`
	H52   float64 `gorm:"column:h52"`
	L52   float64 `gorm:"column:l52"`
	Hl    float64 `gorm:"column:hl"`
	Lb    float64 `gorm:"column:liangbi"`
	Sz    float64 `gorm:"column:shizhi"`
	Hsl   float64 `gorm:"column:huanshoulv"`
	Up    string  `gorm:"column:up"`
	CB    string  `gorm:"column:caibao"`
	CjlD  float64 `gorm:"column:cjlrateday"`
	ZCL   float64 `gorm:"column:zcrate"`
	ZCW   float64 `gorm:"column:zcweek"`
	TAG   int     `gorm:"column:tag"`
}

func CreateStos(stos []*Sto) error {
	err := dbLite.Table("stock_basic").Create(stos).Error
	return err
}

func DeleteALL() error {
	err := dbLite.Debug().Exec("delete from `stock_basic`").Error
	return err
}

func DeleteStoById(id int64) error {
	err := dbLite.Debug().Table("stock_basic").Where("id=?", id).Delete(new(Sto)).Error
	return err
}
func SelectStoByName(name string) (*Sto, error) {
	var sto Sto
	err := dbLite.Table("stock_basic").Where("name=?", name).Find(&sto).Error
	return &sto, err
}
func GetAllStockFromDB(hard string, tpe string) ([]*Sto, error) {
	resp := make([]*Sto, 0)
	query := dbLite.Debug().Table("stock_basic").Where("1=1")
	//A股刷新
	//query := dbLite.Debug().Table("stock_basic").Where("type=3")
	if tpe != "" {
		if t, err := strconv.ParseInt(tpe, 10, 64); err == nil {
			if t > 0 {
				query = query.Where("type=?", t)
			} else {
				query = query.Where("tag=?", 1)
			}
		}
	}
	if hard == "0" || hard == "" {
		query = query.Where("up!=?", time.Now().Format("2006-01-02"))
	}
	err := query.Find(&resp).Error
	return resp, err
}

type SearchDto struct {
	Name      string
	ZclLow    float64
	ZclHigh   float64
	CjlLow    float64
	CjlHigh   float64
	HlHigh    float64
	HlLow     float64
	PeHigh    float64
	PeLow     float64
	Yield     float64
	PriceLow  float64
	PriceHigh float64
	Liangbi   float64
	Tpe       int
	Skip      int
	Size      int
	Sort      string
	SortType  string
}

func Search(dto *SearchDto) ([]*Sto, error) {
	resp := make([]*Sto, 0)
	debug := dbLite.Debug()
	db := debug.Table("stock_basic").Where("1=1")

	if dto.Name != "" {
		err := db.Where("name=?", strings.ToUpper(dto.Name)).Find(&resp).Error
		return resp, err
	}

	if dto.Tpe != 0 {
		if dto.Tpe > 0 {
			db = db.Where("type=?", dto.Tpe)
		} else {
			db.Where("tag=?", 1)
		}
	}
	if dto.Liangbi != 0.0 {
		db = db.Where("liangbi>?", dto.Liangbi)
	}
	if dto.ZclLow != 0.0 {
		db = db.Where("zcrate>=?", dto.ZclLow)
	}
	if dto.ZclHigh != 0.0 {
		db = db.Where("zcrate<=?", dto.ZclHigh)
	}
	if dto.CjlLow != 0.0 {
		db = db.Where("cjlrateday>=?", dto.CjlLow)
	}
	if dto.CjlHigh != 0.0 {
		db = db.Where("cjlrateday<=?", dto.CjlHigh)
	}
	if dto.HlLow != 0.0 {
		db = db.Where("hl>=?", dto.HlLow)
	}
	if dto.HlHigh != 0.0 {
		db = db.Where("hl<=?", dto.HlHigh)
	}
	if dto.PeHigh != 0.0 {
		db = db.Where("pe<=?", dto.PeHigh)
	}
	if dto.PeLow != -10000 {
		db = db.Where("pe>=?", dto.PeLow)
	}
	if dto.Yield != 0.0 {
		db = db.Where("yield>=?", dto.Yield)
	}
	if dto.PriceLow != 0.0 {
		db = db.Where("price>=?", dto.PriceLow)
	}
	if dto.PriceHigh != 0.0 {
		db = db.Where("price<=?", dto.PriceHigh)
	}

	if dto.Sort != "" {
		if dto.SortType == "desc" {
			db = db.Order(fmt.Sprintf("%s %s", dto.Sort, "desc"))
		} else {
			db = db.Order(fmt.Sprintf("%s %s", dto.Sort, "asc"))
		}
	}
	if dto.Size != 0 {
		db = db.Limit(dto.Size)
	}
	if dto.Skip != 0 {
		db = db.Offset(dto.Skip)
	}
	err := db.Find(&resp).Error
	return resp, err
}

func UpdateStoById(sto *Sto) error {
	sto.Up = time.Now().Format("2006-01-02")
	err := dbLite.Table("stock_basic").Where("id=?", sto.ID).Updates(sto).Error
	return err
}

func UpdateByID(id int64, pe, yield float64, chn string, price, h52, l52, hl, liangbi, shizhi, huanshoulv, cjlrateday, zcrate, zcweek float64) error {
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
		"cjlrateday": math.Round(cjlrateday*10000) / 10000,
		"zcrate":     math.Round(zcrate*10000) / 10000,
		"zcweek":     math.Round(zcweek*10000) / 10000,
	}
	err := dbLite.Table("stock_basic").Where("id = ?", id).Updates(update).Error
	return err

}

func UpdateTagByID(id int64, tag string) error {
	update := map[string]interface{}{
		"tag": tag,
	}
	err := dbLite.Table("stock_basic").Where("id = ?", id).Updates(update).Error
	return err
}
