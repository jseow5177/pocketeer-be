package category

import cuc "github.com/jseow5177/pockteer-be/usecase/category"

type categoryHandler struct {
	categoryUseCase cuc.UseCase
}

func NewCategoryHandler(categoryUseCase cuc.UseCase) *categoryHandler {
	return &categoryHandler{
		categoryUseCase: categoryUseCase,
	}
}
