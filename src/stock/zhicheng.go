package stock

import (
	"math"
)

type ZCPoint struct {
	KlineData *KlineData
	Price     float64
	Effect    float64
}

func getZhicheng(kData []*KlineData, ZCL_DIFF_PERCENT float64) float64 {

	if len(kData) == 0 {
		return -1
	}
	lows := make([]float64, 0)
	for _, v := range kData {
		lows = append(lows, v.Low)
	}

	zhicheng := make([]*ZCPoint, 0)
	temp := float64(100000)
	up := true
	//5432 1 2345
	//125.5993 124.9821 130.2088 129.258 131.6623 130.8658 130.6767 129.6015 149.1936
	maxPrice := float64(0)
	minPrice := float64(10000)
	for _, v := range kData {
		//找支撑位
		if v.Low < temp {
			up = false
		} else if v.Low > temp {
			if !up {
				zhicheng = append(zhicheng, &ZCPoint{Price: temp, Effect: 0.1, KlineData: v})
			}
			up = true
		}
		temp = v.Low
		//找最高最低价
		if v.Low < minPrice {
			minPrice = v.Low
		}
		if v.High > maxPrice {
			maxPrice = v.High
		}
	}

	hight := maxPrice - minPrice
	for _, z := range zhicheng {
		z.Effect = (maxPrice - z.Price) / hight
		z.Effect *= z.Effect
	}

	for _, i := range zhicheng {
		orgEff := i.Effect
		maxeff := orgEff
		for _, j := range zhicheng {
			if math.Abs(i.Price-j.Price)/i.Price < ZCL_DIFF_PERCENT {
				//maxeff += math.Abs(float64(iidx-jidx)) * orgEff
				maxeff += orgEff
			}
		}
		i.Effect = maxeff
	}

	closePrice := kData[len(kData)-1].Close
	zhichenglv := float64(0)
	for _, v := range zhicheng {
		diff := math.Abs(closePrice - v.Price)
		if diff < closePrice*ZCL_DIFF_PERCENT {
			zhichenglv += v.Effect
		}
	}
	//log.Println(zhichenglv)
	return zhichenglv
}
