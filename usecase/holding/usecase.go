package holding

import (
	"context"
	"errors"
	"fmt"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/lot"
	"github.com/rs/zerolog/log"
)

var (
	ErrAccountNotInvestment = errors.New("account is not investment type")
)

type holdingUseCase struct {
	txMgr        repo.TxMgr
	accountRepo  repo.AccountRepo
	holdingRepo  repo.HoldingRepo
	lotRepo      repo.LotRepo
	lotUseCase   lot.UseCase
	securityRepo repo.SecurityRepo
	quoteRepo    repo.QuoteRepo
}

func NewHoldingUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	lotUseCase lot.UseCase,
	securityRepo repo.SecurityRepo,
	quoteRepo repo.QuoteRepo,
) UseCase {
	return &holdingUseCase{
		txMgr,
		accountRepo,
		holdingRepo,
		lotRepo,
		lotUseCase,
		securityRepo,
		quoteRepo,
	}
}

func (uc *holdingUseCase) CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error) {
	h, err := req.ToHoldingEntity(req.GetAccountID())
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

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if _, err = uc.holdingRepo.Create(txCtx, h); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new holding to repo, err: %v", err)
			return err
		}

		if len(req.Lots) == 0 {
			return nil
		}

		lotsRes, err := uc.lotUseCase.CreateLots(txCtx, req.ToCreateLotsRequest(h.GetHoldingID()))
		if err != nil {
			return err
		}
		h.SetLots(lotsRes.Lots)

		return nil
	}); err != nil {
		return nil, err
	}

	return &CreateHoldingResponse{
		Holding: h,
	}, nil
}

func (uc *holdingUseCase) CreateHoldings(ctx context.Context, req *CreateHoldingsRequest) (*CreateHoldingsResponse, error) {
	hs, err := req.ToHoldingEntities()
	if err != nil {
		return nil, err
	}

	if len(hs) == 0 {
		return new(CreateHoldingsResponse), nil
	}

	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if !ac.IsInvestment() {
		return nil, ErrAccountNotInvestment
	}

	for _, h := range hs {
		if h.IsDefault() {
			if _, err = uc.securityRepo.Get(ctx, req.ToSecurityFilter(h.GetSymbol())); err != nil {
				log.Ctx(ctx).Error().Msgf("fail to get security from repo, err: %v", err)
				return nil, fmt.Errorf("symbol %v, err: %v", h.GetSymbol(), err)
			}
		}
	}

	errChan := make(chan int, len(req.Holdings))
	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		holdingIDs, err := uc.holdingRepo.CreateMany(txCtx, hs)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new holdings to repo, err: %v", err)
			return err
		}

		return goutil.ParallelizeWork(txCtx, len(req.Holdings), 5, func(ctx context.Context, workNum int) error {
			r := req.Holdings[workNum]

			if len(r.Lots) == 0 {
				return nil
			}

			lotsRes, err := uc.lotUseCase.CreateLots(ctx, r.ToCreateLotsRequest(holdingIDs[workNum]))
			if err != nil {
				errChan <- workNum
				return err
			}
			hs[workNum].SetLots(lotsRes.Lots)

			return nil
		})
	}); err != nil {
		return nil, fmt.Errorf("holding idx %v, err: %v", <-errChan, err)
	}

	return &CreateHoldingsResponse{
		Holdings: hs,
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

	if hu.Symbol != nil && h.IsDefault() {
		if _, err = uc.securityRepo.Get(ctx, &repo.SecurityFilter{
			Symbol: hu.Symbol,
		}); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get security from repo, err: %v", err)
			return nil, err
		}
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
		LotStatus: goutil.Uint32(uint32(entity.LotStatusNormal)),
	})
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to calc lot aggr from repo, err: %v", err)
		return err
	}

	h.SetTotalShares(aggr.TotalShares)
	h.SetTotalCost(goutil.Float64(aggr.GetTotalCost() * config.USDToSGD))

	// Compute avg cost per share
	var avgCostPerShare float64
	if aggr.GetTotalCost() != 0 {
		avgCostPerShare = aggr.GetTotalCost() / aggr.GetTotalShares()
	}
	h.SetAvgCostPerShare(goutil.Float64(avgCostPerShare * config.USDToSGD))

	// Calculate value as Total Shares * Current Price
	// We support only USD holdings now, so convert value from USD to SGD
	// TODO: Have better currency handling
	latestValue := h.GetTotalShares() * q.GetLatestPrice() * config.USDToSGD

	h.SetLatestValue(goutil.Float64(latestValue))
	h.SetQuote(q)

	return nil
}
