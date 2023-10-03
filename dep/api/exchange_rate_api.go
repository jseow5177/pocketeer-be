package api

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type ExchangeRateAPI interface {
	GetExchangeRates(ctx context.Context, erf *ExchangeRateFilter) ([]*entity.ExchangeRate, error)
}

type ExchangeRateFilter struct {
	Date       *string
	Base       *string
	Currencies []string
}

type ExchangeRateFilterOption = func(erf *ExchangeRateFilter)

func WithExchangeRateBase(base *string) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.Base = base
	}
}

func WithExchangeRateCurrencies(currencies ...string) ExchangeRateFilterOption {
	return func(erf *ExchangeRateFilter) {
		erf.Currencies = currencies
	}
}

func NewExchangeRateFilter(date string, opts ...ExchangeRateFilterOption) *ExchangeRateFilter {
	erf := &ExchangeRateFilter{
		Date: goutil.String(date),
	}
	for _, opt := range opts {
		opt(erf)
	}
	return erf
}

func (f *ExchangeRateFilter) GetBase() string {
	if f != nil && f.Base != nil {
		return *f.Base
	}
	return ""
}

func (f *ExchangeRateFilter) GetCurrencies() []string {
	if f != nil && f.Currencies != nil {
		return f.Currencies
	}
	return nil
}

func (f *ExchangeRateFilter) GetDate() string {
	if f != nil && f.Date != nil {
		return *f.Date
	}
	return ""
}
