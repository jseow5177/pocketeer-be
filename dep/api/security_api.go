package api

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type SecurityAPI interface {
	SearchSecurities(ctx context.Context, sf *SecurityFilter) ([]*entity.Security, error)
	GetLatestQuote(ctx context.Context, sf *SecurityFilter) (*entity.Quote, error)
	ListSymbols(ctx context.Context, sf *SecurityFilter) ([]*entity.Security, error)
}

type SecurityFilter struct {
	Keyword  *string `filter:"keyword"`
	Symbol   *string `filter:"symbol"`
	Exchange *string
}

func (f *SecurityFilter) GetKeyword() string {
	if f != nil && f.Keyword != nil {
		return *f.Keyword
	}
	return ""
}

func (f *SecurityFilter) GetSymbol() string {
	if f != nil && f.Symbol != nil {
		return *f.Symbol
	}
	return ""
}

func (f *SecurityFilter) GetExchange() string {
	if f != nil && f.Exchange != nil {
		return *f.Exchange
	}
	return ""
}
