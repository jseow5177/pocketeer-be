package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type ExchangeRateRepo interface {
	GetMany(ctx context.Context, erf *ExchangeRateFilter) ([]*entity.ExchangeRate, error)
	CreateMany(ctx context.Context, ers []*entity.ExchangeRate) ([]string, error)
}

type ExchangeRateFilter struct {
	From         *string `filter:"from"`
	To           *string `filter:"to"`
	TimestampGte *uint64 `filter:"timestamp__gte"`
	TimestampLte *uint64 `filter:"timestamp__lte"`
}

type ExchangeRateFilterOption = func(erf *ExchangeRateFilter)

func WithTimestampGte(timestampGte *uint64) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.TimestampGte = timestampGte
	}
}

func WithTimestampLte(timestampLte *uint64) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.TimestampLte = timestampLte
	}
}

func NewExchangeRateFilter(from, to string, opts ...ExchangeRateFilterOption) *ExchangeRateFilter {
	erf := &ExchangeRateFilter{
		From: goutil.String(from),
		To:   goutil.String(to),
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
