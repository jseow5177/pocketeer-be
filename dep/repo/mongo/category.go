package mongo

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
	now := uint64(time.Now().Unix())

	c.UpdateTime = goutil.Uint64(now)
	c.UpdateTime = goutil.Uint64(now)

	cm := model.ToCategoryModel(c)
	id, err := m.mColl.create(ctx, cm)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *categoryMongo) Update(ctx context.Context, cf *repo.CategoryFilter, c *entity.Category) error {
	c.UpdateTime = goutil.Uint64(uint64(time.Now().Unix()))

	cm := model.ToCategoryModel(c)
	if err := m.mColl.update(ctx, cf, cm); err != nil {
		return err
	}

	return nil
}

func (m *categoryMongo) Get(ctx context.Context, cf *repo.CategoryFilter) (*entity.Category, error) {
	c := new(model.Category)
	if err := m.mColl.get(ctx, cf, &c); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrCategoryNotFound
		}
		return nil, err
	}

	return model.ToCategoryEntity(c), nil
}

func (m *categoryMongo) GetMany(ctx context.Context, cf *repo.CategoryFilter) ([]*entity.Category, error) {
	res, err := m.mColl.getMany(ctx, cf, nil, new(model.Category))
	if err != nil {
		return nil, err
	}

	ecs := make([]*entity.Category, 0, len(res))
	for _, r := range res {
		ecs = append(ecs, model.ToCategoryEntity(r.(*model.Category)))
	}

	return ecs, nil
}
