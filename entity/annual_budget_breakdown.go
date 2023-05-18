package entity

import (
	"fmt"

	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/util"
)

type AnnualBudgetBreakdown struct {
	DefaultBudget  *Budget
	MonthlyBudgets []*Budget
}

func NewAnnualBudgetBreakdown(
	fullYearBudgets []*Budget,
) (*AnnualBudgetBreakdown, error) {
	var defaultBudget *Budget
	monthlyBudgets := make([]*Budget, 0)

	for _, budget := range fullYearBudgets {
		if budget.IsDefaultBudget() {
			defaultBudget = budget
			continue
		}

		monthlyBudgets = append(monthlyBudgets, budget)
	}

	if defaultBudget == nil {
		return nil, errutil.NotFoundError(fmt.Errorf("default budget not found"))
	}

	if len(monthlyBudgets) != len(util.MonthTypes) {
		return nil, errutil.InternalServerError(
			fmt.Errorf(
				"monthly budget len=%d, should be %d",
				len(monthlyBudgets), len(util.MonthTypes),
			),
		)
	}

	return &AnnualBudgetBreakdown{
		DefaultBudget:  defaultBudget,
		MonthlyBudgets: monthlyBudgets,
	}, nil
}

func DefaultAnnualBudgetBreakdown(
	userID,
	categoryID string,
	year uint32,
) *AnnualBudgetBreakdown {
	defaultBudget := NewBudget(userID, categoryID, year, defaultBudgetMonth, true, defaultBudgetAmount)

	monthlyBudgets := make([]*Budget, 0)
	for month := range util.MonthTypes {
		monthlyBudget := NewBudget(userID, categoryID, year, month, false, defaultBudgetAmount)
		monthlyBudgets = append(monthlyBudgets, monthlyBudget)
	}

	return &AnnualBudgetBreakdown{
		DefaultBudget:  defaultBudget,
		MonthlyBudgets: monthlyBudgets,
	}
}

// ****************** Getters
func (e *AnnualBudgetBreakdown) GetDefaultBudget() *Budget {
	return e.DefaultBudget
}

func (e *AnnualBudgetBreakdown) GetMonthlyBudgets() []*Budget {
	return e.MonthlyBudgets
}

func (e *AnnualBudgetBreakdown) GetBudgetType() uint32 {
	return e.GetDefaultBudget().GetBudgetType()
}

// ******************
func (e *AnnualBudgetBreakdown) IsYearlyBudget() bool {
	return e.GetBudgetType() == uint32(BudgetTypeYearly)
}

func (e *AnnualBudgetBreakdown) IsMonthlyBudget() bool {
	return e.GetBudgetType() == uint32(BudgetTypeMonthly)
}

func (e *AnnualBudgetBreakdown) SetBudgetType(budgetType uint32) {
	e.SetBudgetType(budgetType)

	if budgetType == uint32(BudgetTypeYearly) {
		e.setMonthlyBudgetsToDefault()
	}
}

func (e *AnnualBudgetBreakdown) SetDefaultBudget(budgetAmount int64) {
	e.DefaultBudget.SetBudgetAmount(budgetAmount)

	if e.IsYearlyBudget() {
		e.setMonthlyBudgetsToDefault()
	} else {
		e.setFutureBudgetsToDefault()
	}
}

func (e *AnnualBudgetBreakdown) SetMonthlyBudget(
	month uint32,
	budgetAmount int64,
) error {
	if e.IsYearlyBudget() {
		return errutil.BadRequestError(fmt.Errorf("budget type is yearly, cannot update monthly budget"))
	}

	for _, budget := range e.MonthlyBudgets {
		if budget.GetMonth() == month {
			budget.SetBudgetAmount(budgetAmount)
		}
	}

	return nil
}

func (e *AnnualBudgetBreakdown) ToBudgets() []*Budget {
	budgets := make([]*Budget, 0)
	budgets = append(budgets, e.DefaultBudget)
	budgets = append(budgets, e.MonthlyBudgets...)

	return budgets
}

func (e *AnnualBudgetBreakdown) setMonthlyBudgetsToDefault() {
	for _, budget := range e.GetMonthlyBudgets() {
		budget.SetBudgetAmount(e.GetDefaultBudget().GetBudgetAmount())
	}
}

func (e *AnnualBudgetBreakdown) setFutureBudgetsToDefault() {
	currMonth := util.GetCurrMonth()

	for _, budget := range e.GetMonthlyBudgets() {
		if budget.GetMonth() > currMonth {
			budget.SetBudgetAmount(e.GetDefaultBudget().GetBudgetAmount())
		}
	}
}
