package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetCategoryBudgetsByMonth(
		ctx context.Context,
		req *GetCategoryBudgetsByMonthRequest,
	) (*GetCategoryBudgetsByMonthResponse, error)

	GetAnnualBudgetBreakdown(
		ctx context.Context,
		req *GetAnnualBudgetBreakdownRequest,
	) (*GetAnnualBudgetBreakdownResponse, error)

	SetBudget(
		ctx context.Context,
		req *SetBudgetRequest,
	) (*SetBudgetResponse, error)
}

type CategoryBudget struct {
	Category *entity.Category
	Budget   *entity.Budget
}

type GetCategoryBudgetsByMonthRequest struct {
	UserID         *string
	Year           *uint32
	Month          *uint32
	CategoryIDs    []string
	IncludedAmount *bool
}

type GetCategoryBudgetsByMonthResponse struct {
	CategoryBudgets []*CategoryBudget
}

type GetAnnualBudgetBreakdownRequest struct {
	UserID     *string
	CategoryID *string
	Year       *uint32
}

type GetAnnualBudgetBreakdownResponse struct {
	AnnualBudgetBreakdown *entity.AnnualBudgetBreakdown
}

type SetBudgetRequest struct {
	UserID       *string
	CategoryID   *string
	Year         *uint32
	Month        *uint32
	IsDefault    *bool
	BudgetAmount *int64
	BudgetType   *uint32
}

type SetBudgetResponse struct {
	AnnualBudgetBreakdown *entity.AnnualBudgetBreakdown
}

// ***************** Funcs
func (u *CategoryBudget) GetCategory() *entity.Category {
	if u != nil && u.Category != nil {
		return u.Category
	}
	return nil
}

func (u *CategoryBudget) GetBudget() *entity.Budget {
	if u != nil && u.Budget != nil {
		return u.Budget
	}
	return nil
}

func (u *GetCategoryBudgetsByMonthRequest) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *GetCategoryBudgetsByMonthRequest) GetYear() uint32 {
	if u != nil && u.Year != nil {
		return *u.Year
	}
	return 0
}

func (u *GetCategoryBudgetsByMonthRequest) GetMonth() uint32 {
	if u != nil && u.Month != nil {
		return *u.Month
	}
	return 0
}

func (u *GetCategoryBudgetsByMonthRequest) GetCategoryIDs() []string {
	return u.CategoryIDs
}

func (u *GetCategoryBudgetsByMonthRequest) GetIncludedAmount() bool {
	return *u.IncludedAmount
}

func (u *GetCategoryBudgetsByMonthRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:      u.UserID,
		CategoryIDs: u.CategoryIDs,
	}
}

func (u *GetCategoryBudgetsByMonthRequest) ToBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:      u.UserID,
		CategoryIDs: u.CategoryIDs,
		Year:        u.Year,
		Month:       u.Month,
	}
}

func (u *GetCategoryBudgetsByMonthResponse) GetCategoryBudgets() []*CategoryBudget {
	return u.CategoryBudgets
}

// **********************************

func (u *GetAnnualBudgetBreakdownRequest) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *GetAnnualBudgetBreakdownRequest) GetCategoryID() string {
	if u != nil && u.CategoryID != nil {
		return *u.CategoryID
	}
	return ""
}

func (u *GetAnnualBudgetBreakdownRequest) GetYear() uint32 {
	if u != nil && u.Year != nil {
		return *u.Year
	}
	return 0
}

func (u *GetAnnualBudgetBreakdownRequest) ToFullBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:     u.UserID,
		CategoryID: u.CategoryID,
		Year:       u.Year,
	}
}

func (u *GetAnnualBudgetBreakdownResponse) GetAnnualBudgetBreakdown() *entity.AnnualBudgetBreakdown {
	if u != nil && u.AnnualBudgetBreakdown != nil {
		return u.AnnualBudgetBreakdown
	}
	return nil
}

// **********************************

func (u *SetBudgetRequest) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *SetBudgetRequest) GetCategoryID() string {
	if u != nil && u.CategoryID != nil {
		return *u.CategoryID
	}
	return ""
}

func (u *SetBudgetRequest) GetYear() uint32 {
	if u != nil && u.Year != nil {
		return *u.Year
	}
	return 0
}

func (u *SetBudgetRequest) GetMonth() uint32 {
	if u != nil && u.Month != nil {
		return *u.Month
	}
	return 0
}

func (u *SetBudgetRequest) GetIsDefault() bool {
	if u != nil && u.IsDefault != nil {
		return *u.IsDefault
	}
	return false
}

func (u *SetBudgetRequest) GetBudgetAmount() int64 {
	if u != nil && u.BudgetAmount != nil {
		return *u.BudgetAmount
	}
	return 0
}

func (u *SetBudgetRequest) GetBudgetType() uint32 {
	if u != nil && u.BudgetType != nil {
		return *u.BudgetType
	}
	return 0
}

func (u *SetBudgetRequest) ToFullBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:     u.UserID,
		CategoryID: u.CategoryID,
		Year:       u.Year,
	}
}

func (u *SetBudgetResponse) GetAnnualBudgetBreakdown() *entity.AnnualBudgetBreakdown {
	if u != nil && u.AnnualBudgetBreakdown != nil {
		return u.AnnualBudgetBreakdown
	}
	return nil
}
