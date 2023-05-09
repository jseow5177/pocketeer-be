package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/data/entity"
	"github.com/jseow5177/pockteer-be/data/model"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
)

const categoryCollName = "category"

type CategoryMongo struct {
	mColl *MongoColl
}

func NewCategoryMongo(mongo *Mongo) *CategoryMongo {
	return &CategoryMongo{
		mColl: NewMongoColl(mongo, categoryCollName),
	}
}

func (m *CategoryMongo) Create(ctx context.Context, c *entity.Category) (string, error) {
	cm := model.ToCategoryModel(c)

	id, err := m.mColl.create(ctx, cm)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *CategoryMongo) Update(ctx context.Context, c *entity.Category) error {
	return nil
}

func (m *CategoryMongo) Get(ctx context.Context, cf *entity.CategoryFilter) (*entity.Category, error) {
	f := mongoutil.BuildFilter(cf)

	c := new(model.Category)
	if err := m.mColl.get(ctx, f, &c); err != nil {
		return nil, err
	}

	return model.ToCategoryEntity(c), nil
}

func (m *CategoryMongo) GetMany(ctx context.Context, cf *entity.CategoryFilter) ([]*entity.Category, error) {
	f := mongoutil.BuildFilter(cf)

	res, err := m.mColl.getMany(ctx, f, new(model.Category))
	if err != nil {
		return nil, err
	}

	ecs := make([]*entity.Category, 0, len(res))
	for _, r := range res {
		ecs = append(ecs, model.ToCategoryEntity(r.(*model.Category)))
	}

	return ecs, nil
}
