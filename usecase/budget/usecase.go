package budget

import (
	"context"
	"fmt"
	"math"
	"sync"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type budgetUseCase struct {
	budgetRepo      repo.BudgetRepo
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
}

func NewBudgetUseCase(
	budgetRepo repo.BudgetRepo,
	categoryRepo repo.CategoryRepo,
	transactionRepo repo.TransactionRepo,
) UseCase {
	return &budgetUseCase{
		budgetRepo,
		categoryRepo,
		transactionRepo,
	}
}

func (uc *budgetUseCase) CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*CreateBudgetResponse, error) {
	res, err := uc.getBudget(ctx, req.ToGetBudgetRequest())
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

func (uc *budgetUseCase) UpdateBudget(ctx context.Context, req *UpdateBudgetRequest) (*UpdateBudgetResponse, error) {
	res, err := uc.getBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}

	if res.Budget == nil {
		return nil, repo.ErrBudgetNotFound
	}

	return nil, nil
}

// Fetch budget from repo, won't compute used amount.
//
// Timezone is not needed.
func (uc *budgetUseCase) getBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error) {
	q, err := req.ToBudgetQuery()
	if err != nil {
		return nil, err
	}

	bs, err := uc.budgetRepo.GetMany(ctx, q)
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

// Fetch budget from repo and compute the used amount.
func (uc *budgetUseCase) GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error) {
	res, err := uc.getBudget(ctx, req)
	if err != nil {
		return nil, err
	}
	b := res.Budget

	if b == nil {
		return new(GetBudgetResponse), nil
	}

	var startTime, endTime uint64

	if b.IsMonth() {
		startTime, endTime, err = util.GetMonthRangeAsUnix(req.GetBudgetDate(), req.GetTimezone())
	} else if b.IsYear() {
		startTime, endTime, err = util.GetYearRangeAsUnix(req.GetBudgetDate(), req.GetTimezone())
	} else {
		return nil, fmt.Errorf("invalid budget type, budget ID: %v", b.GetBudgetID())
	}
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget date range, err: %v", err)
		return nil, err
	}

	aggrs, err := uc.transactionRepo.CalcTotalAmount(
		ctx,
		"category_id",
		req.ToTransactionFilter(req.GetUserID(), startTime, endTime),
	)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sum transactions by category ID, err: %v", err)
		return nil, err
	}

	var usedAmount float64
	if len(aggrs) > 0 {
		usedAmount = math.Abs(aggrs[0].GetTotalAmount())
	}
	b.SetUsedAmount(goutil.Float64(usedAmount))

	return &GetBudgetResponse{
		Budget: b,
	}, nil
}

func (uc *budgetUseCase) GetBudgets(ctx context.Context, req *GetBudgetsRequest) (*GetBudgetsResponse, error) {
	var (
		mu     sync.Mutex
		g      = new(errgroup.Group)
		wgChan = make(chan struct{}, 10)
	)
	reqs := req.ToGetBudgetRequests(req.GetUserID())
	budgets := make([]*entity.Budget, 0)

	for _, req := range reqs {
		wgChan <- struct{}{}
		req := req
		g.Go(func() error {
			defer func() {
				<-wgChan
			}()

			res, err := uc.GetBudget(ctx, req)
			if err != nil {
				return err
			}

			if res.Budget == nil {
				return nil
			}

			mu.Lock()
			budgets = append(budgets, res.Budget)
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &GetBudgetsResponse{
		Budgets: budgets,
	}, nil
}

func (uc *budgetUseCase) DeleteBudget(ctx context.Context, req *DeleteBudgetRequest) (*DeleteBudgetResponse, error) {
	res, err := uc.getBudget(ctx, req.ToGetBudgetRequest())
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
