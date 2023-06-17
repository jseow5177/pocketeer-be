package entity

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type AccountStatus uint32

const (
	AccountStatusInvalid AccountStatus = iota
	AccountStatusNormal
	AccountStatusDeleted
)

const AccountTypeBitShift = 4

type AccountType uint32

const (
	AccountTypeInvalid AccountType = iota
	AccountTypeAsset
	AccountTypeDebt
)

// Allow 2^4 unique assets
const (
	AssetCash AccountType = AccountTypeAsset<<AccountTypeBitShift | iota
	AssetBankAccount
	AssetInvestment
)

// Allow 2^4 unique debt
const (
	DebtCreditCard AccountType = AccountTypeDebt<<AccountTypeBitShift | iota
	DebtLoan
	DebtMortgage
)

var AccountTypes = map[uint32]string{
	uint32(AssetCash):        "cash",
	uint32(AssetBankAccount): "bank account",
	uint32(AssetInvestment):  "investment",
	uint32(DebtCreditCard):   "credit card",
	uint32(DebtLoan):         "loan",
	uint32(DebtMortgage):     "mortgage",
}

type Account struct {
	UserID        *string
	AccountID     *string
	AccountName   *string
	Balance       *float64
	AccountType   *uint32
	AccountStatus *uint32
	Note          *string
	CreateTime    *uint64
	UpdateTime    *uint64
}

func (ac *Account) GetUserID() string {
	if ac != nil && ac.UserID != nil {
		return *ac.UserID
	}
	return ""
}

func (ac *Account) GetAccountID() string {
	if ac != nil && ac.AccountID != nil {
		return *ac.AccountID
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

func (ac *Account) IsAccountTypeAsset() bool {
	return (ac.GetAccountType() >> AccountTypeBitShift & uint32(AccountTypeAsset)) > 0
}

func (ac *Account) IsAccountTypeDebt() bool {
	return (ac.GetAccountType() >> AccountTypeBitShift & uint32(AccountTypeDebt)) > 0
}

func (ac *Account) SetBalance(balance string) {
	ac.Balance = goutil.Float64(util.MonetaryStrToFloat(balance))
}
