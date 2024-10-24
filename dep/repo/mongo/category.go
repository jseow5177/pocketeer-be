package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const categoryCollName = "category"

type categoryMongo struct {
	mColl *MongoColl
}

func NewCategoryMongo(mongo *Mongo) repo.CategoryRepo {
	return &categoryMongo{
		mColl: NewMongoColl(mongo, categoryCollName),
	}
}

func (m *categoryMongo) Create(ctx context.Context, c *entity.Category) (string, error) {
	cm := model.ToCategoryModelFromEntity(c)
	id, err := m.mColl.create(ctx, cm)
	if err != nil {
		return "", err
	}
	c.SetCategoryID(goutil.String(id))

	return id, nil
}

func (m *categoryMongo) CreateMany(ctx context.Context, cs []*entity.Category) ([]string, error) {
	cms := make([]interface{}, 0)
	for _, c := range cs {
		cms = append(cms, model.ToCategoryModelFromEntity(c))
	}

	ids, err := m.mColl.createMany(ctx, cms)
	if err != nil {
		return nil, err
	}

	for i, c := range cs {
		c.SetCategoryID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *categoryMongo) Update(ctx context.Context, cf *repo.CategoryFilter, cu *entity.CategoryUpdate) error {
	f := mongoutil.BuildFilter(cf)

	cm := model.ToCategoryModelFromUpdate(cu)
	if err := m.mColl.update(ctx, f, cm); err != nil {
		return err
	}

	return nil
}

func (m *categoryMongo) Get(ctx context.Context, cf *repo.CategoryFilter) (*entity.Category, error) {
	f := mongoutil.BuildFilter(cf)

	c := new(model.Category)
	if err := m.mColl.get(ctx, &c, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrCategoryNotFound
		}
		return nil, err
	}

	return model.ToCategoryEntity(c)
}

func (m *categoryMongo) GetMany(ctx context.Context, cf *repo.CategoryFilter) ([]*entity.Category, error) {
	f := mongoutil.BuildFilter(cf)

	res, err := m.mColl.getMany(ctx, new(model.Category), nil, f)
	if err != nil {
		return nil, err
	}

	ecs := make([]*entity.Category, 0, len(res))
	for _, r := range res {
		ec, err := model.ToCategoryEntity(r.(*model.Category))
		if err != nil {
			return nil, err
		}
		ecs = append(ecs, ec)
	}

	return ecs, nil
}
