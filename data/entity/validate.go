package entity

import "errors"

var (
	ErrInvalidCatType = errors.New("invalid cat type")
)

func CheckCatType(catType uint32) error {
	if _, ok := CatTypes[catType]; ok {
		return nil
	}
	return ErrInvalidCatType
}
