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
	InitUser(ctx context.Context, req *InitUserRequest) (*InitUserResponse, error)
	SendOTP(ctx context.Context, req *SendOTPRequest) (*SendOTPResponse, error)
	SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error)
	LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error)
	UpdateUserMeta(ctx context.Context, req *UpdateUserMetaRequest) (*UpdateUserMetaResponse, error)
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

func (m *SignUpRequest) ToOTPFilter() *repo.OTPFilter {
	return &repo.OTPFilter{
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

func (m *SignUpRequest) ToUserUpdate() *entity.UserUpdate {
	return entity.NewUserUpdate(
		entity.WithUpdateUserPassword(m.Password),
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
	User        *entity.User
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *LogInResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type VerifyEmailRequest struct {
	Email *string
	Code  *string
}

func (m *VerifyEmailRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *VerifyEmailRequest) GetCode() string {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return ""
}

func (m *VerifyEmailRequest) ToOTPFilter() *repo.OTPFilter {
	return &repo.OTPFilter{
		Email: m.Email,
	}
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

type VerifyEmailResponse struct {
	AccessToken *string
	User        *entity.User
}

func (m *VerifyEmailResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *VerifyEmailResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

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

type InitUserRequest struct {
	UserID   *string
	Currency *string
}

func (m *InitUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *InitUserRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *InitUserRequest) ToUserUpdate() *entity.UserUpdate {
	return entity.NewUserUpdate(
		entity.WithUpdateUserFlag(goutil.Uint32(uint32(entity.UserFlagDefault))),
		entity.WithUpdateUserMeta(
			entity.NewUserMetaUpdate(
				entity.WithUpdateUserMetaCurrency(m.Currency),
			),
		),
	)
}

func (m *InitUserRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		UserID:     m.UserID,
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
}

type InitUserResponse struct{}

type SendOTPRequest struct {
	Email *string
}

func (m *SendOTPRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SendOTPRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		Email: m.Email,
	}
}

type SendOTPResponse struct {
	Email *string
}

type UpdateUserMetaRequest struct {
	UserID   *string
	Currency *string
}

func (m *UpdateUserMetaRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateUserMetaRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *UpdateUserMetaRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		UserID: m.UserID,
	}
}

func (m *UpdateUserMetaRequest) ToUserUpdate() *entity.UserUpdate {
	return entity.NewUserUpdate(
		entity.WithUpdateUserMeta(
			entity.NewUserMetaUpdate(
				entity.WithUpdateUserMetaCurrency(m.Currency),
			),
		),
	)
}

type UpdateUserMetaResponse struct {
	User *entity.User
}

func (m *UpdateUserMetaResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}
