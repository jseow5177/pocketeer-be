package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type lotUseCase struct {
	lotRepo     repo.LotRepo
	holdingRepo repo.HoldingRepo
}

func NewLotUseCase(
	lotRepo repo.LotRepo,
	holdingRepo repo.HoldingRepo,
) UseCase {
	return &lotUseCase{
		lotRepo,
		holdingRepo,
	}
}

func (uc *lotUseCase) CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get holding from repo, err: %v", err)
		return nil, err
	}

	if !h.IsDefault() {
		return nil, entity.ErrHoldingCannotHaveLots
	}

	// use holding's currency
	l := req.ToLotEntity(h.GetCurrency())

	_, err = uc.lotRepo.Create(ctx, l)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save new lot to repo, err: %v", err)
		return nil, err
	}

	return &CreateLotResponse{
		Lot: l,
	}, nil
}

func (uc *lotUseCase) UpdateLot(ctx context.Context, req *UpdateLotRequest) (*UpdateLotResponse, error) {
	l, err := uc.lotRepo.Get(ctx, req.ToLotFilter())
	if err != nil {
		return nil, err
	}

	lu := l.Update(
		entity.WithUpdateLotCostPerShare(req.CostPerShare),
		entity.WithUpdateLotShares(req.Shares),
		entity.WithUpdateLotTradeDate(req.TradeDate),
	)

	if lu == nil {
		log.Ctx(ctx).Info().Msg("lot has no updates")
		return &UpdateLotResponse{
			l,
		}, nil
	}

	if err = uc.lotRepo.Update(ctx, req.ToLotFilter(), lu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to save lot updates to repo, err: %v", err)
		return nil, err
	}

	return &UpdateLotResponse{
		Lot: l,
	}, nil
}

func (uc *lotUseCase) GetLot(ctx context.Context, req *GetLotRequest) (*GetLotResponse, error) {
	l, err := uc.lotRepo.Get(ctx, req.ToLotFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lot from repo, err: %v", err)
		return nil, err
	}

	return &GetLotResponse{
		Lot: l,
	}, nil
}

func (uc *lotUseCase) GetLots(ctx context.Context, req *GetLotsRequest) (*GetLotsResponse, error) {
	ls, err := uc.lotRepo.GetMany(ctx, req.ToLotFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lots from repo, err: %v", err)
		return nil, err
	}

	return &GetLotsResponse{
		Lots: ls,
	}, nil
}

func (uc *lotUseCase) DeleteLot(ctx context.Context, req *DeleteLotRequest) (*DeleteLotResponse, error) {
	lf := req.ToLotFilter()

	l, err := uc.lotRepo.Get(ctx, lf)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lot from repo, err: %v", err)
		return nil, err
	}

	lu := l.Update(
		entity.WithUpdateLotStatus(goutil.Uint32(uint32(entity.LotStatusDeleted))),
	)

	// mark lot as deleted
	if err := uc.lotRepo.Update(ctx, lf, lu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to mark lot as deleted, err: %v", err)
		return nil, err
	}

	return new(DeleteLotResponse), nil
}
