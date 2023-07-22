package security

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/rs/zerolog/log"
)

type securityUseCase struct {
	securityRepo repo.SecurityRepo
}

func NewSecurityUseCase(securityRepo repo.SecurityRepo) UseCase {
	return &securityUseCase{
		securityRepo,
	}
}

func (uc *securityUseCase) SearchSecurities(ctx context.Context, req *SearchSecuritiesRequest) (*SearchSecuritiesResponse, error) {
	ss, err := uc.securityRepo.GetMany(ctx, req.ToSecurityFilter())
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to search securities, err: %v", err)
		return nil, err
	}

	return &SearchSecuritiesResponse{
		Securities: ss,
	}, nil
}
