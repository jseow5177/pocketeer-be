package repo

import (
	"context"

	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type SnapshotRepo interface {
	CreateMany(ctx context.Context, sps []*entity.Snapshot) ([]string, error)

	GetMany(ctx context.Context, spf *SnapshotFilter) ([]*entity.Snapshot, error)
}

type SnapshotFilter struct {
	UserID       *string `filter:"user_id"`
	SnapshotType *uint32 `filter:"snapshot_type"`
	TimestampGte *uint64 `filter:"timestamp__gte"`
	TimestampLte *uint64 `filter:"timestamp__lte"`
	Paging       *Paging `filter:"-"`
}

type SnapshotFilterOption = func(spf *SnapshotFilter)

func WithSnapshotType(snapshotType *uint32) SnapshotFilterOption {
	return func(spf *SnapshotFilter) {
		spf.SnapshotType = snapshotType
	}
}

func WithSnapshotTimestampGte(timestampGte *uint64) SnapshotFilterOption {
	return func(spf *SnapshotFilter) {
		spf.TimestampGte = timestampGte
	}
}

func WithSnapshotTimestampLte(timestampLte *uint64) SnapshotFilterOption {
	return func(spf *SnapshotFilter) {
		spf.TimestampLte = timestampLte
	}
}

func WithSnapshotPaging(paging *Paging) SnapshotFilterOption {
	return func(spf *SnapshotFilter) {
		spf.Paging = paging
	}
}

func NewSnapshotFilter(userID string, opts ...SnapshotFilterOption) *SnapshotFilter {
	sf := &SnapshotFilter{
		UserID: goutil.String(userID),
	}
	for _, opt := range opts {
		opt(sf)
	}
	return sf
}

func (f *SnapshotFilter) GetUserID() string {
	if f != nil && f.UserID != nil {
		return *f.UserID
	}
	return ""
}

func (f *SnapshotFilter) GetSnapshotType() uint32 {
	if f != nil && f.SnapshotType != nil {
		return *f.SnapshotType
	}
	return 0
}

func (f *SnapshotFilter) GetTimestampGte() uint64 {
	if f != nil && f.TimestampGte != nil {
		return *f.TimestampGte
	}
	return 0
}

func (f *SnapshotFilter) GetTimestampLte() uint64 {
	if f != nil && f.TimestampLte != nil {
		return *f.TimestampLte
	}
	return 0
}
