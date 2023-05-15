package entity

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type BudgtType uint32

const (
	BudgetTypeMonthly BudgtType = 1
	BudgetTypeAnnual  BudgtType = 2
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonthly): "budget_type_monthly",
	uint32(BudgetTypeAnnual):  "budget_type_annual",
}

type BudgetConfig struct {
	UserID         *string
	CatID          *string
	BudgetID       *string
	BudgetType     *uint32
	DefaultBudgets []*BudgetBreakdown
	CustomBudgets  []*BudgetBreakdown
}

type BudgetBreakdown struct {
	Year   uint32
	Month  uint32
	Amount int64
}

func NewDefaultBudgetConfig(
	userID,
	catID string,
) *BudgetConfig {
	return &BudgetConfig{
		UserID:     goutil.String(userID),
		CatID:      goutil.String(catID),
		BudgetType: goutil.Uint32(DefaultBudgetType),
	}
}

// ****************** Getters ******************
func (e *BudgetConfig) GetUserID() string {
	if e != nil && e.UserID != nil {
		return *e.UserID
	}
	return ""
}

func (e *BudgetConfig) GetCatID() string {
	if e != nil && e.CatID != nil {
		return *e.CatID
	}
	return ""
}

func (e *BudgetConfig) GetBudgetID() string {
	if e != nil && e.BudgetID != nil {
		return *e.BudgetID
	}
	return ""
}

func (e *BudgetConfig) GetBudgetType() uint32 {
	if e != nil && e.BudgetType != nil {
		return *e.BudgetType
	}
	return 0
}

func (e *BudgetConfig) GetDefaultBudgets() []*BudgetBreakdown {
	if e != nil && e.DefaultBudgets != nil {
		return e.DefaultBudgets
	}
	return []*BudgetBreakdown{}
}

func (e *BudgetConfig) GetCustomBudgets() []*BudgetBreakdown {
	if e != nil && e.CustomBudgets != nil {
		return e.CustomBudgets
	}
	return []*BudgetBreakdown{}
}

// ****************** Utils ******************
func (e *BudgetConfig) GetBudget(
	year,
	month uint32,
	cat *Category,
) *Budget {
	return &Budget{
		CatID:           goutil.String(e.GetCatID()),
		CatName:         goutil.String(cat.GetCategoryName()),
		BudgetType:      goutil.Uint32(e.GetBudgetType()),
		TransactionType: goutil.Uint32(cat.GetCategoryType()),
		Year:            goutil.Uint32(year),
		Month:           goutil.Uint32(month),
		BudgetAmount:    goutil.Int64(e.calculateBudgetAmount(year, month)),
	}
}

func (e *BudgetConfig) GetDefaultBudget(
	year uint32,
	cat *Category,
) *Budget {
	return &Budget{
		CatID:           goutil.String(e.GetCatID()),
		CatName:         goutil.String(cat.GetCategoryName()),
		BudgetType:      goutil.Uint32(e.GetBudgetType()),
		TransactionType: goutil.Uint32(cat.GetCategoryType()),
		Year:            goutil.Uint32(year),
		BudgetAmount:    goutil.Int64(e.calculateDefaultAmount(year, uint32(util.Constant_DEC))), // we get the latest month
	}
}

func (e *BudgetConfig) HasDefaultBudgetChanged(
	year uint32,
	month uint32,
	newDefaultBudget int64,
) bool {
	defaultBudgetAmount := e.calculateDefaultAmount(year, month)
	return defaultBudgetAmount == newDefaultBudget
}

func (e *BudgetConfig) SetNewDefaultBudget(
	year,
	month uint32,
	newDefaultBudget int64,
) {
	for _, defaultBudget := range e.DefaultBudgets {
		if defaultBudget.Year == year && defaultBudget.Month == month {
			defaultBudget.Amount = newDefaultBudget
			return
		}
	}

	e.DefaultBudgets = append(e.DefaultBudgets, &BudgetBreakdown{
		Year:   year,
		Month:  month,
		Amount: newDefaultBudget,
	})
}

func (e *BudgetConfig) SetNewCustomBudget(
	year,
	month uint32,
	newCustomBudget int64,
) {
	for _, customBudget := range e.CustomBudgets {
		if customBudget.Year == year && customBudget.Month == month {
			customBudget.Amount = newCustomBudget
			return
		}
	}

	e.CustomBudgets = append(e.CustomBudgets, &BudgetBreakdown{
		Year:   year,
		Month:  month,
		Amount: newCustomBudget,
	})
}

func (e *BudgetConfig) SetBudgetType(budgetType uint32) {
	e.BudgetType = goutil.Uint32(budgetType)
}

func (e *BudgetConfig) GetMonthlyBudgets(
	year uint32,
	cat *Category,
) []*Budget {
	if e.GetBudgetType() == uint32(BudgetTypeAnnual) {
		return []*Budget{}
	}

	monthlyBudgets := make([]*Budget, len(util.MonthTypes))
	for month := range util.MonthTypes {
		monthlyBudgets[month] = e.GetBudget(year, month, cat)
	}

	return monthlyBudgets
}

func (e *BudgetConfig) calculateBudgetAmount(
	year, month uint32,
) int64 {
	if e.GetBudgetType() == uint32(BudgetTypeAnnual) {
		return e.calculateDefaultAmount(year, month)
	}

	isAmountCustomised, amount := e.isAmountCustomised(year, month)
	if isAmountCustomised {
		return amount
	}

	return e.calculateDefaultAmount(year, month)
}

func (e *BudgetConfig) isAmountCustomised(
	year, month uint32,
) (bool, int64) {
	dateToBreakdownMap := getDateToBreakdownMap(e.CustomBudgets)
	key := getYearMonthKey(year, month)

	if breakdown, ok := dateToBreakdownMap[key]; ok {
		return true, breakdown.Amount
	}

	return false, 0
}

func (e *BudgetConfig) calculateDefaultAmount(
	year, month uint32,
) int64 {
	if len(e.DefaultBudgets) == 0 {
		return DefaultBudgetAmount
	}

	sortBudgetBreakdowns(e.DefaultBudgets)

	chosen := e.DefaultBudgets[0]
	for _, bd := range e.DefaultBudgets {
		isInputLargerEqual := isDate1LargerEqual(
			year, month,
			bd.Year, bd.Month,
		)

		if isInputLargerEqual {
			chosen = bd
		}
	}

	return chosen.Amount
}
