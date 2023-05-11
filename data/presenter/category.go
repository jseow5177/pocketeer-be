package presenter

import (
	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type Category struct {
	CatID      *string `json:"cat_id,omitempty"`
	CatName    *string `json:"cat_name,omitempty"`
	CatType    *uint32 `json:"cat_type,omitempty"`
	CreateTime *uint64 `json:"create_time,omitempty"`
	UpdateTime *uint64 `json:"update_time,omitempty"`
}

func (m *Category) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return ""
}

func (m *Category) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
}

func (m *Category) GetCatType() uint32 {
	if m != nil && m.CatType != nil {
		return *m.CatType
	}
	return 0
}

func (m *Category) GetCreateTime() uint64 {
	if m != nil && m.CreateTime != nil {
		return *m.CreateTime
	}
	return 0
}

func (m *Category) GetUpdateTime() uint64 {
	if m != nil && m.UpdateTime != nil {
		return *m.UpdateTime
	}
	return 0
}

func ToCategoryPresenter(c *entity.Category) *Category {
	return &Category{
		CatID:      c.CatID,
		CatName:    c.CatName,
		CatType:    c.CatType,
		CreateTime: c.CreateTime,
		UpdateTime: c.UpdateTime,
	}
}

type CreateCategoryRequest struct {
	CatName *string `json:"cat_name,omitempty"`
	CatType *uint32 `json:"cat_type,omitempty"`
}

func (m *CreateCategoryRequest) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
}

func (m *CreateCategoryRequest) GetCatType() uint32 {
	if m != nil && m.CatType != nil {
		return *m.CatType
	}
	return 0
}

func (m *CreateCategoryRequest) ToCategoryEntity(userID string) *entity.Category {
	return &entity.Category{
		UserID:  goutil.String(userID),
		CatName: m.CatName,
		CatType: m.CatType,
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

func (m *CreateCategoryResponse) SetCategory(c *entity.Category) {
	m.Category = ToCategoryPresenter(c)
}

type GetCategoriesRequest struct {
	CatType *uint32 `json:"cat_type"`
}

func (m *GetCategoriesRequest) GetCatType() uint32 {
	if m != nil && m.CatType != nil {
		return *m.CatType
	}
	return 0
}

func (m *GetCategoriesRequest) ToCategoryFilter(userID string) *repo.CategoryFilter {
	return &repo.CategoryFilter{
		UserID:  goutil.String(userID),
		CatType: m.CatType,
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

func (m *GetCategoriesResponse) SetCategories(ecs []*entity.Category) {
	cs := make([]*Category, 0)
	for _, ec := range ecs {
		cs = append(cs, ToCategoryPresenter(ec))
	}
	m.Categories = cs
}

type UpdateCategoryRequest struct {
	CatID   *string `json:"cat_id,omitempty"`
	CatName *string `json:"cat_name,omitempty"`
}

func (m *UpdateCategoryRequest) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return ""
}

func (m *UpdateCategoryRequest) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
}

func (m *UpdateCategoryRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		CatID: m.CatID,
	}
}

func (m *UpdateCategoryRequest) ToCategoryEntity() *entity.Category {
	return &entity.Category{
		CatName: m.CatName,
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

func (m *UpdateCategoryResponse) SetCategory(c *entity.Category) {
	m.Category = ToCategoryPresenter(c)
}

type GetCategoryRequest struct {
	CatID *string `json:"cat_id"`
}

func (m *GetCategoryRequest) GetCatID() string {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return ""
}

func (m *GetCategoryRequest) ToCategoryFilter() *repo.CategoryFilter {
	return &repo.CategoryFilter{
		CatID: m.CatID,
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

func (m *GetCategoryResponse) SetCategory(c *entity.Category) {
	m.Category = ToCategoryPresenter(c)
}
