package server

import "strconv"

func getFloat(s string) float64 {
	float, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return float
}

func getInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
