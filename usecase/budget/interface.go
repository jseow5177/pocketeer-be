package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error)
	GetBudgets(ctx context.Context, req *GetBudgetsRequest) (*GetBudgetsResponse, error)

	CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*CreateBudgetResponse, error)
}

type CreateBudgetRequest struct {
	CategoryID *string
	UserID     *string
	BudgetType *uint32
	Amount     *float64
}

func (m *CreateBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
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

func (m *CreateBudgetRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *CreateBudgetRequest) ToBudgetEntity() *entity.Budget {
	return entity.NewBudget(
		m.GetUserID(),
		m.GetCategoryID(),
		entity.WithBudgetType(m.BudgetType),
		entity.WithBudgetAmount(m.Amount),
	)
}

type CreateBudgetResponse struct {
	Budget *entity.Budget
}

func (m *CreateBudgetResponse) GetBudget() *entity.Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

type GetBudgetsRequest struct{}

type GetBudgetsResponse struct{}

type GetBudgetRequest struct{}

type GetBudgetResponse struct{}
