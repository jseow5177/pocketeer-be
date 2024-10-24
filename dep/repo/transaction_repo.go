package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrTransactionNotFound = errutil.NotFoundError(errors.New("transaction not found"))
)

type TransactionRepo interface {
	Get(ctx context.Context, tf *TransactionFilter) (*entity.Transaction, error)
	GetMany(ctx context.Context, tq *TransactionQuery) ([]*entity.Transaction, error)

	Create(ctx context.Context, t *entity.Transaction) (string, error)
	Update(ctx context.Context, tf *TransactionFilter, t *entity.TransactionUpdate) error
}

type TransactionQuery struct {
	Filters []*TransactionFilter
	Queries []*TransactionQuery
	Op      filter.BoolOp
	Paging  *Paging
}

func (q *TransactionQuery) GetQueries() []filter.Query {
	qs := make([]filter.Query, 0)
	for _, bq := range q.Queries {
		qs = append(qs, bq)
	}
	return qs
}

func (q *TransactionQuery) GetFilters() []interface{} {
	ibfs := make([]interface{}, 0)
	for _, bf := range q.Filters {
		ibfs = append(ibfs, bf)
	}
	return ibfs
}

func (q *TransactionQuery) GetOp() filter.BoolOp {
	return q.Op
}

type TransactionFilter struct {
	UserID             *string  `filter:"user_id"`
	TransactionID      *string  `filter:"_id"`
	AccountID          *string  `filter:"account_id"`
	FromAccountID      *string  `filter:"from_account_id"`
	ToAccountID        *string  `filter:"to_account_id"`
	CategoryID         *string  `filter:"category_id"`
	CategoryIDs        []string `filter:"category_id__in"`
	TransactionStatus  *uint32  `filter:"transaction_status"`
	TransactionType    *uint32  `filter:"transaction_type"`
	TransactionTypes   []uint32 `filter:"transaction_type__in"`
	TransactionTimeGte *uint64  `filter:"transaction_time__gte"`
	TransactionTimeLte *uint64  `filter:"transaction_time__lte"`
}

type TransactionFilterOption = func(tf *TransactionFilter)

func WithTransactionID(transactionID *string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionID = transactionID
	}
}

func WithTransactionAccountID(accountID *string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.AccountID = accountID
	}
}

func WithTransactionFromAccountID(fromAccountID *string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.FromAccountID = fromAccountID
	}
}

func WithTransactionToAccountID(toAccountID *string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.ToAccountID = toAccountID
	}
}

func WithTransactionCategoryID(categoryID *string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.CategoryID = categoryID
	}
}

func WithTransactionCategoryIDs(categoryIDs []string) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.CategoryIDs = categoryIDs
	}
}

func WithTransactionStatus(transactionStatus *uint32) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionStatus = transactionStatus
	}
}

func WithTransactionType(transactionType *uint32) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionType = transactionType
	}
}

func WithTransactionTypes(transactionTypes []uint32) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionTypes = transactionTypes
	}
}

func WithTransactionTimeGte(transactionTimeGte *uint64) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionTimeGte = transactionTimeGte
	}
}

func WithTransactionTimeLte(transactionTimeLte *uint64) TransactionFilterOption {
	return func(tf *TransactionFilter) {
		tf.TransactionTimeLte = transactionTimeLte
	}
}

func NewTransactionFilter(userID string, opts ...TransactionFilterOption) *TransactionFilter {
	tf := &TransactionFilter{
		UserID:            goutil.String(userID),
		TransactionStatus: goutil.Uint32(uint32(entity.TransactionStatusNormal)),
	}
	for _, opt := range opts {
		opt(tf)
	}
	return tf
}

func (f *TransactionFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *TransactionFilter) GetTransactionID() string {
	if f != nil && f.TransactionID != nil {
		return *f.TransactionID
	}
	return ""
}

func (f *TransactionFilter) GetAccountID() string {
	if f != nil && f.AccountID != nil {
		return *f.AccountID
	}
	return ""
}

func (f *TransactionFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *TransactionFilter) GetCategoryIDs() []string {
	if f != nil && f.CategoryIDs != nil {
		return f.CategoryIDs
	}
	return nil
}

func (f *TransactionFilter) GetTransactionStatus() uint32 {
	if f != nil && f.TransactionStatus != nil {
		return *f.TransactionStatus
	}
	return 0
}

func (f *TransactionFilter) GetTransactionType() uint32 {
	if f != nil && f.TransactionType != nil {
		return *f.TransactionType
	}
	return 0
}

func (f *TransactionFilter) GetTransactionTypes() []uint32 {
	if f != nil && f.TransactionTypes != nil {
		return f.TransactionTypes
	}
	return nil
}

func (f *TransactionFilter) GetTransactionTimeGte() uint64 {
	if f != nil && f.TransactionTimeGte != nil {
		return *f.TransactionTimeGte
	}
	return 0
}

func (f *TransactionFilter) GetTransactionTimeLte() uint64 {
	if f != nil && f.TransactionTimeLte != nil {
		return *f.TransactionTimeLte
	}
	return 0
}
