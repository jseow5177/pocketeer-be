package entity

import (
	"time"

	"github.com/jseow5177/pockteer-be/pkg/goutil"
)

type SnapshotUnit uint32

const (
	SnapshotUnitMonth SnapshotUnit = iota
)

var SnapshotUnits = map[uint32]string{
	uint32(SnapshotUnitMonth): "month",
}

type SnapshotType uint32

const (
	SnapshotTypeInvalid SnapshotType = iota
	SnapshotTypeAccount
)

var SnapshotTypes = map[uint32]string{
	uint32(SnapshotTypeAccount): "snapshot account",
}

type Snapshot struct {
	UserID       *string
	SnapshotID   *string
	SnapshotType *uint32
	Record       *string
	Timestamp    *uint64
	CreateTime   *uint64

	// Only `Period` and `Value` will be returned
	Period *string
	Value  *string
}

type SnapshotOption = func(sp *Snapshot)

func WithSnapshotID(snapshotID *string) SnapshotOption {
	return func(sp *Snapshot) {
		sp.SnapshotID = snapshotID
	}
}

func WithSnapshotRecord(record *string) SnapshotOption {
	return func(sp *Snapshot) {
		sp.Record = record
	}
}

func WithSnapshotTimestamp(timestamp *uint64) SnapshotOption {
	return func(sp *Snapshot) {
		sp.SetTimestamp(timestamp)
	}
}

func WithSnapshotCreateTime(createTime *uint64) SnapshotOption {
	return func(sp *Snapshot) {
		sp.SetCreateTime(createTime)
	}
}

func NewSnapshot(userID string, snapshotType uint32, opts ...SnapshotOption) *Snapshot {
	now := uint64(time.Now().UnixMilli())
	ss := &Snapshot{
		UserID:       goutil.String(userID),
		SnapshotType: goutil.Uint32(snapshotType),
		Record:       goutil.String("{}"),
		Timestamp:    goutil.Uint64(now),
		CreateTime:   goutil.Uint64(now),
	}

	for _, opt := range opts {
		opt(ss)
	}

	return ss
}

func (sp *Snapshot) GetUserID() string {
	if sp != nil && sp.UserID != nil {
		return *sp.UserID
	}
	return ""
}

func (sp *Snapshot) SetUserID(userID *string) {
	sp.UserID = userID
}

func (sp *Snapshot) GetSnapshotID() string {
	if sp != nil && sp.SnapshotID != nil {
		return *sp.SnapshotID
	}
	return ""
}

func (sp *Snapshot) SetSnapshotID(snapshotID *string) {
	sp.SnapshotID = snapshotID
}

func (sp *Snapshot) GetSnapshotType() uint32 {
	if sp != nil && sp.SnapshotType != nil {
		return *sp.SnapshotType
	}
	return 0
}

func (sp *Snapshot) SetSnapshotType(snapshotType *uint32) {
	sp.SnapshotType = snapshotType
}

func (sp *Snapshot) GetTimestamp() uint64 {
	if sp != nil && sp.Timestamp != nil {
		return *sp.Timestamp
	}
	return 0
}

func (sp *Snapshot) SetTimestamp(timestamp *uint64) {
	sp.Timestamp = timestamp
}

func (sp *Snapshot) GetCreateTime() uint64 {
	if sp != nil && sp.CreateTime != nil {
		return *sp.CreateTime
	}
	return 0
}

func (sp *Snapshot) SetCreateTime(createTime *uint64) {
	sp.CreateTime = createTime
}

func (sp *Snapshot) GetPeriod() string {
	if sp != nil && sp.Period != nil {
		return *sp.Period
	}
	return ""
}

func (sp *Snapshot) SetPeriod(period *string) {
	sp.Period = period
}

func (sp *Snapshot) GetValue() string {
	if sp != nil && sp.Value != nil {
		return *sp.Value
	}
	return ""
}

func (sp *Snapshot) SetValue(value *string) {
	sp.Value = value
}
