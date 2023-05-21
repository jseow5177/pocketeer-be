package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionRepo interface {
	Get(ctx context.Context, tf *TransactionFilter) (*entity.Transaction, error)
	GetMany(ctx context.Context, tf *TransactionFilter) ([]*entity.Transaction, error)

	Create(ctx context.Context, t *entity.Transaction) (string, error)
	Update(ctx context.Context, tf *TransactionFilter, t *entity.Transaction) error
}

type TransactionFilter struct {
	UserID             *string `filter:"user_id"`
	TransactionID      *string `filter:"_id"`
	CategoryID         *string `filter:"category_id"`
	TransactionType    *uint32 `filter:"transaction_type"`
	TransactionTimeGte *uint64 `filter:"transaction_time__gte"`
	TransactionTimeLte *uint64 `filter:"transaction_time__lte"`
	Paging             *Paging `filter:"-"`
}

func (f *TransactionFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *TransactionFilter) GetTransactionID() string {
	if f != nil && f.TransactionID != nil {
		return *f.TransactionID
	}
	return ""
}

func (f *TransactionFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *TransactionFilter) GetTransactionType() uint32 {
	if f != nil && f.TransactionType != nil {
		return *f.TransactionType
	}
	return 0
}

func (f *TransactionFilter) GetTransactionTimeGte() uint64 {
	if f != nil && f.TransactionTimeGte != nil {
		return *f.TransactionTimeGte
	}
	return 0
}

func (f *TransactionFilter) GetTransactionTimeLte() uint64 {
	if f != nil && f.TransactionTimeLte != nil {
		return *f.TransactionTimeLte
	}
	return 0
}

func (f *TransactionFilter) GetPaging() *Paging {
	if f != nil && f.Paging != nil {
		return f.Paging
	}
	return new(Paging)
}
