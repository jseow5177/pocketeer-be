package budget

import (
	"fmt"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/usecase/util"
)

func toCategoryBudgets(
	categoryIDToBudget map[string]*entity.Budget,
	categories []*entity.Category,
) ([]*CategoryBudget, error) {
	categoryBudgets := make([]*CategoryBudget, 0)
	idToCategoryMap := util.GetIDToCategoryMap(categories)

	for categoryID, budget := range categoryIDToBudget {
		category, ok := idToCategoryMap[categoryID]
		if !ok {
			return nil, errutil.InternalServerError(
				fmt.Errorf(
					"fail to convert to categoryBudgets, cannot find category_id=%s",
					categoryID,
				),
			)
		}

		categoryBudgets = append(
			categoryBudgets,
			&CategoryBudget{
				CategoryName: category.GetCategoryName(),
				Budget:       budget,
			},
		)
	}

	return categoryBudgets, nil
}
