package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
	acm := model.ToAccountModel(ac)
	id, err := m.mColl.create(ctx, acm)
	if err != nil {
		return "", err
	}
	entity.SetAccount(ac, entity.WithAccountID(goutil.String(id)))

	return id, nil
}

func (m *accountMongo) Update(ctx context.Context, acf *repo.AccountFilter, ac *entity.Account) error {
	acm := model.ToAccountModel(ac)
	if err := m.mColl.update(ctx, acf, acm); err != nil {
		return err
	}

	return nil
}

func (m *accountMongo) Get(ctx context.Context, acf *repo.AccountFilter) (*entity.Account, error) {
	ac := new(model.Account)
	if err := m.mColl.get(ctx, acf, &ac); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrAccountNotFound
		}
		return nil, err
	}

	return model.ToAccountEntity(ac), nil
}
