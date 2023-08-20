package holding

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/api/handler/lot"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrEmptyAccountID = errors.New("empty account_id")
)

var CreateHoldingValidator = validator.MustForm(map[string]validator.Validator{
	"account_id": &validator.String{
		Optional: true, // allow reuse in CreateAccount
	},
	"symbol": &validator.String{
		Optional: false,
	},
	"holding_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckHoldingType},
	},
	"total_cost": &validator.String{
		Optional:   true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"latest_value": &validator.String{
		Optional:   true,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"lots": &validator.Slice{
		Optional:  true,
		MaxLen:    5,
		Validator: lot.CreateLotValidator,
	},
})

func (h *holdingHandler) CreateHolding(ctx context.Context, req *presenter.CreateHoldingRequest, res *presenter.CreateHoldingResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	if req.GetAccountID() == "" {
		return ErrEmptyAccountID
	}

	useCaseRes, err := h.holdingUseCase.CreateHolding(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create holding, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
