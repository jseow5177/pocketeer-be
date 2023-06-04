package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var CreateTransactionValidator = validator.MustForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: false,
	},
	"amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckAmount},
	},
	"transaction_type": &validator.UInt32{
		Optional:   false,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": &validator.UInt64{
		Optional: false,
	},
	"note": &validator.String{
		Optional: true,
		MaxLen:   uint32(config.MaxTransactionNoteLength),
	},
})

func (h *transactionHandler) CreateTransaction(ctx context.Context, req *presenter.CreateTransactionRequest, res *presenter.CreateTransactionResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.CreateTransaction(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
