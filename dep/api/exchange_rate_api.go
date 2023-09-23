package api

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type ExchangeRateAPI interface {
	GetExchangeRates(ctx context.Context, erf *ExchangeRateFilter) ([]*entity.ExchangeRate, error)
}

type ExchangeRateFilter struct {
	Date    *string
	Base    *string
	Symbols []string
}

func NewExchangeRateFilter(timestamp *uint64, base string, symbols ...string) *ExchangeRateFilter {
	var date string
	if timestamp != nil {
		ts := util.ToUnixStartOfMonth(*timestamp)
		date = util.FormatDate(time.UnixMilli(int64(ts)))
	}

	return &ExchangeRateFilter{
		Base:    goutil.String(base),
		Date:    goutil.String(date),
		Symbols: symbols,
	}
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
