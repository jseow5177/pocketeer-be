package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UserInitStage uint32

const (
	InitStageOne UserInitStage = iota
	InitStageTwo
	InitStageThree
	InitStageFour
)

type UserMetaUpdate struct {
	InitStage *uint32
}

func (umu *UserMetaUpdate) GetInitStage() uint32 {
	if umu != nil && umu.InitStage != nil {
		return *umu.InitStage
	}
	return 0
}

func (umu *UserMetaUpdate) SetInitStage(initStage *uint32) {
	umu.InitStage = initStage
}

type UserMetaUpdateOption func(uu *UserMetaUpdate)

func WithUpdateUserMetaInitStage(initStage *uint32) UserMetaUpdateOption {
	return func(umu *UserMetaUpdate) {
		umu.SetInitStage(initStage)
	}
}

func NewUserMetaUpdate(opts ...UserMetaUpdateOption) *UserMetaUpdate {
	umu := new(UserMetaUpdate)
	for _, opt := range opts {
		opt(umu)
	}
	return umu
}

type UserMeta struct {
	InitStage *uint32
}

type UserMetaOption = func(um *UserMeta)

func WithUserMetaInitStage(initStage *uint32) UserMetaOption {
	return func(um *UserMeta) {
		um.SetInitStage(initStage)
	}
}

func (um *UserMeta) GetInitStage() uint32 {
	if um != nil && um.InitStage != nil {
		return *um.InitStage
	}
	return 0
}

func (um *UserMeta) SetInitStage(initStage *uint32) {
	um.InitStage = initStage
}

func (um *UserMeta) Update(umu *UserMetaUpdate) *UserMetaUpdate {
	var (
		hasUpdate      bool
		userMetaUpdate = new(UserMetaUpdate)
	)

	if umu.InitStage != nil && um.GetInitStage() != umu.GetInitStage() {
		hasUpdate = true
		um.SetInitStage(umu.InitStage)

		defer func() {
			userMetaUpdate.SetInitStage(um.InitStage)
		}()
	}

	if !hasUpdate {
		return nil
	}

	return userMetaUpdate
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
	Password   *string
	Hash       *string
	UpdateTime *uint64
	Meta       *UserMetaUpdate
}

func (uu *UserUpdate) GetPassword() string {
	if uu != nil && uu.Password != nil {
		return *uu.Password
	}
	return ""
}

func (uu *UserUpdate) SetPassword(password *string) {
	uu.Password = password
}

func (uu *UserUpdate) GetHash() string {
	if uu != nil && uu.Hash != nil {
		return *uu.Hash
	}
	return ""
}

func (uu *UserUpdate) SetHash(hash *string) {
	uu.Hash = hash
}

func (uu *UserUpdate) GetUserFlag() uint32 {
	if uu != nil && uu.UserFlag != nil {
		return *uu.UserFlag
	}
	return 0
}

func (uu *UserUpdate) SetUserFlag(userFlag *uint32) {
	uu.UserFlag = userFlag
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

func (uu *UserUpdate) GetMeta() *UserMetaUpdate {
	if uu != nil && uu.Meta != nil {
		return nil
	}
	return uu.Meta
}

func (uu *UserUpdate) SetMeta(meta *UserMetaUpdate) {
	uu.Meta = meta
}

type UserUpdateOption func(uu *UserUpdate)

func WithUpdateUserPassword(password *string) UserUpdateOption {
	return func(uu *UserUpdate) {
		uu.SetPassword(password)
	}
}

func WithUpdateUserFlag(userFlag *uint32) UserUpdateOption {
	return func(uu *UserUpdate) {
		uu.SetUserFlag(userFlag)
	}
}

func WithUpdateUserStatus(userStatus *uint32) UserUpdateOption {
	return func(uu *UserUpdate) {
		uu.SetUserStatus(userStatus)
	}
}

func WithUpdateUserMeta(meta *UserMetaUpdate) UserUpdateOption {
	return func(uu *UserUpdate) {
		uu.SetMeta(meta)
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
	UserFlag   *uint32
	Hash       *string
	Salt       *string
	CreateTime *uint64
	UpdateTime *uint64
	Meta       *UserMeta
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

func WithUserMeta(userMeta *UserMeta) UserOption {
	return func(u *User) {
		u.Meta = userMeta
	}
}

func NewUser(email, password string, opts ...UserOption) (*User, error) {
	now := uint64(time.Now().UnixMilli())
	u := &User{
		Email:      goutil.String(email),
		UserFlag:   goutil.Uint32(uint32(UserFlagNewUser)),
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
		hash, err := u.getHash(password)
		if err != nil {
			return nil, err
		}
		u.SetHash(goutil.String(string(hash)))
	}

	if err := u.checkOpts(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) checkOpts() error {
	return nil
}

func (u *User) getHash(password string) (string, error) {
	salt := u.GetSalt()
	if salt == "" {
		var err error
		salt, err = u.createSalt()
		if err != nil {
			return "", err
		}
		u.SetSalt(goutil.String(salt))
	}

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

func (u *User) Update(uu *UserUpdate) (*UserUpdate, error) {
	var (
		hasUpdate  bool
		userUpdate = new(UserUpdate)
	)

	if uu.UserStatus != nil && uu.GetUserStatus() != u.GetUserStatus() {
		hasUpdate = true
		u.SetUserStatus(uu.UserStatus)

		defer func() {
			userUpdate.SetUserStatus(u.UserStatus)
		}()
	}

	if uu.UserFlag != nil && uu.GetUserFlag() != u.GetUserFlag() {
		hasUpdate = true
		u.SetUserFlag(uu.UserFlag)

		defer func() {
			userUpdate.SetUserFlag(u.UserFlag)
		}()
	}

	if uu.Password != nil {
		hash, err := u.getHash(uu.GetPassword())
		if err != nil {
			return nil, err
		}

		if u.GetHash() != hash {
			hasUpdate = true
			u.SetHash(goutil.String(hash))

			defer func() {
				userUpdate.SetHash(u.Hash)
			}()
		}
	}

	if uu.Meta != nil {
		umu := u.Meta.Update(uu.Meta)

		if umu != nil {
			hasUpdate = true
			userUpdate.Meta = new(UserMetaUpdate)

			if umu.InitStage != nil {
				u.Meta.SetInitStage(umu.InitStage)

				defer func() {
					userUpdate.Meta.SetInitStage(u.Meta.InitStage)
				}()
			}
		}
	}

	if !hasUpdate {
		return nil, nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	u.SetUpdateTime(now)

	if err := u.checkOpts(); err != nil {
		return nil, err
	}

	userUpdate.SetUpdateTime(now)

	return userUpdate, nil
}

func (u *User) IsSamePassword(password string) (bool, error) {
	hash, err := u.getHash(password)
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
