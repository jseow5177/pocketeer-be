package common

import "github.com/jseow5177/pockteer-be/entity"

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

type TransactionSummary struct {
	Category        *entity.Category
	TransactionType *uint32
	Sum             *float64
	Currency        *string
}

func (m *TransactionSummary) GetCategory() *entity.Category {
	if m != nil && m.Category != nil {
		return m.Category
	}
	return nil
}

func (m *TransactionSummary) GetSum() float64 {
	if m != nil && m.Sum != nil {
		return *m.Sum
	}
	return 0
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
