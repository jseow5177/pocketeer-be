package holding

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
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
	securityAPI api.SecurityAPI
}

func NewHoldingUseCase(
	accountRepo repo.AccountRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	securityAPI api.SecurityAPI,
) UseCase {
	return &holdingUseCase{
		accountRepo,
		holdingRepo,
		lotRepo,
		securityAPI,
	}
}

func (uc *holdingUseCase) CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error) {
	h, err := req.ToHoldingEntity()
	if err != nil {
		return nil, err
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if !ac.IsInvestment() {
		return nil, ErrAccountNotInvestment
	}

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

	if err = uc.calcHoldingValue(ctx, h); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute holding value, err: %v", err)
		return nil, err
	}

	return &GetHoldingResponse{
		h,
	}, nil
}

func (uc *holdingUseCase) GetHoldings(ctx context.Context, req *GetHoldingsRequest) (*GetHoldingsResponse, error) {
	hs, err := uc.holdingRepo.GetMany(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holdings from repo, err: %v", err)
		return nil, err
	}

	if err := goutil.ParallelizeWork(ctx, len(hs), 5, func(ctx context.Context, workNum int) error {
		return uc.calcHoldingValue(ctx, hs[workNum])
	}); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute holdings value, err: %v", err)
		return nil, err
	}

	return &GetHoldingsResponse{
		hs,
	}, nil
}

func (uc *holdingUseCase) UpdateHolding(ctx context.Context, req *UpdateHoldingRequest) (*UpdateHoldingResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		return nil, err
	}

	hu, hasUpdate, err := h.Update(req.ToHoldingUpdate())
	if err != nil {
		return nil, err
	}

	if !hasUpdate {
		log.Ctx(ctx).Info().Msg("holding has no updates")
		return &UpdateHoldingResponse{
			Holding: h,
		}, nil
	}

	if err = uc.holdingRepo.Update(ctx, req.ToHoldingFilter(), hu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save holding updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateHoldingResponse{
		Holding: h,
	}, nil
}

func (uc *holdingUseCase) calcHoldingValue(ctx context.Context, h *entity.Holding) error {
	// value of custom holding is already stored in DB
	if h.IsCustom() {
		return nil
	}

	// Compute total shares and cost
	aggr, err := uc.lotRepo.CalcTotalSharesAndCost(ctx, &repo.LotFilter{
		UserID:    h.UserID,
		HoldingID: h.HoldingID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to calc lot aggr from repo, err: %v", err)
		return err
	}

	// Compute avg cost per share
	var avgCostPerShare float64
	if aggr.GetTotalCost() != 0 {
		avgCostPerShare = util.RoundFloat(aggr.GetTotalCost()/aggr.GetTotalShares(), config.StandardDP)
	}

	h.SetAvgCostPerShare(goutil.Float64(avgCostPerShare))
	h.SetTotalShares(aggr.TotalShares)
	h.SetTotalCost(aggr.TotalCost)

	// Get quote and calculate Total Shares * Current Price
	quote, err := uc.securityAPI.GetLatestQuote(ctx, &api.SecurityFilter{
		Symbol: h.Symbol,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get latest quote, symbol: %v, err: %v", h.GetSymbol(), err)
		return err
	}
	latestValue := util.RoundFloat(h.GetTotalShares()*quote.GetLatestPrice(), config.StandardDP)
	h.SetLatestValue(goutil.Float64(latestValue))
	h.SetQuote(quote)

	return nil
}
