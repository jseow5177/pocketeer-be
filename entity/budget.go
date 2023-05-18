package entity

import "github.com/jseow5177/pockteer-be/pkg/goutil"

type BudgetType uint32

const (
	BudgetTypeMonthly BudgetType = 1
	BudgetTypeYearly  BudgetType = 2
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonthly): "monthly",
	uint32(BudgetTypeYearly):  "yearly",
}

type Budget struct {
	BudgetID     *string
	UserID       *string
	CategoryID   *string
	IsDefault    *bool
	BudgetType   *uint32
	Year         *uint32
	Month        *uint32
	BudgetAmount *int64
}

func NewBudget(
	userID,
	categoryID string,
	year,
	month uint32,
	isDefault bool,
	budgetAmount int64,
) *Budget {
	return &Budget{
		UserID:       goutil.String(userID),
		CategoryID:   goutil.String(categoryID),
		IsDefault:    goutil.Bool(isDefault),
		BudgetType:   goutil.Uint32(uint32(BudgetTypeMonthly)),
		Year:         goutil.Uint32(year),
		Month:        goutil.Uint32(month),
		BudgetAmount: goutil.Int64(budgetAmount),
	}
}

// ****************** Getters | Setters
func (e *Budget) GetBudgetID() string {
	if e != nil && e.BudgetID != nil {
		return *e.BudgetID
	}
	return ""
}

func (e *Budget) GetUserID() string {
	if e != nil && e.UserID != nil {
		return *e.UserID
	}
	return ""
}

func (e *Budget) GetCategoryID() string {
	if e != nil && e.CategoryID != nil {
		return *e.CategoryID
	}
	return ""
}

func (e *Budget) GetIsDefault() bool {
	if e != nil && e.IsDefault != nil {
		return *e.IsDefault
	}
	return false
}

func (e *Budget) GetBudgetType() uint32 {
	if e != nil && e.BudgetType != nil {
		return *e.BudgetType
	}
	return 0
}

func (e *Budget) GetYear() uint32 {
	if e != nil && e.Year != nil {
		return *e.Year
	}
	return 0
}

func (e *Budget) GetMonth() uint32 {
	if e != nil && e.Month != nil {
		return *e.Month
	}
	return 0
}

func (e *Budget) GetBudgetAmount() int64 {
	if e != nil && e.BudgetAmount != nil {
		return *e.BudgetAmount
	}
	return 0
}

func (e *Budget) SetBudgetType(budgetType uint32) {
	e.BudgetType = goutil.Uint32(budgetType)
}

func (e *Budget) SetBudgetAmount(budgetAmount int64) {
	e.BudgetAmount = goutil.Int64(budgetAmount)
}

func (e *Budget) SetIsDefault(isDefault bool) {
	e.IsDefault = goutil.Bool(isDefault)
}

// ******************

func (e *Budget) IsDefaultBudget() bool {
	return e.GetIsDefault()
}
