package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var signMethod = jwt.SigningMethodHS256

type CustomClaims struct {
	UserID *string `json:"user_id,omitempty"`
}

type claims struct {
	*CustomClaims
	jwt.RegisteredClaims
}

type Token struct {
	claims    *CustomClaims
	expiresIn int64
	issuer    string
	secret    string
}

func NewToken(tokenCfg *config.Token, customClaims *CustomClaims) *Token {
	return &Token{
		expiresIn: tokenCfg.ExpiresIn,
		issuer:    tokenCfg.Issuer,
		secret:    tokenCfg.Secret,
		claims:    customClaims,
	}
}

func (t *Token) Sign() (jti string, signedToken string, err error) {
	jti = goutil.NextXID()

	expiresAt := time.Now().Add(time.Duration(t.expiresIn) * time.Second)

	claims := new(claims)
	claims.CustomClaims = t.claims
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    t.issuer,
		ID:        jti,
	}

	token := jwt.NewWithClaims(signMethod, claims)
	signedToken, err = token.SignedString([]byte(t.secret))
	if err != nil {
		return "", "", err
	}

	return jti, signedToken, nil
}
