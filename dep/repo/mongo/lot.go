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

const (
	lotCollName = "lot"
)

type lotMongo struct {
	mColl *MongoColl
}

func NewLotMongo(mongo *Mongo) repo.LotRepo {
	return &lotMongo{
		mColl: NewMongoColl(mongo, lotCollName),
	}
}

func (m *lotMongo) Create(ctx context.Context, l *entity.Lot) (string, error) {
	lm := model.ToLotModelFromEntity(l)
	id, err := m.mColl.create(ctx, lm)
	if err != nil {
		return "", err
	}
	l.SetLotID(goutil.String(id))

	return id, nil
}

func (m *lotMongo) CreateMany(ctx context.Context, ls []*entity.Lot) ([]string, error) {
	lms := make([]interface{}, 0)
	for _, l := range ls {
		lms = append(lms, model.ToLotModelFromEntity(l))
	}
	ids, err := m.mColl.createMany(ctx, lms)
	if err != nil {
		return nil, err
	}

	for i, l := range ls {
		l.SetLotID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *lotMongo) Update(ctx context.Context, lf *repo.LotFilter, lu *entity.LotUpdate) error {
	lm := model.ToLotModelFromUpdate(lu)
	if err := m.mColl.update(ctx, lf, lm); err != nil {
		return err
	}

	return nil
}

func (m *lotMongo) Delete(ctx context.Context, lf *repo.LotFilter) error {
	return m.Update(ctx, lf, entity.NewLotUpdate(
		entity.WithUpdateLotStatus(goutil.Uint32(uint32(entity.LotStatusDeleted))),
	))
}

func (m *lotMongo) Get(ctx context.Context, lf *repo.LotFilter) (*entity.Lot, error) {
	f := mongoutil.BuildFilter(lf)

	lm := new(model.Lot)
	if err := m.mColl.get(ctx, &lm, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrLotNotFound
		}
		return nil, err
	}

	return model.ToLotEntity(lm), nil
}

func (m *lotMongo) GetMany(ctx context.Context, lf *repo.LotFilter) ([]*entity.Lot, error) {
	f := mongoutil.BuildFilter(lf)

	res, err := m.mColl.getMany(ctx, new(model.Lot), nil, f)
	if err != nil {
		return nil, err
	}

	els := make([]*entity.Lot, 0, len(res))
	for _, r := range res {
		els = append(els, model.ToLotEntity(r.(*model.Lot)))
	}

	return els, nil
}
