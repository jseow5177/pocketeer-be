package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepo interface {
	Get(ctx context.Context, uf *UserFilter) (*entity.User, error)
	Create(ctx context.Context, u *entity.User) (string, error)
}

type UserFilter struct {
	UserID     *string `filter:"user_id"`
	UserName   *string `filter:"username"`
	UserStatus *uint32 `filter:"user_status"`
}

func (f *UserFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *UserFilter) GetUserName() string {
	if f != nil && f.UserName != nil {
		return *f.UserName
	}
	return ""
}

func (f *UserFilter) GetCategoryType() uint32 {
	if f != nil && f.UserStatus != nil {
		return *f.UserStatus
	}
	return 0
}
