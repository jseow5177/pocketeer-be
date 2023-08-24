package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/user"
)

type User struct {
	UserID     *string `json:"user_id,omitempty"`
	Email      *string `json:"email,omitempty"`
	Username   *string `json:"username,omitempty"`
	UserFlag   *uint32 `json:"user_flag,omitempty"`
	UserStatus *uint32 `json:"user_status,omitempty"`
	CreateTime *uint64 `json:"create_time,omitempty"`
	UpdateTime *uint64 `json:"update_time,omitempty"`
}

func (u *User) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *User) GetUsername() string {
	if u != nil && u.Username != nil {
		return *u.Username
	}
	return ""
}

func (u *User) GetEmail() string {
	if u != nil && u.Email != nil {
		return *u.Email
	}
	return ""
}

func (u *User) GetUserFlag() uint32 {
	if u != nil && u.UserFlag != nil {
		return *u.UserFlag
	}
	return 0
}

func (u *User) GetUserStatus() uint32 {
	if u != nil && u.UserStatus != nil {
		return *u.UserStatus
	}
	return 0
}

func (u *User) GetCreateTime() uint64 {
	if u != nil && u.CreateTime != nil {
		return *u.CreateTime
	}
	return 0
}

func (u *User) GetUpdateTime() uint64 {
	if u != nil && u.UpdateTime != nil {
		return *u.UpdateTime
	}
	return 0
}

type GetUserRequest struct{}

func (m *GetUserRequest) ToUseCaseReq(userID string) *user.GetUserRequest {
	return &user.GetUserRequest{
		UserID: goutil.String(userID),
	}
}

type GetUserResponse struct {
	User *User `json:"user,omitempty"`
}

func (m *GetUserResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *GetUserResponse) Set(useCaseRes *user.GetUserResponse) {
	m.User = toUser(useCaseRes.User)
}

type SignUpRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
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

func (m *SignUpRequest) ToUseCaseReq() *user.SignUpRequest {
	return &user.SignUpRequest{
		Email:    m.Email,
		Password: m.Password,
	}
}

type SignUpResponse struct {
	User *User `json:"user,omitempty"`
}

func (m *SignUpResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

func (m *SignUpResponse) Set(useCaseRes *user.SignUpResponse) {
	m.User = toUser(useCaseRes.User)
}

type LogInRequest struct {
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
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

func (m *LogInRequest) ToUseCaseReq() *user.LogInRequest {
	return &user.LogInRequest{
		Email:    m.Email,
		Password: m.Password,
	}
}

func (m *LogInResponse) Set(useCaseRes *user.LogInResponse) {
	m.AccessToken = useCaseRes.AccessToken
	m.User = toUser(useCaseRes.User)
}

type LogInResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
	User        *User   `json:"user,omitempty"`
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}

func (m *LogInResponse) GetUser() *User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type VerifyEmailRequest struct {
	Email *string `json:"email,omitempty"`
	Code  *string `json:"code,omitempty"`
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

func (m *VerifyEmailRequest) ToUseCaseReq() *user.VerifyEmailRequest {
	return &user.VerifyEmailRequest{
		Email: m.Email,
		Code:  m.Code,
	}
}

type VerifyEmailResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
	User        *User   `json:"user,omitempty"`
}

func (m *VerifyEmailResponse) Set(useCaseRes *user.VerifyEmailResponse) {
	m.AccessToken = useCaseRes.AccessToken
	m.User = toUser(useCaseRes.User)
}

type InitUserRequest struct{}

func (m *InitUserRequest) ToUseCaseReq(userID string) *user.InitUserRequest {
	return &user.InitUserRequest{
		UserID: goutil.String(userID),
	}
}

type InitUserResponse struct{}

func (m *InitUserResponse) Set(useCaseRes *user.InitUserResponse) {}

type SendOTPRequest struct {
	Email *string `json:"email,omitempty"`
}

func (m *SendOTPRequest) GetEmail() string {
	if m != nil && m.Email != nil {
		return *m.Email
	}
	return ""
}

func (m *SendOTPRequest) ToUseCaseReq() *user.SendOTPRequest {
	return &user.SendOTPRequest{
		Email: m.Email,
	}
}

type SendOTPResponse struct{}

func (m *SendOTPResponse) Set(useCaseRes *user.SendOTPResponse) {}
