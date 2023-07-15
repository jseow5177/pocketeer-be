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

	// TODO: Verify symbol must be unique

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

func (uc *holdingUseCase) calcHoldingValue(ctx context.Context, h *entity.Holding) error {
	// Compute total shares and cost
	aggr, err := uc.lotRepo.CalcTotalSharesAndCost(ctx, &repo.LotFilter{
		UserID:    h.UserID,
		HoldingID: h.HoldingID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to calc lot aggr from repo, err: %v", err)
		return err
	}

	// Compute avg cost
	var avgCost float64
	if aggr.GetTotalCost() != 0 {
		avgCost = util.RoundFloat(aggr.GetTotalCost()/aggr.GetTotalShares(), config.PreciseDP)
	}

	h.SetTotalShares(aggr.TotalShares)
	h.SetAvgCost(goutil.Float64(avgCost))

	// If custom holding, the latest value is Total Shares * Avg Cost
	// Else, get quote and calculate Total Shares * Current Price
	if h.IsCustom() {
		latestValue := util.RoundFloat(h.GetTotalShares()*h.GetAvgCost(), config.StandardDP)
		h.SetLatestValue(goutil.Float64(latestValue))
	} else {
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
	}

	return nil
}
