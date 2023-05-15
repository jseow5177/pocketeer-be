package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
)

// An abstraction for a storage that supports transactions
type TxMgr interface {
	WithTx(ctx context.Context, txFn func(txCtx context.Context) error) error
}

// An abstraction for Transaction storage
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

// An abstraction for Category storage
type CategoryRepo interface {
	Get(ctx context.Context, cf *CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, cf *CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) (string, error)
	Update(ctx context.Context, cf *CategoryFilter, c *entity.Category) error
}

type CategoryFilter struct {
	UserID       *string `filter:"user_id"`
	CategoryID   *string `filter:"_id"`
	CategoryType *uint32 `filter:"category_type"`
}

func (f *CategoryFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *CategoryFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *CategoryFilter) GetCategoryType() uint32 {
	if f != nil && f.CategoryType != nil {
		return *f.CategoryType
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
