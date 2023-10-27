package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type holdingUseCase struct {
	txMgr            repo.TxMgr
	accountRepo      repo.AccountRepo
	holdingRepo      repo.HoldingRepo
	lotRepo          repo.LotRepo
	securityRepo     repo.SecurityRepo
	quoteRepo        repo.QuoteRepo
	exchangeRateRepo repo.ExchangeRateRepo
}

func NewHoldingUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	securityRepo repo.SecurityRepo,
	quoteRepo repo.QuoteRepo,
	exchangeRateRepo repo.ExchangeRateRepo,
) UseCase {
	return &holdingUseCase{
		txMgr,
		accountRepo,
		holdingRepo,
		lotRepo,
		securityRepo,
		quoteRepo,
		exchangeRateRepo,
	}
}

func (uc *holdingUseCase) CreateHolding(ctx context.Context, req *CreateHoldingRequest) (*CreateHoldingResponse, error) {
	ac, err := uc.accountRepo.Get(ctx, req.ToAccountFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account from repo, err: %v", err)
		return nil, err
	}

	if !ac.IsInvestment() {
		return nil, entity.ErrAccountCannotHaveHoldings
	}

	// default to account currency
	currency := ac.GetCurrency()

	if req.GetHoldingType() == uint32(entity.HoldingTypeDefault) {
		s, err := uc.securityRepo.Get(ctx, req.ToSecurityFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get security from repo, err: %v", err)
			return nil, err
		}

		// use symbol's currency
		currency = s.GetCurrency()
	}

	h, err := req.ToHoldingEntity(currency)
	if err != nil {
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		// create holding
		if _, err = uc.holdingRepo.Create(txCtx, h); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to save new holding to repo, err: %v", err)
			return err
		}

		// create lots
		if len(h.Lots) > 0 {
			if _, err = uc.lotRepo.CreateMany(txCtx, h.Lots); err != nil {
				log.Ctx(ctx).Error().Msgf("fail to save new lots to repo, err: %v", err)
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	h.ComputeCostGainAndValue()

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

	if h.IsDefault() {
		q, err := uc.quoteRepo.Get(ctx, req.ToQuoteFilter(h.GetSymbol()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get quote from repo, err: %v", err)
			return nil, err
		}
		h.SetQuote(q)

		ls, err := uc.lotRepo.GetMany(ctx, req.ToLotFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get lots from repo, err: %v", err)
			return nil, err
		}
		h.SetLots(ls)
	}

	h.ComputeCostGainAndValue()

	return &GetHoldingResponse{
		Holding: h,
	}, nil
}

func (uc *holdingUseCase) UpdateHolding(ctx context.Context, req *UpdateHoldingRequest) (*UpdateHoldingResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		return nil, err
	}

	hu, err := h.Update(
		entity.WithUpdateHoldingSymbol(req.Symbol),
		entity.WithUpdateHoldingLatestValue(req.LatestValue),
		entity.WithUpdateHoldingTotalCost(req.TotalCost),
	)
	if err != nil {
		return nil, err
	}

	if len(req.Lots) > 1 {
		// One lot only
	}

	if hu == nil {
		log.Ctx(ctx).Info().Msg("holding has no updates")
		return &UpdateHoldingResponse{
			Holding: h,
		}, nil
	}

	if hu.Symbol != nil && h.IsDefault() {
		if _, err = uc.securityRepo.Get(ctx, repo.NewSecurityFilter(
			repo.WithSecuritySymbol(hu.Symbol),
		)); err != nil {
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

func (uc *holdingUseCase) DeleteHolding(ctx context.Context, req *DeleteHoldingRequest) (*DeleteHoldingResponse, error) {
	hf := req.ToHoldingFilter()

	h, err := uc.holdingRepo.Get(ctx, hf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		hu, err := h.Update(
			entity.WithUpdateHoldingStatus(goutil.Uint32(uint32(entity.HoldingStatusDeleted))),
		)
		if err != nil {
			return err
		}

		// mark holding as deleted
		if err := uc.holdingRepo.Update(txCtx, hf, hu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark holding as deleted, err: %v", err)
			return err
		}

		lu := &entity.LotUpdate{
			LotStatus: goutil.Uint32(uint32(entity.LotStatusDeleted)),
		}

		// mark lots as deleted
		if err := uc.lotRepo.Update(txCtx, req.ToLotFilter(), lu); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to mark lots as deleted, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return new(DeleteHoldingResponse), nil
}
