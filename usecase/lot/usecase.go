package lot

import (
	"context"
	"errors"
	"fmt"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidHolding = errors.New("invalid holding")
)

type lotUseCase struct {
	lotRepo     repo.LotRepo
	holdingRepo repo.HoldingRepo
}

func NewLotUseCase(lotRepo repo.LotRepo, holdingRepo repo.HoldingRepo) UseCase {
	return &lotUseCase{
		lotRepo,
		holdingRepo,
	}
}

func (uc *lotUseCase) DeleteLot(ctx context.Context, req *DeleteLotRequest) (*DeleteLotResponse, error) {
	l, err := uc.lotRepo.Get(ctx, req.ToLotFilter())
	if err != nil && err != repo.ErrLotNotFound {
		log.Ctx(ctx).Error().Msgf("fail to get lot from repo, err: %v", err)
		return nil, err
	}

	if err == repo.ErrLotNotFound {
		return new(DeleteLotResponse), nil
	}

	lu, _ := l.Update(req.ToLotUpdate())

	fmt.Println(lu.GetLotStatus())

	// mark lot as deleted
	if err := uc.lotRepo.Update(ctx, req.ToLotFilter(), lu); err != nil {
		log.Ctx(ctx).Error().Msgf("fail to mark lot as deleted, err: %v", err)
		return nil, err
	}

	return new(DeleteLotResponse), nil
}

func (uc *lotUseCase) CreateLot(ctx context.Context, req *CreateLotRequest) (*CreateLotResponse, error) {
	h, err := uc.holdingRepo.Get(ctx, req.ToHoldingFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get lot from repo, err: %v", err)
		return nil, err
	}

	if h.IsCustom() {
		return nil, ErrInvalidHolding
	}

	l := req.ToLotEntity()
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

	lu, hasUpdate := l.Update(req.ToLotUpdate())
	if !hasUpdate {
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
