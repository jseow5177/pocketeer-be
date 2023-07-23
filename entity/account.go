package entity

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrSetBalanceForbidden = errors.New("set balance forbidden")
	ErrMustSetBalance      = errors.New("balance must be set")
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

var ParentAccountTypes = map[uint32]string{
	uint32(AccountTypeAsset): "asset",
	uint32(AccountTypeDebt):  "debt",
}

var ChildAccountTypes = map[uint32]string{
	uint32(AssetCash):        "cash",
	uint32(AssetBankAccount): "bank account",
	uint32(AssetInvestment):  "investment",
	uint32(DebtCreditCard):   "credit card",
	uint32(DebtLoan):         "loan",
	uint32(DebtMortgage):     "mortgage",
}

type AccountUpdate struct {
	AccountName *string
	Balance     *float64
	Note        *string
	UpdateTime  *uint64
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

func (acu *AccountUpdate) GetUpdateTime() uint64 {
	if acu != nil && acu.UpdateTime != nil {
		return *acu.UpdateTime
	}
	return 0
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

	// Investment
	TotalCost *float64
	Holdings  []*Holding
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

func WithAccountStatus(accountStatus *uint32) AccountOption {
	return func(ac *Account) {
		ac.AccountStatus = accountStatus
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

func NewAccount(userID string, opts ...AccountOption) (*Account, error) {
	now := uint64(time.Now().UnixMilli())
	ac := &Account{
		UserID:        goutil.String(userID),
		AccountName:   goutil.String(""),
		AccountType:   goutil.Uint32(uint32(AssetCash)),
		AccountStatus: goutil.Uint32(uint32(AccountStatusNormal)),
		Note:          goutil.String(""),
		CreateTime:    goutil.Uint64(now),
		UpdateTime:    goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(ac)
	}

	if err := ac.checkOpts(); err != nil {
		return nil, err
	}

	return ac, nil
}

func (ac *Account) checkOpts() error {
	if ac.CanSetBalance() && ac.Balance == nil {
		return ErrMustSetBalance
	}

	if !ac.CanSetBalance() && ac.Balance != nil {
		return ErrSetBalanceForbidden
	}

	return nil
}

func (ac *Account) Update(acu *AccountUpdate) (accountUpdate *AccountUpdate, hasUpdate bool, err error) {
	accountUpdate = new(AccountUpdate)

	if acu.AccountName != nil && acu.GetAccountName() != ac.GetAccountName() {
		hasUpdate = true
		ac.AccountName = acu.AccountName

		defer func() {
			accountUpdate.AccountName = ac.AccountName
		}()
	}

	if acu.Balance != nil && acu.GetBalance() != ac.GetBalance() {
		hasUpdate = true
		ac.Balance = acu.Balance

		defer func() {
			accountUpdate.Balance = ac.Balance
		}()
	}

	if acu.Note != nil && acu.GetNote() != ac.GetNote() {
		hasUpdate = true
		ac.Note = acu.Note

		defer func() {
			accountUpdate.Note = ac.Note
		}()
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	ac.UpdateTime = acu.UpdateTime

	if err = ac.checkOpts(); err != nil {
		return nil, false, err
	}

	accountUpdate.UpdateTime = now

	return
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

func (ac *Account) SetAccountID(accountID *string) {
	ac.AccountID = accountID
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

func (ac *Account) SetBalance(balance *float64) {
	ac.Balance = balance
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

func (ac *Account) GetTotalCost() float64 {
	if ac != nil && ac.TotalCost != nil {
		return *ac.TotalCost
	}
	return 0
}

func (ac *Account) SetTotalCost(totalCost *float64) {
	ac.TotalCost = totalCost
}

func (ac *Account) GetHoldings() []*Holding {
	if ac != nil && ac.Holdings != nil {
		return ac.Holdings
	}
	return nil
}

func (ac *Account) SetHoldings(hs []*Holding) {
	ac.Holdings = hs
}

func (ac *Account) IsAsset() bool {
	return (ac.GetAccountType() >> AccountTypeBitShift & uint32(AccountTypeAsset)) > 0
}

func (ac *Account) IsDebt() bool {
	return (ac.GetAccountType() >> AccountTypeBitShift & uint32(AccountTypeDebt)) > 0
}

func (ac *Account) IsInvestment() bool {
	return ac.GetAccountType() == uint32(AssetInvestment)
}

func (ac *Account) CanSetBalance() bool {
	return ac.GetAccountType() != uint32(AssetInvestment)
}
