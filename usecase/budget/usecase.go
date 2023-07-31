package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/rs/zerolog/log"
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
	res, err := uc.GetBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}

	if res.Budget != nil {
		return nil, repo.ErrBudgetAlreadyExists
	}

	b, err := req.ToBudgetEntity()
	if err != nil {
		return nil, err
	}

	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if _, err := b.CanBudgetUnderCategory(c); err != nil {
		return nil, err
	}

	if _, err := uc.budgetRepo.Create(ctx, b); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new budget to repo, err: %v", err)
		return nil, err
	}

	return &CreateBudgetResponse{
		Budget: b,
	}, nil
}

func (uc *budgetUseCase) GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error) {
	q, paging, err := req.ToBudgetQuery()
	if err != nil {
		return nil, err
	}

	bs, err := uc.budgetRepo.GetMany(ctx, paging, q)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget from repo, err: %v", err)
		return nil, err
	}

	var b *entity.Budget
	if len(bs) > 0 && !bs[0].IsDeleted() {
		b = bs[0]
	}

	return &GetBudgetResponse{
		Budget: b,
	}, nil
}

func (uc *budgetUseCase) GetBudgets(ctx context.Context, req *GetBudgetsRequest) (*GetBudgetsResponse, error) {
	return &GetBudgetsResponse{}, nil
}

func (uc *budgetUseCase) DeleteBudget(ctx context.Context, req *DeleteBudgetRequest) (*DeleteBudgetResponse, error) {
	res, err := uc.GetBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}

	if res.Budget == nil {
		return new(DeleteBudgetResponse), nil
	}

	// create a dummy, deleted budget
	b, err := req.ToBudgetEntity(res.Budget.GetBudgetType())
	if err != nil {
		return nil, err
	}

	if _, err := uc.budgetRepo.Create(ctx, b); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new budget to repo, err: %v", err)
		return nil, err
	}

	return new(DeleteBudgetResponse), nil
}
