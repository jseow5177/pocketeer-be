package entity

import (
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrInvalidHoldingType      = errors.New("invalid holding type")
	ErrInvalidChildAccountType = errors.New("invalid child account type")
	ErrInvalidAccountType      = errors.New("invalid account type")
	ErrInvalidTransactionType  = errors.New("invalid transaction type")
	ErrInvalidCategoryType     = errors.New("invalid category type")
	ErrInvalidBudgetType       = errors.New("invalid budget type")
	ErrMonetaryStr             = errors.New("invalid monetary str")
	ErrInvalidTransactionSumBy = errors.New("invalid transactions sum by")
	ErrNegativeMonetaryStr     = errors.New("monetary str cannot be negative")
)

func CheckHoldingType(holdingType uint32) error {
	if _, ok := HoldingTypes[holdingType]; !ok {
		return ErrInvalidHoldingType
	}
	return nil
}

func CheckAccountType(accountType uint32) error {
	_, isParent := ParentAccountTypes[accountType]
	_, isChild := ChildAccountTypes[accountType]
	if !isParent && !isChild {
		return ErrInvalidAccountType
	}
	return nil
}

func CheckChildAccountType(accountType uint32) error {
	if _, ok := ChildAccountTypes[accountType]; !ok {
		return ErrInvalidChildAccountType
	}
	return nil
}

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

func CheckBudgetType(budgetType uint32) error {
	if _, ok := BudgetTypes[budgetType]; ok {
		return nil
	}
	return ErrInvalidTransactionType
}

func CheckMonetaryStr(str string) error {
	if err := goutil.IsFloat(str, config.AmountDecimalPlaces); err != nil {
		return ErrMonetaryStr
	}
	return nil
}

func CheckPositiveMonetaryStr(str string) error {
	f, err := util.MonetaryStrToFloat(str)
	if err != nil {
		return ErrMonetaryStr
	}
	if f < 0 {
		return ErrNegativeMonetaryStr
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
