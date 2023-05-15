package budget

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/presenter"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/util"
)

type budgetUseCase struct {
	bcRepo  repo.BudgetConfigRepo
	catRepo repo.CategoryRepo
}

func NewBudgetUseCase(
	budgetConfigRepo repo.BudgetConfigRepo,
	catRepo repo.CategoryRepo,
) UseCase {
	return &budgetUseCase{
		bcRepo:  budgetConfigRepo,
		catRepo: catRepo,
	}
}

func (uc *budgetUseCase) GetMonthBudgets(
	ctx context.Context,
	userID string,
	req *presenter.GetMonthBudgetsRequest,
) (budgets []*entity.Budget, err error) {
	categories, err := uc.catRepo.GetMany(ctx, req.ToCategoryFilter(userID))
	if err != nil {
		return nil, err
	}

	catIDs := util.GetCatIDs(categories)
	budgetConfigs, err := uc.bcRepo.GetMany(
		ctx,
		&repo.BudgetConfigFilter{
			UserID: goutil.String(userID),
			CatIDs: catIDs,
		},
	)
	if err != nil {
		return nil, err
	}

	catIDToBudgetConfig := uc.getCatIDToConfigWithDefault(
		userID,
		catIDs,
		budgetConfigs,
	)

	catIDToCatMap := util.GetCatIDToCategoryMap(categories)
	budgets = make([]*entity.Budget, 0)
	for _, cfg := range catIDToBudgetConfig {
		budgets = append(
			budgets,
			cfg.GetBudget(
				req.GetYear(),
				req.GetMonth(),
				catIDToCatMap[cfg.GetCatID()],
			),
		)
	}

	return budgets, nil
}

func (uc *budgetUseCase) getCatIDToConfigWithDefault(
	userID string,
	catIDs []string,
	budgetConfigs []*entity.BudgetConfig,
) map[string]*entity.BudgetConfig {
	catIDToBudgetConfig := make(map[string]*entity.BudgetConfig)
	for _, cfg := range budgetConfigs {
		catIDToBudgetConfig[cfg.GetCatID()] = cfg
	}

	for _, catID := range catIDs {
		if _, ok := catIDToBudgetConfig[catID]; !ok {
			catIDToBudgetConfig[catID] = entity.NewDefaultBudgetConfig(
				userID,
				catID,
			)
		}
	}

	return catIDToBudgetConfig
}

func (uc *budgetUseCase) GetFullYearBudget(
	ctx context.Context,
	userID string,
	req *presenter.GetFullYearBudgetRequest,
) (defaultBudget *entity.Budget, monthlyBudgets []*entity.Budget, err error) {
	cat, err := uc.catRepo.Get(
		ctx,
		req.ToCategoryFilter(userID),
	)
	if err != nil {
		return nil, nil, err
	}

	budgetConfig, err := uc.bcRepo.Get(
		ctx, req.ToBudgetConfigFilter(userID),
	)
	if err != nil {
		return nil, nil, err
	}

	if budgetConfig == nil {
		budgetConfig = entity.NewDefaultBudgetConfig(
			userID,
			cat.GetCategoryID(),
		)
	}

	defaultBudget = budgetConfig.GetDefaultBudget(req.GetYear(), cat)
	monthlyBudgets = budgetConfig.GetMonthlyBudgets(req.GetYear(), cat)

	return defaultBudget, monthlyBudgets, nil
}

func (uc *budgetUseCase) SetBudget(
	ctx context.Context,
	userID string,
	req *presenter.SetBudgetRequest,
) error {
	budgetConfig, err := uc.bcRepo.Get(ctx, req.ToBudgetConfigFilter(userID))
	if err != nil {
		return err
	}

	if budgetConfig == nil {
		budgetConfig = entity.NewDefaultBudgetConfig(userID, req.GetCatID())
		budgetConfig.SetBudgetType(req.GetBudgetType())
	}

	if budgetConfig.HasDefaultBudgetChanged(req.GetCurrYear(), req.GetCurrMonth(), req.GetDefaultBudget()) {
		budgetConfig.SetNewDefaultBudget(req.GetCurrYear(), req.GetCurrMonth(), req.GetDefaultBudget())
	}

	for _, customBudget := range req.GetCustomBudgets() {
		budgetConfig.SetNewCustomBudget(
			customBudget.GetYear(),
			customBudget.GetMonth(),
			customBudget.GetAmount(),
		)
	}

	_, err = uc.bcRepo.Update(ctx, budgetConfig)
	if err != nil {
		return err
	}

	return err
}
