package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
)

type ExchangeRateRepo interface {
	GetMany(ctx context.Context, erf *ExchangeRateFilter) ([]*entity.ExchangeRate, error)
	Get(ctx context.Context, erf *GetExchangeRateFilter) (*entity.ExchangeRate, error)
	CreateMany(ctx context.Context, ers []*entity.ExchangeRate) ([]string, error)
	Create(ctx context.Context, er *entity.ExchangeRate) (string, error)
}

type GetExchangeRateFilter struct {
	From      *string
	To        *string
	Timestamp *uint64
}

func (f *GetExchangeRateFilter) GetFrom() string {
	if f != nil && f.From != nil {
		return *f.From
	}
	return ""
}

func (f *GetExchangeRateFilter) GetTo() string {
	if f != nil && f.To != nil {
		return *f.To
	}
	return ""
}

func (f *GetExchangeRateFilter) GetTimestamp() uint64 {
	if f != nil && f.Timestamp != nil {
		return *f.Timestamp
	}
	return 0
}

func (f *GetExchangeRateFilter) ToExchangeRateFilter() *ExchangeRateFilter {
	minDate, _ := util.ParseDate(config.MinCurrencyDate)

	t := time.UnixMilli(int64(f.GetTimestamp()))

	if t.Before(minDate) {
		t = minDate
	}

	return &ExchangeRateFilter{
		From:         f.From,
		To:           f.To,
		TimestampLte: goutil.Uint64(uint64(t.UnixMilli())),
		Paging: &Paging{
			Limit: goutil.Uint32(1),
			Sorts: []filter.Sort{
				&Sort{
					Field: goutil.String("timestamp"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		},
	}
}

type ExchangeRateFilter struct {
	From         *string `filter:"from"`
	To           *string `filter:"to"`
	TimestampGte *uint64 `filter:"timestamp__gte"`
	TimestampLte *uint64 `filter:"timestamp__lte"`
	Paging       *Paging `filter:"-"`
}

type ExchangeRateFilterOption = func(erf *ExchangeRateFilter)

func WithTimestampGte(timestampGte *uint64) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.TimestampGte = timestampGte
	}
}

func WithTimestampLte(timestampLte *uint64) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.SetTimestampLte(timestampLte)
	}
}

func WithExchangeRateFrom(from *string) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.From = from
	}
}

func WithExchangeRatePaging(paging *Paging) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.Paging = paging
	}
}

func NewExchangeRateFilter(to string, opts ...ExchangeRateFilterOption) *ExchangeRateFilter {
	erf := &ExchangeRateFilter{
		To: goutil.String(to),
	}
	for _, opt := range opts {
		opt(erf)
	}
	return erf
}

func (f *ExchangeRateFilter) GetFrom() string {
	if f != nil && f.From != nil {
		return *f.From
	}
	return ""
}

func (f *ExchangeRateFilter) GetTo() string {
	if f != nil && f.To != nil {
		return *f.To
	}
	return ""
}

func (f *ExchangeRateFilter) GetTimestampGte() uint64 {
	if f != nil && f.TimestampGte != nil {
		return *f.TimestampGte
	}
	return 0
}

func (f *ExchangeRateFilter) GetTimestampLte() uint64 {
	if f != nil && f.TimestampLte != nil {
		return *f.TimestampLte
	}
	return 0
}

func (f *ExchangeRateFilter) SetTimestampLte(timestampLte *uint64) {
	f.TimestampLte = timestampLte
}
