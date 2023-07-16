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

	CalcTotalSharesAndCost(ctx context.Context, lf *LotFilter) (*LotAggr, error)
}

type LotFilter struct {
	LotID     *string `filter:"_id"`
	UserID    *string `filter:"user_id"`
	HoldingID *string `filter:"holding_id"`
	Paging    *Paging `filter:"-"`
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

func (f *LotFilter) GetPaging() *Paging {
	if f != nil && f.Paging != nil {
		return f.Paging
	}
	return nil
}

type LotAggr struct {
	GroupBy     *string
	TotalShares *float64
	TotalCost   *float64
}

func (ag *LotAggr) GetGroupBy() string {
	if ag != nil && ag.GroupBy != nil {
		return *ag.GroupBy
	}
	return ""
}

func (ag *LotAggr) GetTotalShares() float64 {
	if ag != nil && ag.TotalShares != nil {
		return *ag.TotalShares
	}
	return 0
}

func (ag *LotAggr) GetTotalCost() float64 {
	if ag != nil && ag.TotalCost != nil {
		return *ag.TotalCost
	}
	return 0
}
