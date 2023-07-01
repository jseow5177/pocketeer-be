package budget

import (
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

type budgetHandler struct {
	budgetUseCase budget.UseCase
}

func NewBudgetHandler(budgetUseCase budget.UseCase) *budgetHandler {
	return &budgetHandler{
		budgetUseCase,
	}
}
