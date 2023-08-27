package budget

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrEmptyCategoryID = errors.New("empty category_id")
)

var CreateBudgetValidator = validator.OptionalForm(map[string]validator.Validator{
	"category_id": &validator.String{
		Optional: true, // allow reuse in InitUser and CreateCategory
	},
	"budget_date": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckDateStr},
	},
	"amount": &validator.String{
		Optional:   false,
		Validators: []validator.StringFunc{entity.CheckPositiveMonetaryStr},
	},
	"budget_type": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetType},
	},
	"budget_repeat": &validator.UInt32{
		Optional:   true,
		Validators: []validator.UInt32Func{entity.CheckBudgetRepeat},
	},
})

func (h *budgetHandler) CreateBudget(ctx context.Context, req *presenter.CreateBudgetRequest, res *presenter.CreateBudgetResponse) error {
	userID := util.GetUserIDFromCtx(ctx)

	if req.GetCategoryID() == "" {
		return ErrEmptyCategoryID
	}

	useCaseRes, err := h.budgetUseCase.CreateBudget(ctx, req.ToUseCaseReq(userID))
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to create budget, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
