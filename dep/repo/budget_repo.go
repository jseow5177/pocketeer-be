package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type BudgetRepo interface {
	GetMany(ctx context.Context, req *BudgetFilter) ([]*entity.Budget, error)
	Set(ctx context.Context, budgets []*entity.Budget) error
}

type BudgetFilter struct {
	UserID      *string
	CategoryID  *string
	CategoryIDs []string
	Year        *uint32
	Month       *uint32
	BudgetType  *uint32
}
