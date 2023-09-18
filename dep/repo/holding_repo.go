package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrHoldingNotFound = errutil.NotFoundError(errors.New("holding not found"))
)

type HoldingRepo interface {
	Get(ctx context.Context, hf *HoldingFilter) (*entity.Holding, error)
	GetMany(ctx context.Context, hf *HoldingFilter) ([]*entity.Holding, error)

	Create(ctx context.Context, h *entity.Holding) (string, error)
	CreateMany(ctx context.Context, hs []*entity.Holding) ([]string, error)
	Update(ctx context.Context, hf *HoldingFilter, hu *entity.HoldingUpdate) error
	DeleteMany(ctx context.Context, hf *HoldingFilter) error
	Delete(ctx context.Context, hf *HoldingFilter) error
}

type HoldingFilter struct {
	HoldingID     *string  `filter:"_id"`
	HoldingIDs    []string `filter:"_id__in"`
	AccountID     *string  `filter:"account_id"`
	UserID        *string  `filter:"user_id"`
	Symbol        *string  `filter:"symbol"`
	HoldingType   *uint32  `filter:"holding_type"`
	HoldingStatus *uint32  `filter:"holding_status"`
}

type HoldingFilterOption = func(hf *HoldingFilter)

func WithHoldingID(holdingID *string) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.HoldingID = holdingID
	}
}

func WithHoldingAccountID(accountID *string) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.AccountID = accountID
	}
}

func WithHoldingSymbol(symbol *string) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.Symbol = symbol
	}
}

func WithHoldingType(holdingType *uint32) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.HoldingType = holdingType
	}
}

func WithHoldingStatus(holdingStatus *uint32) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.HoldingStatus = holdingStatus
	}
}

func WithHoldingIDs(holdingIDs []string) HoldingFilterOption {
	return func(hf *HoldingFilter) {
		hf.HoldingIDs = holdingIDs
	}
}

func NewHoldingFilter(userID string, opts ...HoldingFilterOption) *HoldingFilter {
	hf := &HoldingFilter{
		UserID:        goutil.String(userID),
		HoldingStatus: goutil.Uint32(uint32(entity.HoldingStatusNormal)),
	}
	for _, opt := range opts {
		opt(hf)
	}
	return hf
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

func (f *HoldingFilter) GetHoldingIDs() []string {
	if f != nil && f.HoldingIDs != nil {
		return f.HoldingIDs
	}
	return nil
}
