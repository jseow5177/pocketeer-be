package aggr

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/account"
	"github.com/jseow5177/pockteer-be/usecase/category"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type UseCase interface {
	GetTransaction(ctx context.Context, req *GetTransactionRequest) (*GetTransactionResponse, error)
	GetTransactions(ctx context.Context, req *GetTransactionsRequest) (*GetTransactionsResponse, error)
	CreateTransaction(ctx context.Context, req *CreateTransactionRequest) (*CreateTransactionResponse, error)
	UpdateTransaction(ctx context.Context, req *UpdateTransactionRequest) (*UpdateTransactionResponse, error)

	GetBudgetWithCategories(ctx context.Context, req *GetBudgetWithCategoriesRequest) (*GetBudgetWithCategoriesResponse, error)
}

type TransactionObj struct {
	Transaction *entity.Transaction
	Category    *entity.Category
	Account     *entity.Account
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

func (m *GetTransactionRequest) ToGetTransactionRequest() *transaction.GetTransactionRequest {
	return &transaction.GetTransactionRequest{
		UserID:        m.UserID,
		TransactionID: m.TransactionID,
	}
}

func (m *GetTransactionRequest) ToGetCategoryRequest(categoryID string) *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		UserID:     m.UserID,
		CategoryID: goutil.String(categoryID),
	}
}

func (m *GetTransactionRequest) ToGetAccountRequest(accountID string) *account.GetAccountRequest {
	return &account.GetAccountRequest{
		UserID:    m.UserID,
		AccountID: goutil.String(accountID),
	}
}

type GetTransactionResponse struct {
	*TransactionObj
}

func (m *GetTransactionResponse) GetTransactionObj() *TransactionObj {
	if m != nil && m.TransactionObj != nil {
		return m.TransactionObj
	}
	return nil
}

type GetTransactionsRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	TransactionType *uint32
	TransactionTime *common.UInt64Filter
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

func (m *GetTransactionsRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *GetTransactionsRequest) GetTransactionTime() *common.UInt64Filter {
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

func (m *GetTransactionsRequest) ToGetTransactionsRequest() *transaction.GetTransactionsRequest {
	return &transaction.GetTransactionsRequest{
		UserID:          m.UserID,
		AccountID:       m.AccountID,
		CategoryID:      m.CategoryID,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
		Paging:          m.Paging,
	}
}

func (m *GetTransactionsRequest) ToGetCategoryRequest() *category.GetCategoryRequest {
	return &category.GetCategoryRequest{
		CategoryID: m.CategoryID,
	}
}

func (m *GetTransactionsRequest) ToGetAccountRequest() *account.GetAccountRequest {
	return &account.GetAccountRequest{
		AccountID: m.AccountID,
	}
}

type GetTransactionsResponse struct {
	TransactionObjs []*TransactionObj
	Paging          *common.Paging
}

func (m *GetTransactionsResponse) GetTransactionObjs() []*TransactionObj {
	if m != nil && m.TransactionObjs != nil {
		return m.TransactionObjs
	}
	return nil
}

func (m *GetTransactionsResponse) GetPaging() *common.Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

type GetBudgetWithCategoriesRequest struct {
	UserID   *string
	BudgetID *string
	Date     time.Time
}

func (m *GetBudgetWithCategoriesRequest) GetDate() time.Time {
	if m != nil {
		return m.Date
	}
	return time.Time{}
}

func (m *GetBudgetWithCategoriesRequest) GetBudgetID() string {
	if m != nil && m.BudgetID != nil {
		return *m.BudgetID
	}
	return ""
}

type GetBudgetWithCategoriesResponse struct {
	Budget     *entity.Budget
	Categories []*entity.Category
}

func (m *GetBudgetWithCategoriesResponse) GetBudget() *entity.Budget {
	if m != nil {
		return m.Budget
	}
	return nil
}

func (m *GetBudgetWithCategoriesResponse) GetCategories() []*entity.Category {
	if m != nil {
		return m.Categories
	}
	return nil
}

type CreateTransactionRequest struct {
	UserID          *string
	AccountID       *string
	CategoryID      *string
	Amount          *string
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

func (m *CreateTransactionRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
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

func (m *CreateTransactionRequest) ToCreateTransactionRequest() *transaction.CreateTransactionRequest {
	return &transaction.CreateTransactionRequest{
		UserID:          m.UserID,
		AccountID:       m.AccountID,
		CategoryID:      m.CategoryID,
		Amount:          m.Amount,
		Note:            m.Note,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
	}
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

type UpdateTransactionRequest struct {
	UserID          *string
	TransactionID   *string
	CategoryID      *string
	Note            *string
	Amount          *string
	TransactionType *uint32
	TransactionTime *uint64
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

func (m *UpdateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *UpdateTransactionRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
}

func (m *UpdateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *UpdateTransactionRequest) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *UpdateTransactionRequest) GetTransactionTime() uint64 {
	if m != nil && m.TransactionTime != nil {
		return *m.TransactionTime
	}
	return 0
}

func (m *UpdateTransactionRequest) ToUpdateTransactionRequest() *transaction.UpdateTransactionRequest {
	return &transaction.UpdateTransactionRequest{
		UserID:          m.UserID,
		TransactionID:   m.TransactionID,
		CategoryID:      m.CategoryID,
		Note:            m.Note,
		Amount:          m.Amount,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
	}
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
