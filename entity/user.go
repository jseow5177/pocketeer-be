package entity

import (
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

func (u *User) IsPasswordCorrect(password string) (bool, error) {
	hash, err := u.hashPassword(password, u.GetSalt())
	if err != nil {
		return false, err
	}
	return u.GetHash() == string(hash), nil
}

func (u *User) SetHash(password string) error {
	newSalt, err := u.createSalt()
	if err != nil {
		return err
	}

	hash, err := u.hashPassword(password, newSalt)
	if err != nil {
		return err
	}

	u.Salt = goutil.String(newSalt)
	u.Hash = goutil.String(string(hash))

	return nil
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
