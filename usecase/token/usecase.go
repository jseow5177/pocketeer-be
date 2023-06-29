package token

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/rs/zerolog/log"
)

var (
	ErrUnsupportedTokenType = errors.New("unsupported token type")
	ErrInvalidToken         = errors.New("invalid token")
)

type tokenUseCase struct {
	tokenCfg *config.Tokens
}

func NewTokenUseCase(tokenCfg *config.Tokens) UseCase {
	return &tokenUseCase{
		tokenCfg,
	}
}

func (uc *tokenUseCase) CreateToken(ctx context.Context, req *CreateTokenRequest) (*CreateTokenResponse, error) {
	tCfg, err := uc.getTokenConfigByType(req.GetTokenType())
	if err != nil {
		return nil, err
	}

	t := entity.NewToken(
		tCfg.Secret,
		tCfg.ExpiresIn,
		entity.WithTokenClaims(req.GetCustomClaims()),
		entity.WithTokenIssuer(goutil.String(tCfg.Issuer)),
	)
	tokenID, token, err := t.Sign()
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to sign token, err: %v", err)
		return nil, err
	}

	return &CreateTokenResponse{
		TokenID: goutil.String(tokenID),
		Token:   goutil.String(token),
	}, nil
}

func (uc *tokenUseCase) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	tCfg, err := uc.getTokenConfigByType(req.GetTokenType())
	if err != nil {
		return nil, err
	}

	tokenID, claims, err := entity.ParseToken(req.GetToken(), tCfg.Secret)
	if err != nil {
		log.Ctx(ctx).Error().Msgf("fail to parse token, err: %v", err)
		return nil, ErrInvalidToken
	}

	return &ValidateTokenResponse{
		TokenID:      goutil.String(tokenID),
		CustomClaims: claims,
	}, nil
}

func (uc *tokenUseCase) getTokenConfigByType(tokenType uint32) (*config.Token, error) {
	switch tokenType {
	case uint32(entity.TokenTypeAccess):
		return uc.tokenCfg.AccessToken, nil
	case uint32(entity.TokenTypeRefresh):
		return uc.tokenCfg.RefreshToken, nil
	}
	return nil, ErrUnsupportedTokenType
}
