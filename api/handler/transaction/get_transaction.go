package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var GetTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
})

func (h *transactionHandler) GetTransaction(ctx context.Context, req *presenter.GetTransactionRequest, res *presenter.GetTransactionResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.aggrUseCase.GetTransaction(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
