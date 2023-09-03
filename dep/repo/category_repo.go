package repo

import (
	"context"
	"errors"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryRepo interface {
	Get(ctx context.Context, cf *CategoryFilter) (*entity.Category, error)
	GetMany(ctx context.Context, cf *CategoryFilter) ([]*entity.Category, error)

	Create(ctx context.Context, c *entity.Category) (string, error)
	CreateMany(ctx context.Context, cs []*entity.Category) ([]string, error)
	Update(ctx context.Context, cf *CategoryFilter, c *entity.CategoryUpdate) error
	Delete(ctx context.Context, cf *CategoryFilter) error
}

type CategoryFilter struct {
	UserID         *string  `filter:"user_id"`
	CategoryID     *string  `filter:"_id"`
	CategoryIDs    []string `filter:"_id__in"`
	CategoryType   *uint32  `filter:"category_type"`
	CategoryStatus *uint32  `filter:"category_status"`
	CategoryName   *string  `filter:"category_name"`
}

type CategoryFilterOption = func(cf *CategoryFilter)

func WithCategoryID(categoryID *string) CategoryFilterOption {
	return func(cf *CategoryFilter) {
		cf.CategoryID = categoryID
	}
}

func WithCategoryIDs(categoryIDs []string) CategoryFilterOption {
	return func(cf *CategoryFilter) {
		cf.CategoryIDs = categoryIDs
	}
}

func WithCategoryType(categoryType *uint32) CategoryFilterOption {
	return func(cf *CategoryFilter) {
		cf.CategoryType = categoryType
	}
}

func WithCategoryStatus(categoryStatus *uint32) CategoryFilterOption {
	return func(cf *CategoryFilter) {
		cf.CategoryStatus = categoryStatus
	}
}

func WithCategoryName(categoryName *string) CategoryFilterOption {
	return func(cf *CategoryFilter) {
		cf.CategoryName = categoryName
	}
}

func NewCategoryFilter(userID string, opts ...CategoryFilterOption) *CategoryFilter {
	cf := &CategoryFilter{
		UserID:         goutil.String(userID),
		CategoryStatus: goutil.Uint32(uint32(entity.CategoryStatusNormal)),
	}
	for _, opt := range opts {
		opt(cf)
	}
	return cf
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

func (f *CategoryFilter) GetCategoryStatus() uint32 {
	if f != nil && f.CategoryStatus != nil {
		return *f.CategoryStatus
	}
	return 0
}
