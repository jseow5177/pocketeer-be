package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
)

var (
	ErrBudgetNotFound      = errors.New("budget not found")
	ErrBudgetAlreadyExists = errors.New("budget already exists")
)

type BudgetRepo interface {
	GetMany(ctx context.Context, bq *BudgetQuery) ([]*entity.Budget, error)

	Create(ctx context.Context, b *entity.Budget) (string, error)
	CreateMany(ctx context.Context, bs []*entity.Budget) ([]string, error)
}

type BudgetQuery struct {
	Filters []*BudgetFilter
	Queries []*BudgetQuery
	Op      filter.BoolOp
	Paging  *Paging
}

func (q *BudgetQuery) GetQueries() []filter.Query {
	qs := make([]filter.Query, 0)
	for _, bq := range q.Queries {
		qs = append(qs, bq)
	}
	return qs
}

func (q *BudgetQuery) GetFilters() []interface{} {
	ibfs := make([]interface{}, 0)
	for _, bf := range q.Filters {
		ibfs = append(ibfs, bf)
	}
	return ibfs
}

func (q *BudgetQuery) GetOp() filter.BoolOp {
	return q.Op
}

type BudgetFilter struct {
	UserID       *string `filter:"user_id"`
	CategoryID   *string `filter:"category_id"`
	StartDate    *uint64 `filter:"start_date"`
	StartDateLte *uint64 `filter:"start_date__lte"`
	EndDate      *uint64 `filter:"end_date"`
	EndDateGte   *uint64 `filter:"end_date__gte"`
	BudgetStatus *uint32 `filter:"budget_status"`
}

func (f *BudgetFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *BudgetFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *BudgetFilter) GetBudgetStatus() uint32 {
	if f != nil && f.BudgetStatus != nil {
		return *f.BudgetStatus
	}
	return 0
}

func (f *BudgetFilter) GetStartDate() uint64 {
	if f != nil && f.StartDate != nil {
		return *f.StartDate
	}
	return 0
}

func (f *BudgetFilter) GetStartDateLte() uint64 {
	if f != nil && f.StartDateLte != nil {
		return *f.StartDateLte
	}
	return 0
}

func (f *BudgetFilter) GetEndDate() uint64 {
	if f != nil && f.EndDate != nil {
		return *f.EndDate
	}
	return 0
}

func (f *BudgetFilter) GetEndDateGte() uint64 {
	if f != nil && f.EndDateGte != nil {
		return *f.EndDateGte
	}
	return 0
}
