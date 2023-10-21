package common

import (
	"sort"

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

type Summary struct {
	Date            *string
	Category        *entity.Category
	TransactionType *uint32
	Sum             *float64
	TotalExpense    *float64
	TotalIncome     *float64
	Currency        *string
	Transactions    []*entity.Transaction
}

type SummaryOption func(ts *Summary)

func WithSummaryDate(date *string) SummaryOption {
	return func(ts *Summary) {
		if date != nil {
			ts.SetDate(date)
		}
	}
}

func WithSummaryCategory(c *entity.Category) SummaryOption {
	return func(ts *Summary) {
		if c != nil {
			ts.SetCategory(c)
		}
	}
}

func WithSummarySum(sum *float64) SummaryOption {
	return func(ts *Summary) {
		if sum != nil {
			ts.SetSum(sum)
		}
	}
}

func WithSummaryTotalExpense(totalExpense *float64) SummaryOption {
	return func(ts *Summary) {
		if totalExpense != nil {
			ts.SetTotalExpense(totalExpense)
		}
	}
}

func WithSummaryTotalIncome(totalIncome *float64) SummaryOption {
	return func(ts *Summary) {
		if totalIncome != nil {
			ts.SetTotalIncome(totalIncome)
		}
	}
}

func WithSummaryTransactionType(transactionType *uint32) SummaryOption {
	return func(ts *Summary) {
		if transactionType != nil {
			ts.SetTransactionType(transactionType)
		}
	}
}

func WithSummaryCurrency(currency *string) SummaryOption {
	return func(ts *Summary) {
		if currency != nil {
			ts.SetCurrency(currency)
		}
	}
}

func WithSummaryTransactions(tss []*entity.Transaction) SummaryOption {
	return func(ts *Summary) {
		if tss != nil {
			ts.SetTransactions(tss)
		}
	}
}

func AscSortSummaryByDate(summary ...*Summary) {
	sort.Slice(summary, func(i, j int) bool {
		return summary[i].GetDate() < summary[j].GetDate()
	})
}

func NewSummary(opts ...SummaryOption) *Summary {
	ts := new(Summary)

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

func (m *Summary) GetDate() string {
	if m != nil && m.Date != nil {
		return *m.Date
	}
	return ""
}

func (m *Summary) SetDate(date *string) {
	m.Date = date
}

func (m *Summary) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *Summary) SetCategory(c *entity.Category) {
	m.Category = c
}

func (m *Summary) GetSum() float64 {
	if m != nil && m.Sum != nil {
		return *m.Sum
	}
	return 0
}

func (m *Summary) SetSum(sum *float64) {
	m.Sum = sum

	if sum != nil {
		s := util.RoundFloatToStandardDP(*sum)
		m.Sum = goutil.Float64(s)
	}
}

func (m *Summary) GetTotalExpense() float64 {
	if m != nil && m.TotalExpense != nil {
		return *m.TotalExpense
	}
	return 0
}

func (m *Summary) SetTotalExpense(totalExpense *float64) {
	m.TotalExpense = totalExpense

	if totalExpense != nil {
		s := util.RoundFloatToStandardDP(*totalExpense)
		m.TotalExpense = goutil.Float64(s)
	}
}

func (m *Summary) GetTotalIncome() float64 {
	if m != nil && m.TotalIncome != nil {
		return *m.TotalIncome
	}
	return 0
}

func (m *Summary) SetTotalIncome(totalIncome *float64) {
	m.TotalIncome = totalIncome

	if totalIncome != nil {
		s := util.RoundFloatToStandardDP(*totalIncome)
		m.TotalIncome = goutil.Float64(s)
	}
}

func (m *Summary) GetTransactionType() uint32 {
	if m != nil && m.TransactionType != nil {
		return *m.TransactionType
	}
	return 0
}

func (m *Summary) SetTransactionType(transactionType *uint32) {
	m.TransactionType = transactionType
}

func (m *Summary) GetCurrency() string {
	if m != nil && m.Currency != nil {
		return *m.Currency
	}
	return ""
}

func (m *Summary) SetCurrency(currency *string) {
	m.Currency = currency
}

func (m *Summary) GetTransactions() []*entity.Transaction {
	if m != nil && m.Transactions != nil {
		return m.Transactions
	}
	return nil
}

func (m *Summary) SetTransactions(ts []*entity.Transaction) {
	m.Transactions = ts
}
