package presenter

type Category struct {
	CatID      *uint64 `json:"cat_id"`
	CatName    *string `json:"cat_name"`
	CreateTime *uint64 `json:"create_time"`
	UpdateTime *uint64 `json:"update_time"`
}

func (m *Category) GetCatID() uint64 {
	if m != nil && m.CatID != nil {
		return *m.CatID
	}
	return 0
}

func (m *Category) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
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

type CreateCategoryRequest struct {
	CatName *string `json:"cat_name"`
}

func (m *CreateCategoryRequest) GetCatName() string {
	if m != nil && m.CatName != nil {
		return *m.CatName
	}
	return ""
}

type CreateCategoryResponse struct {
	Category *Category `json:"category"`
}

func (m *CreateCategoryResponse) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}
