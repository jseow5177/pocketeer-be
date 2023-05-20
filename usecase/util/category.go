package util

import "github.com/jseow5177/pockteer-be/entity"

func GetCategoryIDToBudgetMap(
	budgets []*entity.Budget,
) map[string]*entity.Budget {
	categoryIDToBudgetMap := make(map[string]*entity.Budget)

	for _, budget := range budgets {
		categoryIDToBudgetMap[budget.GetCategoryID()] = budget
	}

	return categoryIDToBudgetMap
}
