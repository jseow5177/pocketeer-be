package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
)

// An abstraction for a storage that supports transactions
type TxMgr interface {
	WithTx(ctx context.Context, txFn func(txCtx context.Context) error) error
}

// An abstraction for Category storage
type CategoryRepo interface {
	Get(ctx context.Context, cf *CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, cf *CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) (string, error)
	Update(ctx context.Context, c *entity.Category) error
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

type BudgetConfigFilter struct {
	UserID          *string  `filter:"user_id"`
	TransactionType *uint32  `filter:"transaction_type"`
	CatIDs          []string `filter:"cat_ids"`
}

type BudgetConfigRepo interface {
	Get(ctx context.Context, filter *BudgetConfigFilter) (*entity.BudgetConfig, error)
	GetMany(ctx context.Context, filter *BudgetConfigFilter) ([]*entity.BudgetConfig, error)
	Update(ctx context.Context, budgetConfig *entity.BudgetConfig) (*entity.BudgetConfig, error)
}
