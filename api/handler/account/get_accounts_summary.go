package account

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetAccountsSummaryValidator = validator.MustForm(map[string]validator.Validator{
	"unit": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckSnapshotUnit},
	},
	"interval": &validator.UInt32{
		Optional: false,
		Min:      goutil.Uint32(1),
	},
})

func (h *accountHandler) GetAccountsSummary(ctx context.Context, req *presenter.GetAccountsSummaryRequest, res *presenter.GetAccountsSummaryResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.accountUseCase.GetAccountsSummary(ctx, req.ToUseCaseReq(user))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get accounts summary from repo, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
