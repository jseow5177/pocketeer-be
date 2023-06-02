package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID     primitive.ObjectID `bson:"user_id,omitempty"`
	Username   *string            `bson:"username,omitempty"`
	UserStatus *uint32            `bson:"user_status,omitempty"`
	Password   *string            `bson:"password,omitempty"`
	Salt       *string            `bson:"salt,omitempty"`
	CreateTime *uint64            `bson:"create_time,omitempty"`
	UpdateTime *uint64            `bson:"update_time,omitempty"`
}

func ToUserModel(u *entity.User) *User {
	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(u.GetUserID()) {
		objID, _ = primitive.ObjectIDFromHex(u.GetUserID())
	}

	return &User{
		UserID:     objID,
		Username:   u.Username,
		UserStatus: u.UserStatus,
		Password:   u.Password,
		Salt:       u.Salt,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

func ToUserEntity(u *User) *entity.User {
	return &entity.User{
		UserID:     goutil.String(u.GetUserID()),
		Username:   u.Username,
		UserStatus: u.UserStatus,
		Password:   u.Password,
		Salt:       u.Salt,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

func (u *User) GetUserID() string {
	if u != nil {
		return u.UserID.Hex()
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

func (u *User) GetPassword() string {
	if u != nil && u.Password != nil {
		return *u.Password
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
