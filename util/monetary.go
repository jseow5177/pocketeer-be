package util

import (
	"math"
	"strconv"

	"github.com/jseow5177/pockteer-be/config"
)

const (
	defaultAmount = 0
)

func MonetaryStrToFloat(val string) float64 {
	af, err := strconv.ParseFloat(val, 64)
	if err != nil {
		af = 0
	}

	p := math.Pow(10, float64(config.AmountDecimalPlaces))
	return math.Round(af*p) / p
}
