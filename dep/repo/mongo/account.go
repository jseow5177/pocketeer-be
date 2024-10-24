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

const accountCollName = "account"

type accountMongo struct {
	mColl *MongoColl
}

func NewAccountMongo(mongo *Mongo) repo.AccountRepo {
	return &accountMongo{
		mColl: NewMongoColl(mongo, accountCollName),
	}
}

func (m *accountMongo) Create(ctx context.Context, ac *entity.Account) (string, error) {
	acm := model.ToAccountModelFromEntity(ac)
	id, err := m.mColl.create(ctx, acm)
	if err != nil {
		return "", err
	}
	ac.SetAccountID(goutil.String(id))

	return id, nil
}

func (m *accountMongo) CreateMany(ctx context.Context, acs []*entity.Account) ([]string, error) {
	acms := make([]interface{}, 0)
	for _, ac := range acs {
		acms = append(acms, model.ToAccountModelFromEntity(ac))
	}

	ids, err := m.mColl.createMany(ctx, acms)
	if err != nil {
		return nil, err
	}

	for i, ac := range acs {
		ac.SetAccountID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *accountMongo) Update(ctx context.Context, acf *repo.AccountFilter, acu *entity.AccountUpdate) error {
	f := mongoutil.BuildFilter(acf)

	acm := model.ToAccountModelFromUpdate(acu)
	if err := m.mColl.update(ctx, f, acm); err != nil {
		return err
	}

	return nil
}

func (m *accountMongo) Get(ctx context.Context, acf *repo.AccountFilter) (*entity.Account, error) {
	f := mongoutil.BuildFilter(acf)

	ac := new(model.Account)
	if err := m.mColl.get(ctx, &ac, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrAccountNotFound
		}
		return nil, err
	}

	return model.ToAccountEntity(ac)
}

func (m *accountMongo) GetMany(ctx context.Context, acf *repo.AccountFilter) ([]*entity.Account, error) {
	f := mongoutil.BuildFilter(acf)

	res, err := m.mColl.getMany(ctx, new(model.Account), nil, f)
	if err != nil {
		return nil, err
	}

	acs := make([]*entity.Account, 0, len(res))
	for _, r := range res {
		ac, err := model.ToAccountEntity(r.(*model.Account))
		if err != nil {
			return nil, err
		}
		acs = append(acs, ac)
	}

	return acs, nil
}
