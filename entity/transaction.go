package entity

import (
	"errors"
	"math"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrMismatchTransactionType   = errors.New("mismatch transaction type")
	ErrInvalidTransactionAccount = errors.New("transaction not allowed under account")
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
	Amount          *float64
	TransactionTime *uint64
	Note            *string
	UpdateTime      *uint64
}

func (tu *TransactionUpdate) GetAmount() float64 {
	if tu != nil && tu.Amount != nil {
		return *tu.Amount
	}
	return 0
}

func (tu *TransactionUpdate) GetTransactionTime() uint64 {
	if tu != nil && tu.TransactionTime != nil {
		return *tu.TransactionTime
	}
	return 0
}

func (tu *TransactionUpdate) GetNote() string {
	if tu != nil && tu.Note != nil {
		return *tu.Note
	}
	return ""
}

func (tu *TransactionUpdate) GetUpdateTime() uint64 {
	if tu != nil && tu.UpdateTime != nil {
		return *tu.UpdateTime
	}
	return 0
}

type TransactionUpdateOption func(acu *TransactionUpdate)

func WithUpdateTransactionAmount(amount *float64) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.Amount = amount
	}
}

func WithUpdateTransactionTime(transactionTime *uint64) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.TransactionTime = transactionTime
	}
}

func WithUpdateTransactionNote(note *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.Note = note
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
	TransactionID   *string
	UserID          *string
	CategoryID      *string
	AccountID       *string
	Amount          *float64
	Note            *string
	TransactionType *uint32
	TransactionTime *uint64
	CreateTime      *uint64
	UpdateTime      *uint64
}

type TransactionOption = func(t *Transaction)

func WithTransactionID(transactionID *string) TransactionOption {
	return func(t *Transaction) {
		t.TransactionID = transactionID
	}
}

func WithTransactionAmount(amount *float64) TransactionOption {
	return func(t *Transaction) {
		t.Amount = amount
	}
}

func WithTransactionNote(note *string) TransactionOption {
	return func(t *Transaction) {
		t.Note = note
	}
}

func WithTransactionType(transactionType *uint32) TransactionOption {
	return func(t *Transaction) {
		t.TransactionType = transactionType
	}
}

func WithTransactionTime(transactionTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.TransactionTime = transactionTime
	}
}

func WithTransactionCreateTime(createTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.CreateTime = createTime
	}
}

func WithTransactionUpdateTime(updateTime *uint64) TransactionOption {
	return func(t *Transaction) {
		t.UpdateTime = updateTime
	}
}

func NewTransaction(userID, accountID, categoryID string, opts ...TransactionOption) *Transaction {
	now := uint64(time.Now().Unix())
	t := &Transaction{
		UserID:          goutil.String(userID),
		CategoryID:      goutil.String(categoryID),
		AccountID:       goutil.String(accountID),
		Amount:          goutil.Float64(0),
		Note:            goutil.String(""),
		TransactionType: goutil.Uint32(uint32(TransactionTypeExpense)),
		TransactionTime: goutil.Uint64(now),
		CreateTime:      goutil.Uint64(now),
		UpdateTime:      goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(t)
	}
	t.checkOpts()
	return t
}

func setTransaction(t *Transaction, opts ...TransactionOption) {
	if t == nil {
		return
	}
	for _, opt := range opts {
		opt(t)
	}
}

func (t *Transaction) checkOpts() {
	if t.GetTransactionType() == uint32(TransactionTypeExpense) {
		if t.GetAmount() > 0 {
			t.Amount = goutil.Float64(-t.GetAmount())
		}
	}

	if t.GetTransactionType() == uint32(TransactionTypeIncome) {
		t.Amount = goutil.Float64(math.Abs(t.GetAmount()))
	}
}

func (t *Transaction) Update(tu *TransactionUpdate) (transactionUpdate *TransactionUpdate, hasUpdate bool) {
	transactionUpdate = new(TransactionUpdate)

	if tu.Amount != nil && tu.GetAmount() != t.GetAmount() {
		hasUpdate = true
		setTransaction(t, WithTransactionAmount(tu.Amount))
	}

	if tu.TransactionTime != nil && tu.GetTransactionTime() != t.GetTransactionTime() {
		hasUpdate = true
		setTransaction(t, WithTransactionNote(tu.Note))
	}

	if tu.Note != nil && tu.GetNote() != t.GetNote() {
		hasUpdate = true
		setTransaction(t, WithTransactionTime(tu.TransactionTime))
	}

	if !hasUpdate {
		return
	}

	now := goutil.Uint64(uint64(time.Now().Unix()))
	setTransaction(t, WithTransactionUpdateTime(now))

	// check
	t.checkOpts()

	transactionUpdate.UpdateTime = now

	if tu.Amount != nil {
		transactionUpdate.Amount = t.Amount
	}

	if tu.Note != nil {
		transactionUpdate.Note = t.Note
	}

	if tu.TransactionTime != nil {
		transactionUpdate.TransactionTime = t.TransactionTime
	}

	return
}

func (t *Transaction) CanTransactionUnderCategory(c *Category) (bool, error) {
	if t.GetTransactionType() != c.GetCategoryType() {
		return false, ErrMismatchTransactionType
	}
	return true, nil
}

func (t *Transaction) CanTransactionUnderAccount(ac *Account) (bool, error) {
	if !ac.CanSetBalance() {
		return false, ErrInvalidTransactionAccount
	}
	return true, nil
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *Transaction) SetTransactionID(transactionID *string) {
	setTransaction(t, WithTransactionID(transactionID))
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

func (t *Transaction) GetAccountID() string {
	if t != nil && t.AccountID != nil {
		return *t.AccountID
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

func GetTransactionTypeByAmount(amount float64) TransactionType {
	if amount <= 0 {
		return TransactionTypeExpense
	}
	return TransactionTypeIncome
}
