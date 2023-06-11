package budget

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetBudget(
		ctx context.Context,
		req *GetBudgetRequest,
	) (*GetBudgetResponse, error)

	GetBudgets(
		ctx context.Context,
		req *GetBudgetsRequest,
	) (*GetBudgetsResponse, error)

	SetBudget(
		ctx context.Context,
		req *SetBudgetRequest,
	) (*SetBudgetResponse, error)
}

type GetBudgetsRequest struct {
	UserID *string
	Date   time.Time
}

type GetBudgetsResponse struct {
	Budgets []*entity.Budget
}

type GetBudgetRequest struct {
	UserID   *string
	BudgetID *string
	Date     time.Time
}

type GetBudgetResponse struct {
	Budget *entity.Budget
}

type SetBudgetRequest struct {
	UserID         *string
	BudgetID       *string
	BudgetName     *string
	BudgetType     *uint32
	BudgetAmount   *float64
	CategoryIDs    []string
	RangeStartDate time.Time
	RangeEndDate   time.Time
}

type SetBudgetResponse struct{}

func (m *GetBudgetsRequest) ToBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID: m.UserID,
	}
}

func (m *GetBudgetRequest) ToBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:   m.UserID,
		BudgetID: m.BudgetID,
	}
}

func (m *SetBudgetRequest) ToBudgetFilter() *repo.BudgetFilter {
	return &repo.BudgetFilter{
		UserID:   m.UserID,
		BudgetID: m.BudgetID,
	}
}

func (m *GetBudgetsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetBudgetsRequest) GetDate() time.Time {
	if m != nil {
		return m.Date
	}
	return time.Time{}
}

func (m *GetBudgetRequest) GetDate() time.Time {
	if m != nil {
		return m.Date
	}
	return time.Time{}
}

func (m *GetBudgetRequest) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *GetBudgetResponse) GetBudget() *entity.Budget {
	if m != nil {
		return m.Budget
	}
	return nil
}

func (m *SetBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

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

func (m *SetBudgetRequest) GetBudgetAmount() float64 {
	if m != nil && m.BudgetAmount != nil {
		return *m.BudgetAmount
	}
	return 0
}

func (m *SetBudgetRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *SetBudgetRequest) GetRangeStartDate() time.Time {
	return m.RangeStartDate
}

func (m *SetBudgetRequest) GetRangeEndDate() time.Time {
	return m.RangeEndDate
}
