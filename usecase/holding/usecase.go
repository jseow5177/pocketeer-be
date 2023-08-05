package holding

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/config"
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
	accountRepo  repo.AccountRepo
	holdingRepo  repo.HoldingRepo
	lotRepo      repo.LotRepo
	securityRepo repo.SecurityRepo
	quoteRepo    repo.QuoteRepo
}

func NewHoldingUseCase(
	accountRepo repo.AccountRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	securityRepo repo.SecurityRepo,
	quoteRepo repo.QuoteRepo,
) UseCase {
	return &holdingUseCase{
		accountRepo,
		holdingRepo,
		lotRepo,
		securityRepo,
		quoteRepo,
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

	if h.IsDefault() {
		if _, err = uc.securityRepo.Get(ctx, req.ToSecurityFilter()); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get security from repo, err: %v", err)
			return nil, err
		}
	}

	if _, err = uc.holdingRepo.Create(ctx, h); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new holding to repo, err: %v", err)
		return nil, err
	}

	return &CreateHoldingResponse{
		Holding: h,
	}, nil
}

func (uc *holdingUseCase) GetHolding(ctx context.Context, req *GetHoldingRequest) (*GetHoldingResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding from repo, err: %v", err)
		return nil, err
	}

	var q *entity.Quote
	if h.IsDefault() {
		q, err = uc.quoteRepo.Get(ctx, req.ToQuoteFilter(h.GetSymbol()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get quote from repo, err: %v", err)
			return nil, err
		}
	}

	if err = uc.calcHoldingValue(ctx, h, q); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute holding value, err: %v", err)
		return nil, err
	}

	return &GetHoldingResponse{
		Holding: h,
	}, nil
}

func (uc *holdingUseCase) GetHoldings(ctx context.Context, req *GetHoldingsRequest) (*GetHoldingsResponse, error) {
	hs, err := uc.holdingRepo.GetMany(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holdings from repo, err: %v", err)
		return nil, err
	}

	if err := goutil.ParallelizeWork(ctx, len(hs), 5, func(ctx context.Context, workNum int) error {
		h := hs[workNum]

		var q *entity.Quote
		if h.IsDefault() {
			q, err = uc.quoteRepo.Get(ctx, req.ToQuoteFilter(h.GetSymbol()))
			if err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get quote from repo, err: %v", err)
				return err
			}
		}
		return uc.calcHoldingValue(ctx, h, q)
	}); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to compute holdings value, err: %v", err)
		return nil, err
	}

	return &GetHoldingsResponse{
		Holdings: hs,
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

func (uc *holdingUseCase) calcHoldingValue(ctx context.Context, h *entity.Holding, q *entity.Quote) error {
	// value of custom holding is already stored in DB
	if h.IsCustom() {
		return nil
	}

	if q == nil {
		q = entity.NewQuote()
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

	// Calculate value as Total Shares * Current Price
	// We support only USD holdings now, so convert value from USD to SGD
	// TODO: Have better currency handling
	latestValue := util.RoundFloat(h.GetTotalShares()*q.GetLatestPrice()*config.USDToSGD, config.StandardDP)

	h.SetLatestValue(goutil.Float64(latestValue))
	h.SetQuote(q)

	return nil
}
