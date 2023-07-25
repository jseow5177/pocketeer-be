package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
)

type budgetUseCase struct {
	budgetRepo   repo.BudgetRepo
	categoryRepo repo.CategoryRepo
}

func NewBudgetUseCase(
	budgetRepo repo.BudgetRepo,
	categoryRepo repo.CategoryRepo,
) UseCase {
	return &budgetUseCase{
		budgetRepo,
		categoryRepo,
	}
}

func (uc *budgetUseCase) CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*CreateBudgetResponse, error) {
	return &CreateBudgetResponse{}, nil
}

func (uc *budgetUseCase) GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error) {
	return &GetBudgetResponse{}, nil
}

func (uc *budgetUseCase) GetBudgets(ctx context.Context, req *GetBudgetsRequest) (*GetBudgetsResponse, error) {
	return &GetBudgetsResponse{}, nil
}
