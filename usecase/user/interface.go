package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/token"
	"github.com/jseow5177/pockteer-be/util"
)

type UseCase interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)
	IsAuthenticated(ctx context.Context, req *IsAuthenticatedRequest) (*IsAuthenticatedResponse, error)
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error)

	SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error)
	LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error)
}

type GetUserRequest struct {
	UserID *string
}

func (m *GetUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetUserRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		UserID: m.UserID,
	}
}

type GetUserResponse struct {
	*entity.User
}

func (m *GetUserResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type SignUpRequest struct {
	Email    *string
	Password *string
}

func (m *SignUpRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SignUpRequest) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *SignUpRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		Email: m.Email,
	}
}

func (m *SignUpRequest) ToUserEntity() (*entity.User, error) {
	username := util.GetEmailPrefix(m.GetEmail())
	return entity.NewUser(
		m.GetEmail(),
		m.GetPassword(),
		entity.WithUsername(goutil.String(username)),
	)
}

type SignUpResponse struct {
	User *entity.User
}

func (m *SignUpResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type LogInRequest struct {
	Email    *string
	Password *string
}

func (m *LogInRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *LogInRequest) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

func (m *LogInRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		Email:      m.Email,
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
}

type LogInResponse struct {
	AccessToken *string
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

type VerifyEmailRequest struct {
	EmailToken *string
}

func (m *VerifyEmailRequest) GetEmailToken() string {
	if m != nil && m.EmailToken != nil {
		return *m.EmailToken
	}
	return ""
}

func (m *VerifyEmailRequest) ToValidateTokenRequest() (*token.ValidateTokenRequest, error) {
	emailToken, err := goutil.Base64Decode(m.GetEmailToken())
	if err != nil {
		return nil, err
	}
	return &token.ValidateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeEmail)),
		Token:     goutil.String(string(emailToken)),
	}, nil
}

func (m *VerifyEmailRequest) ToUserFilter(email string) *repo.UserFilter {
	return &repo.UserFilter{
		Email:      goutil.String(email),
		UserStatus: goutil.Uint32(uint32(entity.UserStatusPending)),
	}
}

func (m *VerifyEmailRequest) ToUserUpdate() *entity.UserUpdate {
	return entity.NewUserUpdate(
		entity.WithUpdateUserStatus(goutil.Uint32(uint32(entity.UserStatusNormal))),
	)
}

type VerifyEmailResponse struct{}

type IsAuthenticatedRequest struct {
	AccessToken *string
}

func (m *IsAuthenticatedRequest) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *IsAuthenticatedRequest) ToValidateTokenRequest() *token.ValidateTokenRequest {
	return &token.ValidateTokenRequest{
		TokenType: goutil.Uint32(uint32(entity.TokenTypeAccess)),
		Token:     m.AccessToken,
	}
}

func (m *IsAuthenticatedRequest) ToUserFilter(userID string) *repo.UserFilter {
	return &repo.UserFilter{
		UserID:     goutil.String(userID),
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
}

type IsAuthenticatedResponse struct {
	User *entity.User
}

func (m *IsAuthenticatedResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}
