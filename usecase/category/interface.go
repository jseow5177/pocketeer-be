package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

type UseCase interface {
	GetCategory(ctx context.Context, req *GetCategoryRequest) (*GetCategoryResponse, error)
	GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error)

	GetCategoryBudget(ctx context.Context, req *GetCategoryBudgetRequest) (*GetCategoryBudgetResponse, error)
	GetCategoriesBudget(ctx context.Context, req *GetCategoriesBudgetRequest) (*GetCategoriesBudgetResponse, error)

	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error)
}

type GetCategoryBudgetRequest struct {
	UserID     *string
	CategoryID *string
	BudgetDate *string
	Timezone   *string
}

func (m *GetCategoryBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetCategoryBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *GetCategoryBudgetRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetCategoryBudgetRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetCategoryBudgetRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

func (m *GetCategoryBudgetRequest) ToTransactionFilter(userID string, startUnix, endUnix uint64) *repo.TransactionFilter {
	return &repo.TransactionFilter{
		UserID:             goutil.String(userID),
		CategoryID:         m.CategoryID,
		TransactionTimeGte: goutil.Uint64(startUnix),
		TransactionTimeLte: goutil.Uint64(endUnix),
	}
}

func (m *GetCategoryBudgetRequest) ToGetBudgetFilter() *repo.GetBudgetFilter {
	return &repo.GetBudgetFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
	}
}

type GetCategoryBudgetResponse struct {
	Category *entity.Category
}

func (m *GetCategoryBudgetResponse) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

type GetCategoryRequest struct {
	UserID     *string
	CategoryID *string
}

func (m *GetCategoryRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetCategoryRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetCategoryRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

type GetCategoryResponse struct {
	Category *entity.Category
}

func (m *GetCategoryResponse) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

type CreateCategoryRequest struct {
	UserID       *string
	CategoryName *string
	CategoryType *uint32
	Budget       *budget.CreateBudgetRequest
}

func (m *CreateCategoryRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateCategoryRequest) GetCategoryName() string {
	if m != nil && m.CategoryName != nil {
		return *m.CategoryName
	}
	return ""
}

func (m *CreateCategoryRequest) GetCategoryType() uint32 {
	if m != nil && m.CategoryType != nil {
		return *m.CategoryType
	}
	return 0
}

func (m *CreateCategoryRequest) GetBudget() *budget.CreateBudgetRequest {
	if m != nil && m.Budget != nil {
		return m.Budget
	}
	return nil
}

func (m *CreateCategoryRequest) ToCategoryEntity() (*entity.Category, error) {
	var b *entity.Budget
	if m.Budget != nil {
		var err error
		b, err = m.Budget.ToBudgetEntity()
		if err != nil {
			return nil, err
		}
	}
	return entity.NewCategory(
		m.GetUserID(),
		m.GetCategoryName(),
		entity.WithCategoryType(m.CategoryType),
		entity.WithCategoryBudget(b),
	)
}

func (m *CreateCategoryRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:       m.UserID,
		CategoryName: m.CategoryName,
		CategoryType: m.CategoryType,
	}
}

type CreateCategoryResponse struct {
	Category *entity.Category
}

func (m *CreateCategoryResponse) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

type UpdateCategoryRequest struct {
	UserID       *string
	CategoryID   *string
	CategoryName *string
}

func (m *UpdateCategoryRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateCategoryRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *UpdateCategoryRequest) GetCategoryName() string {
	if m != nil && m.CategoryName != nil {
		return *m.CategoryName
	}
	return ""
}

func (m *UpdateCategoryRequest) ToCategoryUpdate() *entity.CategoryUpdate {
	return entity.NewCategoryUpdate(
		entity.WithUpdateCategoryName(m.CategoryName),
	)
}

func (m *UpdateCategoryRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

type UpdateCategoryResponse struct {
	Category *entity.Category
}

func (m *UpdateCategoryResponse) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

type GetCategoriesRequest struct {
	UserID       *string
	CategoryType *uint32
	CategoryIDs  []string
}

func (m *GetCategoriesRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetCategoriesRequest) GetCategoryType() uint32 {
	if m != nil && m.CategoryType != nil {
		return *m.CategoryType
	}
	return 0
}

func (m *GetCategoriesRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *GetCategoriesRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:       m.UserID,
		CategoryType: m.CategoryType,
		CategoryIDs:  m.CategoryIDs,
	}
}

type GetCategoriesResponse struct {
	Categories []*entity.Category
}

func (m *GetCategoriesResponse) GetCategories() []*entity.Category {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

type GetCategoriesBudgetRequest struct {
	UserID      *string
	BudgetDate  *string
	CategoryIDs []string
	Timezone    *string
}

func (m *GetCategoriesBudgetRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetCategoriesBudgetRequest) GetBudgetDate() string {
	if m != nil && m.BudgetDate != nil {
		return *m.BudgetDate
	}
	return ""
}

func (m *GetCategoriesBudgetRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *GetCategoriesBudgetRequest) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

func (m *GetCategoriesBudgetRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:      m.UserID,
		CategoryIDs: m.CategoryIDs,
	}
}

func (m *GetCategoriesBudgetRequest) ToGetCategoryBudgetRequest(categoryID string) *GetCategoryBudgetRequest {
	return &GetCategoryBudgetRequest{
		CategoryID: goutil.String(categoryID),
		BudgetDate: m.BudgetDate,
		Timezone:   m.Timezone,
	}
}

type GetCategoriesBudgetResponse struct {
	Categories []*entity.Category
}

func (m *GetCategoriesBudgetResponse) GetCategories() []*entity.Category {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}
