package budget

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
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

func (uc *budgetUseCase) CreateBudget(ctx context.Context, req *CreateBudgetRequest) (*CreateBudgetResponse, error) {
	f := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, f)
	if err != nil && err != repo.ErrBudgetNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get budget from repo, err: %v", err)
		return nil, err
	}

	if b != nil {
		return nil, repo.ErrBudgetAlreadyExists
	}

	u := entity.GetUserFromCtx(ctx)

	b, err = req.ToBudgetEntity(u.Meta.GetCurrency())
	if err != nil {
		return nil, err
	}

	c, err := uc.categoryRepo.Get(ctx, req.ToCategoryFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get category from repo, err: %v", err)
		return nil, err
	}

	if err := b.CanBudgetUnderCategory(c); err != nil {
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
	f := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, f)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget from repo, err: %v", err)
		return nil, err
	}
	obt := b.GetBudgetType() // old budget type

	bu, err := req.ToBudgetUpdate()
	if err != nil {
		return nil, err
	}

	bu, err = b.Update(bu)
	if err != nil {
		return nil, err
	}

	if bu == nil {
		log.Ctx(ctx).Info().Msg("budget has no updates")
		return &UpdateBudgetResponse{
			Budget: b,
		}, nil
	}

	if err = uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// if there is a change in budget type, wipe all old records
		// delete time must be before update time
		if bu.BudgetType != nil {
			if err := uc.budgetRepo.Delete(txCtx, req.ToDeleteBudgetFilter(obt)); err != nil {
				log.Ctx(txCtx).Error().Msgf("fail to delete budget from repo, err: %v", err)
				return err
			}
		}

		b.SetBudgetID(nil)
		b.SetUpdateTime(goutil.Uint64(uint64(time.Now().UnixMilli())))

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
	f := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}

	return &GetBudgetResponse{
		Budget: b,
	}, nil
}

func (uc *budgetUseCase) DeleteBudget(ctx context.Context, req *DeleteBudgetRequest) (*DeleteBudgetResponse, error) {
	gf := req.ToGetBudgetFilter()

	b, err := uc.budgetRepo.Get(ctx, gf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get budget from repo, err: %v", err)
		return nil, err
	}

	df := req.ToDeleteBudgetFilter(b.GetBudgetType())

	if err := uc.budgetRepo.Delete(ctx, df); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete budget from repo, err: %v", err)
		return nil, err
	}

	return new(DeleteBudgetResponse), nil
}
