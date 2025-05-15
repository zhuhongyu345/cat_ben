package stock

import (
	"cat_ben/src/db"
	"log"
	"math"
	"strconv"
	"strings"
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
	if strings.Contains(stock.Name, ":") || strings.Contains(stock.Name, "-") || strings.Contains(stock.Name, "p") {
		_ = db.DeleteStoById(stock.ID)
		log.Println("delete one stock:" + stock.Name)
		return
	}
	basic, err := getWBBasic(stock.Name, stock.Mic)

	//respXQ, err := getDetailFromXQ(stock.Name)
	if err != nil {
		return
	}

	stock.PE = getFloat(basic.PeTtm)
	stock.Yield = getFloat(basic.Yield)
	stock.CHN = basic.Name
	stock.Price = getFloat(basic.Close)
	stock.H52 = getFloat(basic.FiftyTwoWkHigh)
	stock.L52 = getFloat(basic.FiftyTwoWkLow)
	stock.Lb = getFloat(basic.VibrateRatio)
	stock.Sz = getFloat(basic.MarketValue) / float64(100000000)
	stock.Hsl = getFloat(basic.TurnoverRate)
	stock.CB = basic.EstimateEarningsDate
	stock.PEF = getFloat(basic.ForwardPe)

	if err != nil {
		return
	}
	hl := (stock.Price - stock.L52) / (stock.H52 - stock.L52)
	if stock.H52 == stock.L52 {
		hl = 1
	}
	stock.Hl = math.Round(hl*10000) / 10000
	kData, cjl, _ := getKlineFromXQ(stock.Name, "day", 69)
	if len(kData) > 0 {
		ts := kData[len(kData)-1].Ts / 1000
		if time.Now().Unix()-ts > 86400*30 && ts > 0 {
			_ = db.DeleteStoById(stock.ID)
			log.Println("delete one stock:" + stock.Name)
			return
		}
	}
	stock.CjlD = math.Round(cjl*1000) / 1000
	acl := getZhicheng(kData, 0.009)
	stock.ZCL = math.Round(acl*10000) / 10000
	kDataW, _, _ := getKlineFromXQ(stock.Name, "week", 159)
	zcw := getZhicheng(kDataW, 0.012)
	stock.ZCW = math.Round(zcw*10000) / 10000

	err = db.UpdateStoById(stock)
	if err != nil {
		log.Println(err)
	}
	//if stock.Type == 2 {
	//	return
	//}
	//go func() {
	//	basic, err := getWBBasic(stock.Name, stock.Mic)
	//	if err != nil {
	//		log.Printf("wb basic err:%s,%s", stock.Name, err)
	//	}
	//	float, _ := strconv.ParseFloat(basic.ForwardPe, 64)
	//	_ = db.UpdateStoById(&db.Sto{
	//		ID:     stock.ID,
	//		Caibao: basic.EstimateEarningsDate,
	//		PEF:    float,
	//	})
	//}()
}
