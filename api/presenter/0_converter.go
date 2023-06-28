package presenter

import (
	"fmt"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/common"
	"github.com/jseow5177/pockteer-be/usecase/transaction"
)

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
		CategoryID:   goutil.String(c.GetCategoryID()),
		CategoryName: goutil.String(c.GetCategoryName()),
		CategoryType: goutil.Uint32(c.GetCategoryType()),
		CreateTime:   goutil.Uint64(c.GetCreateTime()),
		UpdateTime:   goutil.Uint64(c.GetUpdateTime()),
	}
}

func toUser(u *entity.User) *User {
	if u == nil {
		return nil
	}

	return &User{
		UserID:     goutil.String(u.GetUserID()),
		Username:   goutil.String(u.GetUsername()),
		UserStatus: goutil.Uint32(u.GetUserStatus()),
		CreateTime: goutil.Uint64(u.GetCreateTime()),
		UpdateTime: goutil.Uint64(u.GetUpdateTime()),
	}
}

func toTransaction(t *entity.Transaction) *Transaction {
	if t == nil {
		return nil
	}

	return &Transaction{
		TransactionID:   goutil.String(t.GetTransactionID()),
		Amount:          goutil.String(fmt.Sprint(t.GetAmount())),
		CategoryID:      goutil.String(t.GetCategoryID()),
		AccountID:       goutil.String(t.GetAccountID()),
		Note:            goutil.String(t.GetNote()),
		TransactionType: goutil.Uint32(t.GetTransactionType()),
		TransactionTime: goutil.Uint64(t.GetTransactionTime()),
		CreateTime:      goutil.Uint64(t.GetCreateTime()),
		UpdateTime:      goutil.Uint64(t.GetUpdateTime()),
	}
}

func toAccounts(acs []*entity.Account) []*Account {
	accounts := make([]*Account, len(acs))
	for idx, ac := range acs {
		accounts[idx] = toAccount(ac)
	}
	return accounts
}

func toAccount(ac *entity.Account) *Account {
	if ac == nil {
		return nil
	}

	return &Account{
		AccountID:     goutil.String(ac.GetAccountID()),
		AccountName:   goutil.String(ac.GetAccountName()),
		Balance:       goutil.String(fmt.Sprint(ac.GetBalance())),
		AccountType:   goutil.Uint32(ac.GetAccountType()),
		AccountStatus: goutil.Uint32(ac.GetAccountStatus()),
		Note:          goutil.String(ac.GetNote()),
		CreateTime:    goutil.Uint64(ac.GetCreateTime()),
		UpdateTime:    goutil.Uint64(ac.GetUpdateTime()),
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
