package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/middleware"
	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var UpdateTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"transaction_id": &validator.String{
		Optional: false,
	},
	"category_id": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"amount": &validator.String{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.StringFunc{entity.CheckAmount},
	},
	"transaction_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": &validator.UInt64{
		Optional:  true,
		UnsetZero: true,
	},
	"note": &validator.String{
		Optional: true,
		MaxLen:   uint32(config.MaxTransactionNoteLength),
	},
})

func (h *transactionHandler) UpdateTransaction(ctx context.Context, req *presenter.UpdateTransactionRequest, res *presenter.UpdateTransactionResponse) error {
	userID := middleware.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.UpdateTransaction(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to update transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}