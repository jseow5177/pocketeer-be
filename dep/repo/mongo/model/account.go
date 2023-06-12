package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	AccountID   primitive.ObjectID `bson:"_id,omitempty"`
	UserID      *string            `bson:"user_id,omitempty"`
	AccountName *string            `bson:"account_name,omitempty"`
	Balance     *float64           `bson:"balance,omitempty"`
	AccountType *uint32            `bson:"account_type,omitempty"`
	CreateTime  *uint64            `bson:"create_time,omitempty"`
	UpdateTime  *uint64            `bson:"update_time,omitempty"`
}
