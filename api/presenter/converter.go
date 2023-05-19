package presenter

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/usecase/budget"
)

func toAnnualBudgetBreakdown(
	annualBreakdown *entity.AnnualBudgetBreakdown,
) *AnnualBudgetBreakdown {
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
	return &BasicBudget{
		Year:   goutil.Uint32(entBudget.GetYear()),
		Month:  goutil.Uint32(entBudget.GetMonth()),
		Amount: goutil.Int64(entBudget.GetBudgetAmount()),
	}
}

func cbToBudgets(
	categoryBudgets []*budget.CategoryBudget,
) []*Budget {
	budgets := make([]*Budget, len(categoryBudgets))
	for idx, cb := range categoryBudgets {
		budgets[idx] = cbToBudget(cb)
	}
	return budgets
}

func cbToBudget(
	categoryBudget *budget.CategoryBudget,
) *Budget {
	return &Budget{
		BudgetID:     goutil.String(categoryBudget.GetBudget().GetBudgetID()),
		BudgetType:   goutil.Uint32(categoryBudget.GetBudget().GetBudgetType()),
		Year:         goutil.Uint32(categoryBudget.GetBudget().GetYear()),
		Month:        goutil.Uint32(categoryBudget.GetBudget().GetMonth()),
		BudgetAmount: goutil.Int64(categoryBudget.GetBudget().GetBudgetAmount()),
		Category:     toCategory(categoryBudget.GetCategory()),
	}
}

func toCategory(
	category *entity.Category,
) *Category {
	return &Category{
		CategoryID:   goutil.String(category.GetCategoryID()),
		CategoryName: goutil.String(category.GetCategoryName()),
		CategoryType: goutil.Uint32(category.GetCategoryType()),
		CreateTime:   goutil.Uint64(category.GetCreateTime()),
		UpdateTime:   goutil.Uint64(category.GetUpdateTime()),
	}
}
