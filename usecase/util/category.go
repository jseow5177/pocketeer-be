package util

import "github.com/jseow5177/pockteer-be/entity"

func GetCatIDs(categories []*entity.Category) []string {
	ids := make([]string, len(categories))
	for idx, cat := range categories {
		ids[idx] = cat.GetCategoryID()
	}
	return ids
}

func GetIDToCategoryMap(categories []*entity.Category) map[string]*entity.Category {
	idToCategoryMap := make(map[string]*entity.Category)

	for _, cat := range categories {
		idToCategoryMap[cat.GetCategoryID()] = cat
	}

	return idToCategoryMap
}

func GetCategoryIDToBudgetMap(
	budgets []*entity.Budget,
) map[string]*entity.Budget {
	categoryIDToBudgetMap := make(map[string]*entity.Budget)

	for _, budget := range budgets {
		categoryIDToBudgetMap[budget.GetCategoryID()] = budget
	}

	return categoryIDToBudgetMap
}
