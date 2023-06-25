package entity

import (
	"errors"
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrMismatchTransactionType = errors.New("mismatch transaction type")
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
	AccountID       *string
	CategoryID      *string
	Amount          *float64
	TransactionType *uint32
	TransactionTime *uint64
	Note            *string
}

func (tu *TransactionUpdate) GetAccountID() string {
	if tu != nil && tu.AccountID != nil {
		return *tu.AccountID
	}
	return ""
}

func (tu *TransactionUpdate) GetCategoryID() string {
	if tu != nil && tu.CategoryID != nil {
		return *tu.CategoryID
	}
	return ""
}

func (tu *TransactionUpdate) GetAmount() float64 {
	if tu != nil && tu.Amount != nil {
		return *tu.Amount
	}
	return 0
}

func (tu *TransactionUpdate) GetTransactionType() uint32 {
	if tu != nil && tu.TransactionType != nil {
		return *tu.TransactionType
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

type TransactionUpdateOption func(acu *TransactionUpdate)

func WithUpdateTransactionAccountID(accountID *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.AccountID = accountID
	}
}

func WithUpdateTransactionCategoryID(categoryID *string) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.CategoryID = categoryID
	}
}

func WithUpdateTransactionAmount(amount *float64) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.Amount = amount
	}
}

func WithUpdateTransactionType(transactionType *uint32) TransactionUpdateOption {
	return func(tu *TransactionUpdate) {
		tu.TransactionType = transactionType
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

func WithTransactionAccountID(accountID *string) TransactionOption {
	return func(t *Transaction) {
		t.AccountID = accountID
	}
}

func WithTransactionCategoryID(categoryID *string) TransactionOption {
	return func(t *Transaction) {
		t.CategoryID = categoryID
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
		AccountID:       goutil.String(accountID),
		CategoryID:      goutil.String(categoryID),
		TransactionType: goutil.Uint32(uint32(TransactionTypeExpense)),
		Amount:          goutil.Float64(0),
		CreateTime:      goutil.Uint64(now),
		UpdateTime:      goutil.Uint64(now),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func SetTransaction(t *Transaction, opts ...TransactionOption) {
	for _, opt := range opts {
		opt(t)
	}
}

func (t *Transaction) GetUpdates(tu *TransactionUpdate, mergeUpdate bool) (transactionUpdate *Transaction) {
	var hasUpdate bool

	if tu.CategoryID != nil && tu.GetCategoryID() != t.GetCategoryID() {
		hasUpdate = true
	}

	if tu.AccountID != nil && tu.GetAccountID() != t.GetAccountID() {
		hasUpdate = true
	}

	if tu.Amount != nil && tu.GetAmount() != t.GetAmount() {
		hasUpdate = true
	}

	if tu.TransactionType != nil && tu.GetTransactionType() != t.GetTransactionType() {
		hasUpdate = true
	}

	if tu.TransactionTime != nil && tu.GetTransactionTime() != t.GetTransactionTime() {
		hasUpdate = true
	}

	if tu.Note != nil && tu.GetNote() != t.GetNote() {
		hasUpdate = true
	}

	if hasUpdate {
		transactionUpdate = new(Transaction)
		now := uint64(time.Now().Unix())

		SetTransaction(
			transactionUpdate,
			WithTransactionCategoryID(tu.CategoryID),
			WithTransactionAccountID(tu.AccountID),
			WithTransactionAmount(tu.Amount),
			WithTransactionType(tu.TransactionType),
			WithTransactionTime(tu.TransactionTime),
			WithTransactionNote(tu.Note),
			WithTransactionUpdateTime(goutil.Uint64(now)),
		)

		if mergeUpdate {
			goutil.MergeWithPtrFields(t, transactionUpdate)
		}
		return
	}

	return
}

func (t *Transaction) CanTransactionUnderCategory(c *Category) (bool, error) {
	if t.GetTransactionType() != c.GetCategoryType() {
		return false, ErrMismatchTransactionType
	}
	return true, nil
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
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
