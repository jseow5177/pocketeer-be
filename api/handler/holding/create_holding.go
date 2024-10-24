package holding

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/handler/lot"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

func NewCreateHoldingValidator(optionalAccountID bool) validator.Validator {
	return validator.MustForm(map[string]validator.Validator{
		"account_id": &validator.String{
			Optional: optionalAccountID,
		},
		"symbol": &validator.String{
			Optional: false,
		},
		"currency": &validator.String{
			Optional:   false,
			Validators: []validator.StringFunc{entity.CheckCurrency},
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
			Validator: lot.NewCreateLotValidator(true),
		},
	})
}

func (h *holdingHandler) CreateHolding(ctx context.Context, req *presenter.CreateHoldingRequest, res *presenter.CreateHoldingResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.holdingUseCase.CreateHolding(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create holding, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
