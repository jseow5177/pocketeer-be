package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrInvalidQuote  = errors.New("invalid quote in mem cache")
)

type QuoteRepo interface {
	Get(ctx context.Context, qf *QuoteFilter) (*entity.Quote, error)
}

type QuoteFilter struct {
	Symbol *string
}

func (f *QuoteFilter) GetSymbol() string {
	if f != nil && f.Symbol != nil {
		return *f.Symbol
	}
	return ""
}
