package snapshot

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetAccountSnapshotsValidator = validator.MustForm(map[string]validator.Validator{
	"account_id": &validator.String{
		Optional: true,
	},
	"unit": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckSnapshotUnit},
	},
	"interval": &validator.UInt32{
		Optional: false,
		Min:      goutil.Uint32(1),
	},
})

func (h *snapshotHandler) GetAccountSnapshots(ctx context.Context, req *presenter.GetAccountSnapshotsRequest, res *presenter.GetAccountSnapshotsResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.snapshotUseCase.GetAccountSnapshots(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get account snapshots, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
