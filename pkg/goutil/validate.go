package goutil

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	ErrInvalidFloat = errors.New("invalid float")
)

func FormatFloat(f float64, decimalPlaces int) string {
	ft := fmt.Sprintf("%%.%df", decimalPlaces)
	return fmt.Sprintf(ft, f)
}

func IsFloat(str string, maxDecimalPoints int) error {
	regex := regexp.MustCompile(`^\d+\.\d{0,` + strconv.Itoa(maxDecimalPoints) + `}$`)
	if !regex.MatchString(str) {
		return ErrInvalidFloat
	}

	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return ErrInvalidFloat
	}

	return nil
}
