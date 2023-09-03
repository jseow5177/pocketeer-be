package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type holdingUseCase struct {
	txMgr        repo.TxMgr
	accountRepo  repo.AccountRepo
	holdingRepo  repo.HoldingRepo
	lotRepo      repo.LotRepo
	securityRepo repo.SecurityRepo
	quoteRepo    repo.QuoteRepo
}

func NewHoldingUseCase(
	txMgr repo.TxMgr,
	accountRepo repo.AccountRepo,
	holdingRepo repo.HoldingRepo,
	lotRepo repo.LotRepo,
	securityRepo repo.SecurityRepo,
	quoteRepo repo.QuoteRepo,
) UseCase {
	return &holdingUseCase{
		txMgr,
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
		return nil, entity.ErrAccountCannotHaveHoldings
	}

	if h.IsDefault() {
		if _, err = uc.securityRepo.Get(ctx, req.ToSecurityFilter()); err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get security from repo, err: %v", err)
			return nil, err
		}

		q, err := uc.quoteRepo.Get(ctx, req.ToQuoteFilter())
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get quote from repo, err: %v", err)
			return nil, err
		}
		h.SetQuote(q)
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if _, err = uc.holdingRepo.Create(txCtx, h); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new holding to repo, err: %v", err)
			return err
		}

		if len(req.Lots) == 0 {
			return nil
		}

		ls := req.ToLotEntities()
		for _, l := range ls {
			l.SetHoldingID(h.HoldingID)
		}

		_, err := uc.lotRepo.CreateMany(txCtx, ls)
		if err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to save new lots to repo, err: %v", err)
			return err
		}

		h.SetLots(ls)
		h.ComputeSharesCostAndValue()

		return nil
	}); err != nil {
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

	if !h.IsDefault() {
		return &GetHoldingResponse{
			Holding: h,
		}, nil
	}

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

	h.ComputeSharesCostAndValue()

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

		if !h.IsDefault() {
			return nil
		}

		q, err := uc.quoteRepo.Get(ctx, req.ToQuoteFilter(h.GetSymbol()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get quote from repo, err: %v", err)
			return err
		}
		h.SetQuote(q)

		ls, err := uc.lotRepo.GetMany(ctx, req.ToLotFilter(h.GetHoldingID()))
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to get lots from repo, err: %v", err)
			return err
		}
		h.SetLots(ls)

		h.ComputeSharesCostAndValue()

		return nil
	}); err != nil {
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

	hu, err := h.Update(req.ToHoldingUpdate())
	if err != nil {
		return nil, err
	}

	if hu == nil {
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

func (uc *holdingUseCase) DeleteHolding(ctx context.Context, req *DeleteHoldingRequest) (*DeleteHoldingResponse, error) {
	hf := req.ToHoldingFilter()

	_, err := uc.holdingRepo.Get(ctx, hf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding from repo, err: %v", err)
		return nil, err
	}

	if err := uc.txMgr.WithTx(ctx, func(txCtx context.Context) error {
		if err := uc.holdingRepo.Delete(txCtx, hf); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to delete holding, err: %v", err)
			return err
		}

		if err := uc.lotRepo.Delete(txCtx, req.ToLotFilter()); err != nil {
			log.Ctx(txCtx).Error().Msgf("fail to delete lots, err: %v", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return new(DeleteHoldingResponse), nil
}
