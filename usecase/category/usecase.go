package category

import (
	"context"
	"fmt"
	"math"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

type categoryUseCase struct {
	txMgr           repo.TxMgr
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
	budgetRepo      repo.BudgetRepo
}

func NewCategoryUseCase(
	txMgr repo.TxMgr,
	categoryRepo repo.CategoryRepo,
	transactionRepo repo.TransactionRepo,
	budgetUseCase budget.UseCase,
	budgetRepo repo.BudgetRepo,
) UseCase {
	return &categoryUseCase{
		txMgr,
		categoryRepo,
		transactionRepo,
		budgetRepo,
	}
}

func (uc *categoryUseCase) GetCategory(ctx context.Context, req *GetCategoryRequest) (*GetCategoryResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	return &GetCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) GetCategoryBudget(ctx context.Context, req *GetCategoryBudgetRequest) (*GetCategoryBudgetResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if !c.CanAddBudget() {
		return &GetCategoryBudgetResponse{
			Category: c,
		}, nil
	}

	b, err := uc.getBudgetWithUsage(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget usage, err: %v", err)
		return nil, err
	}

	c.SetBudget(b)

	return &GetCategoryBudgetResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CreateCategoryResponse, error) {
	c, err := req.ToCategoryEntity()
	if err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		_, err = uc.categoryRepo.Create(txCtx, c)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new category to repo, err: %v", err)
			return err
		}

		if c.Budget == nil {
			return nil
		}

		b := c.Budget
		b.SetCategoryID(c.CategoryID)

		_, err := uc.budgetRepo.Create(txCtx, b)
		if err != nil {
			return err
		}

		c.SetBudget(b)

		return nil
	}); err != nil {
		return nil, err
	}

	return &CreateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		return nil, err
	}

	cu, err := c.Update(req.ToCategoryUpdate())
	if err != nil {
		return nil, err
	}

	if cu == nil {
		log.Ctx(ctx).Info().Msg("category has no updates")
		return &UpdateCategoryResponse{
			c,
		}, nil
	}

	if err = uc.categoryRepo.Update(ctx, req.ToCategoryFilter(), cu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteCategoryResponse, error) {
	f := req.ToCategoryFilter()

	_, err := uc.categoryRepo.Get(ctx, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// mark category as deleted
		if err := uc.categoryRepo.Delete(txCtx, f); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark category as deleted, err: %v", err)
			return err
		}

		// hard delete budgets
		if err := uc.budgetRepo.DeleteMany(ctx, req.ToBudgetFilter()); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to delete category budgets, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return new(DeleteCategoryResponse), nil
}

func (uc *categoryUseCase) GetCategories(ctx context.Context, req *GetCategoriesRequest) (*GetCategoriesResponse, error) {
	cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	return &GetCategoriesResponse{
		Categories: cs,
	}, nil
}

func (uc *categoryUseCase) GetCategoriesBudget(ctx context.Context, req *GetCategoriesBudgetRequest) (*GetCategoriesBudgetResponse, error) {
	cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	var (
		catIDs = make([]string, 0)
		cbs    = make(map[string]*entity.Category)
	)
	for _, c := range cs {
		if c.CanAddBudget() {
			if _, ok := cbs[c.GetCategoryID()]; ok {
				// deduplicate
				continue
			}
			catIDs = append(catIDs, c.GetCategoryID())
			cbs[c.GetCategoryID()] = c
		}
	}

	if len(catIDs) == 0 {
		return &GetCategoriesBudgetResponse{
			Categories: cs,
		}, nil
	}

	if err := goutil.ParallelizeWork(ctx, len(catIDs), 10, func(ctx context.Context, workNum int) error {
		b, err := uc.getBudgetWithUsage(ctx, req.ToGetCategoryBudgetRequest(catIDs[workNum]))
		if err != nil {
			return err
		}
		if b != nil {
			cbs[b.GetCategoryID()].SetBudget(b)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &GetCategoriesBudgetResponse{
		Categories: cs,
	}, nil
}

func (uc *categoryUseCase) getBudgetWithUsage(ctx context.Context, req *GetCategoryBudgetRequest) (*entity.Budget, error) {
	f := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, f)
	if err != nil && err != repo.ErrBudgetNotFound {
		return nil, err
	}

	// no budget
	if b == nil {
		return nil, nil
	}

	var start, end uint64
	if b.IsMonth() {
		start, end, err = util.GetMonthRangeAsUnix(req.GetBudgetDate(), req.GetTimezone())
	} else if b.IsYear() {
		start, end, err = util.GetYearRangeAsUnix(req.GetBudgetDate(), req.GetTimezone())
	} else {
		return nil, fmt.Errorf("invalid budget type, budget ID: %v", b.GetBudgetID())
	}
	if err != nil {
		return nil, err
	}

	tf := req.ToTransactionFilter(req.GetUserID(), start, end)
	aggrs, err := uc.transactionRepo.CalcTotalAmount(ctx, "category_id", tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transactions by category ID, err: %v", err)
		return nil, err
	}

	var usedAmount float64
	if len(aggrs) > 0 {
		usedAmount = math.Abs(aggrs[0].GetTotalAmount())
	}

	b.SetUsedAmount(goutil.Float64(usedAmount))

	return b, nil
}

func (uc *categoryUseCase) SumCategoryTransactions(ctx context.Context, req *SumCategoryTransactionsRequest) (*SumCategoryTransactionsResponse, error) {
	cs, err := uc.categoryRepo.GetMany(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get categories from repo, err: %v", err)
		return nil, err
	}

	tf := req.ToTransactionFilter()

	aggrs, err := uc.transactionRepo.Sum(ctx, "category_id", tf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transactions, err: %v", err)
		return nil, err
	}

	sums := make(map[*entity.Category]float64)
	for _, c := range cs {
		if c.IsDeleted() {
			sums[nil] += aggrs[c.GetCategoryID()]
		} else {
			sums[c] += aggrs[c.GetCategoryID()]
		}
	}

	// remove uncategorized if it has no sum
	if sum, ok := sums[nil]; ok && sum == 0 {
		delete(sums, nil)
	}

	res := make([]*common.TransactionSummary, 0)
	for c, sum := range sums {
		res = append(res, &common.TransactionSummary{
			Category: c,
			Sum:      goutil.Float64(util.RoundFloatToStandardDP(sum)),
		})
	}

	return &SumCategoryTransactionsResponse{
		Sums: res,
	}, nil
}
