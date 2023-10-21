package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrExchangeRateNotFound = errors.New("exchange rate not found")
)

type ExchangeRateRepo interface {
	Get(ctx context.Context, erf *ExchangeRateFilter) (*entity.ExchangeRate, error)
	CreateMany(ctx context.Context, ers []*entity.ExchangeRate) ([]string, error)
	Create(ctx context.Context, er *entity.ExchangeRate) (string, error)
}

type ExchangeRateFilter struct {
	From      *string `filter:"from"`
	To        *string `filter:"to"`
	Timestamp *uint64 `filter:"timestamp"`
	Paging    *Paging `filter:"-"`
}

type ExchangeRateFilterOption = func(erf *ExchangeRateFilter)

func WithExchangeRateTimestamp(timestamp *uint64) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.Timestamp = timestamp

		if timestamp != nil {
			minDate, _ := util.ParseDate(config.MinCurrencyDate)

			t := time.UnixMilli(int64(*timestamp))

			if t.Before(minDate) {
				t = minDate
			}

			erf.Timestamp = goutil.Uint64(uint64(t.UnixMilli()))
		}
	}
}

func WithExchangeRateFrom(from *string) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.From = from
	}
}

func WithExchangeRateTo(to *string) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.To = to
	}
}

func WithExchangeRatePaging(paging *Paging) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.Paging = paging
	}
}

func NewExchangeRateFilter(opts ...ExchangeRateFilterOption) *ExchangeRateFilter {
	erf := new(ExchangeRateFilter)
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

func (f *ExchangeRateFilter) GetTimestamp() uint64 {
	if f != nil && f.Timestamp != nil {
		return *f.Timestamp
	}
	return 0
}
