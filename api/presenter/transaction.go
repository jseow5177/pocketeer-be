package presenter

import (
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
	"github.com/jseow5177/pockteer-be/util"
)

type Transaction struct {
	TransactionID     *string   `json:"transaction_id,omitempty"`
	CategoryID        *string   `json:"category_id,omitempty"`
	Category          *Category `json:"category,omitempty"`
	AccountID         *string   `json:"account_id,omitempty"`
	Account           *Account  `json:"account,omitempty"`
	Currency          *string   `json:"currency,omitempty"`
	Amount            *string   `json:"amount,omitempty"`
	Note              *string   `json:"note,omitempty"`
	TransactionStatus *uint32   `json:"transaction_status,omitempty"`
	TransactionType   *uint32   `json:"transaction_type,omitempty"`
	TransactionTime   *uint64   `json:"transaction_time,omitempty"`
	CreateTime        *uint64   `json:"create_time,omitempty"`
	UpdateTime        *uint64   `json:"update_time,omitempty"`
}

func (t *Transaction) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *Transaction) GetAmount() string {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return ""
}

func (t *Transaction) GetCurrency() string {
	if t != nil && t.Currency != nil {
		return *t.Currency
	}
	return ""
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

func (t *Transaction) GetTransactionStatus() uint32 {
	if t != nil && t.TransactionStatus != nil {
		return *t.TransactionStatus
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

func (t *Transaction) GetCategory() *Category {
	if t != nil && t.Category != nil {
		return t.Category
	}
	return nil
}

func (t *Transaction) GetAccount() *Account {
	if t != nil && t.Account != nil {
		return t.Account
	}
	return nil
}

type CreateTransactionRequest struct {
	CategoryID      *string `json:"category_id,omitempty"`
	AccountID       *string `json:"account_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	Currency        *string `json:"currency,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
	Note            *string `json:"note,omitempty"`
}

func (m *CreateTransactionRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *CreateTransactionRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *CreateTransactionRequest) GetAmount() string {
	if m != nil && m.Amount != nil {
		return *m.Amount
	}
	return ""
}

func (m *CreateTransactionRequest) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
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

func (m *CreateTransactionRequest) GetNote() string {
	if m != nil && m.Note != nil {
		return *m.Note
	}
	return ""
}

func (m *CreateTransactionRequest) ToUseCaseReq(userID string) *transaction.CreateTransactionRequest {
	var amount *float64
	if m.Amount != nil {
		a, _ := util.MonetaryStrToFloat(m.GetAmount())
		amount = goutil.Float64(a)
	}

	return &transaction.CreateTransactionRequest{
		UserID:          goutil.String(userID),
		CategoryID:      m.CategoryID,
		AccountID:       m.AccountID,
		Amount:          amount,
		Currency:        m.Currency,
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime,
		Note:            m.Note,
	}
}

type CreateTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *CreateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *CreateTransactionResponse) Set(useCaseRes *transaction.CreateTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.Transaction)
}

type GetTransactionRequest struct {
	TransactionID *string `json:"transaction_id,omitempty"`
}

func (m *GetTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *GetTransactionRequest) ToUseCaseReq(userID string) *transaction.GetTransactionRequest {
	return &transaction.GetTransactionRequest{
		UserID:        goutil.String(userID),
		TransactionID: m.TransactionID,
	}
}

type GetTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *GetTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *GetTransactionResponse) Set(useCaseRes *transaction.GetTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.Transaction)
}

type GetTransactionsRequest struct {
	CategoryID      *string      `json:"category_id,omitempty"`
	CategoryIDs     []string     `json:"category_ids,omitempty"`
	AccountID       *string      `json:"account_id,omitempty"`
	TransactionType *uint32      `json:"transaction_type,omitempty"`
	TransactionTime *RangeFilter `json:"transaction_time,omitempty"`
	Paging          *Paging      `json:"paging,omitempty"`
}

func (m *GetTransactionsRequest) GetCategoryID() string {
	if m != nil && m.CategoryID != nil {
		return *m.CategoryID
	}
	return ""
}

func (m *GetTransactionsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
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

func (m *GetTransactionsRequest) GetTransactionTime() *RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *GetTransactionsRequest) GetPaging() *Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsRequest) ToUseCaseReq(userID string) *transaction.GetTransactionsRequest {
	return &transaction.GetTransactionsRequest{
		UserID:          goutil.String(userID),
		CategoryID:      m.CategoryID,
		AccountID:       m.AccountID,
		CategoryIDs:     m.CategoryIDs,
		TransactionType: m.TransactionType,
		Paging:          m.Paging.toPaging(),
		TransactionTime: m.TransactionTime.toRangeFilter(),
	}
}

type GetTransactionsResponse struct {
	Transactions []*Transaction `json:"transactions,omitempty"`
	Paging       *Paging        `json:"paging,omitempty"`
}

func (m *GetTransactionsResponse) GetTransactions() []*Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *GetTransactionsResponse) GetPaging() *Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionsResponse) Set(useCaseRes *transaction.GetTransactionsResponse) {
	m.Transactions = toTransactions(useCaseRes.Transactions)
	m.Paging = toPaging(useCaseRes.Paging)
}

type UpdateTransactionRequest struct {
	TransactionID   *string `json:"transaction_id,omitempty"`
	CategoryID      *string `json:"category_id,omitempty"`
	AccountID       *string `json:"account_id,omitempty"`
	Amount          *string `json:"amount,omitempty"`
	Note            *string `json:"note,omitempty"`
	TransactionType *uint32 `json:"transaction_type,omitempty"`
	TransactionTime *uint64 `json:"transaction_time,omitempty"`
	Currency        *string `json:"currency,omitempty"`
}

func (t *UpdateTransactionRequest) GetTransactionID() string {
	if t != nil && t.TransactionID != nil {
		return *t.TransactionID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetAccountID() string {
	if t != nil && t.AccountID != nil {
		return *t.AccountID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetCategoryID() string {
	if t != nil && t.CategoryID != nil {
		return *t.CategoryID
	}
	return ""
}

func (t *UpdateTransactionRequest) GetAmount() string {
	if t != nil && t.Amount != nil {
		return *t.Amount
	}
	return ""
}

func (t *UpdateTransactionRequest) GetNote() string {
	if t != nil && t.Note != nil {
		return *t.Note
	}
	return ""
}

func (t *UpdateTransactionRequest) GetTransactionType() uint32 {
	if t != nil && t.TransactionType != nil {
		return *t.TransactionType
	}
	return 0
}

func (t *UpdateTransactionRequest) GetTransactionTime() uint64 {
	if t != nil && t.TransactionTime != nil {
		return *t.TransactionTime
	}
	return 0
}

func (t *UpdateTransactionRequest) GetCurrency() string {
	if t != nil && t.Currency != nil {
		return *t.Currency
	}
	return ""
}

func (m *UpdateTransactionRequest) ToUseCaseReq(userID string) *transaction.UpdateTransactionRequest {
	var amount *float64
	if m.Amount != nil {
		a, _ := util.MonetaryStrToFloat(m.GetAmount())
		amount = goutil.Float64(a)
	}
	return &transaction.UpdateTransactionRequest{
		UserID:          goutil.String(userID),
		AccountID:       m.AccountID,
		TransactionID:   m.TransactionID,
		Note:            m.Note,
		Amount:          amount,
		TransactionTime: m.TransactionTime,
		CategoryID:      m.CategoryID,
		TransactionType: m.TransactionType,
		Currency:        m.Currency,
	}
}

type UpdateTransactionResponse struct {
	Transaction *Transaction `json:"transaction,omitempty"`
}

func (m *UpdateTransactionResponse) GetTransaction() *Transaction {
	if m != nil && m.Transaction != nil {
		return m.Transaction
	}
	return nil
}

func (m *UpdateTransactionResponse) Set(useCaseRes *transaction.UpdateTransactionResponse) {
	m.Transaction = toTransaction(useCaseRes.Transaction)
}

type AggrTransactionsRequest struct {
	TransactionTime  *RangeFilter `json:"transaction_time,omitempty"`
	CategoryIDs      []string     `json:"category_ids,omitempty"`
	BudgetIDs        []string     `json:"budget_ids,omitempty"`
	TransactionTypes []uint32     `json:"transaction_types,omitempty"`
}

func (m *AggrTransactionsRequest) GetTransactionTime() *RangeFilter {
	if m != nil && m.TransactionTime != nil {
		return m.TransactionTime
	}
	return nil
}

func (m *AggrTransactionsRequest) GetCategoryIDs() []string {
	if m != nil && m.CategoryIDs != nil {
		return m.CategoryIDs
	}
	return nil
}

func (m *AggrTransactionsRequest) GetBudgetIDs() []string {
	if m != nil && m.BudgetIDs != nil {
		return m.BudgetIDs
	}
	return nil
}

func (m *AggrTransactionsRequest) GetTransactionTypes() []uint32 {
	if m != nil && m.TransactionTypes != nil {
		return m.TransactionTypes
	}
	return nil
}

func (m *AggrTransactionsRequest) ToUseCaseReq(userID string) *transaction.AggrTransactionsRequest {
	return &transaction.AggrTransactionsRequest{
		UserID:           goutil.String(userID),
		TransactionTypes: m.TransactionTypes,
		CategoryIDs:      m.CategoryIDs,
		TransactionTime:  m.TransactionTime.toRangeFilter(),
	}
}

type TransactionSummary struct {
	Date            *string        `json:"date,omitempty"`
	Category        *Category      `json:"category,omitempty"`
	TransactionType *uint32        `json:"transaction_type,omitempty"`
	Sum             *string        `json:"sum,omitempty"`
	TotalExpense    *string        `json:"total_expense,omitempty"`
	TotalIncome     *string        `json:"total_income,omitempty"`
	Currency        *string        `json:"currency,omitempty"`
	Transactions    []*Transaction `json:"transactions,omitempty"`
}

func (m *TransactionSummary) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

func (m *TransactionSummary) GetCategory() *Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *TransactionSummary) GetSum() string {
	if m != nil && m.Sum != nil {
		return *m.Sum
	}
	return ""
}

func (m *TransactionSummary) GetTotalExpense() string {
	if m != nil && m.TotalExpense != nil {
		return *m.TotalExpense
	}
	return ""
}

func (m *TransactionSummary) GetTotalIncome() string {
	if m != nil && m.TotalIncome != nil {
		return *m.TotalIncome
	}
	return ""
}

func (m *TransactionSummary) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *TransactionSummary) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *TransactionSummary) GetTransactions() []*Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

type Aggr struct {
	Sum *float64 `json:"sum,omitempty"`
}

type AggrTransactionsResponse struct {
	Results map[string]*Aggr `json:"results,omitempty"`
}

func (m *AggrTransactionsResponse) GetResults() map[string]*Aggr {
	if m != nil && m.Results != nil {
		return m.Results
	}
	return nil
}

func (m *AggrTransactionsResponse) Set(useCaseRes *transaction.AggrTransactionsResponse) {
	res := make(map[string]*Aggr)
	for k, v := range useCaseRes.Results {
		res[k] = toAggr(v)
	}
	m.Results = res
}

type DeleteTransactionRequest struct {
	TransactionID *string `json:"transaction_id,omitempty"`
}

func (m *DeleteTransactionRequest) GetTransactionID() string {
	if m != nil && m.TransactionID != nil {
		return *m.TransactionID
	}
	return ""
}

func (m *DeleteTransactionRequest) ToUseCaseReq(userID string) *transaction.DeleteTransactionRequest {
	return &transaction.DeleteTransactionRequest{
		UserID:        goutil.String(userID),
		TransactionID: m.TransactionID,
	}
}

type DeleteTransactionResponse struct{}

func (m *DeleteTransactionResponse) Set(useCaseRes *transaction.DeleteTransactionResponse) {}

type SumTransactionsRequest struct {
	TransactionTime *RangeFilter `json:"transaction_time,omitempty"`
	TransactionType *uint32      `json:"transaction_type,omitempty"`
}

func (m *SumTransactionsRequest) GetTransactionTime() *RangeFilter {
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

func (m *SumTransactionsRequest) ToUseCaseReq(userID string) *transaction.SumTransactionsRequest {
	return &transaction.SumTransactionsRequest{
		UserID:          goutil.String(userID),
		TransactionType: m.TransactionType,
		TransactionTime: m.TransactionTime.toRangeFilter(),
	}
}

type SumTransactionsResponse struct {
	Sums []*TransactionSummary `json:"sums,omitempty"`
}

func (m *SumTransactionsResponse) GetSums() []*TransactionSummary {
	if m != nil && m.Sums != nil {
		return m.Sums
	}
	return nil
}

func (m *SumTransactionsResponse) Set(useCaseRes *transaction.SumTransactionsResponse) {
	m.Sums = toTransactionSummaries(useCaseRes.Sums)
}

type GetTransactionGroupsRequest struct {
	AppMeta *AppMeta `json:"app_meta,omitempty"`
	*GetTransactionsRequest
}

func (m *GetTransactionGroupsRequest) GetAppMeta() *AppMeta {
	if m != nil && m.AppMeta != nil {
		return m.AppMeta
	}
	return nil
}

func (m *GetTransactionGroupsRequest) ToUseCaseReq(userID string) *transaction.GetTransactionGroupsRequest {
	return &transaction.GetTransactionGroupsRequest{
		AppMeta:                m.AppMeta.toAppMeta(),
		GetTransactionsRequest: m.GetTransactionsRequest.ToUseCaseReq(userID),
	}
}

type GetTransactionGroupsResponse struct {
	TransactionGroups []*TransactionSummary `json:"transaction_groups,omitempty"`
	Paging            *Paging               `json:"paging,omitempty"`
}

func (m *GetTransactionGroupsResponse) GetTransactionGroups() []*TransactionSummary {
	if m != nil && m.TransactionGroups != nil {
		return m.TransactionGroups
	}
	return nil
}

func (m *GetTransactionGroupsResponse) GetPaging() *Paging {
	if m != nil && m.Paging != nil {
		return m.Paging
	}
	return nil
}

func (m *GetTransactionGroupsResponse) Set(useCaseRes *transaction.GetTransactionGroupsResponse) {
	m.TransactionGroups = toTransactionSummaries(useCaseRes.TransactionGroups)
	m.Paging = toPaging(useCaseRes.Paging)
}
