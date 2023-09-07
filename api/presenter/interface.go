package presenter

import (
	"fmt"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

type Paging struct {
	Limit *uint32 `json:"limit,omitempty"`
	Page  *uint32 `json:"page,omitempty"`
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
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func (uv *RangeFilter) GetGte() uint64 {
	if uv != nil && uv.Gte != nil {
		return *uv.Gte
	}
	return 0
}

func (uv *RangeFilter) GetLte() uint64 {
	if uv != nil && uv.Lte != nil {
		return *uv.Lte
	}
	return 0
}

func toBudget(b *entity.Budget) *Budget {
	if b == nil {
		return nil
	}

	var amount *string
	if b.Amount != nil {
		amount = goutil.String(fmt.Sprint(b.GetAmount()))
	}

	var usedAmount *string
	if b.UsedAmount != nil {
		usedAmount = goutil.String(fmt.Sprint(b.GetUsedAmount()))
	}

	return &Budget{
		BudgetID:     b.BudgetID,
		CategoryID:   b.CategoryID,
		BudgetType:   b.BudgetType,
		BudgetStatus: b.BudgetStatus,
		Amount:       amount,
		CreateTime:   b.CreateTime,
		UpdateTime:   b.UpdateTime,
		UsedAmount:   usedAmount,
	}
}

func toCategory(c *entity.Category) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		CategoryID:     c.CategoryID,
		CategoryName:   c.CategoryName,
		CategoryType:   c.CategoryType,
		CategoryStatus: c.CategoryStatus,
		CreateTime:     c.CreateTime,
		UpdateTime:     c.UpdateTime,
		Budget:         toBudget(c.Budget),
	}
}

func toCategories(cs []*entity.Category) []*Category {
	categories := make([]*Category, len(cs))
	for idx, c := range cs {
		categories[idx] = toCategory(c)
	}
	return categories
}

func toUserMeta(um *entity.UserMeta) *UserMeta {
	if um == nil {
		return nil
	}

	return &UserMeta{
		InitStage: um.InitStage,
	}
}

func toUser(u *entity.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		UserID:     u.UserID,
		Username:   u.Username,
		Email:      u.Email,
		UserFlag:   u.UserFlag,
		UserStatus: u.UserStatus,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
		Meta:       toUserMeta(u.Meta),
	}
}

func toTransaction(t *entity.Transaction) *Transaction {
	if t == nil {
		return nil
	}

	var amount *string
	if t.Amount != nil {
		amount = goutil.String(fmt.Sprint(t.GetAmount()))
	}

	return &Transaction{
		TransactionID:     t.TransactionID,
		Amount:            amount,
		CategoryID:        t.CategoryID,
		Category:          toCategory(t.Category),
		AccountID:         t.AccountID,
		Account:           toAccount(t.Account),
		Note:              t.Note,
		TransactionStatus: t.TransactionStatus,
		TransactionType:   t.TransactionType,
		TransactionTime:   t.TransactionTime,
		CreateTime:        t.CreateTime,
		UpdateTime:        t.UpdateTime,
	}
}

func toTransactions(ts []*entity.Transaction) []*Transaction {
	transactions := make([]*Transaction, len(ts))
	for idx, t := range ts {
		transactions[idx] = toTransaction(t)
	}
	return transactions
}

func toHolding(h *entity.Holding) *Holding {
	if h == nil {
		return nil
	}

	var latestValue *string
	if h.LatestValue != nil {
		latestValue = goutil.String(fmt.Sprint(h.GetLatestValue()))
	}

	var totalCost *string
	if h.TotalCost != nil {
		totalCost = goutil.String(fmt.Sprint(h.GetTotalCost()))
	}

	var totalShares *string
	if h.TotalShares != nil {
		totalShares = goutil.String(fmt.Sprint(h.GetTotalShares()))
	}

	var avgCostPerShare *string
	if h.AvgCostPerShare != nil {
		avgCostPerShare = goutil.String(fmt.Sprint(h.GetAvgCostPerShare()))
	}

	return &Holding{
		HoldingID:       h.HoldingID,
		AccountID:       h.AccountID,
		Symbol:          h.Symbol,
		HoldingType:     h.HoldingType,
		HoldingStatus:   h.HoldingStatus,
		CreateTime:      h.CreateTime,
		UpdateTime:      h.UpdateTime,
		LatestValue:     latestValue,
		TotalCost:       totalCost,
		TotalShares:     totalShares,
		AvgCostPerShare: avgCostPerShare,
		Quote:           toQuote(h.Quote),
		Lots:            toLots(h.Lots),
	}
}

func toHoldings(hs []*entity.Holding) []*Holding {
	holdings := make([]*Holding, len(hs))
	for idx, h := range hs {
		holdings[idx] = toHolding(h)
	}
	return holdings
}

func toQuote(q *entity.Quote) *Quote {
	if q == nil {
		return nil
	}

	var latestPrice *string
	if q.LatestPrice != nil {
		latestPrice = goutil.String(fmt.Sprint(q.GetLatestPrice()))
	}

	var change *string
	if q.Change != nil {
		change = goutil.String(fmt.Sprint(q.GetChange()))
	}

	var changePercent *string
	if q.ChangePercent != nil {
		changePercent = goutil.String(fmt.Sprint(q.GetChangePercent()))
	}

	var previousClose *string
	if q.PreviousClose != nil {
		previousClose = goutil.String(fmt.Sprint(q.GetPreviousClose()))
	}

	return &Quote{
		LatestPrice:   latestPrice,
		Change:        change,
		ChangePercent: changePercent,
		PreviousClose: previousClose,
		UpdateTime:    q.UpdateTime,
		Currency:      q.Currency,
	}
}

func toLot(l *entity.Lot) *Lot {
	if l == nil {
		return nil
	}

	var shares *string
	if l.Shares != nil {
		shares = goutil.String(fmt.Sprint(l.GetShares()))
	}

	var costPerShare *string
	if l.CostPerShare != nil {
		costPerShare = goutil.String(fmt.Sprint(l.GetCostPerShare()))
	}

	return &Lot{
		LotID:        l.LotID,
		HoldingID:    l.HoldingID,
		Shares:       shares,
		CostPerShare: costPerShare,
		LotStatus:    l.LotStatus,
		TradeDate:    l.TradeDate,
		CreateTime:   l.CreateTime,
		UpdateTime:   l.UpdateTime,
	}
}

func toLots(ls []*entity.Lot) []*Lot {
	lots := make([]*Lot, len(ls))
	for idx, l := range ls {
		lots[idx] = toLot(l)
	}
	return lots
}

func toAccount(ac *entity.Account) *Account {
	if ac == nil {
		return nil
	}

	var balance *string
	if ac.Balance != nil {
		balance = goutil.String(fmt.Sprint(ac.GetBalance()))
	}

	var totalCost *string
	if ac.TotalCost != nil {
		totalCost = goutil.String(fmt.Sprint(ac.GetTotalCost()))
	}

	return &Account{
		AccountID:     ac.AccountID,
		AccountName:   ac.AccountName,
		Balance:       balance,
		AccountType:   ac.AccountType,
		AccountStatus: ac.AccountStatus,
		Note:          ac.Note,
		CreateTime:    ac.CreateTime,
		UpdateTime:    ac.UpdateTime,
		TotalCost:     totalCost,
		Holdings:      toHoldings(ac.Holdings),
	}
}

func toAccounts(acs []*entity.Account) []*Account {
	accounts := make([]*Account, len(acs))
	for idx, ac := range acs {
		accounts[idx] = toAccount(ac)
	}
	return accounts
}

func toPaging(p *common.Paging) *Paging {
	if p == nil {
		return nil
	}

	return &Paging{
		Limit: p.Limit,
		Page:  p.Page,
	}
}

func toAggr(aggr *transaction.Aggr) *Aggr {
	if aggr == nil {
		return nil
	}

	return &Aggr{
		Sum: aggr.Sum,
	}
}

func toSecurity(s *entity.Security) *Security {
	if s == nil {
		return nil
	}

	return &Security{
		Symbol:       s.Symbol,
		SecurityName: s.SecurityName,
		SecurityType: s.SecurityType,
		Region:       s.Region,
		Currency:     s.Currency,
	}
}

func toExchangeRate(er *entity.ExchangeRate) *ExchangeRate {
	if er == nil {
		return nil
	}

	return &ExchangeRate{
		ExchangeRateID: er.ExchangeRateID,
		From:           er.From,
		To:             er.To,
		Timestamp:      er.Timestamp,
		CreateTime:     er.CreateTime,
	}
}

func toExchangeRates(ers []*entity.ExchangeRate) []*ExchangeRate {
	exchangeRates := make([]*ExchangeRate, len(ers))
	for idx, er := range ers {
		exchangeRates[idx] = toExchangeRate(er)
	}
	return exchangeRates
}
