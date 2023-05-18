package budget

import uc "github.com/jseow5177/pockteer-be/usecase/budget"

type budgetHandler struct {
	budgetUseCase uc.UseCase
}

func NewBudgetHandler(budgetUseCase uc.UseCase) *budgetHandler {
	return &budgetHandler{budgetUseCase: budgetUseCase}
}