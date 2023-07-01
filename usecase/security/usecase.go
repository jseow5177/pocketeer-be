package security

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/api"
	"github.com/rs/zerolog/log"
)

type securityUseCase struct {
	securityAPI api.SecurityAPI
}

func NewSecurityUseCase(securityAPI api.SecurityAPI) UseCase {
	return &securityUseCase{
		securityAPI,
	}
}

func (uc *securityUseCase) SearchSecurities(ctx context.Context, req *SearchSecuritiesRequest) (*SearchSecuritiesResponse, error) {
	ss, err := uc.securityAPI.SearchSecurities(ctx, req.GetKeyword())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to search securities, err: %v", err)
		return nil, err
	}

	return &SearchSecuritiesResponse{
		Securities: ss,
	}, nil
}
