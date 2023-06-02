package user

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
)

type UseCase interface {
	GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error)

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

type LogInResponse struct {
	AccessToken *string
}

func (m *LogInResponse) GetAccessToken() string {
	if m != nil && m.AccessToken != nil {
		return *m.AccessToken
	}
	return ""
}
