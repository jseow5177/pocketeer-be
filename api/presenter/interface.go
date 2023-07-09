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

type UInt64Filter struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func (uv *UInt64Filter) GetGte() uint64 {
	if uv != nil && uv.Gte != nil {
		return *uv.Gte
	}
	return 0
}

func (uv *UInt64Filter) GetLte() uint64 {
	if uv != nil && uv.Lte != nil {
		return *uv.Lte
	}
	return 0
}

func toBudgets(
	entities []*entity.Budget,
) []*Budget {
	budgets := make([]*Budget, len(entities))
	for idx, entity := range entities {
		budgets[idx] = toBudget(entity)
	}
	return budgets
}

func toBudget(
	budget *entity.Budget,
) *Budget {
	return &Budget{
		BudgetID:         budget.BudgetID,
		BudgetName:       budget.BudgetName,
		BudgetType:       budget.BudgetType,
		CategoryIDs:      budget.CategoryIDs,
		BudgetBreakdowns: toBudgetBreakdowns(budget.BreakdownMap),
	}
}

func toBudgetBreakdowns(
	breakdownMap entity.BreakdownMap,
) []*BudgetBreakdown {
	breakdowns := make([]*BudgetBreakdown, 0)

	for _, bd := range breakdownMap {
		if bd == nil {
			continue
		}

		breakdowns = append(
			breakdowns,
			&BudgetBreakdown{
				Year:   bd.Year,
				Month:  bd.Month,
				Amount: bd.Amount,
			},
		)
	}
	return breakdowns
}

func toCategories(cs []*entity.Category) []*Category {
	categories := make([]*Category, len(cs))
	for idx, c := range cs {
		categories[idx] = toCategory(c)
	}
	return categories
}

func toCategory(c *entity.Category) *Category {
	if c == nil {
		return nil
	}

	return &Category{
		CategoryID:   c.CategoryID,
		CategoryName: c.CategoryName,
		CategoryType: c.CategoryType,
		CreateTime:   c.CreateTime,
		UpdateTime:   c.UpdateTime,
	}
}

func toUser(u *entity.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		UserID:     u.UserID,
		Username:   u.Username,
		UserStatus: u.UserStatus,
		CreateTime: u.CreateTime,
		UpdateTime: u.UpdateTime,
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
		TransactionID:   t.TransactionID,
		Amount:          amount,
		CategoryID:      t.CategoryID,
		AccountID:       t.AccountID,
		Note:            t.Note,
		TransactionType: t.TransactionType,
		TransactionTime: t.TransactionTime,
		CreateTime:      t.CreateTime,
		UpdateTime:      t.UpdateTime,
	}
}

func toHolding(h *entity.Holding) *Holding {
	if h == nil {
		return nil
	}

	return &Holding{
		HoldingID:     h.HoldingID,
		AccountID:     h.AccountID,
		Symbol:        h.Symbol,
		HoldingType:   h.HoldingType,
		HoldingStatus: h.HoldingStatus,
		CreateTime:    h.CreateTime,
		UpdateTime:    h.UpdateTime,
		LatestValue:   h.LatestValue,
		AvgCost:       h.AvgCost,
		TotalShares:   h.TotalShares,
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

func toAccount(ac *entity.Account) *Account {
	if ac == nil {
		return nil
	}

	var balance *string
	if ac.Balance != nil {
		balance = goutil.String(fmt.Sprint(ac.GetBalance()))
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
	}
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

func toSecurity(security *entity.Security) *Security {
	return &Security{
		Symbol:       security.Symbol,
		SecurityName: security.SecurityName,
		SecurityType: security.SecurityType,
		Region:       security.Region,
		Currency:     security.Currency,
	}
}
