package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
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
		if mongo.IsDuplicateKeyError(err) {
			return "", repo.ErrHoldingAlreadyExists
		}
		return "", err
	}
	h.SetHoldingID(goutil.String(id))

	return id, nil
}

func (m *holdingMongo) Update(ctx context.Context, hf *repo.HoldingFilter, hu *entity.HoldingUpdate) error {
	hm := model.ToHoldingModelFromUpdate(hu)
	if err := m.mColl.update(ctx, hf, hm); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return repo.ErrHoldingAlreadyExists
		}
		return err
	}

	return nil
}

func (m *holdingMongo) Get(ctx context.Context, hf *repo.HoldingFilter) (*entity.Holding, error) {
	h := new(model.Holding)
	if err := m.mColl.get(ctx, hf, &h); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.ErrHoldingNotFound
		}
		return nil, err
	}

	return model.ToHoldingEntity(h)
}

func (m *holdingMongo) GetMany(ctx context.Context, hf *repo.HoldingFilter) ([]*entity.Holding, error) {
	res, err := m.mColl.getMany(ctx, hf, nil, new(model.Holding))
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
