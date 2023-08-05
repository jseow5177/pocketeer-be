package budget

import (
	"context"
	"time"

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
	UpdateBudget(ctx context.Context, req *UpdateBudgetRequest) (*UpdateBudgetResponse, error)
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
	startDate, endDate, err := entity.GetBudgetDateRange(
		m.GetBudgetDate(),
		m.GetBudgetType(),
		m.GetBudgetRepeat(),
	)
	if err != nil {
		return nil, err
	}

	return entity.NewBudget(
		m.GetUserID(),
		m.GetCategoryID(),
		entity.WithBudgetAmount(m.Amount),
		entity.WithBudgetType(m.BudgetType),
		entity.WithBudgetStartDate(goutil.Uint64(startDate)),
		entity.WithBudgetEndDate(goutil.Uint64(endDate)),
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
	DeleteTime   *uint64
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

func (m *DeleteBudgetRequest) GetDeleteTime() uint64 {
	if m != nil && m.DeleteTime != nil {
		return *m.DeleteTime
	}
	return 0
}

func (m *DeleteBudgetRequest) ToGetBudgetRequest() *GetBudgetRequest {
	return &GetBudgetRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

func (m *DeleteBudgetRequest) ToBudgetEntity(budgetType uint32, deleteTime uint64) (*entity.Budget, error) {
	startDate, endDate, err := entity.GetBudgetDateRange(
		m.GetBudgetDate(),
		budgetType,
		m.GetBudgetRepeat(),
	)
	if err != nil {
		return nil, err
	}

	t := uint64(time.Now().UnixMilli())
	if deleteTime != 0 {
		t = deleteTime
	}

	return entity.NewBudget(
		m.GetUserID(),
		m.GetCategoryID(),
		entity.WithBudgetAmount(goutil.Float64(0)),
		entity.WithBudgetType(goutil.Uint32(budgetType)),
		entity.WithBudgetStatus(goutil.Uint32(uint32(entity.BudgetStatusDeleted))),
		entity.WithBudgetStartDate(goutil.Uint64(startDate)),
		entity.WithBudgetEndDate(goutil.Uint64(endDate)),
		entity.WithBudgetCreateTime(goutil.Uint64(t)),
		entity.WithBudgetUpdateTime(goutil.Uint64(t)),
	)
}

type DeleteBudgetResponse struct{}

type GetBudgetRequest struct {
	UserID     *string
	CategoryID *string
	BudgetDate *string
	Timezone   *string
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

func (m *GetBudgetRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetBudgetRequest) ToTransactionFilter(userID string, startUnix, endUnix uint64) *repo.TransactionFilter {
	return &repo.TransactionFilter{
		UserID:             goutil.String(userID),
		CategoryID:         m.CategoryID,
		TransactionTimeGte: goutil.Uint64(startUnix),
		TransactionTimeLte: goutil.Uint64(endUnix),
	}
}

func (m *GetBudgetRequest) ToBudgetQuery() (*repo.BudgetQuery, error) {
	t, err := util.ParseDateStrToInt(m.GetBudgetDate())
	if err != nil {
		return nil, err
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
		Paging: &repo.Paging{
			Limit: goutil.Uint32(1),
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("update_time"),
					Order: goutil.String(config.OrderDesc),
				},
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

type UpdateBudgetRequest struct {
	CategoryID   *string
	UserID       *string
	BudgetDate   *string
	BudgetType   *uint32
	BudgetRepeat *uint32
	Amount       *float64
}

func (m *UpdateBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
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

func (m *UpdateBudgetRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *UpdateBudgetRequest) ToGetBudgetRequest() *GetBudgetRequest {
	return &GetBudgetRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

func (m *UpdateBudgetRequest) ToDeleteBudgetRequest(deleteTime uint64) *DeleteBudgetRequest {
	return &DeleteBudgetRequest{
		UserID:       m.UserID,
		CategoryID:   m.CategoryID,
		BudgetDate:   m.BudgetDate,
		BudgetRepeat: goutil.Uint32(uint32(entity.BudgetRepeatAllTime)),
		DeleteTime:   goutil.Uint64(deleteTime),
	}
}

func (m *UpdateBudgetRequest) ToBudgetUpdate() (*entity.BudgetUpdate, error) {
	startDate, endDate, err := entity.GetBudgetDateRange(
		m.GetBudgetDate(),
		m.GetBudgetType(),
		m.GetBudgetRepeat(),
	)
	if err != nil {
		return nil, err
	}
	return entity.NewBudgetUpdate(
		entity.WithUpdateBudgetAmount(m.Amount),
		entity.WithUpdateBudgetType(m.BudgetType),
		entity.WithUpdateBudgetStartDate(goutil.Uint64(startDate)),
		entity.WithUpdateBudgetEndDate(goutil.Uint64(endDate)),
	), nil
}

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
	BudgetDate  *string
	CategoryIDs []string
	Timezone    *string
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

func (m *GetBudgetsRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetBudgetsRequest) ToGetBudgetRequests(userID string) []*GetBudgetRequest {
	reqs := make([]*GetBudgetRequest, 0)
	for _, catID := range m.CategoryIDs {
		reqs = append(reqs, &GetBudgetRequest{
			UserID:     goutil.String(userID),
			CategoryID: goutil.String(catID),
			BudgetDate: m.BudgetDate,
			Timezone:   m.Timezone,
		})
	}
	return reqs
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
