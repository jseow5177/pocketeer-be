package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type BudgetType uint32

const (
	BudgetTypeMonth BudgetType = iota
	BudgetTypeYear
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonth): "month",
	uint32(BudgetTypeYear):  "year",
}

type BudgetStatus uint32

const (
	BudgetStatusNormal BudgetStatus = iota
	BudgetStatusDeleted
)

type Budget struct {
	BudgetID     *string
	CategoryID   *string
	UserID       *string
	BudgetType   *uint32
	Amount       *float64
	BudgetStatus *uint32
	CreateTime   *uint64
	UpdateTime   *uint64
}

type BudgetOption = func(b *Budget)

func WithBudgetID(budgetID *string) BudgetOption {
	return func(b *Budget) {
		b.BudgetID = budgetID
	}
}

func WithBudgetType(budgetType *uint32) BudgetOption {
	return func(b *Budget) {
		b.BudgetType = budgetType
	}
}

func WithBudgetStatus(budgetStatus *uint32) BudgetOption {
	return func(b *Budget) {
		b.BudgetStatus = budgetStatus
	}
}

func WithBudgetAmount(amount *float64) BudgetOption {
	return func(b *Budget) {
		b.Amount = amount
	}
}

func WithBudgetCreateTime(createTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.CreateTime = createTime
	}
}

func WithBudgetUpdateTime(updateTime *uint64) BudgetOption {
	return func(b *Budget) {
		b.UpdateTime = updateTime
	}
}

func NewBudget(userID, categoryID string, opts ...BudgetOption) *Budget {
	now := uint64(time.Now().UnixMilli())
	b := &Budget{
		UserID:       goutil.String(userID),
		CategoryID:   goutil.String(categoryID),
		BudgetType:   goutil.Uint32(uint32(BudgetTypeMonth)),
		BudgetStatus: goutil.Uint32(uint32(BudgetStatusNormal)),
		Amount:       goutil.Float64(0),
		CreateTime:   goutil.Uint64(now),
		UpdateTime:   goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(b)
	}

	b.checkOpts()

	return b
}

func (b *Budget) checkOpts() {}

func (b *Budget) GetBudgetID() string {
	if b != nil && b.BudgetID != nil {
		return *b.BudgetID
	}
	return ""
}

func (b *Budget) SetBudgetID(budgetID *string) {
	b.BudgetID = budgetID
}

func (b *Budget) GetUserID() string {
	if b != nil && b.UserID != nil {
		return *b.UserID
	}
	return ""
}

func (b *Budget) GetCategoryID() string {
	if b != nil && b.CategoryID != nil {
		return *b.CategoryID
	}
	return ""
}

func (b *Budget) GetBudgetType() uint32 {
	if b != nil && b.BudgetType != nil {
		return *b.BudgetType
	}
	return 0
}

func (b *Budget) GetBudgetStatus() uint32 {
	if b != nil && b.BudgetStatus != nil {
		return *b.BudgetStatus
	}
	return 0
}

func (b *Budget) GetAmount() float64 {
	if b != nil && b.Amount != nil {
		return *b.Amount
	}
	return 0
}

func (b *Budget) GetCreateTime() uint64 {
	if b != nil && b.CreateTime != nil {
		return *b.CreateTime
	}
	return 0
}

func (b *Budget) GetUpdateTime() uint64 {
	if b != nil && b.UpdateTime != nil {
		return *b.UpdateTime
	}
	return 0
}
