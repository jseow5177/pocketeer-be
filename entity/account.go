package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
	uint32(DebtCreditCard):   "credit card",
	uint32(DebtLoan):         "loan",
	uint32(DebtMortgage):     "mortgage",
}

type AccountUpdate struct {
	AccountName *string
	Balance     *float64
	Note        *string
}

func (acu *AccountUpdate) GetAccountName() string {
	if acu != nil && acu.AccountName != nil {
		return *acu.AccountName
	}
	return ""
}

func (acu *AccountUpdate) GetBalance() float64 {
	if acu != nil && acu.Balance != nil {
		return *acu.Balance
	}
	return 0
}

func (acu *AccountUpdate) GetNote() string {
	if acu != nil && acu.Note != nil {
		return *acu.Note
	}
	return ""
}

type AccountUpdateOption func(acu *AccountUpdate)

func WithUpdateAccountName(accountName *string) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.AccountName = accountName
	}
}

func WithUpdateAccountBalance(balance *float64) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.Balance = balance
	}
}

func WithUpdateAccountNote(note *string) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.Note = note
	}
}

func NewAccountUpdate(opts ...AccountUpdateOption) *AccountUpdate {
	au := new(AccountUpdate)
	for _, opt := range opts {
		opt(au)
	}
	return au
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

type AccountOption = func(ac *Account)

func WithAccountID(accountID *string) AccountOption {
	return func(ac *Account) {
		ac.AccountID = accountID
	}
}

func WithAccountName(accountName *string) AccountOption {
	return func(ac *Account) {
		ac.AccountName = accountName
	}
}

func WithAccountBalance(balance *float64) AccountOption {
	return func(ac *Account) {
		ac.Balance = balance
	}
}

func WithAccountType(accountType *uint32) AccountOption {
	return func(ac *Account) {
		ac.AccountType = accountType
	}
}

func WithAccountNote(note *string) AccountOption {
	return func(ac *Account) {
		ac.Note = note
	}
}

func WithAccountCreateTime(createTime *uint64) AccountOption {
	return func(ac *Account) {
		ac.CreateTime = createTime
	}
}

func WithAccountUpdateTime(updateTime *uint64) AccountOption {
	return func(ac *Account) {
		ac.UpdateTime = updateTime
	}
}

func NewAccount(userID string, opts ...AccountOption) *Account {
	now := uint64(time.Now().Unix())
	ac := &Account{
		UserID:        goutil.String(userID),
		AccountType:   goutil.Uint32(uint32(AssetCash)),
		Balance:       goutil.Float64(0),
		AccountStatus: goutil.Uint32(uint32(AccountStatusNormal)),
		CreateTime:    goutil.Uint64(now),
		UpdateTime:    goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(ac)
	}
	return ac
}

func SetAccount(ac *Account, opts ...AccountOption) {
	for _, opt := range opts {
		opt(ac)
	}
}

func (ac *Account) GetUpdates(acu *AccountUpdate, mergeUpdate bool) (accountUpdate *Account) {
	var hasUpdate bool

	if acu.AccountName != nil && acu.GetAccountName() != ac.GetAccountName() {
		hasUpdate = true
	}

	if acu.Balance != nil && acu.GetBalance() != ac.GetBalance() {
		hasUpdate = true
	}

	if acu.Note != nil && acu.GetNote() != ac.GetNote() {
		hasUpdate = true
	}

	if hasUpdate {
		accountUpdate = new(Account)
		now := uint64(time.Now().Unix())

		SetAccount(
			accountUpdate,
			WithAccountName(acu.AccountName),
			WithAccountBalance(acu.Balance),
			WithAccountNote(acu.Note),
			WithAccountUpdateTime(goutil.Uint64(now)),
		)

		if mergeUpdate {
			goutil.MergeWithPtrFields(ac, accountUpdate)
		}
		return
	}

	return
}

func (ac *Account) AddBalance(amount float64) (accountUpdate *Account) {
	newBalance := ac.GetBalance() + amount
	return ac.GetUpdates(NewAccountUpdate(
		WithUpdateAccountBalance(goutil.Float64(newBalance)),
	), true)
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
