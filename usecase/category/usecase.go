package category

import (
	"context"
	"fmt"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

type categoryUseCase struct {
	txMgr            repo.TxMgr
	categoryRepo     repo.CategoryRepo
	transactionRepo  repo.TransactionRepo
	budgetRepo       repo.BudgetRepo
	exchangeRateRepo repo.ExchangeRateRepo
}

func NewCategoryUseCase(
	txMgr repo.TxMgr,
	categoryRepo repo.CategoryRepo,
	transactionRepo repo.TransactionRepo,
	budgetUseCase budget.UseCase,
	budgetRepo repo.BudgetRepo,
	exchangeRateRepo repo.ExchangeRateRepo,
) UseCase {
	return &categoryUseCase{
		txMgr,
		categoryRepo,
		transactionRepo,
		budgetRepo,
		exchangeRateRepo,
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

	_, err = uc.categoryRepo.Create(ctx, c)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new category to repo, err: %v", err)
		return nil, err
	}

	return &CreateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) UpdateCategory(ctx context.Context, req *UpdateCategoryRequest) (*UpdateCategoryResponse, error) {
	cf := req.ToCategoryFilter()
	c, err := uc.categoryRepo.Get(ctx, cf)
	if err != nil {
		return nil, err
	}

	cu, err := c.Update(
		entity.WithUpdateCategoryName(req.CategoryName),
	)
	if err != nil {
		return nil, err
	}

	if cu == nil {
		log.Ctx(ctx).Info().Msg("category has no updates")
		return &UpdateCategoryResponse{
			c,
		}, nil
	}

	if err := uc.categoryRepo.Update(ctx, cf, cu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save category updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateCategoryResponse{
		Category: c,
	}, nil
}

func (uc *categoryUseCase) DeleteCategory(ctx context.Context, req *DeleteCategoryRequest) (*DeleteCategoryResponse, error) {
	cf := req.ToCategoryFilter()

	c, err := uc.categoryRepo.Get(ctx, cf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		cu, err := c.Update(
			entity.WithUpdateCategoryStatus(goutil.Uint32(uint32(entity.CategoryStatusDeleted))),
		)
		if err != nil {
			return err
		}

		// mark category as deleted
		if err := uc.categoryRepo.Update(txCtx, cf, cu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save category updates to repo, err: %v", err)
			return err
		}

		// hard delete budgets
		bf := req.ToBudgetFilter()
		if err := uc.budgetRepo.DeleteMany(ctx, bf); err != nil {
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
		catID := catIDs[workNum]
		gReq := req.ToGetCategoryBudgetRequest(catID)

		b, err := uc.getBudgetWithUsage(ctx, gReq)
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
	bf := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, bf)
	if err != nil && err != repo.ErrBudgetNotFound {
		return nil, err
	}

	// no budget
	if b == nil {
		return nil, nil
	}

	// get budget date range
	var start, end uint64
	if b.IsMonth() {
		start, end, err = util.GetMonthRangeAsUnix(req.GetBudgetDate(), req.AppMeta.GetTimezone())
	} else if b.IsYear() {
		start, end, err = util.GetYearRangeAsUnix(req.GetBudgetDate(), req.AppMeta.GetTimezone())
	} else {
		return nil, fmt.Errorf("invalid budget type, budget ID: %v", b.GetBudgetID())
	}
	if err != nil {
		return nil, err
	}

	tq := req.ToTransactionQuery(req.GetUserID(), start, end)
	ts, err := uc.transactionRepo.GetMany(ctx, tq)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	var usedAmount float64
	for _, t := range ts {
		amount := t.GetAmount()

		if t.GetCurrency() != b.GetCurrency() {
			erf := req.ToExchangeRateFilter(
				b.GetCurrency(),
				t.GetCurrency(),
				t.GetTransactionTime(),
			)
			er, err := uc.exchangeRateRepo.Get(ctx, erf)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get exchange rate from repo, err: %v", err)
				return nil, err
			}

			amount *= er.GetRate()
		}

		usedAmount += amount
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

	var (
		categoryIDs   = make([]string, 0)
		categoryMap   = make(map[string]*entity.Category)
		sumByCategory = make(map[*entity.Category]float64)
	)
	for _, c := range cs {
		categoryIDs = append(categoryIDs, c.GetCategoryID())
		categoryMap[c.GetCategoryID()] = c

		if c.IsDeleted() {
			sumByCategory[nil] = 0
		} else {
			sumByCategory[c] = 0
		}
	}

	tq := req.ToTransactionQuery(categoryIDs)
	ts, err := uc.transactionRepo.GetMany(ctx, tq)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions from repo, err: %v", err)
		return nil, err
	}

	u := entity.GetUserFromCtx(ctx)

	for _, t := range ts {
		amount := t.GetAmount()

		if t.GetCurrency() != u.Meta.GetCurrency() {
			erf := req.ToExchangeRateFilter(
				u.Meta.GetCurrency(),
				t.GetCurrency(),
				t.GetTransactionTime(),
			)
			er, err := uc.exchangeRateRepo.Get(ctx, erf)
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get exchange rate from repo, err: %v", err)
				return nil, err
			}

			amount *= er.GetRate()
		}

		if c := categoryMap[t.GetCategoryID()]; c.IsDeleted() {
			sumByCategory[nil] += amount
		} else {
			sumByCategory[c] += amount
		}
	}

	sums := make([]*common.TransactionSummary, 0)
	for c, sum := range sumByCategory {
		// remove uncategorized if it has no sum
		if c == nil && sum == 0 {
			continue
		}

		sums = append(sums, &common.TransactionSummary{
			Category: c,
			Sum:      goutil.Float64(util.RoundFloatToStandardDP(sum)),
			Currency: u.Meta.Currency,
		})
	}

	return &SumCategoryTransactionsResponse{
		Sums: sums,
	}, nil
}
