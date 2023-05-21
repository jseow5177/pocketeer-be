package presenter

//go:generate easytags $GOFILE

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/category"
)

type Category struct {
	CategoryID   *string `json:"category_id,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
	CategoryType *uint32 `json:"category_type,omitempty"`
	CreateTime   *uint64 `json:"create_time,omitempty"`
	UpdateTime   *uint64 `json:"update_time,omitempty"`
}

func (c *Category) GetCategoryID() string {
	if c != nil && c.CategoryID != nil {
		return *c.CategoryID
	}
	return ""
}

func (c *Category) GetCategoryName() string {
	if c != nil && c.CategoryName != nil {
		return *c.CategoryName
	}
	return ""
}

func (c *Category) GetCategoryType() uint32 {
	if c != nil && c.CategoryType != nil {
		return *c.CategoryType
	}
	return 0
}

func (c *Category) GetCreateTime() uint64 {
	if c != nil && c.CreateTime != nil {
		return *c.CreateTime
	}
	return 0
}

func (c *Category) GetUpdateTime() uint64 {
	if c != nil && c.UpdateTime != nil {
		return *c.UpdateTime
	}
	return 0
}

type CreateCategoryRequest struct {
	CategoryName *string `json:"category_name,omitempty"`
	CategoryType *uint32 `json:"category_type,omitempty"`
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

func (m *CreateCategoryRequest) ToUseCaseReq(userID string) *category.CreateCategoryRequest {
	return &category.CreateCategoryRequest{
		UserID:       goutil.String(userID),
		CategoryName: m.CategoryName,
		CategoryType: m.CategoryType,
	}
}

type CreateCategoryResponse struct {
	Category *Category `json:"category,omitempty"`
}

func (m *CreateCategoryResponse) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *CreateCategoryResponse) Set(useCaseRes *category.CreateCategoryResponse) {
	m.Category = toCategory(useCaseRes.Category)
}

type GetCategoriesRequest struct {
	CategoryType *uint32 `json:"category_type"`
}

func (m *GetCategoriesRequest) GetCategoryTypee() uint32 {
	if m != nil && m.CategoryType != nil {
		return *m.CategoryType
	}
	return 0
}

func (m *GetCategoriesRequest) ToUseCaseReq(userID string) *category.GetCategoriesRequest {
	return &category.GetCategoriesRequest{
		UserID:       goutil.String(userID),
		CategoryType: m.CategoryType,
	}
}

type GetCategoriesResponse struct {
	Categories []*Category `json:"categories"`
}

func (m *GetCategoriesResponse) GetCategories() []*Category {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

func (m *GetCategoriesResponse) Set(useCaseRes *category.GetCategoriesResponse) {
	cs := make([]*Category, 0)
	for _, c := range useCaseRes.Categories {
		cs = append(cs, toCategory(c))
	}
	m.Categories = cs
}

type UpdateCategoryRequest struct {
	CategoryID   *string `json:"category_id,omitempty"`
	CategoryName *string `json:"category_name,omitempty"`
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

func (m *UpdateCategoryRequest) ToUseCaseReq(userID string) *category.UpdateCategoryRequest {
	return &category.UpdateCategoryRequest{
		UserID:       goutil.String(userID),
		CategoryName: m.CategoryName,
	}
}

type UpdateCategoryResponse struct {
	Category *Category `json:"category,omitempty"`
}

func (m *UpdateCategoryResponse) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *UpdateCategoryResponse) Set(useCaseRes *category.UpdateCategoryResponse) {
	m.Category = toCategory(useCaseRes.Category)
}

type GetCategoryRequest struct {
	CategoryID *string `json:"category_id"`
}

func (m *GetCategoryRequest) GetCatID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetCategoryRequest) ToUseCaseReq(userID string) *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
	}
}

type GetCategoryResponse struct {
	Category *Category `json:"category,omitempty"`
}

func (m *GetCategoryResponse) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *GetCategoryResponse) Set(useCaseRes *category.GetCategoryResponse) {
	m.Category = toCategory(useCaseRes.Category)
}
