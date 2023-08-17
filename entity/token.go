package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var signTokenMethod = jwt.SigningMethodHS256

type CustomClaims struct {
	UserID *string `json:"user_id,omitempty"`
	Email  *string `json:"email,omitempty"`
}

func (cc *CustomClaims) GetUserID() string {
	if cc != nil && cc.UserID != nil {
		return *cc.UserID
	}
	return ""
}

func (cc *CustomClaims) GetEmail() string {
	if cc != nil && cc.Email != nil {
		return *cc.Email
	}
	return ""
}

type claims struct {
	*CustomClaims
	jwt.RegisteredClaims
}

type TokenType uint32

const (
	TokenTypeAccess TokenType = iota
	TokenTypeRefresh
)

type Token struct {
	Claims    *CustomClaims
	ExpiresIn *int64
	Issuer    *string
	Secret    *string
}

type TokenOption = func(t *Token)

func WithTokenIssuer(issuer *string) TokenOption {
	return func(t *Token) {
		t.Issuer = issuer
	}
}

func WithTokenClaims(claims *CustomClaims) TokenOption {
	return func(t *Token) {
		t.Claims = claims
	}
}

func (t *Token) GetClaims() *CustomClaims {
	if t != nil && t.Claims != nil {
		return t.Claims
	}
	return nil
}

func (t *Token) GetExpiresIn() int64 {
	if t != nil && t.ExpiresIn != nil {
		return *t.ExpiresIn
	}
	return 0
}

func (t *Token) GetIssuer() string {
	if t != nil && t.Issuer != nil {
		return *t.Issuer
	}
	return ""
}

func (t *Token) GetSecret() string {
	if t != nil && t.Secret != nil {
		return *t.Secret
	}
	return ""
}

func NewToken(secret string, expiresIn int64, opts ...TokenOption) *Token {
	t := &Token{
		Secret:    goutil.String(secret),
		ExpiresIn: goutil.Int64(expiresIn),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func ParseToken(plainToken, secret string) (jti string, customClaims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(plainToken, new(claims), func(_ *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{signTokenMethod.Name}))
	if err != nil || !token.Valid {
		return "", nil, err
	}

	cms := token.Claims.(*claims)

	return cms.ID, cms.CustomClaims, nil
}

func (t *Token) Sign() (jti string, signedToken string, err error) {
	jti = goutil.NextXID()

	expiresAt := time.Now().Add(time.Duration(t.GetExpiresIn()) * time.Second)

	claims := new(claims)
	claims.CustomClaims = t.GetClaims()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    t.GetIssuer(),
		ID:        jti,
	}

	token := jwt.NewWithClaims(signTokenMethod, claims)
	signedToken, err = token.SignedString([]byte(t.GetSecret()))
	if err != nil {
		return "", "", err
	}

	return jti, signedToken, nil
}
