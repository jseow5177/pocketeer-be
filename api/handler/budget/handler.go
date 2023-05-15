package budget

import "github.com/jseow5177/pockteer-be/dep/repo"

type BudgetHandler struct {
	categoryRepo     repo.CategoryRepo
	budgetConfigRepo repo.BudgetConfigRepo
}

func NewCategoryHandler(
	categoryRepo repo.CategoryRepo,
	budgetConfigRepo repo.BudgetConfigRepo,
) *BudgetHandler {
	return &BudgetHandler{
		categoryRepo:     categoryRepo,
		budgetConfigRepo: budgetConfigRepo,
	}
}
