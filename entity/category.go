package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type CategoryUpdate struct {
	CategoryName   *string
	CategoryStatus *uint32
	UpdateTime     *uint64
}

func (cu *CategoryUpdate) GetCategoryName() string {
	if cu != nil && cu.CategoryName != nil {
		return *cu.CategoryName
	}
	return ""
}

func (cu *CategoryUpdate) SetCategoryName(categoryName *string) {
	cu.CategoryName = categoryName
}

func (cu *CategoryUpdate) GetCategoryStatus() uint32 {
	if cu != nil && cu.CategoryStatus != nil {
		return *cu.CategoryStatus
	}
	return 0
}

func (cu *CategoryUpdate) SetCategoryStatus(categoryStatus *uint32) {
	cu.CategoryStatus = categoryStatus
}

func (cu *CategoryUpdate) GetUpdateTime() uint64 {
	if cu != nil && cu.UpdateTime != nil {
		return *cu.UpdateTime
	}
	return 0
}

func (cu *CategoryUpdate) SetUpdateTime(updateTime *uint64) {
	cu.UpdateTime = updateTime
}

type CategoryUpdateOption func(cu *CategoryUpdate)

func WithUpdateCategoryName(categoryName *string) CategoryUpdateOption {
	return func(cu *CategoryUpdate) {
		cu.SetCategoryName(categoryName)
	}
}

func WithUpdateCategoryStatus(categoryStatus *uint32) CategoryUpdateOption {
	return func(cu *CategoryUpdate) {
		cu.SetCategoryStatus(categoryStatus)
	}
}

func NewCategoryUpdate(opts ...CategoryUpdateOption) *CategoryUpdate {
	cu := new(CategoryUpdate)
	for _, opt := range opts {
		opt(cu)
	}
	return cu
}

type CategoryStatus uint32

const (
	CategoryStatusInvalid CategoryStatus = iota
	CategoryStatusNormal
	CategoryStatusDeleted
)

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
		c.SetCategoryID(categoryID)
	}
}

func WithCategoryType(categoryType *uint32) CategoryOption {
	return func(c *Category) {
		c.SetCategoryType(categoryType)
	}
}

func WithCategoryStatus(categoryStatus *uint32) CategoryOption {
	return func(c *Category) {
		c.SetCategoryStatus(categoryStatus)
	}
}

func WithCategoryCreateTime(createTime *uint64) CategoryOption {
	return func(c *Category) {
		c.SetCreateTime(createTime)
	}
}

func WithCategoryUpdateTime(updateTime *uint64) CategoryOption {
	return func(c *Category) {
		c.SetUpdateTime(updateTime)
	}
}

func NewCategory(userID, categoryName string, opts ...CategoryOption) (*Category, error) {
	now := uint64(time.Now().UnixMilli())
	c := &Category{
		UserID:         goutil.String(userID),
		CategoryName:   goutil.String(categoryName),
		CategoryType:   goutil.Uint32(uint32(TransactionTypeExpense)),
		CategoryStatus: goutil.Uint32(uint32(CategoryStatusNormal)),
		CreateTime:     goutil.Uint64(now),
		UpdateTime:     goutil.Uint64(now),
	}

	for _, opt := range opts {
		opt(c)
	}

	if err := c.checkOpts(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Category) checkOpts() error {
	return nil
}

func (c *Category) Update(cu *CategoryUpdate) (*CategoryUpdate, error) {
	var (
		hasUpdate      bool
		categoryUpdate = new(CategoryUpdate)
	)

	if cu.CategoryName != nil && cu.GetCategoryName() != c.GetCategoryName() {
		hasUpdate = true
		c.SetCategoryName(cu.CategoryName)

		defer func() {
			categoryUpdate.SetCategoryName(c.CategoryName)
		}()
	}

	if cu.CategoryStatus != nil && cu.GetCategoryStatus() != c.GetCategoryStatus() {
		hasUpdate = true
		c.SetCategoryStatus(cu.CategoryStatus)

		defer func() {
			categoryUpdate.SetCategoryStatus(c.CategoryStatus)
		}()
	}

	if !hasUpdate {
		return nil, nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	c.SetUpdateTime(now)

	// check
	if err := c.checkOpts(); err != nil {
		return nil, err
	}

	categoryUpdate.SetUpdateTime(now)

	return categoryUpdate, nil
}

func (c *Category) GetUserID() string {
	if c != nil && c.UserID != nil {
		return *c.UserID
	}
	return ""
}

func (c *Category) SetUserID(userID *string) {
	c.UserID = userID
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
