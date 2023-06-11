package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/aggr"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/util"
)

//go:generate easytags $GOFILE
type CategoryBudget struct {
	Budget     *Budget     `json:"budget,omitempty"`
	Categories []*Category `json:"categories,omitempty"`
}

type Budget struct {
	BudgetID         *string            `json:"budget_id,omitempty"`
	BudgetName       *string            `json:"budget_name,omitempty"`
	BudgetType       *uint32            `json:"budget_type,omitempty"`
	CategoryIDs      []string           `json:"category_ids,omitempty"`
	BudgetBreakdowns []*BudgetBreakdown `json:"budget_breakdowns,omitempty"`
}

type BudgetBreakdown struct {
	Amount *float64 `json:"amount,omitempty"`
	Year   *int     `json:"year,omitempty"`
	Month  *int     `json:"month,omitempty"`
}

type GetBudgetRequest struct {
	BudgetID *string `json:"budget_id,omitempty"`
	Date     *string `json:"date,omitempty"`
}

type GetBudgetResponse struct {
	CategoryBudget *CategoryBudget `json:"category_budget,omitempty"`
}

type GetBudgetsRequest struct {
	Date *string `json:"date,omitempty"`
}

type GetBudgetsResponse struct {
	Budgets []*Budget `json:"budgets,omitempty"`
}

type SetBudgetRequest struct {
	BudgetID       *string  `json:"budget_id,omitempty"`
	BudgetName     *string  `json:"budget_name,omitempty"`
	BudgetType     *uint32  `json:"budget_type,omitempty"`
	BudgetAmount   *string  `json:"budget_amount,omitempty"`
	CategoryIDs    []string `json:"category_ids,omitempty"`
	RangeStartDate *string  `json:"range_start_date,omitempty"`
	RangeEndDate   *string  `json:"range_end_date,omitempty"`
}

type SetBudgetResponse struct{}

func (m *GetBudgetRequest) ToUseCaseReq(userID string) *aggr.GetBudgetWithCategoriesRequest {
	date, _ := util.DateStrToDate(m.GetDate())

	return &aggr.GetBudgetWithCategoriesRequest{
		UserID:   goutil.String(userID),
		BudgetID: m.BudgetID,
		Date:     date,
	}
}

func (m *GetBudgetResponse) Set(usecaseRes *aggr.GetBudgetWithCategoriesResponse) {
	m.CategoryBudget = &CategoryBudget{
		Budget:     toBudget(usecaseRes.GetBudget()),
		Categories: toCategories(usecaseRes.GetCategories()),
	}
}

func (m *GetBudgetsRequest) ToUseCaseReq(userID string) *budget.GetBudgetsRequest {
	date, _ := util.DateStrToDate(m.GetDate())

	return &budget.GetBudgetsRequest{
		UserID: goutil.String(userID),
		Date:   date,
	}
}

func (m *GetBudgetsResponse) Set(usecaseRes *budget.GetBudgetsResponse) {
	m.Budgets = toBudgets(usecaseRes.Budgets)
}

func (m *SetBudgetRequest) ToUseCaseReq(userID string) *budget.SetBudgetRequest {
	rangeStartDate, _ := util.DateStrToDate(m.GetRangeStartDate())
	rangeEndDate, _ := util.DateStrToDate(m.GetRangeEndDate())

	return &budget.SetBudgetRequest{
		UserID:         goutil.String(userID),
		BudgetID:       m.BudgetID,
		BudgetName:     m.BudgetName,
		BudgetType:     m.BudgetType,
		BudgetAmount:   goutil.Float64(util.MonetaryStrToFloat(m.GetBudgetAmount())),
		CategoryIDs:    m.CategoryIDs,
		RangeStartDate: rangeStartDate,
		RangeEndDate:   rangeEndDate,
	}
}

func (m *SetBudgetResponse) Set(usecaseRes *budget.SetBudgetResponse) {}

// Getters ----------------------------------------------------------------------
// CategoryBudget getters
func (m *CategoryBudget) GetBudget() *Budget {
	if m != nil {
		return m.Budget
	}
	return nil
}

func (m *CategoryBudget) GetCategories() []*Category {
	if m != nil {
		return m.Categories
	}
	return nil
}

// Budget getters
func (m *Budget) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *Budget) GetBudgetName() string {
	if m != nil && m.BudgetName != nil {
		return *m.BudgetName
	}
	return ""
}

func (m *Budget) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *Budget) GetCategoryIDs() []string {
	if m != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *Budget) GetBudgetBreakdowns() []*BudgetBreakdown {
	if m != nil {
		return m.BudgetBreakdowns
	}
	return nil
}

// BudgetBreakdown getters
func (m *BudgetBreakdown) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *BudgetBreakdown) GetYear() int {
	if m != nil && m.Year != nil {
		return *m.Year
	}
	return 0
}

func (m *BudgetBreakdown) GetMonth() int {
	if m != nil && m.Month != nil {
		return *m.Month
	}
	return 0
}

// GetBudgetRequest getters
func (m *GetBudgetRequest) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *GetBudgetRequest) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

// GetBudgetResponse getter
func (m *GetBudgetResponse) GetCategoryBudget() *CategoryBudget {
	if m != nil {
		return m.CategoryBudget
	}
	return nil
}

// GetBudgetsRequest getter
func (m *GetBudgetsRequest) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

// GetBudgetsResponse getter
func (m *GetBudgetsResponse) GetBudgets() []*Budget {
	if m != nil {
		return m.Budgets
	}
	return nil
}

// SetBudgetRequest getters
func (m *SetBudgetRequest) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *SetBudgetRequest) GetBudgetName() string {
	if m != nil && m.BudgetName != nil {
		return *m.BudgetName
	}
	return ""
}

func (m *SetBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *SetBudgetRequest) GetBudgetAmount() string {
	if m != nil && m.BudgetAmount != nil {
		return *m.BudgetAmount
	}
	return ""
}

func (m *SetBudgetRequest) GetCategoryIDs() []string {
	if m != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *SetBudgetRequest) GetRangeStartDate() string {
	if m != nil && m.RangeStartDate != nil {
		return *m.RangeStartDate
	}
	return ""
}

func (m *SetBudgetRequest) GetRangeEndDate() string {
	if m != nil && m.RangeEndDate != nil {
		return *m.RangeEndDate
	}
	return ""
}
