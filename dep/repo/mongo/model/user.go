package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID     primitive.ObjectID `bson:"_id,omitempty"`
	Email      *string            `bson:"email,omitempty"`
	Username   *string            `bson:"username,omitempty"`
	UserStatus *uint32            `bson:"user_status,omitempty"`
	Hash       *string            `bson:"hash,omitempty"`
	Salt       *string            `bson:"salt,omitempty"`
	CreateTime *uint64            `bson:"create_time,omitempty"`
	UpdateTime *uint64            `bson:"update_time,omitempty"`
}

func ToUserModelFromEntity(u *entity.User) *User {
	if u == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(u.GetUserID()) {
		objID, _ = primitive.ObjectIDFromHex(u.GetUserID())
	}

	var encodedHash *string
	if u.Hash != nil {
		encodedHash = goutil.String(goutil.Base64Encode([]byte(u.GetHash())))
	}

	var encodedSalt *string
	if u.Salt != nil {
		encodedSalt = goutil.String(goutil.Base64Encode([]byte(u.GetSalt())))
	}

	return &User{
		UserID:     objID,
		Email:      u.Email,
		Username:   u.Username,
		UserStatus: u.UserStatus,
		Hash:       encodedHash,
		Salt:       encodedSalt,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
	}
}

func ToUserModelFromUpdate(uu *entity.UserUpdate) *User {
	if uu == nil {
		return nil
	}

	return &User{
		UserStatus: uu.UserStatus,
		UpdateTime: uu.UpdateTime,
	}
}

func ToUserEntity(u *User) (*entity.User, error) {
	if u == nil {
		return nil, nil
	}

	var decodedHash *string
	if u.Hash != nil {
		b, err := goutil.Base64Decode(u.GetHash())
		if err != nil {
			return nil, err
		}
		decodedHash = goutil.String(string(b))
	}

	var decodedSalt *string
	if u.Salt != nil {
		b, err := goutil.Base64Decode(u.GetSalt())
		if err != nil {
			return nil, err
		}
		decodedSalt = goutil.String(string(b))
	}

	return entity.NewUser(
		u.GetEmail(),
		"",
		entity.WithUserID(goutil.String(u.GetUserID())),
		entity.WithUserHash(decodedHash),
		entity.WithUserSalt(decodedSalt),
		entity.WithUserStatus(u.UserStatus),
		entity.WithUserCreateTime(u.CreateTime),
		entity.WithUserUpdateTime(u.UpdateTime),
		entity.WithUsername(u.Username),
	)
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

func (u *User) GetEmail() string {
	if u != nil && u.Email != nil {
		return *u.Email
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
