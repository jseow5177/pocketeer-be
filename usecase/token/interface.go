package token

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	CreateToken(ctx context.Context, req *CreateTokenRequest) (*CreateTokenResponse, error)
	ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error)
}

type CreateTokenRequest struct {
	TokenType    *uint32
	CustomClaims *entity.CustomClaims
}

func (m *CreateTokenRequest) GetTokenType() uint32 {
	if m != nil && m.TokenType != nil {
		return *m.TokenType
	}
	return 0
}

func (m *CreateTokenRequest) GetCustomClaims() *entity.CustomClaims {
	if m != nil && m.CustomClaims != nil {
		return m.CustomClaims
	}
	return nil
}

type CreateTokenResponse struct {
	TokenID *string
	Token   *string
}

func (m *CreateTokenResponse) GetTokenID() string {
	if m != nil && m.TokenID != nil {
		return *m.TokenID
	}
	return ""
}

func (m *CreateTokenResponse) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

type ValidateTokenRequest struct {
	TokenType *uint32
	Token     *string
}

func (m *ValidateTokenRequest) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *ValidateTokenRequest) GetTokenType() uint32 {
	if m != nil && m.TokenType != nil {
		return *m.TokenType
	}
	return 0
}

type ValidateTokenResponse struct {
	TokenID      *string
	CustomClaims *entity.CustomClaims
}

func (m *ValidateTokenResponse) GetTokenID() string {
	if m != nil && m.TokenID != nil {
		return *m.TokenID
	}
	return ""
}

func (m *ValidateTokenResponse) GetCustomClaims() *entity.CustomClaims {
	if m != nil && m.CustomClaims != nil {
		return m.CustomClaims
	}
	return nil
}
