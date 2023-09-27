package entity

import (
	"errors"
	"net/mail"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrInvalidBudgetRepeat     = errutil.ValidationError(errors.New("invalid budget repeat"))
	ErrInvalidCurrency         = errutil.ValidationError(errors.New("invalid currency"))
	ErrInvalidEmail            = errutil.ValidationError(errors.New("invalid email"))
	ErrInvalidDate             = errutil.ValidationError(errors.New("invalid date"))
	ErrInvalidTimezone         = errutil.ValidationError(errors.New("invalid timezone"))
	ErrInvalidHoldingType      = errutil.ValidationError(errors.New("invalid holding type"))
	ErrInvalidChildAccountType = errutil.ValidationError(errors.New("invalid child account type"))
	ErrInvalidAccountType      = errutil.ValidationError(errors.New("invalid account type"))
	ErrInvalidTransactionType  = errutil.ValidationError(errors.New("invalid transaction type"))
	ErrInvalidCategoryType     = errutil.ValidationError(errors.New("invalid category type"))
	ErrInvalidBudgetType       = errutil.ValidationError(errors.New("invalid budget type"))
	ErrInvalidMonetaryStr      = errutil.ValidationError(errors.New("invalid monetary str"))
	ErrInvalidTransactionSumBy = errutil.ValidationError(errors.New("invalid transactions sum by"))
	ErrInvalidBudgetRepeats    = errutil.ValidationError(errors.New("invalid budget periods"))
	ErrMustBePositive          = errutil.ValidationError(errors.New("must be positive"))
)

func CheckCurrency(currency string) error {
	if _, ok := Currencies[currency]; !ok {
		return ErrInvalidCurrency
	}
	return nil
}

func CheckEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func CheckBudgetRepeat(budgetRepeat uint32) error {
	if _, ok := BudgetRepeats[budgetRepeat]; !ok {
		return ErrInvalidBudgetRepeat
	}
	return nil
}

func CheckDateStr(date string) error {
	if _, err := util.ParseDate(date); err != nil {
		return ErrInvalidDate
	}
	return nil
}

func CheckTimezone(timezone string) error {
	if _, err := time.LoadLocation(timezone); err != nil {
		return ErrInvalidTimezone
	}
	return nil
}

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
	if _, ok := TransactionTypes[transactionType]; !ok {
		return ErrInvalidTransactionType
	}
	return nil
}

func CheckBudgetType(budgetType uint32) error {
	if _, ok := BudgetTypes[budgetType]; !ok {
		return ErrInvalidTransactionType
	}
	return nil
}

func CheckMonetaryStr(str string) error {
	if _, err := util.MonetaryStrToFloat(str); err != nil {
		return ErrInvalidMonetaryStr
	}
	return nil
}

func CheckPositiveMonetaryStr(str string) error {
	f, err := util.MonetaryStrToFloat(str)
	if err != nil {
		return ErrInvalidMonetaryStr
	}
	if f < 0 {
		return ErrMustBePositive
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

func AppMetaValidator() validator.Validator {
	return validator.MustForm(map[string]validator.Validator{
		"timezone": &validator.String{
			Optional:   false,
			Validators: []validator.StringFunc{CheckTimezone},
		},
	})
}
