package stock

import (
	"cat_ben/src/db"
	"math"
	"strconv"
	"time"
)

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
	byName, _ := db.SelectStoByName(name)
	resp := make(map[string]interface{})
	resp["basic"] = byName
	resp["pes"] = pes
	resp["prices"] = prices
	resp["times"] = times
	return resp
}
