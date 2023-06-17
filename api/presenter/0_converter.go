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

func toCategories(
	entities []*entity.Category,
) []*Category {
	categories := make([]*Category, len(entities))
	for idx, entity := range entities {
		categories[idx] = toCategory(entity)
	}
	return categories
}

func toCategory(category *entity.Category) *Category {
	if category == nil {
		return nil
	}

	return &Category{
		CategoryID:   category.CategoryID,
		CategoryName: category.CategoryName,
		CategoryType: category.CategoryType,
		CreateTime:   category.CreateTime,
		UpdateTime:   category.UpdateTime,
	}
}

func toUser(user *entity.User) *User {
	if user == nil {
		return nil
	}

	return &User{
		UserID:     user.UserID,
		Username:   user.Username,
		UserStatus: user.UserStatus,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}

func toTransaction(transaction *entity.Transaction, category *entity.Category) *Transaction {
	if transaction == nil {
		return nil
	}

	var ammount *string
	if transaction.Amount != nil {
		ammount = goutil.String(fmt.Sprint(transaction.GetAmount()))
	}

	return &Transaction{
		TransactionID:   transaction.TransactionID,
		Category:        toCategory(category),
		Amount:          ammount,
		Note:            transaction.Note,
		TransactionType: transaction.TransactionType,
		TransactionTime: transaction.TransactionTime,
		CreateTime:      transaction.CreateTime,
		UpdateTime:      transaction.UpdateTime,
	}
}

func toAccount(account *entity.Account) *Account {
	if account == nil {
		return nil
	}

	var balance *string
	if account.Balance != nil {
		balance = goutil.String(fmt.Sprint(account.GetBalance()))
	}

	return &Account{
		AccountID:     account.AccountID,
		AccountName:   account.AccountName,
		Balance:       balance,
		AccountType:   account.AccountType,
		AccountStatus: account.AccountStatus,
		Note:          account.Note,
		CreateTime:    account.CreateTime,
		UpdateTime:    account.UpdateTime,
	}
}

func toPaging(paging *common.Paging) *Paging {
	if paging == nil {
		return nil
	}

	return &Paging{
		Limit: paging.Limit,
		Page:  paging.Page,
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
