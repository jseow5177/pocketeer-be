package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepo interface {
	Get(ctx context.Context, af *AccountFilter) (*entity.Account, error)

	Create(ctx context.Context, ac *entity.Account) (string, error)
	Update(ctx context.Context, af *AccountFilter, acu *entity.AccountUpdate) error
}

type AccountFilter struct {
	UserID    *string `filter:"user_id"`
	AccountID *string `filter:"_id"`
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
