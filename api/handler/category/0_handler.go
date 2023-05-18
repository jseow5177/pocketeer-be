package category

import cuc "github.com/jseow5177/pockteer-be/usecase/category"

type CategoryHandler struct {
	categoryUseCase cuc.UseCase
}

func NewCategoryHandler(categoryUseCase cuc.UseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
	}
}
