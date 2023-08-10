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
	UserStatusPending
)

var UserStatuses = map[uint32]string{
	uint32(UserStatusNormal):  "normal",
	uint32(UserStatusDeleted): "deleted",
	uint32(UserStatusPending): "pending",
}

type UserUpdate struct {
	UserStatus *uint32
	UpdateTime *uint64
}

func (uu *UserUpdate) GetUserStatus() uint32 {
	if uu != nil && uu.UserStatus != nil {
		return *uu.UserStatus
	}
	return 0
}

func (uu *UserUpdate) SetUserStatus(userStatus *uint32) {
	uu.UserStatus = userStatus
}

func (uu *UserUpdate) GetUpdateTime() uint64 {
	if uu != nil && uu.UpdateTime != nil {
		return *uu.UpdateTime
	}
	return 0
}

func (uu *UserUpdate) SetUpdateTime(updateTime *uint64) {
	uu.UpdateTime = updateTime
}

type UserUpdateOption func(uu *UserUpdate)

func WithUpdateUserStatus(userStatus *uint32) UserUpdateOption {
	return func(uu *UserUpdate) {
		uu.SetUserStatus(userStatus)
	}
}

func NewUserUpdate(opts ...UserUpdateOption) *UserUpdate {
	uu := new(UserUpdate)
	for _, opt := range opts {
		opt(uu)
	}
	return uu
}

type User struct {
	UserID     *string
	Email      *string
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
		u.SetUserID(userID)
	}
}

func WithUsername(username *string) UserOption {
	return func(u *User) {
		u.SetUsername(username)
	}
}

func WithUserEmail(email *string) UserOption {
	return func(u *User) {
		u.SetEmail(email)
	}
}

func WithUserHash(hash *string) UserOption {
	return func(u *User) {
		u.SetHash(hash)
	}
}

func WithUserSalt(salt *string) UserOption {
	return func(u *User) {
		u.SetSalt(salt)
	}
}

func WithUserStatus(userStatus *uint32) UserOption {
	return func(u *User) {
		u.SetUserStatus(userStatus)
	}
}

func WithUserCreateTime(createTime *uint64) UserOption {
	return func(u *User) {
		u.SetCreateTime(createTime)
	}
}

func WithUserUpdateTime(updateTime *uint64) UserOption {
	return func(u *User) {
		u.SetUpdateTime(updateTime)
	}
}

func NewUser(email, password string, opts ...UserOption) (*User, error) {
	now := uint64(time.Now().UnixMilli())
	u := &User{
		Email:      goutil.String(email),
		UserStatus: goutil.Uint32(uint32(UserStatusPending)),
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
			u.SetSalt(goutil.String(salt))
		}

		hash, err := u.hashPassword(password, salt)
		if err != nil {
			return nil, err
		}
		u.SetHash(goutil.String(string(hash)))
	}

	u.checkOpts()

	return u, nil
}

func (u *User) checkOpts() {}

func (u *User) Update(uu *UserUpdate) (userUpdate *UserUpdate, hasUpdate bool) {
	userUpdate = new(UserUpdate)

	if uu.UserStatus != nil && uu.GetUserStatus() != u.GetUserStatus() {
		hasUpdate = true
		u.SetUserStatus(uu.UserStatus)

		defer func() {
			userUpdate.SetUserStatus(u.UserStatus)
		}()
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	u.SetUpdateTime(now)

	u.checkOpts()

	userUpdate.SetUpdateTime(now)

	return
}

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
	u.UserID = userID
}

func (u *User) GetUsername() string {
	if u != nil && u.Username != nil {
		return *u.Username
	}
	return ""
}

func (u *User) SetUsername(username *string) {
	u.Username = username
}

func (u *User) GetEmail() string {
	if u != nil && u.Email != nil {
		return *u.Email
	}
	return ""
}

func (u *User) SetEmail(email *string) {
	u.Email = email
}

func (u *User) GetUserStatus() uint32 {
	if u != nil && u.UserStatus != nil {
		return *u.UserStatus
	}
	return 0
}

func (u *User) SetUserStatus(userStatus *uint32) {
	u.UserStatus = userStatus
}

func (u *User) GetHash() string {
	if u != nil && u.Hash != nil {
		return *u.Hash
	}
	return ""
}

func (u *User) SetHash(hash *string) {
	u.Hash = hash
}

func (u *User) GetSalt() string {
	if u != nil && u.Salt != nil {
		return *u.Salt
	}
	return ""
}

func (u *User) SetSalt(salt *string) {
	u.Salt = salt
}

func (u *User) GetCreateTime() uint64 {
	if u != nil && u.CreateTime != nil {
		return *u.CreateTime
	}
	return 0
}

func (u *User) SetCreateTime(createTime *uint64) {
	u.CreateTime = createTime
}

func (u *User) GetUpdateTime() uint64 {
	if u != nil && u.UpdateTime != nil {
		return *u.UpdateTime
	}
	return 0
}

func (u *User) SetUpdateTime(updateTime *uint64) {
	u.UpdateTime = updateTime
}

func (u *User) IsNormal() bool {
	return u.GetUserStatus() == uint32(UserStatusNormal)
}

func (u *User) IsPendingVerification() bool {
	return u.GetUserStatus() == uint32(UserStatusPending)
}
