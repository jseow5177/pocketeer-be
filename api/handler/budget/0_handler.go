package budget

import (
	"github.com/jseow5177/pockteer-be/usecase/aggr"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

type budgetHandler struct {
	budgetUseCase budget.UseCase
	aggrUsecase   aggr.UseCase
}

func NewBudgetHandler(budgetUseCase budget.UseCase, aggrUsecase aggr.UseCase) *budgetHandler {
	return &budgetHandler{
		budgetUseCase,
		aggrUsecase,
	}
}
