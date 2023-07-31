package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type UseCase interface {
	GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error)
	GetBudgets(ctx context.Context, req *GetBudgetsRequest) (*GetBudgetsResponse, error)

	CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*CreateBudgetResponse, error)
	DeleteBudget(ctx context.Context, req *DeleteBudgetRequest) (*DeleteBudgetResponse, error)
}

type CreateBudgetRequest struct {
	CategoryID   *string
	UserID       *string
	BudgetDate   *string
	BudgetType   *uint32
	BudgetRepeat *uint32
	Amount       *float64
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

func (m *CreateBudgetRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *CreateBudgetRequest) ToBudgetEntity() (*entity.Budget, error) {
	return entity.NewBudget(
		m.GetUserID(),
		m.GetCategoryID(),
		entity.WithBudgetDate(m.BudgetDate),
		entity.WithBudgetAmount(m.Amount),
		entity.WithBudgetType(m.BudgetType),
		entity.WithBudgetRepeat(m.BudgetRepeat),
	)
}

func (m *CreateBudgetRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

func (m *CreateBudgetRequest) ToGetBudgetRequest() *GetBudgetRequest {
	return &GetBudgetRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
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

type DeleteBudgetRequest struct {
	UserID       *string
	BudgetDate   *string
	CategoryID   *string
	BudgetRepeat *uint32
}

func (m *DeleteBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *DeleteBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *DeleteBudgetRequest) ToGetBudgetRequest() *GetBudgetRequest {
	return &GetBudgetRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

func (m *DeleteBudgetRequest) ToBudgetEntity(budgetType uint32) (*entity.Budget, error) {
	return entity.NewBudget(
		m.GetUserID(),
		m.GetCategoryID(),
		entity.WithBudgetDate(m.BudgetDate),
		entity.WithBudgetAmount(goutil.Float64(0)),
		entity.WithBudgetType(goutil.Uint32(budgetType)),
		entity.WithBudgetRepeat(m.BudgetRepeat),
		entity.WithBudgetStatus(goutil.Uint32(uint32(entity.BudgetStatusDeleted))),
	)
}

type DeleteBudgetResponse struct{}

type GetBudgetRequest struct {
	UserID     *string
	CategoryID *string
	BudgetDate *string
}

func (m *GetBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
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

func (m *GetBudgetRequest) ToBudgetQuery() (*repo.BudgetQuery, *repo.Paging, error) {
	t, err := util.ParseDateStrToInt(m.GetBudgetDate())
	if err != nil {
		return nil, nil, err
	}

	return &repo.BudgetQuery{
			Queries: []*repo.BudgetQuery{
				{
					Filters: []*repo.BudgetFilter{
						{
							StartDateLte: goutil.Uint64(t),
							EndDateGte:   goutil.Uint64(t),
						},
						{
							StartDateLte: goutil.Uint64(t),
							EndDate:      goutil.Uint64(0),
						},
						{
							StartDate: goutil.Uint64(0),
							EndDate:   goutil.Uint64(0),
						},
					},
					Op: filter.Or,
				},
				{
					Filters: []*repo.BudgetFilter{
						{
							UserID:     m.UserID,
							CategoryID: m.CategoryID,
						},
					},
				},
			},
			Op: filter.And,
		},
		&repo.Paging{
			Limit: goutil.Uint32(1),
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("create_time"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		}, nil
}

type GetBudgetResponse struct {
	Budget *entity.Budget
}

func (m *GetBudgetResponse) GetBudget() *entity.Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

type UpdateBudgetRequest struct{}

type UpdateBudgetResponse struct {
	Budget *entity.Budget
}

func (m *UpdateBudgetResponse) GetBudget() *entity.Budget {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

type GetBudgetsRequest struct {
	UserID      *string
	CategoryIDs []string
	BudgetDate  *string
}

func (m *GetBudgetsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
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

type GetBudgetsResponse struct {
	Budgets []*entity.Budget
}

func (m *GetBudgetsResponse) GetBudgets() []*entity.Budget {
	if m != nil && m.Budgets != nil {
		return m.Budgets
	}
	return nil
}
