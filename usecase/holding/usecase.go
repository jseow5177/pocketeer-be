package holding

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrAccountNotInvestment = errors.New("account is not investment type")
)

type holdingUseCase struct {
	accountRepo repo.AccountRepo
	holdingRepo repo.HoldingRepo
	lotRepo     repo.LotRepo
}

func NewHoldingUseCase(accountRepo repo.AccountRepo, holdingRepo repo.HoldingRepo, lotRepo repo.LotRepo) UseCase {
	return &holdingUseCase{
		accountRepo,
		holdingRepo,
		lotRepo,
	}
}

func (uc *holdingUseCase) CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter(req.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if !ac.IsInvestment() {
		return nil, ErrAccountNotInvestment
	}

	h := req.ToHoldingEntity()

	// TODO: Check if symbol exists
	log.Ctx(ctx).Info().Msgf("checking if symbol exists: %v", h.GetSymbol())

	if _, err = uc.holdingRepo.Create(ctx, h); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new holding to repo, err: %v", err)
		return nil, err
	}

	return &CreateHoldingResponse{
		h,
	}, nil
}

func (uc *holdingUseCase) GetHolding(ctx context.Context, req *GetHoldingRequest) (*GetHoldingResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding from repo, err: %v", err)
		return nil, err
	}

	// Compute total shares and cost
	aggr, err := uc.lotRepo.CalcTotalSharesAndCost(ctx, req.ToLotFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to calc lot aggr from repo, err: %v", err)
		return nil, err
	}

	var avgCost float64
	if aggr.GetTotalCost() != 0 {
		avgCost = util.RoundFloat(aggr.GetTotalCost()/aggr.GetTotalShares(), config.PreciseDP)
	}

	h.SetTotalShares(aggr.TotalShares)
	h.SetAvgCost(goutil.Float64(avgCost))

	// TODO: Get latest value

	return &GetHoldingResponse{
		h,
	}, nil
}
