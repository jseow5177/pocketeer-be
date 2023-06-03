package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)

	SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error)
	LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error)
}

type GetUserRequest struct {
	UserID     *string
	UserName   *string
	UserStatus *uint32
}

func (m *GetUserRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetUserRequest) GetUserName() string {
	if m != nil && m.UserName != nil {
		return *m.UserName
	}
	return ""
}

func (m *GetUserRequest) GetUserStatus() uint32 {
	if m != nil && m.UserStatus != nil {
		return *m.UserStatus
	}
	return 0
}

func (m *GetUserRequest) ToUserFilter() *repo.UserFilter {
	return &repo.UserFilter{
		UserID:     m.UserID,
		UserName:   m.UserName,
		UserStatus: m.UserStatus,
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
	Username *string
	Password *string
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

func (m *SignUpRequest) ToGetUserRequest() *GetUserRequest {
	return &GetUserRequest{
		UserName:   m.Username,
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
}

func (m *SignUpRequest) ToUserEntity() *entity.User {
	return &entity.User{
		Username:   m.Username,
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
}

type SignUpResponse struct {
	*entity.User
}

func (m *SignUpResponse) GetUser() *entity.User {
	if m != nil && m.User != nil {
		return m.User
	}
	return nil
}

type LogInRequest struct {
	Username *string
	Password *string
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

func (m *LogInRequest) ToGetUserRequest() *GetUserRequest {
	return &GetUserRequest{
		UserName:   m.Username,
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
