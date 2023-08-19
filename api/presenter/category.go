package presenter

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
	Budget       *Budget `json:"budget,omitempty"`
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

func (c *Category) GetBudget() *Budget {
	if c != nil && c.Budget != nil {
		return c.Budget
	}
	return nil
}

type CreateCategoriesRequest struct {
	Categories []*CreateCategoryRequest `json:"categories,omitempty"`
}

func (m *CreateCategoriesRequest) GetCategories() []*CreateCategoryRequest {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

func (m *CreateCategoriesRequest) ToUseCaseReq(userID string) *category.CreateCategoriesRequest {
	cs := make([]*category.CreateCategoryRequest, 0)
	for _, c := range m.Categories {
		cs = append(cs, c.ToUseCaseReq(userID))
	}
	return &category.CreateCategoriesRequest{
		Categories: cs,
	}
}

type CreateCategoriesResponse struct {
	Categories []*Category `json:"categories,omitempty"`
}

func (m *CreateCategoriesResponse) Set(useCaseRes *category.CreateCategoriesResponse) {
	m.Categories = toCategories(useCaseRes.Categories)
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
	CategoryType *uint32  `json:"category_type,omitempty"`
	CategoryIDs  []string `json:"category_ids,omitempty"`
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

func (m *GetCategoriesRequest) ToUseCaseReq(userID string) *category.GetCategoriesRequest {
	return &category.GetCategoriesRequest{
		UserID:       goutil.String(userID),
		CategoryType: m.CategoryType,
		CategoryIDs:  m.CategoryIDs,
	}
}

type GetCategoriesResponse struct {
	Categories []*Category `json:"categories,omitempty"`
}

func (m *GetCategoriesResponse) GetCategories() []*Category {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

func (m *GetCategoriesResponse) Set(useCaseRes *category.GetCategoriesResponse) {
	m.Categories = toCategories(useCaseRes.Categories)
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
		CategoryID:   m.CategoryID,
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
	CategoryID *string `json:"category_id,omitempty"`
}

func (m *GetCategoryRequest) GetCategoryID() string {
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

type GetCategoryBudgetRequest struct {
	CategoryID *string `json:"category_id,omitempty"`
	BudgetDate *string `json:"budget_date,omitempty"`
	Timezone   *string `json:"timezone,omitempty"`
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

func (m *GetCategoryBudgetRequest) ToUseCaseReq(userID string) *category.GetCategoryBudgetRequest {
	return &category.GetCategoryBudgetRequest{
		UserID:     goutil.String(userID),
		CategoryID: m.CategoryID,
		BudgetDate: m.BudgetDate,
		Timezone:   m.Timezone,
	}
}

type GetCategoryBudgetResponse struct {
	Category *Category `json:"category,omitempty"`
}

func (m *GetCategoryBudgetResponse) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *GetCategoryBudgetResponse) Set(res *category.GetCategoryBudgetResponse) {
	m.Category = toCategory(res.Category)
}

type GetCategoriesBudgetRequest struct {
	CategoryIDs []string `json:"category_ids,omitempty"`
	BudgetDate  *string  `json:"budget_date,omitempty"`
	Timezone    *string  `json:"timezone,omitempty"`
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

func (m *GetCategoriesBudgetRequest) ToUseCaseReq(userID string) *category.GetCategoriesBudgetRequest {
	return &category.GetCategoriesBudgetRequest{
		UserID:      goutil.String(userID),
		Timezone:    m.Timezone,
		BudgetDate:  m.BudgetDate,
		CategoryIDs: m.CategoryIDs,
	}
}

type GetCategoriesBudgetResponse struct {
	Categories []*Category `json:"categories,omitempty"`
}

func (m *GetCategoriesBudgetResponse) GetCategories() []*Category {
	if m != nil && m.Categories != nil {
		return m.Categories
	}
	return nil
}

func (m *GetCategoriesBudgetResponse) Set(res *category.GetCategoriesBudgetResponse) {
	m.Categories = toCategories(res.Categories)
}
