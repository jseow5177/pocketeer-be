package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepo interface {
	Get(ctx context.Context, af *AccountFilter) (*entity.Account, error)
	GetMany(ctx context.Context, acf *AccountFilter) ([]*entity.Account, error)

	Create(ctx context.Context, ac *entity.Account) (string, error)
	Update(ctx context.Context, acf *AccountFilter, acu *entity.AccountUpdate) error
}

type AccountFilter struct {
	UserID            *string  `filter:"user_id"`
	AccountID         *string  `filter:"_id"`
	AccountType       *uint32  `filter:"account_type"`
	AccountTypeBitPos []uint32 `filter:"account_type__bitsAllSet"` // bit pos
}

type AccountFilterOption = func(acf *AccountFilter)

func WitAccountID(accountID *string) AccountFilterOption {
	return func(acf *AccountFilter) {
		acf.AccountID = accountID
	}
}

func WitAccountType(accountType *uint32) AccountFilterOption {
	return func(acf *AccountFilter) {
		if accountType == nil {
			return
		}
		if _, ok := entity.ParentAccountTypes[*accountType]; ok {
			acf.AccountTypeBitPos = goutil.GetSetBitPos(*accountType << entity.AccountTypeBitShift)
		}
		if _, ok := entity.ChildAccountTypes[*accountType]; ok {
			acf.AccountType = accountType
		}
	}
}

func NewAccountFilter(userID string, opts ...AccountFilterOption) *AccountFilter {
	acf := new(AccountFilter)
	for _, opt := range opts {
		opt(acf)
	}
	return acf
}

func (f *AccountFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *AccountFilter) GetAccountID() string {
	if f != nil && f.AccountID != nil {
		return *f.AccountID
	}
	return ""
}

func (f *AccountFilter) GetAccountType() uint32 {
	if f != nil && f.AccountType != nil {
		return *f.AccountType
	}
	return 0
}

func (f *AccountFilter) GetAccountTypeBitPos() []uint32 {
	if f != nil && f.AccountTypeBitPos != nil {
		return f.AccountTypeBitPos
	}
	return nil
}
