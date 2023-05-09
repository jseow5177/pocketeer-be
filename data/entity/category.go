package entity

type CatType uint32

const (
	CatTypeExpense CatType = 1
	CatTypeIncome  CatType = 2
)

var CatTypes = map[uint32]string{
	uint32(CatTypeExpense): "expense",
	uint32(CatTypeIncome):  "income",
}

type CategoryFilter struct {
	CatID   *string `filter:"_id"`
	CatType *uint32 `filter:"cat_type"`
}

func (f *CategoryFilter) GetCatID() string {
	if f != nil && f.CatID != nil {
		return *f.CatID
	}
	return ""
}

func (f *CategoryFilter) GetCatType() uint32 {
	if f != nil && f.CatType != nil {
		return *f.CatType
	}
	return 0
}

type Category struct {
	CatID      *string
	CatName    *string
	CatType    *uint32
	CreateTime *uint64
	UpdateTime *uint64
}

func (c *Category) GetCatID() string {
	if c != nil && c.CatID != nil {
		return *c.CatID
	}
	return ""
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
