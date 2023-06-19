package category

import "github.com/jseow5177/pockteer-be/usecase/category"

type categoryHandler struct {
	categoryUseCase category.UseCase
}

func NewCategoryHandler(categoryUseCase category.UseCase) *categoryHandler {
	return &categoryHandler{
		categoryUseCase,
	}
}
