package stock

import (
	"cat_ben/src/db"
	"log"
	"math"
	"strconv"
	"time"
)

// from https://xueqiu.com/S/INTC

var sleep = 0

func getSleep() int {
	if sleep > 0 {
		return sleep
	}
	value, err := db.GetValue(`sleep`)
	if err != nil {
		sleep = 1000
	}
	atoi, err := strconv.Atoi(value)
	if err != nil {
		sleep = 1000
	}
	sleep = atoi
	log.Printf("sleep:%d", sleep)
	return sleep
}

func FlushBasic(hard string, tpe string) {
	log.Println("flush start")
	defer func() {
		if r := recover(); r != nil {
			log.Println("recover:", r)
		}
	}()
	stocks, _ := db.GetAllStockFromDB(hard, tpe)
	begin := time.Now().Unix()
	count := 0
	for _, stock := range stocks {
		time.Sleep(time.Duration(getSleep()) * time.Millisecond)
		FlushOne(stock)
		count++
	}

	log.Print("flush finish cost:")
	log.Println(time.Now().Unix() - begin)
	log.Print("count:")
	log.Println(count)
}

func FlushOne(stock *db.Sto) {
	pe, yield, chn, price, h52, l52, liangbi, shizhi, huanshoulv, err := getDetailFromXQ(stock.Name)
	if err != nil {
		return
	}
	hl := (price - l52) / (h52 - l52)
	if h52 == l52 {
		hl = 1
	}
	hl = math.Round(hl*10000) / 10000
	kData, cjl, _ := getKlineFromXQ(stock.Name, "day", 69)
	if len(kData) > 0 {
		ts := kData[len(kData)-1].Ts / 1000
		if time.Now().Unix()-ts > 86400*30 && ts > 0 {
			err := db.DeleteStoById(stock.ID)
			if err != nil {
				log.Printf("delete one stock err,%s:%s", stock.Name, err)
			}
			log.Println("delete one stock:" + stock.Name)
			return
		}
	}

	zcl := getZhicheng(kData, 0.009)
	kDataW, _, _ := getKlineFromXQ(stock.Name, "week", 159)
	zcw := getZhicheng(kDataW, 0.012)

	err = db.UpdateByID(stock.ID, pe, yield, chn, price, h52, l52, hl, liangbi, shizhi, huanshoulv, cjl, zcl, zcw)
	if err != nil {
		log.Println(err)
	}
	if stock.Type == 2 {
		return
	}
	go func() {
		basic, err := getWBBasic(stock.Name, stock.Mic)
		if err != nil {
			log.Printf("wb basic err:%s,%s", stock.Name, err)
		}
		float, _ := strconv.ParseFloat(basic.ForwardPe, 64)
		_ = db.UpdateStoById(&db.Sto{
			ID:     stock.ID,
			Caibao: basic.EstimateEarningsDate,
			PEF:    float,
		})
	}()
}
