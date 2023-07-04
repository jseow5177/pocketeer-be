package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrLotNotFound = errors.New("lot not found")
)

type LotRepo interface {
	Get(ctx context.Context, lf *LotFilter) (*entity.Lot, error)
	GetMany(ctx context.Context, lf *LotFilter) ([]*entity.Lot, error)

	Create(ctx context.Context, l *entity.Lot) (string, error)
	Update(ctx context.Context, lf *LotFilter, lu *entity.LotUpdate) error
}

type LotFilter struct {
	UserID    *string `filter:"user_id"`
	LotID     *string `filter:"_id"`
	HoldingID *string `filter:"holding_id"`
}

func (f *LotFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *LotFilter) GetLotID() string {
	if f != nil && f.LotID != nil {
		return *f.LotID
	}
	return ""
}

func (f *LotFilter) GetHoldingID() string {
	if f != nil && f.HoldingID != nil {
		return *f.HoldingID
	}
	return ""
}
