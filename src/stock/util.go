package stock

import "strconv"

func getFloat(s string) float64 {
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return float
}
