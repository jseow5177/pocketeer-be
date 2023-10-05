package entity

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrSetBalanceForbidden       = errors.New("set balance forbidden")
	ErrMustSetBalance            = errors.New("balance must be set")
	ErrAccountCannotHaveHoldings = errors.New("account cannot have holdings")
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

type AccountUpdateOption func(ac *Account) bool

func WithUpdateAccountName(accountName *string) AccountUpdateOption {
	return func(ac *Account) bool {
		if accountName != nil && ac.GetAccountName() != *accountName {
			ac.SetAccountName(accountName)
			return true
		}
		return false
	}
}

func WithUpdateAccountBalance(balance *float64) AccountUpdateOption {
	return func(ac *Account) bool {
		if balance != nil && ac.GetBalance() != *balance {
			ac.SetBalance(balance)
			return true
		}
		return false
	}
}

func WithUpdateAccountNote(note *string) AccountUpdateOption {
	return func(ac *Account) bool {
		if note != nil && ac.GetNote() != *note {
			ac.SetNote(note)
			return true
		}
		return false
	}
}

func WithUpdateAccountStatus(accountStatus *uint32) AccountUpdateOption {
	return func(ac *Account) bool {
		if accountStatus != nil && ac.GetAccountStatus() != *accountStatus {
			ac.SetAccountStatus(accountStatus)
			return true
		}
		return false
	}
}

type Account struct {
	UserID        *string
	AccountID     *string
	AccountName   *string
	Currency      *string
	Balance       *float64
	AccountType   *uint32
	AccountStatus *uint32
	Note          *string
	CreateTime    *uint64
	UpdateTime    *uint64

	// Investment
	Gain        *float64
	PercentGain *float64
	Holdings    []*Holding
}

type AccountOption = func(ac *Account)

func WithAccountID(accountID *string) AccountOption {
	return func(ac *Account) {
		ac.SetAccountID(accountID)
	}
}

func WithAccountName(accountName *string) AccountOption {
	return func(ac *Account) {
		ac.SetAccountName(accountName)
	}
}

func WithAccountCurrency(currency *string) AccountOption {
	return func(ac *Account) {
		ac.SetCurrency(currency)
	}
}

func WithAccountBalance(balance *float64) AccountOption {
	return func(ac *Account) {
		ac.SetBalance(balance)
	}
}

func WithAccountStatus(accountStatus *uint32) AccountOption {
	return func(ac *Account) {
		ac.SetAccountStatus(accountStatus)
	}
}

func WithAccountType(accountType *uint32) AccountOption {
	return func(ac *Account) {
		ac.SetAccountType(accountType)
	}
}

func WithAccountNote(note *string) AccountOption {
	return func(ac *Account) {
		ac.SetNote(note)
	}
}

func WithAccountCreateTime(createTime *uint64) AccountOption {
	return func(ac *Account) {
		ac.SetCreateTime(createTime)
	}
}

func WithAccountUpdateTime(updateTime *uint64) AccountOption {
	return func(ac *Account) {
		ac.SetUpdateTime(updateTime)
	}
}

func WithAccountHoldings(holdings []*Holding) AccountOption {
	return func(ac *Account) {
		ac.SetHoldings(holdings)
	}
}

func NewAccount(userID string, opts ...AccountOption) (*Account, error) {
	now := uint64(time.Now().UnixMilli())
	ac := &Account{
		UserID:        goutil.String(userID),
		AccountName:   goutil.String(""),
		AccountType:   goutil.Uint32(uint32(AssetCash)),
		Currency:      goutil.String(string(CurrencySGD)),
		AccountStatus: goutil.Uint32(uint32(AccountStatusNormal)),
		Note:          goutil.String(""),
		CreateTime:    goutil.Uint64(now),
		UpdateTime:    goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(ac)
	}

	if err := ac.validate(); err != nil {
		return nil, err
	}

	return ac, nil
}

func (ac *Account) validate() error {
	if ac.CanSetBalance() && ac.Balance == nil {
		return ErrMustSetBalance
	}

	if !ac.CanSetBalance() && ac.Balance != nil {
		return ErrSetBalanceForbidden
	}

	if !ac.IsInvestment() && len(ac.Holdings) > 0 {
		return ErrAccountCannotHaveHoldings
	}

	return nil
}

type AccountUpdate struct {
	AccountName   *string
	Balance       *float64
	Note          *string
	UpdateTime    *uint64
	AccountStatus *uint32
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

func (acu *AccountUpdate) GetAccountStatus() uint32 {
	if acu != nil && acu.AccountStatus != nil {
		return *acu.AccountStatus
	}
	return 0
}

func (acu *AccountUpdate) GetUpdateTime() uint64 {
	if acu != nil && acu.UpdateTime != nil {
		return *acu.UpdateTime
	}
	return 0
}

func (acu *AccountUpdate) GetNote() string {
	if acu != nil && acu.Note != nil {
		return *acu.Note
	}
	return ""
}

func (ac *Account) ToAccountUpdate() *AccountUpdate {
	return &AccountUpdate{
		AccountName:   ac.AccountName,
		Balance:       ac.Balance,
		Note:          ac.Note,
		AccountStatus: ac.AccountStatus,
		UpdateTime:    ac.UpdateTime,
	}
}

func (ac *Account) Update(acus ...AccountUpdateOption) (*AccountUpdate, error) {
	if len(acus) == 0 {
		return nil, nil
	}

	var hasUpdate bool
	for _, acu := range acus {
		if ok := acu(ac); ok {
			hasUpdate = true
		}
	}

	if !hasUpdate {
		return nil, nil
	}

	// check
	if err := ac.validate(); err != nil {
		return nil, err
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	ac.SetUpdateTime(now)

	return ac.ToAccountUpdate(), nil
}

func (ac *Account) GetUserID() string {
	if ac != nil && ac.UserID != nil {
		return *ac.UserID
	}
	return ""
}

func (ac *Account) SetUserID(userID *string) {
	ac.UserID = userID
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

func (ac *Account) SetAccountName(accountName *string) {
	ac.AccountName = accountName
}

func (ac *Account) GetCurrency() string {
	if ac != nil && ac.Currency != nil {
		return *ac.Currency
	}
	return ""
}

func (ac *Account) SetCurrency(currency *string) {
	ac.Currency = currency
}

func (ac *Account) GetBalance() float64 {
	if ac != nil && ac.Balance != nil {
		return *ac.Balance
	}
	return 0
}

func (ac *Account) SetBalance(balance *float64) {
	ac.Balance = balance

	if balance != nil {
		b := util.RoundFloatToStandardDP(*balance)
		ac.Balance = goutil.Float64(b)
	}
}

func (ac *Account) GetAccountStatus() uint32 {
	if ac != nil && ac.AccountStatus != nil {
		return *ac.AccountStatus
	}
	return 0
}

func (ac *Account) SetAccountStatus(accountStatus *uint32) {
	ac.AccountStatus = accountStatus
}

func (ac *Account) GetAccountType() uint32 {
	if ac != nil && ac.AccountType != nil {
		return *ac.AccountType
	}
	return 0
}

func (ac *Account) SetAccountType(accountType *uint32) {
	ac.AccountType = accountType
}

func (ac *Account) GetNote() string {
	if ac != nil && ac.Note != nil {
		return *ac.Note
	}
	return ""
}

func (ac *Account) SetNote(note *string) {
	ac.Note = note
}

func (ac *Account) GetCreateTime() uint64 {
	if ac != nil && ac.CreateTime != nil {
		return *ac.CreateTime
	}
	return 0
}

func (ac *Account) SetCreateTime(createTime *uint64) {
	ac.CreateTime = createTime
}

func (ac *Account) GetUpdateTime() uint64 {
	if ac != nil && ac.UpdateTime != nil {
		return *ac.UpdateTime
	}
	return 0
}

func (ac *Account) SetUpdateTime(updateTime *uint64) {
	ac.UpdateTime = updateTime
}

func (ac *Account) GetGain() float64 {
	if ac != nil && ac.Gain != nil {
		return *ac.Gain
	}
	return 0
}

func (ac *Account) SetGain(gain *float64) {
	ac.Gain = gain

	if gain != nil {
		g := util.RoundFloatToStandardDP(*gain)
		ac.Gain = goutil.Float64(g)
	}
}

func (ac *Account) GetPercentGain() float64 {
	if ac != nil && ac.PercentGain != nil {
		return *ac.PercentGain
	}
	return 0
}

func (ac *Account) SetPercentGain(percentGain *float64) {
	ac.PercentGain = percentGain

	if percentGain != nil {
		pg := util.RoundFloatToStandardDP(*percentGain)
		ac.PercentGain = goutil.Float64(pg)
	}
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
