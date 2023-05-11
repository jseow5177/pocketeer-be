package entity

type Category struct {
	UserID     *string
	CatID      *string
	CatName    *string
	CatType    *uint32
	CreateTime *uint64
	UpdateTime *uint64
}

func (c *Category) GetUserID() string {
	if c != nil && c.UserID != nil {
		return *c.UserID
	}
	return ""
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
