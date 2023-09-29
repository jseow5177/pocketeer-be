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
	FromDate   *string
	ToDate     *string
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

func NewExchangeRateFilter(fromDate, toDate string, opts ...ExchangeRateFilterOption) *ExchangeRateFilter {
	erf := &ExchangeRateFilter{
		FromDate: goutil.String(fromDate),
		ToDate:   goutil.String(toDate),
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

func (f *ExchangeRateFilter) GetFromDate() string {
	if f != nil && f.FromDate != nil {
		return *f.FromDate
	}
	return ""
}

func (f *ExchangeRateFilter) GetToDate() string {
	if f != nil && f.ToDate != nil {
		return *f.ToDate
	}
	return ""
}
