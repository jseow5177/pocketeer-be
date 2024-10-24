package lot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

func NewCreateLotValidator(optionalHoldingID bool) validator.Validator {
	return validator.MustForm(map[string]validator.Validator{
		"holding_id": &validator.String{
			Optional: optionalHoldingID,
		},
		"shares": &validator.String{
			Optional:   false,
			Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
		},
		"cost_per_share": &validator.String{
			Optional:   false,
			Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
		},
		"trade_date": &validator.UInt64{
			Optional: true,
		},
	})
}

func (h *lotHandler) CreateLot(ctx context.Context, req *presenter.CreateLotRequest, res *presenter.CreateLotResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.lotUseCase.CreateLot(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create lot, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
