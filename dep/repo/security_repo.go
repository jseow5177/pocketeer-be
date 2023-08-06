package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrSecurityNotFound = errors.New("security not found")
)

type SecurityRepo interface {
	GetMany(ctx context.Context, sf *SecurityFilter) ([]*entity.Security, error)
	Get(ctx context.Context, sf *SecurityFilter) (*entity.Security, error)

	CreateMany(ctx context.Context, ss []*entity.Security) error
}

type SecurityFilter struct {
	SymbolRegex *string `filter:"symbol__regex"`
	Symbol      *string `filter:"symbol"`
	Paging      *Paging `filter:"-"`
}

func (f *SecurityFilter) GetSymbol() string {
	if f != nil && f.Symbol != nil {
		return *f.Symbol
	}
	return ""
}

func (f *SecurityFilter) GetSymbolRegex() string {
	if f != nil && f.SymbolRegex != nil {
		return *f.SymbolRegex
	}
	return ""
}

func (f *SecurityFilter) GetPaging() *Paging {
	if f != nil && f.Paging != nil {
		return f.Paging
	}
	return nil
}
