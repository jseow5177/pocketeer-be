package entity

import (
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
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
	if err := goutil.IsFloat(amount, AmountDecimalPlaces); err != nil {
		return ErrInvalidAmount
	}
	return nil
}

func PagingValidator(optional bool) validator.Validator {
	return &validator.Form{
		Optional: optional,
		Validators: map[string]validator.Validator{
			"limit": &validator.UInt32{
				Optional:  true,
				UnsetZero: true,
				Max:       goutil.Uint32(config.MaxPagingLimit),
			},
			"page": &validator.UInt32{
				Optional:  true,
				UnsetZero: true,
				Min:       goutil.Uint32(config.MinPagingPage),
			},
		},
	}
}

func UInt64FilterValidator(optional bool) validator.Validator {
	return &validator.Form{
		Optional: optional,
		Validators: map[string]validator.Validator{
			"gte": &validator.UInt64{
				Optional:  true,
				UnsetZero: true,
			},
			"lte": &validator.UInt64{
				Optional:  true,
				UnsetZero: true,
			},
		},
	}
}
