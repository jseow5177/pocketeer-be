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
	Update(ctx context.Context, tf *TransactionFilter, t *entity.TransactionUpdate) error

	Sum(ctx context.Context, sumBy string, tf *TransactionFilter) (map[string]float64, error)
}

type TransactionFilter struct {
	UserID             *string  `filter:"user_id"`
	TransactionID      *string  `filter:"_id"`
	AccountID          *string  `filter:"account_id"`
	CategoryID         *string  `filter:"category_id"`
	CategoryIDs        []string `filter:"category_id__in"`
	TransactionType    *uint32  `filter:"transaction_type"`
	TransactionTypes   []uint32 `filter:"transaction_type__in"`
	TransactionTimeGte *uint64  `filter:"transaction_time__gte"`
	TransactionTimeLte *uint64  `filter:"transaction_time__lte"`
	Paging             *Paging  `filter:"-"`
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

func (f *TransactionFilter) GetAccountID() string {
	if f != nil && f.AccountID != nil {
		return *f.AccountID
	}
	return ""
}

func (f *TransactionFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *TransactionFilter) GetCategoryIDs() []string {
	if f != nil && f.CategoryIDs != nil {
		return f.CategoryIDs
	}
	return nil
}

func (f *TransactionFilter) GetTransactionType() uint32 {
	if f != nil && f.TransactionType != nil {
		return *f.TransactionType
	}
	return 0
}

func (f *TransactionFilter) GetTransactionTypes() []uint32 {
	if f != nil && f.TransactionTypes != nil {
		return f.TransactionTypes
	}
	return nil
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
	return nil
}
