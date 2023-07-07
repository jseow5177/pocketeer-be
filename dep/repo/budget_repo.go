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
	Set(ctx context.Context, budget *entity.Budget) error
}

type BudgetFilter struct {
	UserID    *string  `filter:"user_id"`
	BudgetID  *string  `filter:"_id"`
	BudgetIDs []string `filter:"_id__in"`
}
