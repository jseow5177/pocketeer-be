package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/usecase/token"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
)

type AuthMiddleware struct {
	tokenUseCase token.UseCase
}

func NewAuthMiddleware(tokenUseCase token.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		tokenUseCase: tokenUseCase,
	}
}

func (am *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// get token from auth header
		authHeader := r.Header.Get("Authorization")
		accessToken := am.stripBearerPrefix(authHeader)

		if accessToken == "" {
			httputil.ReturnServerResponse(w, nil, ErrUserNotAuthenticated)
			return
		}

		res, err := am.tokenUseCase.ValidateAccessToken(ctx, &token.ValidateAccessTokenRequest{
			AccessToken: goutil.String(accessToken),
		})
		if err != nil {
			httputil.ReturnServerResponse(w, nil, ErrUserNotAuthenticated)
			return
		}

		r = r.WithContext(util.SetUserIDToCtx(ctx, res.GetUserID()))

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) stripBearerPrefix(authHeader string) string {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:]
	}
	return ""
}
