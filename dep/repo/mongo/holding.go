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

const holdingCollName = "holding"

type holdingMongo struct {
	mColl *MongoColl
}

func NewHoldingMongo(mongo *Mongo) repo.HoldingRepo {
	return &holdingMongo{
		mColl: NewMongoColl(mongo, holdingCollName),
	}
}

func (m *holdingMongo) Create(ctx context.Context, h *entity.Holding) (string, error) {
	hm := model.ToHoldingModelFromEntity(h)
	id, err := m.mColl.create(ctx, hm)
	if err != nil {
		return "", err
	}
	h.SetHoldingID(goutil.String(id))

	return id, nil
}

func (m *holdingMongo) CreateMany(ctx context.Context, hs []*entity.Holding) ([]string, error) {
	hms := make([]interface{}, 0)
	for _, h := range hs {
		hms = append(hms, model.ToHoldingModelFromEntity(h))
	}

	ids, err := m.mColl.createMany(ctx, hms)
	if err != nil {
		return nil, err
	}

	for i, h := range hs {
		h.SetHoldingID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *holdingMongo) Update(ctx context.Context, hf *repo.HoldingFilter, hu *entity.HoldingUpdate) error {
	f := mongoutil.BuildFilter(hf)

	hm := model.ToHoldingModelFromUpdate(hu)
	if err := m.mColl.update(ctx, f, hm); err != nil {
		return err
	}

	return nil
}

func (m *holdingMongo) Get(ctx context.Context, hf *repo.HoldingFilter) (*entity.Holding, error) {
	f := mongoutil.BuildFilter(hf)

	h := new(model.Holding)
	if err := m.mColl.get(ctx, &h, f); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrHoldingNotFound
		}
		return nil, err
	}

	return model.ToHoldingEntity(h)
}

func (m *holdingMongo) Delete(ctx context.Context, hf *repo.HoldingFilter) error {
	return m.Update(ctx, hf, entity.NewHoldingUpdate(
		entity.WithUpdateHoldingStatus(goutil.Uint32(uint32(entity.HoldingStatusDeleted))),
	))
}

func (m *holdingMongo) DeleteMany(ctx context.Context, hf *repo.HoldingFilter) error {
	hm := model.ToHoldingModelFromUpdate(entity.NewHoldingUpdate(
		entity.WithUpdateHoldingStatus(goutil.Uint32(uint32(entity.HoldingStatusDeleted))),
	))

	f := mongoutil.BuildFilter(hf)

	return m.mColl.updateMany(ctx, f, hm)
}

func (m *holdingMongo) GetMany(ctx context.Context, hf *repo.HoldingFilter) ([]*entity.Holding, error) {
	f := mongoutil.BuildFilter(hf)

	res, err := m.mColl.getMany(ctx, new(model.Holding), hf.Paging, f)
	if err != nil {
		return nil, err
	}

	ehs := make([]*entity.Holding, 0, len(res))
	for _, r := range res {
		eh, err := model.ToHoldingEntity(r.(*model.Holding))
		if err != nil {
			return nil, err
		}
		ehs = append(ehs, eh)
	}

	return ehs, nil
}
