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

func (uc *tokenUseCase) CreateAuthTokenPair(ctx context.Context, req *CreateAuthTokenPairRequest) (*CreateAuthTokenPairResponse, error) {
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

	return &CreateAuthTokenPairResponse{
		TokenPair: &TokenPair{
			AccessToken:  goutil.String(accessToken),
			RefreshToken: goutil.String(refreshToken),
		},
	}, nil
}

func (uc *tokenUseCase) ValidateAccessToken(ctx context.Context, req *ValidateAccessTokenRequest) (*ValidateAccessTokenResponse, error) {
	_, claims, err := entity.ParseToken(req.GetAccessToken(), uc.accessTokenCfg.Secret)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("invalid access token, err: %v", err)
		return nil, err
	}

	// TODO: SANITY CHECK!! CHECK IF USER EXISTS

	return &ValidateAccessTokenResponse{
		UserID: claims.UserID,
	}, nil
}
