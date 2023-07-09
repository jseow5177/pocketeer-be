package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type CategoryUpdate struct {
	CategoryName *string
	UpdateTime   *uint64
}

func (cu *CategoryUpdate) GetCategoryName() string {
	if cu != nil && cu.CategoryName != nil {
		return *cu.CategoryName
	}
	return ""
}

func (cu *CategoryUpdate) GetUpdateTime() uint64 {
	if cu != nil && cu.UpdateTime != nil {
		return *cu.UpdateTime
	}
	return 0
}

type CategoryUpdateOption func(cu *CategoryUpdate)

func WithUpdateCategoryName(categoryName *string) CategoryUpdateOption {
	return func(cu *CategoryUpdate) {
		cu.CategoryName = categoryName
	}
}

func NewCategoryUpdate(opts ...CategoryUpdateOption) *CategoryUpdate {
	cu := new(CategoryUpdate)
	for _, opt := range opts {
		opt(cu)
	}
	return cu
}

type Category struct {
	UserID       *string
	CategoryID   *string
	CategoryName *string
	CategoryType *uint32
	CreateTime   *uint64
	UpdateTime   *uint64
}

type CategoryOption = func(c *Category)

func WithCategoryID(categoryID *string) CategoryOption {
	return func(c *Category) {
		c.CategoryID = categoryID
	}
}

func WithCategoryName(categoryName *string) CategoryOption {
	return func(c *Category) {
		c.CategoryName = categoryName
	}
}

func WithCategoryType(categoryType *uint32) CategoryOption {
	return func(c *Category) {
		c.CategoryType = categoryType
	}
}

func WithCategoryCreateTime(createTime *uint64) CategoryOption {
	return func(c *Category) {
		c.CreateTime = createTime
	}
}

func WithCategoryUpdateTime(updateTime *uint64) CategoryOption {
	return func(c *Category) {
		c.UpdateTime = updateTime
	}
}

func NewCategory(userID string, opts ...CategoryOption) *Category {
	now := uint64(time.Now().Unix())
	c := &Category{
		UserID:       goutil.String(userID),
		CategoryType: goutil.Uint32(uint32(TransactionTypeExpense)),
		CreateTime:   goutil.Uint64(now),
		UpdateTime:   goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(c)
	}
	c.checkOpts()
	return c
}

func setCategory(c *Category, opts ...CategoryOption) {
	if c == nil {
		return
	}

	for _, opt := range opts {
		opt(c)
	}
}

func (c *Category) checkOpts() {}

func (c *Category) Update(cu *CategoryUpdate) (categoryUpdate *CategoryUpdate, hasUpdate bool) {
	categoryUpdate = new(CategoryUpdate)

	if cu.CategoryName != nil && cu.GetCategoryName() != c.GetCategoryName() {
		hasUpdate = true
		setCategory(c, WithCategoryName(cu.CategoryName))
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().Unix()))
	setCategory(c, WithCategoryUpdateTime(now))

	// check
	c.checkOpts()

	categoryUpdate.UpdateTime = now

	if cu.CategoryName != nil {
		categoryUpdate.CategoryName = c.CategoryName
	}

	return
}

func (c *Category) GetUserID() string {
	if c != nil && c.UserID != nil {
		return *c.UserID
	}
	return ""
}

func (c *Category) GetCategoryID() string {
	if c != nil && c.CategoryID != nil {
		return *c.CategoryID
	}
	return ""
}

func (c *Category) SetCategoryID(categoryID *string) {
	setCategory(c, WithCategoryID(categoryID))
}

func (c *Category) GetCategoryName() string {
	if c != nil && c.CategoryName != nil {
		return *c.CategoryName
	}
	return ""
}

func (f *Category) GetCategoryType() uint32 {
	if f != nil && f.CategoryType != nil {
		return *f.CategoryType
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
