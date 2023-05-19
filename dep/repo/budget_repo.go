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
	UserID      *string  `filter:"user_id"`
	CategoryID  *string  `filter:"category_id"`
	CategoryIDs []string `filter:"category_id__in"`
	Year        *uint32  `filter:"year"`
	Month       *uint32  `filter:"month"`
	BudgetType  *uint32  `filter:"budget_type"`
}
