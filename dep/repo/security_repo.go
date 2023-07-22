package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type SecurityRepo interface {
	GetMany(ctx context.Context, sf *SecurityFilter) ([]*entity.Security, error)

	CreateMany(ctx context.Context, ss []*entity.Security) error
}

type SecurityFilter struct {
	SymbolRegex *string `filter:"symbol__regex"`
	Paging      *Paging `filter:"-"`
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
