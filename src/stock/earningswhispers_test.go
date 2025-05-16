package stock

import (
	"cat_ben/src/db"
	"encoding/json"
	"testing"
	"time"
)

func TestFlushAllIcon(t *testing.T) {
	fromDB, err := db.GetAllStockFromDB("1", "")
	if err != nil {
		return
	}
	for _, v := range fromDB {
		SaveImgLocal(v.Name)
	}

}

func TestSaveImgLocal(t *testing.T) {
	SaveImgLocal("AAPL")
}

func TestGetByDate(t *testing.T) {
	date, err := GetByDate(time.Now().Format("20060102"))
	if err != nil {
		return
	}

	marshal, _ := json.Marshal(date)
	t.Log(string(marshal))

}
