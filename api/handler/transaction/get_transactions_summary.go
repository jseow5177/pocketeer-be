package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetTransactionsSummaryValidator = validator.MustForm(map[string]validator.Validator{
	"app_meta": entity.AppMetaValidator(),
	"unit": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckSnapshotUnit},
	},
	"interval": &validator.UInt32{
		Optional: false,
		Min:      goutil.Uint32(1),
	},
})

func (h *transactionHandler) GetTransactionsSummary(ctx context.Context, req *presenter.GetTransactionsSummaryRequest, res *presenter.GetTransactionsSummaryResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.GetTransactionsSummary(ctx, req.ToUseCaseReq(user))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transactions summary, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
