package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Account struct {
	AccountID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID        *string            `bson:"user_id,omitempty"`
	AccountName   *string            `bson:"account_name,omitempty"`
	Balance       *float64           `bson:"balance,omitempty"`
	AccountType   *uint32            `bson:"account_type,omitempty"`
	AccountStatus *uint32            `bson:"account_status,omitempty"`
	Note          *string            `bson:"note,omitempty"`
	CreateTime    *uint64            `bson:"create_time,omitempty"`
	UpdateTime    *uint64            `bson:"update_time,omitempty"`
}

func ToAccountModelFromEntity(ac *entity.Account) *Account {
	if ac == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(ac.GetAccountID()) {
		objID, _ = primitive.ObjectIDFromHex(ac.GetAccountID())
	}

	return &Account{
		AccountID:     objID,
		UserID:        ac.UserID,
		AccountName:   ac.AccountName,
		Balance:       ac.Balance,
		Note:          ac.Note,
		AccountType:   ac.AccountType,
		AccountStatus: ac.AccountStatus,
		CreateTime:    ac.CreateTime,
		UpdateTime:    ac.UpdateTime,
	}
}

func ToAccountModelFromUpdate(acu *entity.AccountUpdate) *Account {
	if acu == nil {
		return nil
	}

	return &Account{
		AccountName: acu.AccountName,
		Balance:     acu.Balance,
		Note:        acu.Note,
		UpdateTime:  acu.UpdateTime,
	}
}

func ToAccountEntity(ac *Account) (*entity.Account, error) {
	if ac == nil {
		return nil, nil
	}

	return entity.NewAccount(
		ac.GetUserID(),
		entity.WithAccountID(goutil.String(ac.GetAccountID())),
		entity.WithAccountName(ac.AccountName),
		entity.WithAccountBalance(ac.Balance),
		entity.WithAccountStatus(ac.AccountStatus),
		entity.WithAccountType(ac.AccountType),
		entity.WithAccountNote(ac.Note),
		entity.WithAccountCreateTime(ac.CreateTime),
		entity.WithAccountUpdateTime(ac.UpdateTime),
	)
}

func (ac *Account) GetUserID() string {
	if ac != nil && ac.UserID != nil {
		return *ac.UserID
	}
	return ""
}

func (ac *Account) GetAccountID() string {
	if ac != nil {
		return ac.AccountID.Hex()
	}
	return ""
}

func (ac *Account) GetAccountName() string {
	if ac != nil && ac.AccountName != nil {
		return *ac.AccountName
	}
	return ""
}

func (ac *Account) GetBalance() float64 {
	if ac != nil && ac.Balance != nil {
		return *ac.Balance
	}
	return 0
}

func (ac *Account) GetAccountStatus() uint32 {
	if ac != nil && ac.AccountStatus != nil {
		return *ac.AccountStatus
	}
	return 0
}

func (ac *Account) GetAccountType() uint32 {
	if ac != nil && ac.AccountType != nil {
		return *ac.AccountType
	}
	return 0
}

func (ac *Account) GetNote() string {
	if ac != nil && ac.Note != nil {
		return *ac.Note
	}
	return ""
}

func (ac *Account) GetCreateTime() uint64 {
	if ac != nil && ac.CreateTime != nil {
		return *ac.CreateTime
	}
	return 0
}

func (ac *Account) GetUpdateTime() uint64 {
	if ac != nil && ac.UpdateTime != nil {
		return *ac.UpdateTime
	}
	return 0
}
