package presenter

import (
	"github.com/jseow5177/pockteer-be/usecase/user"
)

type User struct {
	UserID     *string `json:"user_id,omitempty"`
	Username   *string `json:"username,omitempty"`
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

type GetUserRequest struct {
	UserID *string `json:"user_id,omitempty"`
}

func (m *GetUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetUserRequest) ToUseCaseReq() *user.GetUserRequest {
	return &user.GetUserRequest{
		UserID: m.UserID,
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
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (m *SignUpRequest) GetUsername() string {
	if m != nil && m.Username != nil {
		return *m.Username
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
		Username: m.Username,
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
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
}

func (m *LogInRequest) GetUsername() string {
	if m != nil && m.Username != nil {
		return *m.Username
	}
	return ""
}

func (m *LogInRequest) GetPassword() string {
	if m != nil && m.Password != nil {
		return *m.Password
	}
	return ""
}

type LogInResponse struct {
	AccessToken *string `json:"access_token,omitempty"`
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}
