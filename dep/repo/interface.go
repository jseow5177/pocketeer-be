package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
)

// An abstraction for a storage that supports transactions
type TxMgr interface {
	WithTx(ctx context.Context, txFn func(txCtx context.Context) error) error
}

// An abstraction for Transaction storage
type TransactionRepo interface {
	Get(ctx context.Context, tf *TransactionFilter) (*entity.Transaction, error)
	//GetMany(ctx context.Context, tf *TransactionFilter) ([]*entity.Transaction, error)

	Create(ctx context.Context, t *entity.Transaction) (string, error)
	//Update(ctx context.Context, tf *TransactionFilter, t *entity.Transaction) error
}

type TransactionFilter struct {
	UserID          *string `filter:"user_id"`
	TransactionID   *string `filter:"_id"`
	TransactionType *uint32 `filter:"transaction_type"`
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

func (f *TransactionFilter) GetTransactionType() uint32 {
	if f != nil && f.TransactionType != nil {
		return *f.TransactionType
	}
	return 0
}

// An abstraction for Category storage
type CategoryRepo interface {
	Get(ctx context.Context, cf *CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, cf *CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) (string, error)
	Update(ctx context.Context, cf *CategoryFilter, c *entity.Category) error
}

type CategoryFilter struct {
	UserID  *string `filter:"user_id"`
	CatID   *string `filter:"_id"`
	CatType *uint32 `filter:"cat_type"`
}

func (f *CategoryFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *CategoryFilter) GetCatID() string {
	if f != nil && f.CatID != nil {
		return *f.CatID
	}
	return ""
}

func (f *CategoryFilter) GetCatType() uint32 {
	if f != nil && f.CatType != nil {
		return *f.CatType
	}
	return 0
}
