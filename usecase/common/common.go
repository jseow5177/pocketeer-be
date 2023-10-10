package common

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/util"
)

type Paging struct {
	Limit *uint32
	Page  *uint32
}

func (p *Paging) GetLimit() uint32 {
	if p != nil && p.Limit != nil {
		return *p.Limit
	}
	return 0
}

func (p *Paging) GetPage() uint32 {
	if p != nil && p.Page != nil {
		return *p.Page
	}
	return 0
}

type RangeFilter struct {
	Gte *uint64
	Lte *uint64
}

func (m *RangeFilter) GetGte() uint64 {
	if m != nil && m.Gte != nil {
		return *m.Gte
	}
	return 0
}

func (m *RangeFilter) GetLte() uint64 {
	if m != nil && m.Lte != nil {
		return *m.Lte
	}
	return 0
}

type AppMeta struct {
	Timezone *string
}

func (m *AppMeta) GetTimezone() string {
	if m != nil && m.Timezone != nil {
		return *m.Timezone
	}
	return ""
}

type TransactionSummary struct {
	Date            *string
	Category        *entity.Category
	TransactionType *uint32
	Sum             *float64 // TotalExpense + TotalIncome
	TotalExpense    *float64
	TotalIncome     *float64
	Currency        *string
	Transactions    []*entity.Transaction
}

type TransactionSummaryOption func(ts *TransactionSummary)

func WithSummaryDate(date *string) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if date != nil {
			ts.SetDate(date)
		}
	}
}

func WithSummaryCategory(c *entity.Category) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if c != nil {
			ts.SetCategory(c)
		}
	}
}

func WithSummarySum(sum *float64) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if sum != nil {
			ts.SetSum(sum)
		}
	}
}

func WithSummaryTotalExpense(totalExpense *float64) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if totalExpense != nil {
			ts.SetTotalExpense(totalExpense)
		}
	}
}

func WithSummaryTotalIncome(totalIncome *float64) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if totalIncome != nil {
			ts.SetTotalIncome(totalIncome)
		}
	}
}

func WithSummaryTransactionType(transactionType *uint32) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if transactionType != nil {
			ts.SetTransactionType(transactionType)
		}
	}
}

func WithSummaryCurrency(currency *string) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if currency != nil {
			ts.SetCurrency(currency)
		}
	}
}

func WithSummaryTransactions(tss []*entity.Transaction) TransactionSummaryOption {
	return func(ts *TransactionSummary) {
		if tss != nil {
			ts.SetTransactions(tss)
		}
	}
}

func NewTransactionSummary(opts ...TransactionSummaryOption) *TransactionSummary {
	ts := new(TransactionSummary)

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

func (m *TransactionSummary) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

func (m *TransactionSummary) SetDate(date *string) {
	m.Date = date
}

func (m *TransactionSummary) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *TransactionSummary) SetCategory(c *entity.Category) {
	m.Category = c
}

func (m *TransactionSummary) GetSum() float64 {
	if m != nil && m.Sum != nil {
		return *m.Sum
	}
	return 0
}

func (m *TransactionSummary) SetSum(sum *float64) {
	m.Sum = sum

	if sum != nil {
		s := util.RoundFloatToStandardDP(*sum)
		m.Sum = goutil.Float64(s)
	}
}

func (m *TransactionSummary) GetTotalExpense() float64 {
	if m != nil && m.TotalExpense != nil {
		return *m.TotalExpense
	}
	return 0
}

func (m *TransactionSummary) SetTotalExpense(totalExpense *float64) {
	m.TotalExpense = totalExpense

	if totalExpense != nil {
		s := util.RoundFloatToStandardDP(*totalExpense)
		m.TotalExpense = goutil.Float64(s)
	}
}

func (m *TransactionSummary) GetTotalIncome() float64 {
	if m != nil && m.TotalIncome != nil {
		return *m.TotalIncome
	}
	return 0
}

func (m *TransactionSummary) SetTotalIncome(totalIncome *float64) {
	m.TotalIncome = totalIncome

	if totalIncome != nil {
		s := util.RoundFloatToStandardDP(*totalIncome)
		m.TotalIncome = goutil.Float64(s)
	}
}

func (m *TransactionSummary) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *TransactionSummary) SetTransactionType(transactionType *uint32) {
	m.TransactionType = transactionType
}

func (m *TransactionSummary) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *TransactionSummary) SetCurrency(currency *string) {
	m.Currency = currency
}

func (m *TransactionSummary) GetTransactions() []*entity.Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *TransactionSummary) SetTransactions(ts []*entity.Transaction) {
	m.Transactions = ts
}
