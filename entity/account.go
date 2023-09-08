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

func (acu *AccountUpdate) SetAccountName(accountName *string) {
	acu.AccountName = accountName
}

func (acu *AccountUpdate) GetBalance() float64 {
	if acu != nil && acu.Balance != nil {
		return *acu.Balance
	}
	return 0
}

func (acu *AccountUpdate) SetBalance(balance *float64) {
	acu.Balance = balance

	if balance != nil {
		b := util.RoundFloatToStandardDP(*balance)
		acu.Balance = goutil.Float64(b)
	}
}

func (acu *AccountUpdate) GetNote() string {
	if acu != nil && acu.Note != nil {
		return *acu.Note
	}
	return ""
}

func (acu *AccountUpdate) SetNote(note *string) {
	acu.Note = note
}

func (acu *AccountUpdate) GetUpdateTime() uint64 {
	if acu != nil && acu.UpdateTime != nil {
		return *acu.UpdateTime
	}
	return 0
}

func (acu *AccountUpdate) SetUpdateTime(updateTime *uint64) {
	acu.UpdateTime = updateTime
}

func (acu *AccountUpdate) GetAccountStatus() uint32 {
	if acu != nil && acu.AccountStatus != nil {
		return *acu.AccountStatus
	}
	return 0
}

func (acu *AccountUpdate) SetAccountStatus(accountStatus *uint32) {
	acu.AccountStatus = accountStatus
}

type AccountUpdateOption func(acu *AccountUpdate)

func WithUpdateAccountName(accountName *string) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.SetAccountName(accountName)
	}
}

func WithUpdateAccountBalance(balance *float64) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.SetBalance(balance)
	}
}

func WithUpdateAccountNote(note *string) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.SetNote(note)
	}
}

func WithUpdateAccountStatus(accountStatus *uint32) AccountUpdateOption {
	return func(acu *AccountUpdate) {
		acu.SetAccountStatus(accountStatus)
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
	Currency      *string
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

	if !ac.IsInvestment() && len(ac.Holdings) > 0 {
		return ErrAccountCannotHaveHoldings
	}

	if err := CheckCurrency(ac.GetCurrency()); err != nil {
		return err
	}

	return nil
}

func (ac *Account) Update(acu *AccountUpdate) (*AccountUpdate, error) {
	var (
		hasUpdate     bool
		accountUpdate = new(AccountUpdate)
	)

	if acu.AccountName != nil && acu.GetAccountName() != ac.GetAccountName() {
		hasUpdate = true
		ac.SetAccountName(acu.AccountName)

		defer func() {
			accountUpdate.SetAccountName(ac.AccountName)
		}()
	}

	if acu.Balance != nil && acu.GetBalance() != ac.GetBalance() {
		hasUpdate = true
		ac.SetBalance(acu.Balance)

		defer func() {
			accountUpdate.SetBalance(ac.Balance)
		}()
	}

	if acu.Note != nil && acu.GetNote() != ac.GetNote() {
		hasUpdate = true
		ac.SetNote(acu.Note)

		defer func() {
			accountUpdate.SetNote(ac.Note)
		}()
	}

	if !hasUpdate {
		return nil, nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	ac.SetUpdateTime(now)

	if err := ac.checkOpts(); err != nil {
		return nil, err
	}

	accountUpdate.SetUpdateTime(now)

	return accountUpdate, nil
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

func (ac *Account) GetTotalCost() float64 {
	if ac != nil && ac.TotalCost != nil {
		return *ac.TotalCost
	}
	return 0
}

func (ac *Account) SetTotalCost(totalCost *float64) {
	ac.TotalCost = totalCost

	if totalCost != nil {
		tc := util.RoundFloatToStandardDP(*totalCost)
		ac.TotalCost = goutil.Float64(tc)
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

func (ac *Account) ComputeCostAndBalance() {
	if !ac.IsInvestment() {
		return
	}

	var (
		totalCost    float64
		totalBalance float64
	)
	for _, h := range ac.Holdings {
		totalCost += h.GetTotalCost()
		totalBalance += h.GetLatestValue()
	}

	ac.SetBalance(goutil.Float64(totalBalance))
	ac.SetTotalCost(goutil.Float64(totalCost))
}
