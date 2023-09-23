package model

import (
	"github.com/jseow5177/pockteer-be/entity"
	"github.com/jseow5177/pockteer-be/pkg/goutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Snapshot struct {
	SnapshotID   primitive.ObjectID `bson:"_id,omitempty"`
	UserID       *string            `bson:"user_id,omitempty"`
	SnapshotType *uint32            `bson:"snapshot_type,omitempty"`
	Record       *string            `bson:"record,omitempty"`
	Timestamp    *uint64            `bson:"timestamp,omitempty"`
	CreateTime   *uint64            `bson:"create_time,omitempty"`
}

func ToSnapshotModelFromEntity(ss *entity.Snapshot) *Snapshot {
	if ss == nil {
		return nil
	}

	objID := primitive.NilObjectID
	if primitive.IsValidObjectID(ss.GetSnapshotID()) {
		objID, _ = primitive.ObjectIDFromHex(ss.GetSnapshotID())
	}

	return &Snapshot{
		SnapshotID:   objID,
		UserID:       ss.UserID,
		SnapshotType: ss.SnapshotType,
		Record:       ss.Record,
		Timestamp:    ss.Timestamp,
		CreateTime:   ss.CreateTime,
	}
}

func ToSnapshotEntity(ss *Snapshot) *entity.Snapshot {
	if ss == nil {
		return nil
	}

	return entity.NewSnapshot(
		ss.GetUserID(),
		ss.GetSnapshotType(),
		entity.WithSnapshotID(goutil.String(ss.GetSnapshotID())),
		entity.WithSnapshotRecord(ss.Record),
		entity.WithSnapshotTimestamp(ss.Timestamp),
		entity.WithSnapshotCreateTime(ss.CreateTime),
	)
}

func (ss *Snapshot) GetUserID() string {
	if ss != nil && ss.UserID != nil {
		return *ss.UserID
	}
	return ""
}

func (ss *Snapshot) GetSnapshotID() string {
	if ss != nil {
		return ss.SnapshotID.Hex()
	}
	return ""
}

func (ss *Snapshot) GetSnapshotType() uint32 {
	if ss != nil && ss.SnapshotType != nil {
		return *ss.SnapshotType
	}
	return 0
}

func (ss *Snapshot) GetTimestamp() uint64 {
	if ss != nil && ss.Timestamp != nil {
		return *ss.Timestamp
	}
	return 0
}

func (ss *Snapshot) GetCreateTime() uint64 {
	if ss != nil && ss.CreateTime != nil {
		return *ss.CreateTime
	}
	return 0
}
