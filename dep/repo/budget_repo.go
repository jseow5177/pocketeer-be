package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrBudgetNotFound = errors.New("budget not found")
)

type BudgetRepo interface {
	Get(ctx context.Context, bf *BudgetFilter) (*entity.Budget, error)
	GetMany(ctx context.Context, bf *BudgetFilter) ([]*entity.Budget, error)

	Create(ctx context.Context, b *entity.Budget) (string, error)
}

type BudgetFilter struct{}
