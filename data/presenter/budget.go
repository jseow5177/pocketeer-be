package presenter

import (
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

//go:generate easytags $GOFILE

type GetMonthBudgetsRequest struct {
	Year              *uint32  `json:"year"`
	Month             *uint32  `json:"month"`
	CatIDs            []string `json:"cat_ids"`
	IncludeUsedAmount *bool    `json:"include_used"`
}

type GetMonthBudgetsResponse struct {
	Budgets []*Budget `json:"budgets"`
}

type Budget struct {
	CatID           *string `json:"cat_id"`
	CatName         *string `json:"cat_name"`
	BudgetType      *uint32 `json:"budget_type"`
	TransactionType *uint32 `json:"transaction_type"`
	Year            *uint32 `json:"year"`
	Month           *uint32 `json:"month"`
	BudgetAmount    *int64  `json:"budget_amount"`
	Used            *int64  `json:"used"`
}

type GetFullYearBudgetRequest struct {
	CatID *string `json:"cat_id"`
	Year  *uint32 `json:"year"`
}

type GetFullYearBudgetResponse struct {
	FullYearBudget *FullYearBudget `json:"budget_breakdown"`
}

type FullYearBudget struct {
	CatID           *string        `json:"cat_id"`
	BudgetType      *uint32        `json:"budget_type"`
	TransactionType *uint32        `json:"transaction_type"`
	Year            *uint32        `json:"year"`
	DefaultBudget   *int64         `json:"default_budget"`
	MonthlyBudgets  []*BasicBudget `json:"monthly_budgets"`
}

type BasicBudget struct {
	Year   *uint32 `json:"year"`
	Month  *uint32 `json:"month"`
	Amount *int64  `json:"amount"`
}

type SetBudgetRequest struct {
	CatID         *string        `json:"cat_id"`
	BudgetType    *uint32        `json:"budget_type"`
	CurrYear      *uint32        `json:"curr_year"`
	CurrMonth     *uint32        `json:"curr_month"`
	DefaultBudget *int64         `json:"default_budget"`
	CustomBudgets []*BasicBudget `json:"custom_budgets"`
}

func (m *GetMonthBudgetsRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID: goutil.String(userID),
	}
}

func (m *GetMonthBudgetsRequest) GetCatIDs() []string {
	if m != nil {
		return m.CatIDs
	}
	return []string{}
}

func (m *GetMonthBudgetsRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *GetMonthBudgetsRequest) GetMonth() uint32 {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

func (m *GetFullYearBudgetRequest) ToBudgetConfigFilter(userID string) *repo.BudgetConfigFilter {
	return &repo.BudgetConfigFilter{
		UserID: goutil.String(userID),
		CatIDs: []string{m.GetCatID()},
	}
}

func (m *GetFullYearBudgetRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID: goutil.String(userID),
		CategoryID:  goutil.String(m.GetCatID()),
	}
}

func (m *GetFullYearBudgetRequest) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return ""
}

func (m *GetFullYearBudgetRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *SetBudgetRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID: goutil.String(userID),
		CategoryID:  goutil.String(m.GetCatID()),
	}
}

func (m *SetBudgetRequest) ToBudgetConfigFilter(userID string) *repo.BudgetConfigFilter {
	return &repo.BudgetConfigFilter{
		UserID: goutil.String(userID),
		CatIDs: []string{m.GetCatID()},
	}
}

func (m *SetBudgetRequest) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return ""
}

func (m *SetBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *SetBudgetRequest) GetCurrYear() uint32 {
	if m != nil && m.CurrYear != nil {
		return *m.CurrYear
	}
	return 0
}

func (m *SetBudgetRequest) GetCurrMonth() uint32 {
	if m != nil && m.CurrMonth != nil {
		return *m.CurrMonth
	}
	return 0
}

func (m *SetBudgetRequest) GetDefaultBudget() int64 {
	if m != nil && m.DefaultBudget != nil {
		return *m.DefaultBudget
	}
	return 0
}

func (m *SetBudgetRequest) GetCustomBudgets() []*BasicBudget {
	if m != nil && m.CustomBudgets != nil {
		return *&m.CustomBudgets
	}
	return []*BasicBudget{}
}

type SetBudgetResponse struct{}

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
