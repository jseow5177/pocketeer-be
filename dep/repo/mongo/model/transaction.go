package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	TransactionID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID            *string            `bson:"user_id,omitempty"`
	CategoryID        *string            `bson:"category_id,omitempty"`
	AccountID         *string            `bson:"account_id,omitempty"`
	Currency          *string            `bson:"currency,omitempty"`
	Amount            *float64           `bson:"amount,omitempty"`
	Note              *string            `bson:"note,omitempty"`
	TransactionStatus *uint32            `bson:"transaction_status,omitempty"`
	TransactionType   *uint32            `bson:"transaction_type,omitempty"`
	TransactionTime   *uint64            `bson:"transaction_time,omitempty"`
	CreateTime        *uint64            `bson:"create_time,omitempty"`
	UpdateTime        *uint64            `bson:"update_time,omitempty"`
}

func ToTransactionModelFromEntity(t *entity.Transaction) *Transaction {
	if t == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(t.GetTransactionID()) {
		objID, _ = primitive.ObjectIDFromHex(t.GetTransactionID())
	}

	return &Transaction{
		TransactionID:     objID,
		UserID:            t.UserID,
		CategoryID:        t.CategoryID,
		AccountID:         t.AccountID,
		Currency:          t.Currency,
		Amount:            t.Amount,
		Note:              t.Note,
		TransactionStatus: t.TransactionStatus,
		TransactionType:   t.TransactionType,
		TransactionTime:   t.TransactionTime,
		CreateTime:        t.CreateTime,
		UpdateTime:        t.UpdateTime,
	}
}

func ToTransactionModelFromUpdate(tu *entity.TransactionUpdate) *Transaction {
	if tu == nil {
		return nil
	}

	return &Transaction{
		Amount:            tu.Amount,
		Note:              tu.Note,
		TransactionTime:   tu.TransactionTime,
		TransactionStatus: tu.TransactionStatus,
		UpdateTime:        tu.UpdateTime,
		AccountID:         tu.AccountID,
		CategoryID:        tu.CategoryID,
		TransactionType:   tu.TransactionType,
		Currency:          tu.Currency,
	}
}

func ToTransactionEntity(t *Transaction) (*entity.Transaction, error) {
	if t == nil {
		return nil, nil
	}

	return entity.NewTransaction(
		t.GetUserID(),
		t.GetAccountID(),
		t.GetCategoryID(),
		entity.WithTransactionID(goutil.String(t.GetTransactionID())),
		entity.WithTransactionAmount(t.Amount),
		entity.WithTransactionNote(t.Note),
		entity.WithTransactionType(t.TransactionType),
		entity.WithTransactionTime(t.TransactionTime),
		entity.WithTransactionCreateTime(t.CreateTime),
		entity.WithTransactionUpdateTime(t.UpdateTime),
		entity.WithTransactionStatus(t.TransactionStatus),
		entity.WithTransactionCurrency(t.Currency),
	)
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

func (t *Transaction) GetAccountID() string {
	if t != nil && t.AccountID != nil {
		return *t.AccountID
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

func (t *Transaction) GetNote() string {
	if t != nil && t.Note != nil {
		return *t.Note
	}
	return ""
}

func (t *Transaction) GetTransactionStatus() uint32 {
	if t != nil && t.TransactionStatus != nil {
		return *t.TransactionStatus
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
