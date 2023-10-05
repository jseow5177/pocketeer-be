package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var DeleteTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
})

func (h *transactionHandler) DeleteTransaction(ctx context.Context, req *presenter.DeleteTransactionRequest, res *presenter.DeleteTransactionResponse) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.DeleteTransaction(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to delete transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
