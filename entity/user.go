package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UserStatus uint32

const (
	UserStatusInvalid UserStatus = iota
	UserStatusNormal
	UserStatusDeleted
)

var UserStatuses = map[uint32]string{
	uint32(UserStatusNormal):  "normal",
	uint32(UserStatusDeleted): "deleted",
}

type User struct {
	UserID     *string
	Username   *string
	UserStatus *uint32
	Hash       *string
	Salt       *string
	CreateTime *uint64
	UpdateTime *uint64
}

type UserOption = func(u *User)

func WithUserID(userID *string) UserOption {
	return func(u *User) {
		u.UserID = userID
	}
}

func WithUsername(username *string) UserOption {
	return func(u *User) {
		u.Username = username
	}
}

func WithHash(hash *string) UserOption {
	return func(u *User) {
		u.Hash = hash
	}
}

func WithSalt(salt *string) UserOption {
	return func(u *User) {
		u.Salt = salt
	}
}

func WithUserStatus(userStatus *uint32) UserOption {
	return func(u *User) {
		u.UserStatus = userStatus
	}
}

func WithUserCreateTime(createTime *uint64) UserOption {
	return func(u *User) {
		u.CreateTime = createTime
	}
}

func WithUserUpdateTime(updateTime *uint64) UserOption {
	return func(u *User) {
		u.UpdateTime = updateTime
	}
}

func NewUser(username, password string, opts ...UserOption) (*User, error) {
	now := uint64(time.Now().Unix())
	u := &User{
		Username:   goutil.String(username),
		UserStatus: goutil.Uint32(uint32(UserStatusNormal)),
		Hash:       goutil.String(""),
		Salt:       goutil.String(""),
		CreateTime: goutil.Uint64(now),
		UpdateTime: goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(u)
	}

	if password != "" {
		salt := u.GetSalt()
		if salt == "" {
			var err error
			salt, err = u.createSalt()
			if err != nil {
				return nil, err
			}
			setUser(u, WithSalt(goutil.String(salt)))
		}

		hash, err := u.hashPassword(password, salt)
		if err != nil {
			return nil, err
		}
		setUser(u, WithHash(goutil.String(string(hash))))
	}

	u.checkOpts()

	return u, nil
}

func setUser(u *User, opts ...UserOption) {
	for _, opt := range opts {
		opt(u)
	}
}

func (u *User) checkOpts() {}

func (u *User) IsPasswordCorrect(password string) (bool, error) {
	hash, err := u.hashPassword(password, u.GetSalt())
	if err != nil {
		return false, err
	}
	return u.GetHash() == string(hash), nil
}

func (u *User) hashPassword(password, salt string) ([]byte, error) {
	return goutil.HMACSha256(password, []byte(salt))
}

func (u *User) createSalt() (string, error) {
	salt, err := goutil.RandByte(config.SaltByteSize)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func (u *User) GetUserID() string {
	if u != nil && u.UserID != nil {
		return *u.UserID
	}
	return ""
}

func (u *User) SetUserID(userID *string) {
	setUser(u, WithUserID(userID))
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

func (u *User) GetHash() string {
	if u != nil && u.Hash != nil {
		return *u.Hash
	}
	return ""
}

func (u *User) GetSalt() string {
	if u != nil && u.Salt != nil {
		return *u.Salt
	}
	return ""
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
