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
	acm := model.ToAccountModelFromEntity(ac)
	id, err := m.mColl.create(ctx, acm)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", repo.ErrAccountAlreadyExists
		}
		return "", err
	}
	ac.SetAccountID(goutil.String(id))

	return id, nil
}

func (m *accountMongo) Update(ctx context.Context, acf *repo.AccountFilter, acu *entity.AccountUpdate) error {
	acm := model.ToAccountModelFromUpdate(acu)
	if err := m.mColl.update(ctx, acf, acm); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return repo.ErrAccountAlreadyExists
		}
		return err
	}

	return nil
}

func (m *accountMongo) Get(ctx context.Context, acf *repo.AccountFilter) (*entity.Account, error) {
	ac := new(model.Account)
	if err := m.mColl.get(ctx, &ac, acf); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrAccountNotFound
		}
		return nil, err
	}

	return model.ToAccountEntity(ac)
}

func (m *accountMongo) GetMany(ctx context.Context, acf *repo.AccountFilter) ([]*entity.Account, error) {
	res, err := m.mColl.getMany(ctx, new(model.Account), nil, acf)
	if err != nil {
		return nil, err
	}

	eacs := make([]*entity.Account, 0, len(res))
	for _, r := range res {
		eac, err := model.ToAccountEntity(r.(*model.Account))
		if err != nil {
			return nil, err
		}
		eacs = append(eacs, eac)
	}

	return eacs, nil
}
