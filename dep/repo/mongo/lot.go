package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

const lotCollName = "lot"

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
	l.SetLotID(id)

	return id, nil
}

func (m *lotMongo) Update(ctx context.Context, lf *repo.LotFilter, lu *entity.LotUpdate) error {
	lm := model.ToLotModelFromUpdate(lu)
	if err := m.mColl.update(ctx, lf, lm); err != nil {
		return err
	}

	return nil
}

func (m *lotMongo) Get(ctx context.Context, lf *repo.LotFilter) (*entity.Lot, error) {
	lm := new(model.Lot)
	if err := m.mColl.get(ctx, lf, &lm); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrLotNotFound
		}
		return nil, err
	}

	return model.ToLotEntity(lm), nil
}

func (m *lotMongo) GetMany(ctx context.Context, lf *repo.LotFilter) ([]*entity.Lot, error) {
	res, err := m.mColl.getMany(ctx, lf, nil, new(model.Lot))
	if err != nil {
		return nil, err
	}

	els := make([]*entity.Lot, 0, len(res))
	for _, r := range res {
		els = append(els, model.ToLotEntity(r.(*model.Lot)))
	}

	return els, nil
}
