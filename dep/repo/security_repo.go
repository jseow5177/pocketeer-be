package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type SecurityRepo interface {
	CreateMany(ctx context.Context, ss []*entity.Security) error
}

type SecurityFilter struct {
	SymbolKeyword *string `filter:"symbol"`
}

func (f *SecurityFilter) GetSymbolKeyword() string {
	if f != nil && f.SymbolKeyword != nil {
		return *f.SymbolKeyword
	}
	return ""
}
