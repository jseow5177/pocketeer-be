package entity

import (
	"errors"
	"math"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

var (
	ErrMismatchTransactionType   = errors.New("mismatch transaction type")
	ErrInvalidTransactionAccount = errors.New("transaction not allowed under account")
)

type TransactionStatus uint32

const (
	TransactionStatusInvalid TransactionStatus = iota
	TransactionStatusNormal
	TransactionStatusDeleted
)

type TransactionType uint32

const (
	TransactionTypeInvalid TransactionType = iota
	TransactionTypeExpense
	TransactionTypeIncome
	TransactionTypeTransfer
)

var TransactionTypes = map[uint32]string{
	uint32(TransactionTypeExpense):  "expense",
	uint32(TransactionTypeIncome):   "income",
	uint32(TransactionTypeTransfer): "transfer",
}

type TransactionUpdate struct {
	Amount            *float64
	TransactionTime   *uint64
	Note              *string
	TransactionStatus *uint32
	UpdateTime        *uint64
	AccountID         *string
	CategoryID        *string
	TransactionType   *uint32
}

func (tu *TransactionUpdate) GetAmount() float64 {
	if tu != nil && tu.Amount != nil {
		return *tu.Amount
	}
	return 0
}

func (tu *TransactionUpdate) SetAmount(amount *float64) {
	tu.Amount = amount

	if amount != nil {
		am := util.RoundFloatToStandardDP(*amount)
		tu.Amount = goutil.Float64(am)
	}
}

func (tu *TransactionUpdate) GetTransactionTime() uint64 {
	if tu != nil && tu.TransactionTime != nil {
		return *tu.TransactionTime
	}
	return 0
}

func (tu *TransactionUpdate) SetTransactionTime(transactionTime *uint64) {
	tu.TransactionTime = transactionTime
}

func (tu *TransactionUpdate) GetTransactionStatus() uint32 {
	if tu != nil && tu.TransactionStatus != nil {
		return *tu.TransactionStatus
	}
	return 0
}

func (tu *TransactionUpdate) SetTransactionStatus(transactionStatus *uint32) {
	tu.TransactionStatus = transactionStatus
}

func (tu *TransactionUpdate) GetNote() string {
	if tu != nil && tu.Note != nil {
		return *tu.Note
	}
	return ""
}

func (tu *TransactionUpdate) SetNote(note *string) {
	tu.Note = note
}

func (tu *TransactionUpdate) GetUpdateTime() uint64 {
	if tu != nil && tu.UpdateTime != nil {
		return *tu.UpdateTime
	}
	return 0
}

func (tu *TransactionUpdate) SetUpdateTime(updateTime *uint64) {
	tu.UpdateTime = updateTime
}

func (tu *TransactionUpdate) GetAccountID() string {
	if tu != nil && tu.AccountID != nil {
		return *tu.AccountID
	}
	return ""
}

func (tu *TransactionUpdate) SetAccountID(accountID *string) {
	tu.AccountID = accountID
}

func (tu *TransactionUpdate) GetTransactionType() uint32 {
	if tu != nil && tu.TransactionType != nil {
		return *tu.TransactionType
	}
	return 0
}

func (tu *TransactionUpdate) SetTransactionType(transactionType *uint32) {
	tu.TransactionType = transactionType
}

func (tu *TransactionUpdate) GetCategoryID() string {
	if tu != nil && tu.CategoryID != nil {
		return *tu.CategoryID
	}
	return ""
}

func (tu *TransactionUpdate) SetCategoryID(categoryID *string) {
	tu.CategoryID = categoryID
}

type TransactionUpdateOption func(acu *TransactionUpdate)

func WithUpdateTransactionAccountID(accountID *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetAccountID(accountID)
	}
}

func WithUpdateTransactionCategoryID(categoryID *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetCategoryID(categoryID)
	}
}

func WithUpdateTransactionAmount(amount *float64) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetAmount(amount)
	}
}

func WithUpdateTransactionTime(transactionTime *uint64) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetTransactionTime(transactionTime)
	}
}

func WithUpdateTransactionNote(note *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetNote(note)
	}
}

func WithUpdateTransactionStatus(transactionStatus *uint32) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetTransactionStatus(transactionStatus)
	}
}

func WithUpdateTransactionType(transactionType *uint32) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.SetTransactionType(transactionType)
	}
}

func NewTransactionUpdate(opts ...TransactionUpdateOption) *TransactionUpdate {
	tu := new(TransactionUpdate)
	for _, opt := range opts {
		opt(tu)
	}
	return tu
}

type Transaction struct {
	TransactionID     *string
	UserID            *string
	CategoryID        *string
	AccountID         *string
	Amount            *float64
	Note              *string
	TransactionStatus *uint32
	TransactionType   *uint32
	TransactionTime   *uint64
	CreateTime        *uint64
	UpdateTime        *uint64

	Category *Category
	Account  *Account
}

type TransactionOption = func(t *Transaction)

func WithTransactionID(transactionID *string) TransactionOption {
	return func(t *Transaction) {
		t.SetTransactionID(transactionID)
	}
}

func WithTransactionAmount(amount *float64) TransactionOption {
	return func(t *Transaction) {
		t.SetAmount(amount)
	}
}

func WithTransactionNote(note *string) TransactionOption {
	return func(t *Transaction) {
		t.SetNote(note)
	}
}

func WithTransactionStatus(transactionStatus *uint32) TransactionOption {
	return func(t *Transaction) {
		t.SetTransactionStatus(transactionStatus)
	}
}

func WithTransactionType(transactionType *uint32) TransactionOption {
	return func(t *Transaction) {
		t.SetTransactionType(transactionType)
	}
}

func WithTransactionTime(transactionTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.SetTransactionTime(transactionTime)
	}
}

func WithTransactionCreateTime(createTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.SetCreateTime(createTime)
	}
}

func WithTransactionUpdateTime(updateTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.SetUpdateTime(updateTime)
	}
}

func NewTransaction(userID, accountID, categoryID string, opts ...TransactionOption) *Transaction {
	now := uint64(time.Now().UnixMilli())
	t := &Transaction{
		UserID:            goutil.String(userID),
		CategoryID:        goutil.String(categoryID),
		AccountID:         goutil.String(accountID),
		Amount:            goutil.Float64(0),
		Note:              goutil.String(""),
		TransactionStatus: goutil.Uint32(uint32(TransactionStatusNormal)),
		TransactionType:   goutil.Uint32(uint32(TransactionTypeExpense)),
		TransactionTime:   goutil.Uint64(now),
		CreateTime:        goutil.Uint64(now),
		UpdateTime:        goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(t)
	}
	t.checkOpts()
	return t
}

func (t *Transaction) checkOpts() {
	if t.IsExpense() {
		if t.GetAmount() > 0 {
			t.SetAmount(goutil.Float64(-t.GetAmount()))
		}
	}

	if t.IsIncome() {
		t.SetAmount(goutil.Float64(math.Abs(t.GetAmount())))
	}
}

func (t *Transaction) Update(tu *TransactionUpdate) *TransactionUpdate {
	var (
		hasUpdate         bool
		transactionUpdate = new(TransactionUpdate)
	)

	if tu.TransactionType != nil && tu.GetTransactionType() != t.GetTransactionType() {
		hasUpdate = true
		t.SetTransactionType(tu.TransactionType)

		defer func() {
			// a change in transaction type may need to reset amount
			transactionUpdate.SetAmount(t.Amount)
			transactionUpdate.SetTransactionType(t.TransactionType)
		}()
	}

	if tu.CategoryID != nil && tu.GetCategoryID() != t.GetCategoryID() {
		hasUpdate = true
		t.SetCategoryID(tu.CategoryID)

		defer func() {
			transactionUpdate.SetCategoryID(t.CategoryID)
		}()
	}

	if tu.Amount != nil && tu.GetAmount() != t.GetAmount() {
		hasUpdate = true
		t.SetAmount(tu.Amount)

		defer func() {
			transactionUpdate.SetAmount(t.Amount)
		}()
	}

	if tu.TransactionTime != nil && tu.GetTransactionTime() != t.GetTransactionTime() {
		hasUpdate = true
		t.SetTransactionTime(tu.TransactionTime)

		defer func() {
			transactionUpdate.SetTransactionTime(t.TransactionTime)
		}()
	}

	if tu.Note != nil && tu.GetNote() != t.GetNote() {
		hasUpdate = true
		t.SetNote(tu.Note)

		defer func() {
			transactionUpdate.SetNote(t.Note)
		}()
	}

	if tu.TransactionStatus != nil && tu.GetTransactionStatus() != t.GetTransactionStatus() {
		hasUpdate = true
		t.SetTransactionStatus(tu.TransactionStatus)

		defer func() {
			transactionUpdate.SetTransactionStatus(t.TransactionStatus)
		}()
	}

	if tu.AccountID != nil && tu.GetAccountID() != t.GetAccountID() {
		hasUpdate = true
		t.SetAccountID(tu.AccountID)

		defer func() {
			transactionUpdate.SetAccountID(t.AccountID)
		}()
	}

	if !hasUpdate {
		return nil
	}

	now := goutil.Uint64(uint64(time.Now().UnixMilli()))
	t.SetUpdateTime(now)

	// check
	t.checkOpts()

	transactionUpdate.SetUpdateTime(now)

	return transactionUpdate
}

func (t *Transaction) CanTransactionUnderCategory(c *Category) error {
	if t.GetTransactionType() != c.GetCategoryType() {
		return ErrMismatchTransactionType
	}
	return nil
}

func (t *Transaction) CanTransactionUnderAccount(ac *Account) error {
	if !ac.CanSetBalance() {
		return ErrInvalidTransactionAccount
	}
	return nil
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *Transaction) SetTransactionID(transactionID *string) {
	t.TransactionID = transactionID
}

func (t *Transaction) GetUserID() string {
	if t != nil && t.UserID != nil {
		return *t.UserID
	}
	return ""
}

func (t *Transaction) SetUserID(userID *string) {
	t.UserID = userID
}

func (t *Transaction) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *Transaction) SetCategoryID(categoryID *string) {
	t.CategoryID = categoryID
}

func (t *Transaction) GetAccountID() string {
	if t != nil && t.AccountID != nil {
		return *t.AccountID
	}
	return ""
}

func (t *Transaction) SetAccountID(accountID *string) {
	t.AccountID = accountID
}

func (t *Transaction) GetAmount() float64 {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return 0
}

func (t *Transaction) SetAmount(amount *float64) {
	t.Amount = amount

	if amount != nil {
		am := util.RoundFloatToStandardDP(*amount)
		t.Amount = goutil.Float64(am)
	}
}

func (t *Transaction) GetNote() string {
	if t != nil && t.Note != nil {
		return *t.Note
	}
	return ""
}

func (t *Transaction) SetNote(note *string) {
	t.Note = note
}

func (t *Transaction) GetTransactionStatus() uint32 {
	if t != nil && t.TransactionStatus != nil {
		return *t.TransactionStatus
	}
	return 0
}

func (t *Transaction) SetTransactionStatus(transactionStatus *uint32) {
	t.TransactionStatus = transactionStatus
}

func (t *Transaction) GetTransactionType() uint32 {
	if t != nil && t.TransactionType != nil {
		return *t.TransactionType
	}
	return 0
}

func (t *Transaction) SetTransactionType(transactionType *uint32) {
	t.TransactionType = transactionType
}

func (t *Transaction) GetTransactionTime() uint64 {
	if t != nil && t.TransactionTime != nil {
		return *t.TransactionTime
	}
	return 0
}

func (t *Transaction) SetTransactionTime(transactionTime *uint64) {
	t.TransactionTime = transactionTime
}

func (t *Transaction) GetCreateTime() uint64 {
	if t != nil && t.CreateTime != nil {
		return *t.CreateTime
	}
	return 0
}

func (t *Transaction) SetCreateTime(createTime *uint64) {
	t.CreateTime = createTime
}

func (t *Transaction) GetUpdateTime() uint64 {
	if t != nil && t.UpdateTime != nil {
		return *t.UpdateTime
	}
	return 0
}

func (t *Transaction) SetUpdateTime(updateTime *uint64) {
	t.UpdateTime = updateTime
}

func (t *Transaction) GetCategory() *Category {
	if t != nil && t.Category != nil {
		return t.Category
	}
	return nil
}

func (t *Transaction) SetCategory(c *Category) {
	t.Category = c
}

func (t *Transaction) GetAccount() *Account {
	if t != nil && t.Account != nil {
		return t.Account
	}
	return nil
}

func (t *Transaction) SetAccount(ac *Account) {
	t.Account = ac
}

func (t *Transaction) IsExpense() bool {
	return t.GetTransactionType() == uint32(TransactionTypeExpense)
}

func (t *Transaction) IsIncome() bool {
	return t.GetTransactionType() == uint32(TransactionTypeIncome)
}

func GetTransactionTypeByAmount(amount float64) TransactionType {
	if amount <= 0 {
		return TransactionTypeExpense
	}
	return TransactionTypeIncome
}
