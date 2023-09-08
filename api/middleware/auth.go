package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/usecase/user"
	"github.com/rs/zerolog/log"
)

var (
	ErrUserNotAuthenticated  = errors.New("user not authenticated")
	ErrAdminNotAuthenticated = errors.New("admin not authenticated")
)

type AdminAuthMiddleware struct {
	adminCfg *config.ServerAdmin
}

func NewAdminAuthMiddleware(adminCfg *config.ServerAdmin) *AdminAuthMiddleware {
	return &AdminAuthMiddleware{
		adminCfg,
	}
}

func (am *AdminAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		s := am.stripBasicPrefix(authHeader)

		creds, err := goutil.Base64Decode(s, base64.StdPadding)
		if err != nil {
			log.Ctx(ctx).Error().Msgf("fail to decode credentials, str: %v, err: %v", s, err)
			httputil.ReturnServerResponse(w, nil, errutil.UnauthorizedError(ErrAdminNotAuthenticated))
			return
		}

		up := strings.Split(string(creds), ":")
		if len(up) != 2 || (up[0] != am.adminCfg.Username || up[1] != am.adminCfg.Password) {
			log.Ctx(ctx).Error().Msgf("invalid credentials: %v", up)
			httputil.ReturnServerResponse(w, nil, errutil.UnauthorizedError(ErrAdminNotAuthenticated))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (am *AdminAuthMiddleware) stripBasicPrefix(authHeader string) string {
	if len(authHeader) > 5 && strings.ToUpper(authHeader[0:6]) == "BASIC " {
		return authHeader[6:]
	}
	return ""
}

type UserAuthMiddleware struct {
	userUseCase user.UseCase
}

func NewUserAuthMiddleware(userUseCase user.UseCase) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		userUseCase,
	}
}

func (am *UserAuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// get token from auth header
		authHeader := r.Header.Get("Authorization")
		accessToken := am.stripBasicPrefix(authHeader)

		if accessToken == "" {
			log.Ctx(ctx).Error().Msg("token is empty")
			httputil.ReturnServerResponse(w, nil, errutil.UnauthorizedError(ErrUserNotAuthenticated))
			return
		}

		res, err := am.userUseCase.IsAuthenticated(ctx, &user.IsAuthenticatedRequest{
			AccessToken: goutil.String(accessToken),
		})
		if err != nil {
			// TODO: Return NotAuthenticated only when user is not found or token is invalid
			// Else, return 500
			httputil.ReturnServerResponse(w, nil, errutil.UnauthorizedError(ErrUserNotAuthenticated))
			return
		}

		r = r.WithContext(entity.SetUserToCtx(ctx, res.User))

		next.ServeHTTP(w, r)
	})
}

func (am *UserAuthMiddleware) stripBasicPrefix(authHeader string) string {
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
		return authHeader[7:]
	}
	return ""
}
