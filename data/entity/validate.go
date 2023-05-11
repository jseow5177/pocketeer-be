package entity

import (
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidCategoryType    = errors.New("invalid category type")
	ErrInvalidAmount          = errors.New("invalid amount")
)

func CheckCategoryType(categoryType uint32) error {
	if err := CheckTransactionType(categoryType); err != nil {
		return ErrInvalidCategoryType
	}
	return nil
}

func CheckTransactionType(transactionType uint32) error {
	if _, ok := TransactionTypes[transactionType]; ok {
		return nil
	}
	return ErrInvalidTransactionType
}

func CheckAmount(amount string) error {
	if err := goutil.IsFloat(amount, config.AmountDecimalPlaces); err != nil {
		return ErrInvalidAmount
	}
	return nil
}
