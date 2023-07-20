package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepo interface {
	Get(ctx context.Context, cf *CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, cf *CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) (string, error)
	Update(ctx context.Context, cf *CategoryFilter, c *entity.CategoryUpdate) error
}

type CategoryFilter struct {
	UserID       *string  `filter:"user_id"`
	CategoryID   *string  `filter:"_id"`
	CategoryIDs  []string `filter:"_id__in"`
	CategoryType *uint32  `filter:"category_type"`
	CategoryName *string  `filter:"category_name"`
}

func (f *CategoryFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *CategoryFilter) GetCategoryID() string {
	if f != nil && f.CategoryID != nil {
		return *f.CategoryID
	}
	return ""
}

func (f *CategoryFilter) GetCategoryName() string {
	if f != nil && f.CategoryName != nil {
		return *f.CategoryName
	}
	return ""
}

func (f *CategoryFilter) GetCategoryIDs() []string {
	if f != nil && f.CategoryIDs != nil {
		return f.CategoryIDs
	}
	return nil
}

func (f *CategoryFilter) GetCategoryType() uint32 {
	if f != nil && f.CategoryType != nil {
		return *f.CategoryType
	}
	return 0
}
