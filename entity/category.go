package entity

type CatType uint32

const (
	CatTypeExpense CatType = 1
	CatTypeIncome  CatType = 2
)

type CategoryFilter struct {
	CatID   *uint64
	CatType *uint32
}

func (f *CategoryFilter) GetCatID() uint64 {
	if f != nil && f.CatID != nil {
		return *f.CatID
	}
	return 0
}

func (f *CategoryFilter) GetCatType() uint32 {
	if f != nil && f.CatType != nil {
		return *f.CatType
	}
	return 0
}

type Category struct {
	CatID      *uint64
	CatName    *string
	CatType    *uint32
	CreateTime *uint64
	UpdateTime *uint64
}

func (c *Category) GetCatID() uint64 {
	if c != nil && c.CatID != nil {
		return *c.CatID
	}
	return 0
}

func (c *Category) GetCatName() string {
	if c != nil && c.CatName != nil {
		return *c.CatName
	}
	return ""
}

func (f *Category) GetCatType() uint32 {
	if f != nil && f.CatType != nil {
		return *f.CatType
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
