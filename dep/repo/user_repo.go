package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrUserNotFound = errutil.NotFoundError(errors.New("user not found"))
)

type UserRepo interface {
	Get(ctx context.Context, uf *UserFilter) (*entity.User, error)
	GetMany(ctx context.Context, uf *UserFilter) ([]*entity.User, error)

	Create(ctx context.Context, u *entity.User) (string, error)
	Update(ctx context.Context, uf *UserFilter, uu *entity.UserUpdate) error
}

type UserFilter struct {
	UserID     *string `filter:"_id"`
	Email      *string `filter:"email"`
	UserStatus *uint32 `filter:"user_status"`
	Paging     *Paging `filter:"-"`
}

type UserFilterOption = func(uf *UserFilter)

func WithUserID(userID *string) UserFilterOption {
	return func(uf *UserFilter) {
		uf.UserID = userID
	}
}

func WithUserEmail(email *string) UserFilterOption {
	return func(uf *UserFilter) {
		uf.Email = email
	}
}

func WithUserStatus(userStatus *uint32) UserFilterOption {
	return func(uf *UserFilter) {
		uf.UserStatus = userStatus
	}
}

func WithUserPaging(paging *Paging) UserFilterOption {
	return func(uf *UserFilter) {
		uf.Paging = paging
	}
}

func NewUserFilter(opts ...UserFilterOption) *UserFilter {
	uf := &UserFilter{
		UserStatus: goutil.Uint32(uint32(entity.UserStatusNormal)),
	}
	for _, opt := range opts {
		opt(uf)
	}
	return uf
}

func (f *UserFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *UserFilter) GetEmail() string {
	if f != nil && f.Email != nil {
		return *f.Email
	}
	return ""
}

func (f *UserFilter) GetCategoryType() uint32 {
	if f != nil && f.UserStatus != nil {
		return *f.UserStatus
	}
	return 0
}
