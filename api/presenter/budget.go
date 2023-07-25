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

type CreateBudgetRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	BudgetType *uint32 `json:"budget_type,omitempty"`
	Amount     *string `json:"amount,omitempty"`
}

func (m *CreateBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateBudgetRequest) GetBudgetType() uint32 {
	if m != nil && m.BudgetType != nil {
		return *m.BudgetType
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
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		Amount:     amount,
		BudgetType: m.BudgetType,
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

type GetBudgetRequest struct{}

type GetBudgetResponse struct{}

type GetBudgetsRequest struct{}

type GetBudgetsResponse struct{}
