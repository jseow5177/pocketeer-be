package token

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

type tokenUseCase struct {
	accessTokenCfg  *config.Token
	refreshTokenCfg *config.Token
}

func NewTokenUseCase(accessTokenCfg, refreshTokenCfg *config.Token) UseCase {
	return &tokenUseCase{
		accessTokenCfg,
		refreshTokenCfg,
	}
}

func (uc *tokenUseCase) CreateTokenPair(ctx context.Context, req *CreateTokenPairRequest) (*CreateTokenPairResponse, error) {
	// sign access token
	at := entity.NewToken(uc.accessTokenCfg, req.ToCustomClaims())
	_, accessToken, err := at.Sign()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sign access token, err: %v", err)
		return nil, err
	}

	// sign refresh token
	rt := entity.NewToken(uc.accessTokenCfg, req.ToCustomClaims())
	_, refreshToken, err := rt.Sign()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sign refresh token, err: %v", err)
		return nil, err
	}

	return &CreateTokenPairResponse{
		TokenPair: &TokenPair{
			AccessToken:  goutil.String(accessToken),
			RefreshToken: goutil.String(refreshToken),
		},
	}, nil
}
