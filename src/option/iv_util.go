package option

import "math"

func GetVOL(isCall bool, stockPrice, actPrice, optionPrice, expireDays float64) float64 {

	low := float64(0)
	vol := float64(50)
	high := float64(100)

	for {
		tempPrice := GetPriceBlackScholes(isCall, stockPrice, actPrice, expireDays, vol)
		if tempPrice > optionPrice {
			high = vol
		} else if tempPrice < optionPrice {
			low = vol
		}
		vol = (low + high) / 2
		if high-low < 0.0001 {
			break
		}
	}
	return math.Round(vol*10000.0) / 10000.0
}

func GetPriceBlackScholes(isCall bool, stockPrice, actPrice, expireDays, volatility float64) float64 {
	T := expireDays / float64(365)
	riskFreeRate := float64(0)
	d1 := (math.Log(stockPrice/actPrice) + (riskFreeRate+volatility*volatility/2)*T) / (volatility * math.Sqrt(T))
	d2 := d1 - volatility*math.Sqrt(T)
	if isCall {
		return stockPrice*CND(d1) - actPrice*math.Exp(-riskFreeRate*T)*CND(d2)
	} else {
		return actPrice*math.Exp(-riskFreeRate*T)*CND(-d2) - stockPrice*CND(-d1)
	}
}

func CND(X float64) float64 {
	var L, K, w float64
	a1 := 0.31938153
	a2 := -0.356563782
	a3 := 1.781477937
	a4 := -1.821255978
	a5 := 1.330274429
	L = math.Abs(X)
	K = 1.0 / (1.0 + 0.2316419*L)
	w = 1.0 - 1.0/math.Sqrt(2.0*math.Pi)*math.Exp(-L*L/2)*(a1*K+a2*K*K+a3*math.Pow(K, 3)+a4*math.Pow(K, 4)+a5*math.Pow(K, 5))
	if X < 0.0 {
		w = 1.0 - w
	}
	return w
}
