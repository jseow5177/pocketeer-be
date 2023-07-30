package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/mongo"
)

const securityCollName = "security"

type securityMongo struct {
	mColl *MongoColl
}

func NewSecurityMongo(mongo *Mongo) repo.SecurityRepo {
	return &securityMongo{
		mColl: NewMongoColl(mongo, securityCollName),
	}
}

func (m *securityMongo) CreateMany(ctx context.Context, ss []*entity.Security) error {
	ssms := make([]interface{}, 0, len(ss))
	for _, s := range ss {
		ssms = append(ssms, model.ToSecurityModelFromEntity(s))
	}

	ids, err := m.mColl.createMany(ctx, ssms)
	if err != nil {
		return err
	}

	for i, s := range ss {
		s.SetSecurityID(goutil.String(ids[i]))
	}

	return nil
}

func (m *securityMongo) Update(ctx context.Context, sf *repo.SecurityFilter, su *entity.SecurityUpdate) error {
	sm := model.ToSecurityModelFromUpdate(su)
	if err := m.mColl.update(ctx, sf, sm); err != nil {
		return err
	}

	return nil
}

func (m *securityMongo) Get(ctx context.Context, sf *repo.SecurityFilter) (*entity.Security, error) {
	s := new(model.Security)
	if err := m.mColl.get(ctx, &s, sf); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrSecurityNotFound
		}
		return nil, err
	}

	return model.ToSecurityEntity(s), nil
}

func (m *securityMongo) GetMany(ctx context.Context, sf *repo.SecurityFilter) ([]*entity.Security, error) {
	res, err := m.mColl.getMany(ctx, new(model.Security), sf.Paging, sf)
	if err != nil {
		return nil, err
	}

	ess := make([]*entity.Security, 0, len(res))
	for _, r := range res {
		ess = append(ess, model.ToSecurityEntity(r.(*model.Security)))
	}

	return ess, nil
}
