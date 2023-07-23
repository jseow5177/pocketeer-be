package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var DeleteTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
})

func (h *transactionHandler) DeleteTransaction(ctx context.Context, req *presenter.DeleteTransactionRequest, res *presenter.DeleteTransactionResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.DeleteTransaction(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
