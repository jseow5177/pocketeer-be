package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	TransactionID   primitive.ObjectID `bson:"_id,omitempty"`
	UserID          *string            `bson:"user_id,omitempty"`
	CategoryID      *string            `bson:"category_id,omitempty"`
	AccountID       *string            `bson:"account_id,omitempty"`
	Amount          *float64           `bson:"amount,omitempty"`
	Note            *string            `bson:"note,omitempty"`
	TransactionType *uint32            `bson:"transaction_type,omitempty"`
	TransactionTime *uint64            `bson:"transaction_time,omitempty"`
	CreateTime      *uint64            `bson:"create_time,omitempty"`
	UpdateTime      *uint64            `bson:"update_time,omitempty"`
}

func ToTransactionModel(t *entity.Transaction) *Transaction {
	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(t.GetTransactionID()) {
		objID, _ = primitive.ObjectIDFromHex(t.GetTransactionID())
	}

	return &Transaction{
		TransactionID:   objID,
		UserID:          t.UserID,
		CategoryID:      t.CategoryID,
		AccountID:       t.AccountID,
		Amount:          t.Amount,
		Note:            t.Note,
		TransactionType: t.TransactionType,
		TransactionTime: t.TransactionTime,
		CreateTime:      t.CreateTime,
		UpdateTime:      t.UpdateTime,
	}
}

func ToTransactionEntity(t *Transaction) *entity.Transaction {
	et := &entity.Transaction{
		TransactionID:   goutil.String(t.GetTransactionID()),
		UserID:          t.UserID,
		CategoryID:      t.CategoryID,
		AccountID:       t.AccountID,
		Note:            t.Note,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		TransactionTime: t.TransactionTime,
		CreateTime:      t.CreateTime,
		UpdateTime:      t.UpdateTime,
	}

	return et
}

func (t *Transaction) GetTransactionID() string {
	if t != nil {
		return t.TransactionID.Hex()
	}
	return ""
}

func (t *Transaction) GetUserID() string {
	if t != nil && t.UserID != nil {
		return *t.UserID
	}
	return ""
}

func (t *Transaction) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *Transaction) GetAmount() float64 {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return 0
}

func (t *Transaction) GetTransactionType() uint32 {
	if t != nil && t.TransactionType != nil {
		return *t.TransactionType
	}
	return 0
}

func (t *Transaction) GetTransactionTime() uint64 {
	if t != nil && t.TransactionTime != nil {
		return *t.TransactionTime
	}
	return 0
}

func (t *Transaction) GetCreateTime() uint64 {
	if t != nil && t.CreateTime != nil {
		return *t.CreateTime
	}
	return 0
}

func (t *Transaction) GetUpdateTime() uint64 {
	if t != nil && t.UpdateTime != nil {
		return *t.UpdateTime
	}
	return 0
}
