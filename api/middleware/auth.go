package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/usecase/user"
	"github.com/jseow5177/pockteer-be/util"
	"github.com/rs/zerolog/log"
)

var (
	ErrUserNotAuthenticated = errors.New("user not authenticated")
)

type AuthMiddleware struct {
	userUseCase user.UseCase
}

func NewAuthMiddleware(userUseCase user.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		userUseCase,
	}
}

func (am *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authErr := errutil.UnauthorizedError(ErrUserNotAuthenticated)

		// get token from auth header
		authHeader := r.Header.Get("Authorization")
		accessToken := am.stripBearerPrefix(authHeader)

		if accessToken == "" {
			log.Ctx(ctx).Error().Msg("token is empty")
			httputil.ReturnServerResponse(w, nil, authErr)
			return
		}

		res, err := am.userUseCase.IsAuthenticated(ctx, &user.IsAuthenticatedRequest{
			AccessToken: goutil.String(accessToken),
		})
		if err != nil {
			httputil.ReturnServerResponse(w, nil, authErr)
			return
		}

		r = r.WithContext(util.SetUserIDToCtx(ctx, res.User.GetUserID()))

		next.ServeHTTP(w, r)
	})
}

func (am *AuthMiddleware) stripBearerPrefix(authHeader string) string {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:]
	}
	return ""
}
