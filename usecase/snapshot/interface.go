package snapshot

import (
	"context"
	"time"

	"github.com/jseow5177/pockteer-be/dep/repo"
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type UseCase interface {
	GetAccountSnapshots(ctx context.Context, req *GetAccountSnapshotsRequest) (*GetAccountSnapshotsResponse, error)
}

type GetAccountSnapshotsRequest struct {
	UserID    *string
	AccountID *string
	Unit      *uint32
	Interval  *uint32
}

func (m *GetAccountSnapshotsRequest) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GetAccountSnapshotsRequest) GetAccountID() string {
	if m != nil && m.AccountID != nil {
		return *m.AccountID
	}
	return ""
}

func (m *GetAccountSnapshotsRequest) GetUnit() uint32 {
	if m != nil && m.Unit != nil {
		return *m.Unit
	}
	return 0
}

func (m *GetAccountSnapshotsRequest) GetInterval() uint32 {
	if m != nil && m.Interval != nil {
		return *m.Interval
	}
	return 0
}

func (m *GetAccountSnapshotsRequest) ToAccountFilter() *repo.AccountFilter {
	return repo.NewAccountFilter(
		m.GetUserID(),
		repo.WithAccountID(m.AccountID),
	)
}

func (m *GetAccountSnapshotsRequest) ToSnapshotFilter() *repo.SnapshotFilter {
	now := time.Now()

	t := now
	switch m.GetUnit() {
	case uint32(entity.SnapshotUnitMonth):
		t = now.AddDate(0, -int(m.GetInterval()), 0)
	}

	return repo.NewSnapshotFilter(
		m.GetUserID(),
		repo.WithSnapshotType(goutil.Uint32(uint32(entity.SnapshotTypeAccount))),
		repo.WithSnapshotTimestampGte(goutil.Uint64(uint64(t.UnixMilli()))),
		repo.WithSnapshotTimestampLte(goutil.Uint64(uint64(now.UnixMilli()))),
	)
}

type GetAccountSnapshotsResponse struct {
	Snapshots []*entity.Snapshot
}

func (m *GetAccountSnapshotsResponse) GetSnapshots() []*entity.Snapshot {
	if m != nil && m.Snapshots != nil {
		return m.Snapshots
	}
	return nil
}
