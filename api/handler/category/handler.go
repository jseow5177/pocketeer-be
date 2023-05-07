package category

import "github.com/jseow5177/pockteer-be/dep/repo"

type CategoryHandler struct {
	categoryRepo repo.CategoryRepo
}

func NewCategoryHandler(categoryRepo repo.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}
