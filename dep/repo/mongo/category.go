package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/model"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

const categoryCollName = "category"

type CategoryMongo struct {
	coll *MongoColl
}

func NewCategoryMongo(mongo *Mongo) *CategoryMongo {
	return &CategoryMongo{
		coll: NewMongoColl(mongo, categoryCollName),
	}
}

func (m *CategoryMongo) Create(ctx context.Context, c *entity.Category) error {
	cm := model.ToCategoryModel(c)

	id, err := m.coll.create(ctx, cm)
	if err != nil {
		return err
	}
	c.CatID = goutil.String(id)

	return nil
}

func (m *CategoryMongo) Update(ctx context.Context, c *entity.Category) error {
	return nil
}

func (m *CategoryMongo) Get(ctx context.Context, f *entity.CategoryFilter) (*entity.Category, error) {
	return nil, nil
}

func (m *CategoryMongo) GetMany(ctx context.Context, f *entity.CategoryFilter) ([]*entity.Category, error) {
	return nil, nil
}
