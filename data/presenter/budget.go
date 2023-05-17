package presenter

import (
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

//go:generate easytags $GOFILE

// ***************** Common structs
type Budget struct {
	CategoryID      *string `json:"category_id"`
	CategoryName    *string `json:"category_name"`
	BudgetType      *uint32 `json:"budget_type"`
	TransactionType *uint32 `json:"transaction_type"`
	Year            *uint32 `json:"year"`
	Month           *uint32 `json:"month"`
	BudgetAmount    *int64  `json:"budget_amount"`
	Used            *int64  `json:"used"`
}

type BasicBudget struct {
	Year   *uint32 `json:"year"`
	Month  *uint32 `json:"month"`
	Amount *int64  `json:"amount"`
}

type FullYearBudget struct {
	CategoryID      *string        `json:"category_id"`
	BudgetType      *uint32        `json:"budget_type"`
	TransactionType *uint32        `json:"transaction_type"`
	Year            *uint32        `json:"year"`
	DefaultBudget   *int64         `json:"default_budget"`
	MonthlyBudgets  []*BasicBudget `json:"monthly_budgets"`
}

// ***************** Request | Response
type GetCategoryBudgetsByMonthRequest struct {
	Year              *uint32  `json:"year"`
	Month             *uint32  `json:"month"`
	CategoryIDs       []string `json:"category_ids"`
	IncludeUsedAmount *bool    `json:"include_used"`
}

type GetCategoryBudgetsByMonthResponse struct {
	Budgets []*Budget `json:"budgets"`
}

type GetBudgetBreakdownByYearRequest struct {
	CategoryID *string `json:"category_id"`
	Year       *uint32 `json:"year"`
}

type GetBudgetBreakdownByYearResponse struct {
	FullYearBudget *FullYearBudget `json:"full_year_budget"`
}

type SetBudgetRequest struct {
	CategoryID   *string `json:"category_id"`
	Year         *uint32 `json:"year"`
	Month        *uint32 `json:"month"`
	IsDefault    *bool   `json:"is_default"`
	BudgetAmount *int64  `json:"budget_amount"`
	BudgetType   *uint32 `json:"budget_type"`
}

type SetBudgetResponse struct{}

// ***************** Funcs
func (m *Budget) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *Budget) GetCategoryName() string {
	if m != nil && m.CategoryName != nil {
		return *m.CategoryName
	}
	return ""
}

func (m *Budget) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *Budget) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *Budget) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *Budget) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *Budget) GetBudgetAmount() int64 {
	if m != nil && m.BudgetAmount != nil {
		return *m.BudgetAmount
	}
	return 0
}

func (m *Budget) GetUsed() int64 {
	if m != nil && m.Used != nil {
		return *m.Used
	}
	return 0
}

func (m *BasicBudget) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *BasicBudget) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *BasicBudget) GetAmount() int64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *FullYearBudget) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *FullYearBudget) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *FullYearBudget) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *FullYearBudget) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *FullYearBudget) GetDefaultBudget() int64 {
	if m != nil && m.DefaultBudget != nil {
		return *m.DefaultBudget
	}
	return 0
}

func (m *FullYearBudget) GetMonthlyBudgets() []*BasicBudget {
	if m != nil && m.MonthlyBudgets != nil {
		return m.MonthlyBudgets
	}
	return []*BasicBudget{}
}

func (m *GetCategoryBudgetsByMonthRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *GetCategoryBudgetsByMonthRequest) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *GetCategoryBudgetsByMonthRequest) GetCategoryIDs() []string {
	return m.CategoryIDs
}

func (m *GetCategoryBudgetsByMonthRequest) GetIncludeUsedAmount() bool {
	if m != nil && m.IncludeUsedAmount != nil {
		return *m.IncludeUsedAmount
	}
	return false
}

func (m *GetCategoryBudgetsByMonthRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID: goutil.String(userID),
	}
}

func (m *GetCategoryBudgetsByMonthRequest) ToBudgetFilter(userID string) *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:      goutil.String(userID),
		CategoryIDs: m.GetCategoryIDs(),
		Year:        m.Year,
		Month:       m.Month,
	}
}

func (m *GetBudgetBreakdownByYearRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetBudgetBreakdownByYearRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *SetBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *SetBudgetRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *SetBudgetRequest) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *SetBudgetRequest) GetIsDefault() bool {
	if m != nil && m.IsDefault != nil {
		return *m.IsDefault
	}
	return false
}

func (m *SetBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *SetBudgetRequest) GetBudgetAmount() int64 {
	if m != nil && m.BudgetAmount != nil {
		return *m.BudgetAmount
	}
	return 0
}

func (m *SetBudgetRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:     goutil.String(userID),
		CategoryID: goutil.String(m.GetCategoryID()),
	}
}
