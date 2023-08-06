package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
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
