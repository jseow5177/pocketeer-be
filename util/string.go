package util

import (
	"math"
	"strconv"

	"github.com/jseow5177/pockteer-be/config"
)

func MonetaryStrToFloat(val string) (float64, error) {
	return StrToFloat(val, config.StandardDP)
}

func StrToFloat(val string, dp int) (float64, error) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	return RoundFloat(f, dp), nil
}

func RoundFloat(f float64, dp int) float64 {
	p := math.Pow(10, float64(dp))

	return math.Round(f*p) / p
}

func RoundFloatToStandardDP(f float64) float64 {
	return RoundFloat(f, config.StandardDP)
}
