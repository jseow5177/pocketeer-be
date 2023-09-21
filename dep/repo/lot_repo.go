package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrLotNotFound = errutil.NotFoundError(errors.New("lot not found"))
)

type LotRepo interface {
	Get(ctx context.Context, lf *LotFilter) (*entity.Lot, error)
	GetMany(ctx context.Context, lf *LotFilter) ([]*entity.Lot, error)

	Create(ctx context.Context, l *entity.Lot) (string, error)
	CreateMany(ctx context.Context, ls []*entity.Lot) ([]string, error)
	Update(ctx context.Context, lf *LotFilter, lu *entity.LotUpdate) error
	Delete(ctx context.Context, lf *LotFilter) error
	DeleteMany(ctx context.Context, lf *LotFilter) error
}

type LotFilter struct {
	LotID      *string  `filter:"_id"`
	UserID     *string  `filter:"user_id"`
	HoldingID  *string  `filter:"holding_id"`
	HoldingIDs []string `filter:"holding_id__in"`
	LotStatus  *uint32  `filter:"lot_status"`
	Paging     *Paging  `filter:"-"`
}

type LotFilterOption = func(lf *LotFilter)

func WitLotID(lotID *string) LotFilterOption {
	return func(lf *LotFilter) {
		lf.LotID = lotID
	}
}

func WithLotHoldingID(holdingID *string) LotFilterOption {
	return func(lf *LotFilter) {
		lf.HoldingID = holdingID
	}
}

func WithLotStatus(lotStatus *uint32) LotFilterOption {
	return func(lf *LotFilter) {
		lf.LotStatus = lotStatus
	}
}

func WithLotPaging(paging *Paging) LotFilterOption {
	return func(lf *LotFilter) {
		lf.Paging = paging
	}
}

func WithLotHoldingIDs(holdingIDs []string) LotFilterOption {
	return func(lf *LotFilter) {
		lf.HoldingIDs = holdingIDs
	}
}

func NewLotFilter(userID string, opts ...LotFilterOption) *LotFilter {
	lf := &LotFilter{
		UserID:    goutil.String(userID),
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
	}
	for _, opt := range opts {
		opt(lf)
	}
	return lf
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

func (f *LotFilter) GetLotStatus() uint32 {
	if f != nil && f.LotStatus != nil {
		return *f.LotStatus
	}
	return 0
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

func (f *LotFilter) GetHoldingIDs() []string {
	if f != nil && f.HoldingIDs != nil {
		return f.HoldingIDs
	}
	return nil
}
