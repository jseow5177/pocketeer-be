package util

import (
	"math"
	"strconv"

	"github.com/jseow5177/pockteer-be/config"
)

func MonetaryStrToFloat(val string) (float64, error) {
	af, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}

	p := math.Pow(10, float64(config.AmountDecimalPlaces))

	return math.Round(af*p) / p, nil
}
