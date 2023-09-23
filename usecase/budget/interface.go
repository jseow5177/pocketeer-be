package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error)

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
	Currency     *string
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

func (m *CreateBudgetRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
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
	startDate, endDate, err := entity.GetBudgetStartEnd(
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
		entity.WithBudgetCurrency(m.Currency),
		entity.WithBudgetAmount(m.Amount),
		entity.WithBudgetType(goutil.Uint32(m.GetBudgetType())),
		entity.WithBudgetStartDate(goutil.Uint64(startDate)),
		entity.WithBudgetEndDate(goutil.Uint64(endDate)),
	)
}

func (m *CreateBudgetRequest) ToCategoryFilter() *repo.CategoryFilter {
	return repo.NewCategoryFilter(
		m.GetUserID(),
		repo.WithCategoryID(m.CategoryID),
	)
}

func (m *CreateBudgetRequest) ToGetBudgetFilter() *repo.GetBudgetFilter {
	return &repo.GetBudgetFilter{
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

func (m *DeleteBudgetRequest) ToGetBudgetFilter() *repo.GetBudgetFilter {
	return &repo.GetBudgetFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

func (m *DeleteBudgetRequest) ToDeleteBudgetFilter(budgetType uint32) *repo.DeleteBudgetFilter {
	return &repo.DeleteBudgetFilter{
		UserID:       m.UserID,
		CategoryID:   m.CategoryID,
		BudgetDate:   m.BudgetDate,
		BudgetRepeat: m.BudgetRepeat,
		BudgetType:   goutil.Uint32(uint32(budgetType)),
	}
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

func (m *GetBudgetRequest) ToGetBudgetFilter() *repo.GetBudgetFilter {
	return &repo.GetBudgetFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
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

func (m *UpdateBudgetRequest) ToGetBudgetFilter() *repo.GetBudgetFilter {
	return &repo.GetBudgetFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

func (m *UpdateBudgetRequest) ToDeleteBudgetFilter(budgetType uint32) *repo.DeleteBudgetFilter {
	return &repo.DeleteBudgetFilter{
		UserID:       m.UserID,
		CategoryID:   m.CategoryID,
		BudgetDate:   m.BudgetDate,
		BudgetRepeat: goutil.Uint32(uint32(entity.BudgetRepeatAllTime)),
		BudgetType:   goutil.Uint32(uint32(budgetType)),
	}
}

func (m *UpdateBudgetRequest) ToBudgetUpdate() (*entity.BudgetUpdate, error) {
	startDate, endDate, err := entity.GetBudgetStartEnd(
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
