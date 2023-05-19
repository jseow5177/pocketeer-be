package budget

import (
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

func toBudgets(categoryBudgets []*budget.CategoryBudget) []*presenter.Budget {
	budgets := make([]*presenter.Budget, len(categoryBudgets))

	for idx, b := range categoryBudgets {
		budget := &presenter.Budget{
			CategoryName: goutil.String(b.CategoryName),
		}

		if b.Budget != nil {
			budget.CategoryID = goutil.String(b.Budget.GetCategoryID())
			budget.BudgetType = goutil.Uint32(b.Budget.GetBudgetType())
			budget.Year = goutil.Uint32(b.Budget.GetYear())
			budget.Month = goutil.Uint32(b.Budget.GetMonth())
			budget.BudgetAmount = goutil.Int64(b.Budget.GetBudgetAmount())
		}

		budgets[idx] = budget
	}

	return budgets
}
