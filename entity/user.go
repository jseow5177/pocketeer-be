package entity

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type ctxKey string

const ctxKeyUser ctxKey = "ctx:user"

func SetUserToCtx(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, user)
}

func GetUserFromCtx(ctx context.Context) *User {
	val := ctx.Value(ctxKeyUser)
	if val == nil {
		return nil
	}

	if user, ok := val.(*User); ok {
		return user
	}

	return nil
}

type UserMeta struct {
	Currency *string
	HideInfo *bool
}

func (um *UserMeta) GetCurrency() string {
	if um != nil && um.Currency != nil {
		return *um.Currency
	}
	return ""
}

func (um *UserMeta) SetCurrency(currency *string) {
	um.Currency = currency
}

func (um *UserMeta) GetHideInfo() bool {
	if um != nil && um.HideInfo != nil {
		return *um.HideInfo
	}
	return false
}

func (um *UserMeta) SetHideInfo(hideInfo *bool) {
	um.HideInfo = hideInfo
}

type UserFlag uint32

const (
	UserFlagDefault UserFlag = iota
	UserFlagNewUser
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
	UserFlag   *uint32
	UserStatus *uint32
	UpdateTime *uint64
	Currency   *string
	HideInfo   *bool
	Hash       *string
	Salt       *string
}

func (uu *UserUpdate) GetUserFlag() uint32 {
	if uu != nil && uu.UserFlag != nil {
		return *uu.UserFlag
	}
	return 0
}

func (uu *UserUpdate) GetUserStatus() uint32 {
	if uu != nil && uu.UserStatus != nil {
		return *uu.UserStatus
	}
	return 0
}

func (uu *UserUpdate) GetUpdateTime() uint64 {
	if uu != nil && uu.UpdateTime != nil {
		return *uu.UpdateTime
	}
	return 0
}

func (uu *UserUpdate) GetCurrency() string {
	if uu != nil && uu.Currency != nil {
		return *uu.Currency
	}
	return ""
}

func (uu *UserUpdate) GetHideInfo() bool {
	if uu != nil && uu.HideInfo != nil {
		return *uu.HideInfo
	}
	return false
}

func (uu *UserUpdate) GetHash() string {
	if uu != nil && uu.Hash != nil {
		return *uu.Hash
	}
	return ""
}

func (uu *UserUpdate) GetSalt() string {
	if uu != nil && uu.Salt != nil {
		return *uu.Salt
	}
	return ""
}

type UserUpdateOption func(u *User)

func WithUpdateUserFlag(userFlag *uint32) UserUpdateOption {
	return func(u *User) {
		if userFlag != nil {
			u.SetUserFlag(userFlag)
		}
	}
}

func WithUpdateUserStatus(userStatus *uint32) UserUpdateOption {
	return func(u *User) {
		if userStatus != nil {
			u.SetUserStatus(userStatus)
		}
	}
}

func WithUpdateUserPassword(password *string) UserUpdateOption {
	return func(u *User) {
		if password != nil {
			u.SetPassword(password)
		}
	}
}

func WithUpdateUserCurrency(currency *string) UserUpdateOption {
	return func(u *User) {
		if u.Meta != nil && currency != nil {
			u.Meta.SetCurrency(currency)
		}
	}
}

func WithUpdateUserHideInfo(hideInfo *bool) UserUpdateOption {
	return func(u *User) {
		if u.Meta != nil && hideInfo != nil {
			u.Meta.SetHideInfo(hideInfo)
		}
	}
}

type User struct {
	UserID     *string
	Email      *string
	Username   *string
	UserStatus *uint32
	UserFlag   *uint32
	Password   *string
	Hash       *string
	Salt       *string
	CreateTime *uint64
	UpdateTime *uint64
	Meta       *UserMeta
}

type UserOption func(u *User)

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

func WithUserPassword(password *string) UserOption {
	return func(u *User) {
		u.SetPassword(password)
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

func WithUserFlag(userFlag *uint32) UserOption {
	return func(u *User) {
		u.SetUserFlag(userFlag)
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

func WithUserCurrency(currency *string) UserOption {
	return func(u *User) {
		if u.Meta != nil {
			u.Meta.Currency = currency
		}
	}
}

func WithUserHideInfo(hideInfo *bool) UserOption {
	return func(u *User) {
		if u.Meta != nil {
			u.Meta.HideInfo = hideInfo
		}
	}
}

func NewUser(email string, opts ...UserOption) (*User, error) {
	now := uint64(time.Now().UnixMilli())
	u := &User{
		Email:      goutil.String(email),
		UserFlag:   goutil.Uint32(uint32(UserFlagNewUser)),
		UserStatus: goutil.Uint32(uint32(UserStatusPending)),
		Password:   goutil.String(""),
		Hash:       goutil.String(""),
		Salt:       goutil.String(""),
		CreateTime: goutil.Uint64(now),
		UpdateTime: goutil.Uint64(now),
		Meta: &UserMeta{
			Currency: goutil.String(""),
			HideInfo: goutil.Bool(false),
		},
	}

	for _, opt := range opts {
		opt(u)
	}

	if u.GetPassword() != "" {
		salt, err := u.createSalt()
		if err != nil {
			return nil, err
		}
		u.SetSalt(goutil.String(salt))

		hash, err := u.createHash(u.GetPassword(), salt)
		if err != nil {
			return nil, err
		}
		u.SetHash(goutil.String(string(hash)))
	}

	if err := u.validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) validate() error {
	return nil
}

func (u *User) createHash(password, salt string) (string, error) {
	hash, err := goutil.HMACSha256(password, []byte(salt))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (u *User) createSalt() (string, error) {
	salt, err := goutil.RandByte(config.SaltByteSize)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func (u *User) Clone() (*User, error) {
	return NewUser(
		u.GetEmail(),
		WithUserID(goutil.String(u.GetUserID())),
		WithUserHash(u.Hash),
		WithUserSalt(u.Salt),
		WithUserStatus(u.UserStatus),
		WithUserCreateTime(u.CreateTime),
		WithUserUpdateTime(u.UpdateTime),
		WithUsername(u.Username),
		WithUserFlag(u.UserFlag),
		WithUserCurrency(u.Meta.Currency),
		WithUserHideInfo(u.Meta.HideInfo),
	)
}

func (u *User) ToUserUpdate(old *User) *UserUpdate {
	var (
		hasUpdate bool

		uu = &UserUpdate{
			UpdateTime: u.UpdateTime,
		}
	)

	if old.GetUserStatus() != u.GetUserStatus() {
		hasUpdate = true
		uu.UserStatus = u.UserStatus
	}

	if old.GetUserFlag() != u.GetUserFlag() {
		hasUpdate = true
		uu.UserFlag = u.UserFlag
	}

	if old.GetHash() != u.GetHash() {
		hasUpdate = true
		uu.Hash = u.Hash
	}

	if old.GetSalt() != u.GetSalt() {
		hasUpdate = true
		uu.Salt = u.Salt
	}

	if old.Meta.GetCurrency() != u.Meta.GetCurrency() {
		hasUpdate = true
		uu.Currency = u.Meta.Currency
	}

	if old.Meta.GetHideInfo() != u.Meta.GetHideInfo() {
		hasUpdate = true
		uu.HideInfo = u.Meta.HideInfo
	}

	if hasUpdate {
		return uu
	}

	return nil
}

func (u *User) Update(uus ...UserUpdateOption) (*UserUpdate, error) {
	if len(uus) == 0 {
		return nil, nil
	}

	old, err := u.Clone()
	if err != nil {
		return nil, err
	}

	for _, uu := range uus {
		uu(u)
	}

	if u.GetPassword() != "" {
		salt, err := u.createSalt()
		if err != nil {
			return nil, err
		}
		u.SetSalt(goutil.String(salt))

		hash, err := u.createHash(u.GetPassword(), salt)
		if err != nil {
			return nil, err
		}
		u.SetHash(goutil.String(string(hash)))
	}

	// check
	if err := u.validate(); err != nil {
		return nil, err
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	u.SetUpdateTime(now)

	return u.ToUserUpdate(old), nil
}

func (u *User) IsSamePassword(password string) (bool, error) {
	hash, err := u.createHash(password, u.GetSalt())
	if err != nil {
		return false, err
	}
	return u.GetHash() == hash, nil
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

func (u *User) GetUserFlag() uint32 {
	if u != nil && u.UserFlag != nil {
		return *u.UserFlag
	}
	return 0
}

func (u *User) SetUserFlag(userFlag *uint32) {
	u.UserFlag = userFlag
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

func (u *User) GetPassword() string {
	if u != nil && u.Password != nil {
		return *u.Password
	}
	return ""
}

func (u *User) SetPassword(password *string) {
	u.Password = password
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

func (u *User) GetMeta() *UserMeta {
	if u != nil && u.Meta != nil {
		return nil
	}
	return u.Meta
}

func (u *User) SetMeta(meta *UserMeta) {
	u.Meta = meta
}

func (u *User) IsNormal() bool {
	return u.GetUserStatus() == uint32(UserStatusNormal)
}

func (u *User) IsPendingVerification() bool {
	return u.GetUserStatus() == uint32(UserStatusPending)
}

func (u *User) IsNew() bool {
	return (u.GetUserFlag() & uint32(UserFlagNewUser)) > 0
}
