package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/pkg/filter"
)

// An abstraction for a storage that supports transactions
type TxMgr interface {
	WithTx(ctx context.Context, txFn func(txCtx context.Context) error) error
}

type Sort struct {
	Field *string
	Order *string
}

func (s *Sort) GetField() *string {
	if s != nil && s.Field != nil {
		return s.Field
	}
	return nil
}

func (s *Sort) GetOrder() *string {
	if s != nil && s.Order != nil {
		return s.Order
	}
	return nil
}

type Paging struct {
	Limit *uint32
	Page  *uint32
	Sorts []filter.Sort
}

func (p *Paging) GetLimit() *uint32 {
	if p != nil && p.Limit != nil {
		return p.Limit
	}
	return nil
}

func (p *Paging) GetPage() *uint32 {
	if p != nil && p.Page != nil {
		return p.Page
	}
	return nil
}

func (p *Paging) GetSorts() []filter.Sort {
	if p != nil && p.Sorts != nil {
		return p.Sorts
	}
	return nil
}
