package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
)

var (
	ErrQuoteNotFound = errutil.NotFoundError(errors.New("quote not found"))
	ErrInvalidQuote  = errors.New("invalid quote in mem cache")
)

type QuoteRepo interface {
	Get(ctx context.Context, qf *QuoteFilter) (*entity.Quote, error)
	Upsert(ctx context.Context, qf *QuoteFilter, q *entity.Quote) error
}

type QuoteFilter struct {
	Symbol  *string  `filter:"symbol"`
	Symbols []string `filter:"symbol__in"`
	Paging  *Paging  `filter:"-"`
}

type QuoteFilterOption = func(qf *QuoteFilter)

func WithQuoteSymbols(symbols []string) QuoteFilterOption {
	return func(qf *QuoteFilter) {
		qf.Symbols = symbols
	}
}

func WithQuoteSymbol(symbol *string) QuoteFilterOption {
	return func(qf *QuoteFilter) {
		qf.Symbol = symbol
	}
}

func WithQuoteRatePaging(paging *Paging) QuoteFilterOption {
	return func(qf *QuoteFilter) {
		qf.Paging = paging
	}
}

func NewQuoteFilter(opts ...QuoteFilterOption) *QuoteFilter {
	qf := new(QuoteFilter)
	for _, opt := range opts {
		opt(qf)
	}
	return qf
}

func (f *QuoteFilter) GetSymbol() string {
	if f != nil && f.Symbol != nil {
		return *f.Symbol
	}
	return ""
}

func (f *QuoteFilter) GetSymbols() []string {
	if f != nil && f.Symbols != nil {
		return f.Symbols
	}
	return nil
}
