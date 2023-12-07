package mongo

import (
	"context"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/dep/repo/mongo/model"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"github.com/jseow5177/pockteer-be/pkg/mongoutil"
)

const snapshotCollName = "snapshot"

type snapshotMongo struct {
	mColl *MongoColl
}

func NewSnapshotMongo(mongo *Mongo) repo.SnapshotRepo {
	return &snapshotMongo{
		mColl: NewMongoColl(mongo, snapshotCollName),
	}
}

func (m *snapshotMongo) CreateMany(ctx context.Context, sps []*entity.Snapshot) ([]string, error) {
	spms := make([]interface{}, 0, len(sps))
	for _, sp := range sps {
		spms = append(spms, model.ToSnapshotModelFromEntity(sp))
	}

	ids, err := m.mColl.createMany(ctx, spms)
	if err != nil {
		return nil, err
	}

	for i, sp := range sps {
		sp.SetSnapshotID(goutil.String(ids[i]))
	}

	return ids, nil
}

func (m *snapshotMongo) GetMany(ctx context.Context, spf *repo.SnapshotFilter) ([]*entity.Snapshot, error) {
	f := mongoutil.BuildFilter(spf)

	res, err := m.mColl.getMany(ctx, new(model.Snapshot), spf.Paging, f)
	if err != nil {
		return nil, err
	}

	sps := make([]*entity.Snapshot, 0, len(res))
	for _, r := range res {
		sps = append(sps, model.ToSnapshotEntity(r.(*model.Snapshot)))
	}

	return sps, nil
}
