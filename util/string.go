package util

import (
	"math"
	"strconv"
	"strings"

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

	return roundFloat(f, dp), nil
}

func roundFloat(f float64, dp int) float64 {
	p := math.Pow(10, float64(dp))

	return math.Round(f*p) / p
}

func RoundFloatToStandardDP(f float64) float64 {
	return roundFloat(f, config.StandardDP)
}

func RoundFloatToPreciseDP(f float64) float64 {
	return roundFloat(f, config.PreciseDP)
}

func GetEmailPrefix(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}
