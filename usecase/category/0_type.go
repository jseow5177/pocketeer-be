package category

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetCategory(ctx context.Context, req *GetCategoryRequest) (*GetCategoryResponse, error)
	GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error)

	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error)
	UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error)
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

func (m *GetCategoryRequest) GetCatID() string {
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

func (m *CreateCategoryRequest) ToCategoryEntity() *entity.Category {
	return &entity.Category{
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

func (m *UpdateCategoryRequest) ToGetCategoryRequest() *GetCategoryRequest {
	return &GetCategoryRequest{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
	}
}

func (m *UpdateCategoryRequest) ToCategoryEntity() *entity.Category {
	return &entity.Category{
		CategoryName: m.CategoryName,
	}
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

func (m *GetCategoriesRequest) GetCategoryTypee() uint32 {
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
