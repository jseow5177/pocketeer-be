package api

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type ExchangeRateAPI interface {
	GetExchangeRates(ctx context.Context, erf *ExchangeRateFilter) ([]*entity.ExchangeRate, error)
}

type ExchangeRateFilter struct {
	Base    *string
	Symbols []string
	Date    *string
}

func (f *ExchangeRateFilter) GetBase() string {
	if f != nil && f.Base != nil {
		return *f.Base
	}
	return ""
}

func (f *ExchangeRateFilter) GetSymbols() []string {
	if f != nil && f.Symbols != nil {
		return f.Symbols
	}
	return nil
}

func (f *ExchangeRateFilter) GetDate() string {
	if f != nil && f.Date != nil {
		return *f.Date
	}
	return ""
}
