package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type CategoryStatus uint32

const (
	CategoryStatusInvalid CategoryStatus = iota
	CategoryStatusNormal
	CategoryStatusDeleted
)

type CategoryUpdateOption func(c *Category)

func WithUpdateCategoryName(categoryName *string) CategoryUpdateOption {
	return func(c *Category) {
		if categoryName != nil {
			c.SetCategoryName(categoryName)
		}
	}
}

func WithUpdateCategoryStatus(categoryStatus *uint32) CategoryUpdateOption {
	return func(c *Category) {
		if categoryStatus != nil {
			c.SetCategoryStatus(categoryStatus)
		}
	}
}

type Category struct {
	UserID         *string
	CategoryID     *string
	CategoryName   *string
	CategoryType   *uint32
	CategoryStatus *uint32
	CreateTime     *uint64
	UpdateTime     *uint64

	Budget *Budget
}

type CategoryOption = func(c *Category)

func WithCategoryID(categoryID *string) CategoryOption {
	return func(c *Category) {
		if categoryID != nil {
			c.SetCategoryID(categoryID)
		}
	}
}

func WithCategoryType(categoryType *uint32) CategoryOption {
	return func(c *Category) {
		if categoryType != nil {
			c.SetCategoryType(categoryType)
		}
	}
}

func WithCategoryStatus(categoryStatus *uint32) CategoryOption {
	return func(c *Category) {
		if categoryStatus != nil {
			c.SetCategoryStatus(categoryStatus)
		}
	}
}

func WithCategoryCreateTime(createTime *uint64) CategoryOption {
	return func(c *Category) {
		if createTime != nil {
			c.SetCreateTime(createTime)
		}
	}
}

func WithCategoryUpdateTime(updateTime *uint64) CategoryOption {
	return func(c *Category) {
		if updateTime != nil {
			c.SetUpdateTime(updateTime)
		}
	}
}

func WithCategoryBudget(budget *Budget) CategoryOption {
	return func(c *Category) {
		if budget != nil {
			c.SetBudget(budget)
		}
	}
}

func (c *Category) Clone() (*Category, error) {
	return NewCategory(
		c.GetUserID(),
		c.GetCategoryName(),
		WithCategoryID(goutil.String(c.GetCategoryID())),
		WithCategoryType(c.CategoryType),
		WithCategoryStatus(c.CategoryStatus),
		WithCategoryCreateTime(c.CreateTime),
		WithCategoryUpdateTime(c.UpdateTime),
	)
}

func NewCategory(UserID, categoryName string, opts ...CategoryOption) (*Category, error) {
	now := uint64(time.Now().UnixMilli())
	c := &Category{
		CategoryID:     goutil.String(""),
		UserID:         goutil.String(UserID),
		CategoryName:   goutil.String(categoryName),
		CategoryType:   goutil.Uint32(uint32(TransactionTypeExpense)),
		CategoryStatus: goutil.Uint32(uint32(CategoryStatusNormal)),
		CreateTime:     goutil.Uint64(now),
		UpdateTime:     goutil.Uint64(now),
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Category) validate() error {
	if !c.CanAddBudget() && c.Budget != nil {
		return ErrBudgetNotAllowed
	}

	return nil
}

type CategoryUpdate struct {
	CategoryName   *string
	CategoryStatus *uint32
	UpdateTime     *uint64
}

func (c *Category) ToCategoryUpdate(old *Category) *CategoryUpdate {
	var (
		hasUpdate bool

		cu = &CategoryUpdate{
			UpdateTime: c.UpdateTime,
		}
	)

	if old.GetCategoryName() != c.GetCategoryName() {
		hasUpdate = true
		cu.CategoryName = c.CategoryName
	}

	if old.GetCategoryStatus() != c.GetCategoryStatus() {
		hasUpdate = true
		cu.CategoryStatus = c.CategoryStatus
	}

	if hasUpdate {
		return cu
	}

	return nil
}

func (c *Category) Update(cus ...CategoryUpdateOption) (*CategoryUpdate, error) {
	if len(cus) == 0 {
		return nil, nil
	}

	old, err := c.Clone()
	if err != nil {
		return nil, err
	}

	for _, cu := range cus {
		cu(c)
	}

	// check
	if err := c.validate(); err != nil {
		return nil, err
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	c.SetUpdateTime(now)

	return c.ToCategoryUpdate(old), nil
}

func (c *Category) GetUserID() string {
	if c != nil && c.UserID != nil {
		return *c.UserID
	}
	return ""
}

func (c *Category) SetUserID(UserID *string) {
	c.UserID = UserID
}

func (c *Category) GetCategoryID() string {
	if c != nil && c.CategoryID != nil {
		return *c.CategoryID
	}
	return ""
}

func (c *Category) SetCategoryID(categoryID *string) {
	c.CategoryID = categoryID
}

func (c *Category) GetCategoryName() string {
	if c != nil && c.CategoryName != nil {
		return *c.CategoryName
	}
	return ""
}

func (c *Category) SetCategoryName(categoryName *string) {
	c.CategoryName = categoryName
}

func (f *Category) GetCategoryType() uint32 {
	if f != nil && f.CategoryType != nil {
		return *f.CategoryType
	}
	return 0
}

func (c *Category) SetCategoryType(categoryType *uint32) {
	c.CategoryType = categoryType
}

func (c *Category) GetCategoryStatus() uint32 {
	if c != nil && c.CategoryStatus != nil {
		return *c.CategoryStatus
	}
	return 0
}

func (c *Category) SetCategoryStatus(categoryStatus *uint32) {
	c.CategoryStatus = categoryStatus
}

func (c *Category) GetCreateTime() uint64 {
	if c != nil && c.CreateTime != nil {
		return *c.CreateTime
	}
	return 0
}

func (c *Category) SetCreateTime(createTime *uint64) {
	c.CreateTime = createTime
}

func (c *Category) GetUpdateTime() uint64 {
	if c != nil && c.UpdateTime != nil {
		return *c.UpdateTime
	}
	return 0
}

func (c *Category) SetUpdateTime(updateTime *uint64) {
	c.UpdateTime = updateTime
}

func (c *Category) GetBudget() *Budget {
	if c != nil && c.Budget != nil {
		return c.Budget
	}
	return nil
}

func (c *Category) SetBudget(b *Budget) {
	c.Budget = b
}

func (c *Category) CanAddBudget() bool {
	return c.GetCategoryType() == uint32(TransactionTypeExpense)
}

func (c *Category) IsDeleted() bool {
	return c.GetCategoryStatus() == uint32(CategoryStatusDeleted)
}
