package token

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateTokenPair(ctx context.Context, req *CreateTokenPairRequest) (*CreateTokenPairResponse, error)
}

type TokenPair struct {
	AccessToken  *string
	RefreshToken *string // TODO
}

type CreateTokenPairRequest struct {
	UserID *string
}

func (m *CreateTokenPairRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateTokenPairRequest) ToCustomClaims() *entity.CustomClaims {
	return &entity.CustomClaims{
		UserID: m.UserID,
	}
}

type CreateTokenPairResponse struct {
	*TokenPair
}

func (m *CreateTokenPairResponse) GetTokenPair() *TokenPair {
	if m != nil && m.TokenPair != nil {
		return m.TokenPair
	}
	return nil
}
