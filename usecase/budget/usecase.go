package budget

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type budgetUseCase struct {
	txMgr           repo.TxMgr
	budgetRepo      repo.BudgetRepo
	categoryRepo    repo.CategoryRepo
	transactionRepo repo.TransactionRepo
}

func NewBudgetUseCase(
	txMgr repo.TxMgr,
	budgetRepo repo.BudgetRepo,
	categoryRepo repo.CategoryRepo,
	transactionRepo repo.TransactionRepo,
) UseCase {
	return &budgetUseCase{
		txMgr,
		budgetRepo,
		categoryRepo,
		transactionRepo,
	}
}

func (uc *budgetUseCase) CreateBudgets(ctx context.Context, req *CreateBudgetsRequest) (*CreateBudgetsResponse, error) {
	var (
		bs = make([]*entity.Budget, 0)
		cs = make(map[string]*entity.Category)
	)

	for _, r := range req.Budgets {
		res, err := uc.GetBudget(ctx, r.ToGetBudgetRequest())
		if err != nil {
			return nil, err
		}

		if res.Budget != nil {
			return nil, repo.ErrBudgetAlreadyExists
		}

		b, err := r.ToBudgetEntity()
		if err != nil {
			return nil, err
		}

		// for simplicity, do not allow multiple budgets for the same category in a batch
		_, ok := cs[b.GetCategoryID()]
		if ok {
			return nil, repo.ErrBudgetAlreadyExists
		}

		c, err := uc.categoryRepo.Get(ctx, r.ToCategoryFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
			return nil, err
		}
		cs[b.GetCategoryID()] = c

		if _, err := b.CanBudgetUnderCategory(c); err != nil {
			return nil, err
		}

		bs = append(bs, b)
	}

	if _, err := uc.budgetRepo.CreateMany(ctx, bs); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new budgets to repo, err: %v", err)
		return nil, err
	}

	return &CreateBudgetsResponse{
		Budgets: bs,
	}, nil
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

func (uc *budgetUseCase) UpdateBudget(ctx context.Context, req *UpdateBudgetRequest) (*UpdateBudgetResponse, error) {
	now := uint64(time.Now().UnixMilli())

	res, err := uc.GetBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}
	b := res.Budget

	if b == nil {
		return nil, repo.ErrBudgetNotFound
	}

	bu, err := req.ToBudgetUpdate()
	if err != nil {
		return nil, err
	}

	bu, hasUpdate, err := b.Update(bu)
	if err != nil {
		return nil, err
	}

	if !hasUpdate {
		log.Ctx(ctx).Info().Msg("budget has no updates")
		return &UpdateBudgetResponse{
			Budget: b,
		}, nil
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// if there is a change in budget type, wipe all old records
		// delete time must be before update time
		if bu.BudgetType != nil {
			if _, err := uc.DeleteBudget(txCtx, req.ToDeleteBudgetRequest(now)); err != nil {
				return err
			}
		}

		b.BudgetID = nil
		if _, err := uc.budgetRepo.Create(txCtx, b); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to create new updated budget to repo, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &UpdateBudgetResponse{
		Budget: b,
	}, nil
}

func (uc *budgetUseCase) GetBudget(ctx context.Context, req *GetBudgetRequest) (*GetBudgetResponse, error) {
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

	if b == nil {
		return new(GetBudgetResponse), nil
	}

	if req.GetTimezone() != "" {
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
	}

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
	res, err := uc.GetBudget(ctx, req.ToGetBudgetRequest())
	if err != nil {
		return nil, err
	}

	if res.Budget == nil {
		return new(DeleteBudgetResponse), nil
	}

	// create a dummy, deleted budget
	b, err := req.ToBudgetEntity(res.Budget.GetBudgetType(), req.GetDeleteTime())
	if err != nil {
		return nil, err
	}

	if _, err := uc.budgetRepo.Create(ctx, b); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save dummy budget to repo, err: %v", err)
		return nil, err
	}

	return new(DeleteBudgetResponse), nil
}
