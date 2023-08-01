package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/util"
)

type Budget struct {
	BudgetID     *string `json:"budget_id,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	BudgetType   *uint32 `json:"budget_type,omitempty"`
	BudgetStatus *uint32 `json:"budget_status,omitempty"`
	Amount       *string `json:"amount,omitempty"`
	CreateTime   *uint64 `json:"create_time,omitempty"`
	UpdateTime   *uint64 `json:"update_time,omitempty"`
	UsedAmount   *string `json:"used_amount,omitempty"`
}

func (b *Budget) GetBudgetID() string {
	if b != nil && b.BudgetID != nil {
		return *b.BudgetID
	}
	return ""
}

func (b *Budget) GetCategoryID() string {
	if b != nil && b.CategoryID != nil {
		return *b.CategoryID
	}
	return ""
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil && b.BudgetType != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) GetBudgetStatus() uint32 {
	if b != nil && b.BudgetStatus != nil {
		return *b.BudgetStatus
	}
	return 0
}

func (b *Budget) GetAmount() string {
	if b != nil && b.Amount != nil {
		return *b.Amount
	}
	return ""
}

func (b *Budget) GetUsedAmount() string {
	if b != nil && b.UsedAmount != nil {
		return *b.UsedAmount
	}
	return ""
}

func (b *Budget) GetCreateTime() uint64 {
	if b != nil && b.CreateTime != nil {
		return *b.CreateTime
	}
	return 0
}

func (b *Budget) GetUpdateTime() uint64 {
	if b != nil && b.UpdateTime != nil {
		return *b.UpdateTime
	}
	return 0
}

type GetBudgetsRequest struct {
	CategoryIDs []string `json:"category_ids,omitempty"`
	BudgetDate  *string  `json:"budget_date,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
}

func (m *GetBudgetsRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *GetBudgetsRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *GetBudgetsRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetBudgetsRequest) ToUseCaseReq(userID string) *budget.GetBudgetsRequest {
	return &budget.GetBudgetsRequest{
		UserID:      goutil.String(userID),
		Timezone:    m.Timezone,
		BudgetDate:  m.BudgetDate,
		CategoryIDs: m.CategoryIDs,
	}
}

type GetBudgetsResponse struct {
	Budgets []*Budget `json:"budgets,omitempty"`
}

func (m *GetBudgetsResponse) GetBudgets() []*Budget {
	if m != nil && m.Budgets != nil {
		return m.Budgets
	}
	return nil
}

func (m *GetBudgetsResponse) Set(res *budget.GetBudgetsResponse) {
	m.Budgets = toBudgets(res.Budgets)
}

type GetBudgetRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	BudgetDate *string `json:"budget_date,omitempty"`
	Timezone   *string `json:"timezone,omitempty"`
}

func (m *GetBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *GetBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetBudgetRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetBudgetRequest) ToUseCaseReq(userID string) *budget.GetBudgetRequest {
	return &budget.GetBudgetRequest{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
		Timezone:   m.Timezone,
	}
}

type GetBudgetResponse struct {
	Budget *Budget `json:"budget,omitempty"`
}

func (m *GetBudgetResponse) GetBudget() *Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

func (m *GetBudgetResponse) Set(res *budget.GetBudgetResponse) {
	m.Budget = toBudget(res.Budget)
}

type CreateBudgetRequest struct {
	BudgetDate   *string `json:"budget_date,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	BudgetType   *uint32 `json:"budget_type,omitempty"`
	BudgetRepeat *uint32 `json:"budget_repeat,omitempty"`
	Amount       *string `json:"amount,omitempty"`
}

func (m *CreateBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *CreateBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *CreateBudgetRequest) GetBudgetRepeat() uint32 {
	if m != nil && m.BudgetRepeat != nil {
		return *m.BudgetRepeat
	}
	return 0
}

func (m *CreateBudgetRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
}

func (m *CreateBudgetRequest) ToUseCaseReq(userID string) *budget.CreateBudgetRequest {
	var amount *float64
	if m.Amount != nil {
		a, _ := util.MonetaryStrToFloat(m.GetAmount())
		amount = goutil.Float64(a)
	}
	return &budget.CreateBudgetRequest{
		UserID:       goutil.String(userID),
		CategoryID:   m.CategoryID,
		Amount:       amount,
		BudgetType:   m.BudgetType,
		BudgetDate:   m.BudgetDate,
		BudgetRepeat: m.BudgetRepeat,
	}
}

type CreateBudgetResponse struct {
	Budget *Budget `json:"budget,omitempty"`
}

func (m *CreateBudgetResponse) GetBudget() *Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

func (m *CreateBudgetResponse) Set(res *budget.CreateBudgetResponse) {
	m.Budget = toBudget(res.Budget)
}

type DeleteBudgetRequest struct {
	BudgetDate   *string `json:"budget_date,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	BudgetRepeat *uint32 `json:"budget_repeat,omitempty"`
}

func (m *DeleteBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *DeleteBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *DeleteBudgetRequest) GetBudgetRepeat() uint32 {
	if m != nil && m.BudgetRepeat != nil {
		return *m.BudgetRepeat
	}
	return 0
}

func (m *DeleteBudgetRequest) ToUseCaseReq(userID string) *budget.DeleteBudgetRequest {
	return &budget.DeleteBudgetRequest{
		UserID:       goutil.String(userID),
		CategoryID:   m.CategoryID,
		BudgetDate:   m.BudgetDate,
		BudgetRepeat: m.BudgetRepeat,
	}
}

type DeleteBudgetResponse struct{}

func (m *DeleteBudgetResponse) Set(res *budget.DeleteBudgetResponse) {}

type UpdateBudgetRequest struct {
	BudgetDate   *string `json:"budget_date,omitempty"`
	CategoryID   *string `json:"category_id,omitempty"`
	BudgetType   *uint32 `json:"budget_type,omitempty"`
	BudgetRepeat *uint32 `json:"budget_repeat,omitempty"`
	Amount       *string `json:"amount,omitempty"`
}

func (m *UpdateBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *UpdateBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *UpdateBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
	}
	return 0
}

func (m *UpdateBudgetRequest) GetBudgetRepeat() uint32 {
	if m != nil && m.BudgetRepeat != nil {
		return *m.BudgetRepeat
	}
	return 0
}

func (m *UpdateBudgetRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
}

type UpdateBudgetResponse struct {
	Budget *Budget `json:"budget,omitempty"`
}

func (m *UpdateBudgetResponse) GetBudget() *Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}
