package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var GetTransactionGroupsValidator = validator.MustForm(map[string]validator.Validator{
	"app_meta": entity.AppMetaValidator(),
	"category_id": &validator.String{
		Optional: true,
	},
	"account_id": &validator.String{
		Optional:  true,
		UnsetZero: true,
	},
	"transaction_type": &validator.UInt32{
		Optional:   true,
		UnsetZero:  true,
		Validators: []validator.UInt32Func{entity.CheckTransactionType},
	},
	"transaction_time": validator.MustForm(map[string]validator.Validator{
		"gte": &validator.UInt64{
			Optional: false,
		},
		"lte": &validator.UInt64{
			Optional: false,
		},
	}),
	"paging": entity.PagingValidator(true),
})

func (h *transactionHandler) GetTransactionGroups(
	ctx context.Context,
	req *presenter.GetTransactionGroupsRequest,
	res *presenter.GetTransactionGroupsResponse,
) error {
	user := entity.GetUserFromCtx(ctx)

	useCaseRes, err := h.transactionUseCase.GetTransactionGroups(ctx, req.ToUseCaseReq(user.GetUserID()))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to get transaction, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
