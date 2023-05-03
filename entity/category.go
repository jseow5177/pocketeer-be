package entity

type Category struct {
	CategoryID   *uint64
	CategoryName *string
}

func (c *Category) GetCategoryID() uint64 {
	if c != nil && c.CategoryID != nil {
		return *c.CategoryID
	}
	return 0
}

func (c *Category) GetCategoryName() string {
	if c != nil && c.CategoryName != nil {
		return *c.CategoryName
	}
	return ""
}
