package aggr

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetBudgetWithCategories(ctx context.Context, req *GetBudgetWithCategoriesRequest) (*GetBudgetWithCategoriesResponse, error)
}

type GetBudgetWithCategoriesRequest struct {
	UserID   *string
	BudgetID *string
	Date     time.Time
}

type GetBudgetWithCategoriesResponse struct {
	Budget     *entity.Budget
	Categories []*entity.Category
}

func (m *GetBudgetWithCategoriesRequest) GetDate() time.Time {
	if m != nil {
		return m.Date
	}
	return time.Time{}
}

func (m *GetBudgetWithCategoriesRequest) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

func (m *GetBudgetWithCategoriesResponse) GetBudget() *entity.Budget {
	if m != nil {
		return m.Budget
	}
	return nil
}

func (m *GetBudgetWithCategoriesResponse) GetCategories() []*entity.Category {
	if m != nil {
		return m.Categories
	}
	return nil
}
