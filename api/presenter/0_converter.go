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

func toCategory(
	category *entity.Category,
) *Category {
	if category == nil {
		return nil
	}

	return &Category{
		CategoryID:   goutil.String(category.GetCategoryID()),
		CategoryName: goutil.String(category.GetCategoryName()),
		CategoryType: goutil.Uint32(category.GetCategoryType()),
		CreateTime:   goutil.Uint64(category.GetCreateTime()),
		UpdateTime:   goutil.Uint64(category.GetUpdateTime()),
	}
}

func toUser(
	user *entity.User,
) *User {
	if user == nil {
		return nil
	}

	return &User{
		UserID:     goutil.String(user.GetUserID()),
		Username:   goutil.String(user.GetUsername()),
		UserStatus: goutil.Uint32(user.GetUserStatus()),
		CreateTime: goutil.Uint64(user.GetCreateTime()),
		UpdateTime: goutil.Uint64(user.GetUpdateTime()),
	}
}

func toTransaction(transaction *entity.Transaction, category *entity.Category) *Transaction {
	return &Transaction{
		TransactionID:   goutil.String(transaction.GetTransactionID()),
		Category:        toCategory(category),
		Amount:          goutil.String(fmt.Sprint(transaction.GetAmount())),
		Note:            goutil.String(transaction.GetNote()),
		TransactionType: goutil.Uint32(transaction.GetTransactionType()),
		TransactionTime: goutil.Uint64(transaction.GetTransactionTime()),
		CreateTime:      goutil.Uint64(transaction.GetCreateTime()),
		UpdateTime:      goutil.Uint64(transaction.GetUpdateTime()),
	}
}

func toPaging(paging *common.Paging) *Paging {
	return &Paging{
		Limit: goutil.Uint32(paging.GetLimit()),
		Page:  goutil.Uint32(paging.GetPage()),
	}
}

func toAggr(aggr *transaction.Aggr) *Aggr {
	return &Aggr{
		Sum: aggr.Sum,
	}
}
