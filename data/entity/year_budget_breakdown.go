package entity

type YearBudgetBreakdown struct {
	DefaultBudget  *Budget
	MonthlyBudgets []*Budget
}

func NewYearBudgetBreakdown(
	fullYearBudgets []*Budget,
) *YearBudgetBreakdown {
	return &YearBudgetBreakdown{}
}

func DefaultYearBudgetBreakdown(
	userID,
	categoryID string,
	year uint32,
) *YearBudgetBreakdown {
	return &YearBudgetBreakdown{}
}

func (e *YearBudgetBreakdown) SetBudgetType(budgetType uint32) {}

func (e *YearBudgetBreakdown) SetDefaultBudget(budgetAmount int64) {}

func (e *YearBudgetBreakdown) SetMonthlyBudget(monthlyBudget int64) {}

func (e *YearBudgetBreakdown) ToBudgets() []*Budget {
	return []*Budget{}
}
