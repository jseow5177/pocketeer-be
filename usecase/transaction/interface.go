package transaction

import (
	"context"

	"github.com/jseow5177/pockteer-be/config"
	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/filter"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
)

type UseCase interface {
	GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error)
	GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error)
	GetTransactionGroups(ctx context.Context, req *GetTransactionGroupsRequest) (*GetTransactionGroupsResponse, error)

	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error)
	DeleteTransaction(ctx context.Context, req *DeleteTransactionRequest) (*DeleteTransactionResponse, error)

	SumTransactions(ctx context.Context, req *SumTransactionsRequest) (*SumTransactionsResponse, error)
}

type GetTransactionRequest struct {
	UserID        *string
	TransactionID *string
}

func (m *GetTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *GetTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *GetTransactionRequest) ToCategoryFilter(categoryID string) *repo.CategoryFilter {
	return repo.NewCategoryFilter(
		m.GetUserID(),
		repo.WithCategoryID(goutil.String(categoryID)),
		repo.WithCategoryStatus(nil),
	)
}

func (m *GetTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
		repo.WithAccountStatus(nil), // deleted account can be shown
	)
}

type GetTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *GetTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type CreateTransactionRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	FromAccountID   *string
	ToAccountID     *string
	Currency        *string
	Amount          *float64
	Note            *string
	TransactionType *uint32
	TransactionTime *uint64
}

func (m *CreateTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *CreateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateTransactionRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (t *CreateTransactionRequest) GetFromAccountID() string {
	if t != nil && t.FromAccountID != nil {
		return *t.FromAccountID
	}
	return ""
}

func (t *CreateTransactionRequest) GetToAccountID() string {
	if t != nil && t.ToAccountID != nil {
		return *t.ToAccountID
	}
	return ""
}

func (m *CreateTransactionRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *CreateTransactionRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *CreateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *CreateTransactionRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *CreateTransactionRequest) GetTransactionTime() uint64 {
	if m != nil && m.TransactionTime != nil {
		return *m.TransactionTime
	}
	return 0
}

func (m *CreateTransactionRequest) ToTransactionEntity() (*entity.Transaction, error) {
	return entity.NewTransaction(
		m.GetUserID(),
		entity.WithTransactionAmount(m.Amount),
		entity.WithTransactionAccountID(m.AccountID),
		entity.WithTransactionCategoryID(m.CategoryID),
		entity.WithTransactionFromAccountID(m.FromAccountID),
		entity.WithTransactionToAccountID(m.ToAccountID),
		entity.WithTransactionNote(m.Note),
		entity.WithTransactionType(m.TransactionType),
		entity.WithTransactionTime(m.TransactionTime),
		entity.WithTransactionCurrency(m.Currency),
	)
}

func (m *CreateTransactionRequest) ToCategoryFilter() *repo.CategoryFilter {
	return repo.NewCategoryFilter(
		m.GetUserID(),
		repo.WithCategoryID(m.CategoryID),
	)
}

func (m *CreateTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
	)
}

type CreateTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *CreateTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type GetTransactionsRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	CategoryIDs     []string
	TransactionType *uint32
	TransactionTime *common.RangeFilter
	Paging          *common.Paging
}

func (m *GetTransactionsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetTransactionsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetTransactionsRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetTransactionsRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *GetTransactionsRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *GetTransactionsRequest) GetTransactionTime() *common.RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *GetTransactionsRequest) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsRequest) ToTransactionQuery() *repo.TransactionQuery {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(common.RangeFilter)
	}

	paging := m.Paging
	if m.Paging == nil {
		paging = new(common.Paging)
	}

	return &repo.TransactionQuery{
		Queries: []*repo.TransactionQuery{
			{
				Filters: []*repo.TransactionFilter{
					repo.NewTransactionFilter(
						m.GetUserID(),
						repo.WithTransactionAccountID(m.AccountID),
					),
					repo.NewTransactionFilter(
						m.GetUserID(),
						repo.WithTransactionFromAccountID(m.AccountID),
					),
					repo.NewTransactionFilter(
						m.GetUserID(),
						repo.WithTransactionToAccountID(m.AccountID),
					),
				},
				Op: filter.Or,
			},
			{
				Filters: []*repo.TransactionFilter{
					repo.NewTransactionFilter(
						m.GetUserID(),
						repo.WithTransactionCategoryID(m.CategoryID),
						repo.WithTransactionCategoryIDs(m.CategoryIDs),
						repo.WithTransactionType(m.TransactionType),
						repo.WithTransactionTimeGte(tt.Gte),
						repo.WithTransactionTimeLte(tt.Lte),
					),
				},
			},
		},
		Op: filter.And,
		Paging: &repo.Paging{
			Limit: paging.Limit,
			Page:  paging.Page,
			Sorts: []filter.Sort{
				&repo.Sort{
					Field: goutil.String("transaction_time"),
					Order: goutil.String(config.OrderDesc),
				},
				&repo.Sort{
					Field: goutil.String("create_time"),
					Order: goutil.String(config.OrderDesc),
				},
			},
		},
	}
}

func (m *GetTransactionsRequest) ToCategoryFilter(categoryIDs []string, categoryStatus uint32) *repo.CategoryFilter {
	return repo.NewCategoryFilter(
		m.GetUserID(),
		repo.WithCategoryIDs(m.CategoryIDs),
		repo.WithCategoryStatus(goutil.Uint32(categoryStatus)),
	)
}

func (m *GetTransactionsRequest) ToAccountFilter(accountIDs []string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountIDs(accountIDs),
		repo.WithAccountStatus(nil), // any status is ok
	)
}

type GetTransactionsResponse struct {
	Transactions []*entity.Transaction
	Paging       *common.Paging
}

func (m *GetTransactionsResponse) GetTransactions() []*entity.Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *GetTransactionsResponse) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

type GetTransactionGroupsRequest struct {
	AppMeta *common.AppMeta
	*GetTransactionsRequest
}

func (m *GetTransactionGroupsRequest) GetAppMeta() *common.AppMeta {
	if m != nil && m.AppMeta != nil {
		return m.AppMeta
	}
	return nil
}

type GetTransactionGroupsResponse struct {
	TransactionGroups []*common.TransactionSummary
	Paging            *common.Paging
}

func (m *GetTransactionGroupsResponse) GetTransactionGroups() []*common.TransactionSummary {
	if m != nil && m.TransactionGroups != nil {
		return m.TransactionGroups
	}
	return nil
}

func (m *GetTransactionGroupsResponse) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

type UpdateTransactionRequest struct {
	UserID          *string
	TransactionID   *string
	AccountID       *string
	CategoryID      *string
	FromAccountID   *string
	ToAccountID     *string
	TransactionType *uint32
	Note            *string
	Amount          *float64
	TransactionTime *uint64
	Currency        *string
}

func (m *UpdateTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetFromAccountID() string {
	if t != nil && t.FromAccountID != nil {
		return *t.FromAccountID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetToAccountID() string {
	if t != nil && t.ToAccountID != nil {
		return *t.ToAccountID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *UpdateTransactionRequest) GetAmount() float64 {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return 0
}

func (m *UpdateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionTime() uint64 {
	if m != nil && m.TransactionTime != nil {
		return *m.TransactionTime
	}
	return 0
}

func (m *UpdateTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *UpdateTransactionRequest) ToCategoryFilter() *repo.CategoryFilter {
	return repo.NewCategoryFilter(
		m.GetUserID(),
		repo.WithCategoryID(m.CategoryID),
	)
}

func (m *UpdateTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
	)
}

type UpdateTransactionResponse struct {
	Transaction *entity.Transaction
}

func (m *UpdateTransactionResponse) GetTransaction() *entity.Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

type SumTransactionsRequest struct {
	UserID          *string
	TransactionTime *common.RangeFilter
	TransactionType *uint32
}

func (m *SumTransactionsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *SumTransactionsRequest) GetTransactionTime() *common.RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *SumTransactionsRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *SumTransactionsRequest) ToTransactionQuery() *repo.TransactionQuery {
	tt := m.TransactionTime
	if tt == nil {
		tt = new(common.RangeFilter)
	}

	return &repo.TransactionQuery{
		Filters: []*repo.TransactionFilter{
			repo.NewTransactionFilter(
				m.GetUserID(),
				repo.WithTransactionType(m.TransactionType),
				repo.WithTransactionTimeGte(tt.Gte),
				repo.WithTransactionTimeLte(tt.Lte),
			),
		},
		Op: filter.And,
	}
}

type SumTransactionsResponse struct {
	Sums []*common.TransactionSummary
}

func (m *SumTransactionsResponse) GetSums() []*common.TransactionSummary {
	if m != nil && m.Sums != nil {
		return m.Sums
	}
	return nil
}

type DeleteTransactionRequest struct {
	UserID        *string
	TransactionID *string
}

func (m *DeleteTransactionRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *DeleteTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *DeleteTransactionRequest) ToTransactionFilter() *repo.TransactionFilter {
	return repo.NewTransactionFilter(
		m.GetUserID(),
		repo.WithTransactionID(m.TransactionID),
	)
}

func (m *DeleteTransactionRequest) ToAccountFilter(accountID string) *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(goutil.String(accountID)),
	)
}

type DeleteTransactionResponse struct{}
