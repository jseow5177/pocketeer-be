package util

import (
	"math"
	"strconv"
)

func MonetaryStrToFloat(val string) (float64, error) {
	return StrToFloat(val, 2)
}

func StrToFloat(val string, dp int) (float64, error) {
	af, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	p := math.Pow(10, float64(dp))

	return math.Round(af*p) / p, nil
}
