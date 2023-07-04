package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrHoldingNotFound = errors.New("holding not found")
)

type HoldingRepo interface {
	Get(ctx context.Context, hf *HoldingFilter) (*entity.Holding, error)
	GetMany(ctx context.Context, hf *HoldingFilter) ([]*entity.Holding, error)

	Create(ctx context.Context, h *entity.Holding) (string, error)
}

type HoldingFilter struct {
	HoldingID *string `filter:"_id"`
	AccountID *string `filter:"account_id"`
	UserID    *string `filter:"user_id"`
}

func (f *HoldingFilter) GetHoldingID() string {
	if f != nil && f.HoldingID != nil {
		return *f.HoldingID
	}
	return ""
}

func (f *HoldingFilter) GetAccountID() string {
	if f != nil && f.AccountID != nil {
		return *f.AccountID
	}
	return ""
}

func (f *HoldingFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}
