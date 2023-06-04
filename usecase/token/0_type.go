package token

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateAuthTokenPair(ctx context.Context, req *CreateAuthTokenPairRequest) (*CreateAuthTokenPairResponse, error)
	ValidateAccessToken(ctx context.Context, req *ValidateAccessTokenRequest) (*ValidateAccessTokenResponse, error)
}

type TokenPair struct {
	AccessToken  *string
	RefreshToken *string // TODO
}

type ValidateAccessTokenRequest struct {
	AccessToken *string
}

func (m *ValidateAccessTokenRequest) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

type ValidateAccessTokenResponse struct {
	UserID *string
}

func (m *ValidateAccessTokenResponse) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

type CreateAuthTokenPairRequest struct {
	UserID *string
}

func (m *CreateAuthTokenPairRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateAuthTokenPairRequest) ToCustomClaims() *entity.CustomClaims {
	return &entity.CustomClaims{
		UserID: m.UserID,
	}
}

type CreateAuthTokenPairResponse struct {
	*TokenPair
}

func (m *CreateAuthTokenPairResponse) GetTokenPair() *TokenPair {
	if m != nil && m.TokenPair != nil {
		return m.TokenPair
	}
	return nil
}
