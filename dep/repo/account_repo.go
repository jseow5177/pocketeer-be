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
	Get(ctx context.Context, acf *AccountFilter) (*entity.Account, error)
	GetMany(ctx context.Context, acf *AccountFilter) ([]*entity.Account, error)

	Create(ctx context.Context, ac *entity.Account) (string, error)
	CreateMany(ctx context.Context, acs []*entity.Account) ([]string, error)
	Update(ctx context.Context, acf *AccountFilter, acu *entity.AccountUpdate) error
}

type AccountFilter struct {
	UserID            *string  `filter:"user_id"`
	AccountID         *string  `filter:"_id"`
	AccountIDs        []string `filter:"_id__in"`
	AccountName       *string  `filter:"account_name"`
	AccountType       *uint32  `filter:"account_type"`
	AccountTypeBitPos []uint32 `filter:"account_type__bitsAllSet"` // bit pos
	AccountStatus     *uint32  `filter:"account_status"`
}

type AccountFilterOption = func(acf *AccountFilter)

func WithAccountID(accountID *string) AccountFilterOption {
	return func(acf *AccountFilter) {
		acf.AccountID = accountID
	}
}

func WithAccountName(accountName *string) AccountFilterOption {
	return func(acf *AccountFilter) {
		acf.AccountName = accountName
	}
}

func WithAccountIDs(accountIDs []string) AccountFilterOption {
	return func(acf *AccountFilter) {
		acf.AccountIDs = accountIDs
	}
}

func WithAccountStatus(accountStatus *uint32) AccountFilterOption {
	return func(acf *AccountFilter) {
		acf.AccountStatus = accountStatus
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
	acf := &AccountFilter{
		UserID:        goutil.String(userID),
		AccountStatus: goutil.Uint32(uint32(entity.AccountStatusNormal)),
	}
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

func (f *AccountFilter) GetAccountName() string {
	if f != nil && f.AccountName != nil {
		return *f.AccountName
	}
	return ""
}

func (f *AccountFilter) GetAccountID() string {
	if f != nil && f.AccountID != nil {
		return *f.AccountID
	}
	return ""
}

func (f *AccountFilter) GetCategoryIDs() []string {
	if f != nil && f.AccountIDs != nil {
		return f.AccountIDs
	}
	return nil
}

func (f *AccountFilter) GetAccountType() uint32 {
	if f != nil && f.AccountType != nil {
		return *f.AccountType
	}
	return 0
}

func (f *AccountFilter) GetAccountStatus() uint32 {
	if f != nil && f.AccountStatus != nil {
		return *f.AccountStatus
	}
	return 0
}

func (f *AccountFilter) GetAccountTypeBitPos() []uint32 {
	if f != nil && f.AccountTypeBitPos != nil {
		return f.AccountTypeBitPos
	}
	return nil
}
