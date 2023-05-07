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
	Get(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) error
	Update(ctx context.Context, c *entity.Category) error
}
