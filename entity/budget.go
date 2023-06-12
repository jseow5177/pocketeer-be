package entity

import (
	"fmt"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type BudgetType uint32

const (
	BudgetTypeMonthly BudgetType = 1
	BudgetTypeYearly  BudgetType = 2
)

var BudgetTypes = map[uint32]string{
	uint32(BudgetTypeMonthly): "monthly",
	uint32(BudgetTypeYearly):  "yearly",
}

type BudgetStatus uint32

const (
	BudgetStatusNormal  BudgetStatus = 1
	BudgetStatusDeleted BudgetStatus = 2
)

var BudgetStatuses = map[uint32]string{
	uint32(BudgetStatusNormal):  "normal",
	uint32(BudgetStatusDeleted): "deleted",
}

type Budget struct {
	BudgetID     *string
	UserID       *string
	BudgetName   *string
	BudgetType   *uint32
	CategoryIDs  []string
	Status       *uint32
	BreakdownMap BreakdownMap
	UpdateTime   *uint64
}

type BreakdownMap map[DateInfo]*BudgetBreakdown

type BudgetBreakdown struct {
	Amount *float64
	Year   *int
	Month  *int
}

type DateInfo struct {
	Year  int
	Month int
}

func NewBudget(
	userID string,
	budgetType uint32,
) *Budget {
	return &Budget{
		UserID:       goutil.String(userID),
		BudgetType:   goutil.Uint32(budgetType),
		CategoryIDs:  make([]string, 0),
		BreakdownMap: make(BreakdownMap),
		Status:       goutil.Uint32(uint32(BudgetStatusNormal)),
	}
}

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

func (e *Budget) GetBudgetName() string {
	if e != nil && e.BudgetName != nil {
		return *e.BudgetName
	}
	return ""
}

func (e *Budget) GetBudgetType() uint32 {
	if e != nil && e.BudgetType != nil {
		return *e.BudgetType
	}
	return 0
}

func (e *Budget) GetCategoryIDs() []string {
	if e == nil {
		return []string{}
	}
	return e.CategoryIDs
}

func (e *Budget) GetBreakdownMap() BreakdownMap {
	if e == nil {
		return map[DateInfo]*BudgetBreakdown{}
	}
	return e.BreakdownMap
}

func (e *Budget) SetBudgetType(budgetType uint32) error {
	if e.BudgetType == nil {
		e.BudgetType = goutil.Uint32(budgetType)
	} else if e.GetBudgetType() != budgetType {
		return fmt.Errorf("cannot change budget type")
	}

	return nil
}

func (e *Budget) SetBudgetName(budgetName string) {
	e.BudgetName = goutil.String(budgetName)
}

func (e *Budget) SetCategoryIDs(categoryIDs []string) {
	e.CategoryIDs = categoryIDs
}

func (e *Budget) SetBudgetAmount(
	budgetAmount float64,
	rangeStartDate time.Time,
	rangeEndDate time.Time,
) {
	if e.GetBudgetType() == uint32(BudgetTypeMonthly) {
		e.setMonthlyBudget(budgetAmount, rangeStartDate, rangeEndDate)
	} else {
		e.setYearlyBudget(budgetAmount, rangeStartDate, rangeEndDate)
	}
}

func (e *Budget) IsBreakdownAvailable(
	date time.Time,
) bool {
	year, month := date.Year(), int(date.Month())

	dateInfo := DateInfo{}
	if e.GetBudgetType() == uint32(BudgetTypeMonthly) {
		dateInfo.Year = year
		dateInfo.Month = month
	} else {
		dateInfo.Year = year
	}

	_, exist := e.BreakdownMap[dateInfo]

	return exist
}
func (e *Budget) FilterBreakdownByDate(
	date time.Time,
) {
	year, month := date.Year(), int(date.Month())

	dateInfo := DateInfo{}
	if e.GetBudgetType() == uint32(BudgetTypeMonthly) {
		dateInfo.Year = year
		dateInfo.Month = month
	} else {
		dateInfo.Year = year
	}

	// Set breakdownMap to only 1 breakdown (the filtered breakdown)
	e.BreakdownMap = map[DateInfo]*BudgetBreakdown{
		dateInfo: e.BreakdownMap[dateInfo],
	}
}

func (e *Budget) setMonthlyBudget(
	budgetAmount float64,
	rangeStartDate time.Time,
	rangeEndDate time.Time,
) {
	date := rangeStartDate

	for date.Equal(rangeEndDate) || date.Before(rangeEndDate) {
		year, month := date.Year(), int(date.Month())
		e.setBudgetBreakdown(budgetAmount, year, month)

		date = date.AddDate(0, 1, 0)
	}
}

func (e *Budget) setYearlyBudget(
	budgetAmount float64,
	rangeStartDate time.Time,
	rangeEndDate time.Time,
) {
	date := rangeStartDate

	for date.Before(rangeEndDate) {
		year, _ := date.Year(), date.Month()
		e.setBudgetBreakdown(budgetAmount, year, 0)

		date = date.AddDate(1, 0, 0)
	}
}

func (e *Budget) setBudgetBreakdown(
	budgetAmount float64,
	year,
	month int,
) {
	dateInfo := DateInfo{
		Year:  year,
		Month: month,
	}

	_, ok := e.BreakdownMap[dateInfo]
	if !ok {
		e.BreakdownMap[dateInfo] = &BudgetBreakdown{
			Year:  goutil.Int(year),
			Month: goutil.Int(month),
		}
	}

	e.BreakdownMap[dateInfo].Amount = goutil.Float64(budgetAmount)
}
