package security

import (
	"context"

	"github.com/jseow5177/pockteer-be/api/presenter"
	"github.com/jseow5177/pockteer-be/pkg/validator"
	"github.com/rs/zerolog/log"
)

var SearchSecuritiesValidator = validator.MustForm(map[string]validator.Validator{
	"keyword": &validator.String{},
})

func (h *securityHandler) SearchSecurities(ctx context.Context, req *presenter.SearchSecuritiesRequest, res *presenter.SearchSecuritiesResponse) error {
	useCaseRes, err := h.securityUseCase.SearchSecurities(ctx, req.ToUseCaseReq())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to search securities, err: %v", err)
		return err
	}

	res.Set(useCaseRes)

	return nil
}
