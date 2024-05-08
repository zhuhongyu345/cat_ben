package stock

import (
	"cat_ben/src/db"
	"log"
	"math"
	"strconv"
	"time"
)

func Search(name, zclLow, zclHigh, cjlLow, cjlHigh, hlLow, hlHight, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType string) []*db.Sto {
	stocks, err := db.Search(name, zclLow, zclHigh, cjlLow, cjlHigh, hlLow, hlHight, peHigh, peLow, yield, priceLow, priceHigh, liangbi, tpe, skip, size, sort, sortType)
	if err != nil {
		log.Printf("db.Search err:%s", err)
	}
	return stocks
}

func GetHistory(name string, period string, count string) map[string]interface{} {
	count64, _ := strconv.ParseInt(count, 10, 64)
	kData, _, _ := getKlineFromXQ(name, period, count64)
	times := make([]string, 0)
	pes := make([]float64, 0)
	prices := make([]float64, 0)
	for _, v := range kData {
		times = append(times, time.UnixMilli(v.Ts).Format("2006-01-02"))
		pes = append(pes, math.Round(v.Pe*10000)/10000)
		prices = append(prices, math.Round(v.Close*10000)/10000)
	}
	resp := make(map[string]interface{})
	resp["pes"] = pes
	resp["prices"] = prices
	resp["times"] = times
	return resp

}
