package presenter

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
	"github.com/jseow5177/pockteer-be/usecase/common"
)

func toAnnualBudgetBreakdown(
	annualBreakdown *entity.AnnualBudgetBreakdown,
) *AnnualBudgetBreakdown {
	if annualBreakdown == nil {
		return nil
	}

	return &AnnualBudgetBreakdown{
		CategoryID:     goutil.String(annualBreakdown.GetDefaultBudget().GetCategoryID()),
		BudgetType:     goutil.Uint32(annualBreakdown.GetBudgetType()),
		Year:           goutil.Uint32(annualBreakdown.GetDefaultBudget().GetYear()),
		DefaultBudget:  goutil.Int64(annualBreakdown.GetDefaultBudget().GetBudgetAmount()),
		MonthlyBudgets: toBasicBudgets(annualBreakdown.GetMonthlyBudgets()),
	}
}

func toBasicBudgets(
	entBudgets []*entity.Budget,
) []*BasicBudget {
	budgets := make([]*BasicBudget, len(entBudgets))
	for idx, b := range entBudgets {
		budgets[idx] = toBasicBudget(b)
	}
	return budgets
}

func toBasicBudget(
	entBudget *entity.Budget,
) *BasicBudget {
	if entBudget == nil {
		return nil
	}

	return &BasicBudget{
		Year:   goutil.Uint32(entBudget.GetYear()),
		Month:  goutil.Uint32(entBudget.GetMonth()),
		Amount: goutil.Int64(entBudget.GetBudgetAmount()),
	}
}

func toCategoryBudgets(
	categoryBudgets []*budget.CategoryBudget,
) []*CategoryBudget {
	cbs := make([]*CategoryBudget, len(categoryBudgets))
	for idx, cb := range categoryBudgets {
		cbs[idx] = toCategoryBudget(cb)
	}
	return cbs
}

func toCategoryBudget(
	categoryBudget *budget.CategoryBudget,
) *CategoryBudget {
	if categoryBudget == nil {
		return nil
	}

	return &CategoryBudget{
		Budget:   toBudget(categoryBudget.Budget),
		Category: toCategory(categoryBudget.Category),
	}
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

func toBudget(
	budget *entity.Budget,
) *Budget {
	if budget == nil {
		return nil
	}

	return &Budget{
		BudgetID:     goutil.String(budget.GetBudgetID()),
		CategoryID:   goutil.String(budget.GetCategoryID()),
		BudgetType:   goutil.Uint32(budget.GetBudgetType()),
		Year:         goutil.Uint32(budget.GetYear()),
		Month:        goutil.Uint32(budget.GetMonth()),
		BudgetAmount: goutil.Int64(budget.GetBudgetAmount()),
	}
}

func toTransaction(transaction *entity.Transaction, category *entity.Category) *Transaction {
	return &Transaction{
		TransactionID:   goutil.String(transaction.GetTransactionID()),
		Category:        toCategory(category),
		Amount:          goutil.String(transaction.GetAmount()),
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
