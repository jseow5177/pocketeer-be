package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrHoldingNotFound      = errors.New("holding not found")
	ErrHoldingAlreadyExists = errors.New("holding already exists")
)

type HoldingRepo interface {
	Get(ctx context.Context, hf *HoldingFilter) (*entity.Holding, error)
	GetMany(ctx context.Context, hf *HoldingFilter) ([]*entity.Holding, error)

	Create(ctx context.Context, h *entity.Holding) (string, error)
	Update(ctx context.Context, hf *HoldingFilter, hu *entity.HoldingUpdate) error
}

type HoldingFilter struct {
	HoldingID   *string `filter:"_id"`
	AccountID   *string `filter:"account_id"`
	UserID      *string `filter:"user_id"`
	Symbol      *string `filter:"symbol"`
	HoldingType *uint32 `filter:"holding_type"`
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

func (f *HoldingFilter) GetSymbol() string {
	if f != nil && f.Symbol != nil {
		return *f.Symbol
	}
	return ""
}

func (f *HoldingFilter) GetHoldingType() uint32 {
	if f != nil && f.HoldingType != nil {
		return *f.HoldingType
	}
	return 0
}
