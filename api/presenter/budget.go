package presenter

import (
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

//go:generate easytags $GOFILE

// ***************** Common structs
type CategoryBudget struct {
	Budget   *Budget   `json:"budget"`
	Category *Category `json:"category"`
}

type Budget struct {
	BudgetID     *string `json:"budget_id"`
	CategoryID   *string `json:"category_id"`
	BudgetType   *uint32 `json:"budget_type"`
	Year         *uint32 `json:"year"`
	Month        *uint32 `json:"month"`
	BudgetAmount *int64  `json:"budget_amount"`
	Used         *int64  `json:"used"`
}

type BasicBudget struct {
	Year   *uint32 `json:"year"`
	Month  *uint32 `json:"month"`
	Amount *int64  `json:"amount"`
}

type AnnualBudgetBreakdown struct {
	CategoryID     *string        `json:"category_id"`
	BudgetType     *uint32        `json:"budget_type"`
	Year           *uint32        `json:"year"`
	DefaultBudget  *int64         `json:"default_budget"`
	MonthlyBudgets []*BasicBudget `json:"monthly_budgets"`
}

// ***************** Request | Response
type GetCategoryBudgetsByMonthRequest struct {
	Year              *uint32  `json:"year"`
	Month             *uint32  `json:"month"`
	CategoryIDs       []string `json:"category_ids"`
	IncludeUsedAmount *bool    `json:"include_used"`
}

type GetCategoryBudgetsByMonthResponse struct {
	CategoryBudgets []*CategoryBudget `json:"category_budgets"`
}

type GetAnnualBudgetBreakdownRequest struct {
	CategoryID *string `json:"category_id"`
	Year       *uint32 `json:"year"`
}

type GetAnnualBudgetBreakdownResponse struct {
	AnnualBudgetBreakdown *AnnualBudgetBreakdown `json:"annual_budget_breakdown"`
}

type SetBudgetRequest struct {
	CategoryID   *string `json:"category_id"`
	Year         *uint32 `json:"year"`
	Month        *uint32 `json:"month"`
	IsDefault    *bool   `json:"is_default"`
	BudgetAmount *int64  `json:"budget_amount"`
	BudgetType   *uint32 `json:"budget_type"`
}

type SetBudgetResponse struct {
	AnnualBudgetBreakdown *AnnualBudgetBreakdown
}

// ***************** Funcs
func (m *Budget) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *Budget) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *Budget) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
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

// *****************
func (m *CategoryBudget) GetCategory() *Category {
	return m.Category
}

func (m *CategoryBudget) GetBudget() *Budget {
	return m.Budget
}

// *****************
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

// *****************
func (m *AnnualBudgetBreakdown) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *AnnualBudgetBreakdown) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *AnnualBudgetBreakdown) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *AnnualBudgetBreakdown) GetDefaultBudget() int64 {
	if m != nil && m.DefaultBudget != nil {
		return *m.DefaultBudget
	}
	return 0
}

func (m *AnnualBudgetBreakdown) GetMonthlyBudgets() []*BasicBudget {
	if m != nil && m.MonthlyBudgets != nil {
		return m.MonthlyBudgets
	}
	return []*BasicBudget{}
}

// *****************
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
		UserID:      goutil.String(userID),
		CategoryIDs: m.GetCategoryIDs(),
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

func (m *GetCategoryBudgetsByMonthRequest) ToUseCaseReq(userID string) *budget.GetCategoryBudgetsByMonthRequest {
	return &budget.GetCategoryBudgetsByMonthRequest{
		UserID:         goutil.String(userID),
		Year:           m.Year,
		Month:          m.Month,
		CategoryIDs:    m.CategoryIDs,
		IncludedAmount: m.IncludeUsedAmount,
	}
}

func (m *GetCategoryBudgetsByMonthResponse) Set(useCaseRes *budget.GetCategoryBudgetsByMonthResponse) {
	m.CategoryBudgets = toCategoryBudgets(useCaseRes.GetCategoryBudgets())
}

// *****************

func (m *GetAnnualBudgetBreakdownRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetAnnualBudgetBreakdownRequest) GetYear() uint32 {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *GetAnnualBudgetBreakdownRequest) ToUseCaseRes(userID string) *budget.GetAnnualBudgetBreakdownRequest {
	return &budget.GetAnnualBudgetBreakdownRequest{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		Year:       m.Year,
	}
}

func (m *GetAnnualBudgetBreakdownRequest) ToFullBudgetFilter(userID string) *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		Year:       m.Year,
	}
}

func (m *GetAnnualBudgetBreakdownResponse) Set(useCaseRes *budget.GetAnnualBudgetBreakdownResponse) {
	m.AnnualBudgetBreakdown = toAnnualBudgetBreakdown(useCaseRes.GetAnnualBudgetBreakdown())
}

// *****************
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

func (m *SetBudgetRequest) ToFullBudgetFilter(userID string) *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		Year:       m.Year,
	}
}

func (m *SetBudgetRequest) ToUseCaseReq(userID string) *budget.SetBudgetRequest {
	return &budget.SetBudgetRequest{
		UserID:       goutil.String(userID),
		CategoryID:   m.CategoryID,
		Year:         m.Year,
		Month:        m.Month,
		IsDefault:    m.IsDefault,
		BudgetAmount: m.BudgetAmount,
		BudgetType:   m.BudgetType,
	}
}

func (m *SetBudgetResponse) Set(useCaseRes *budget.SetBudgetResponse) {
	m.AnnualBudgetBreakdown = toAnnualBudgetBreakdown(useCaseRes.GetAnnualBudgetBreakdown())
}
